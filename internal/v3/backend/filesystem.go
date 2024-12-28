package backend

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dadav/gorge/internal/config"
	"github.com/dadav/gorge/internal/log"
	"github.com/dadav/gorge/internal/model"
	"github.com/dadav/gorge/internal/v3/utils"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/hashicorp/go-version"
)

type FilesystemBackend struct {
	muModules  sync.RWMutex
	Modules    map[string]*gen.Module
	ModulesDir string
	muReleases sync.RWMutex
	Releases   map[string][]*gen.Release
}

var _ Backend = (*FilesystemBackend)(nil)

const (
	defaultVersion = "0.0.0"
	metadataFile   = "metadata.json"
	readmeFile     = "README.md"
	tarGzExt       = ".tar.gz"
)

func NewFilesystemBackend(path string) *FilesystemBackend {
	return &FilesystemBackend{
		Modules:    map[string]*gen.Module{},
		ModulesDir: path,
		Releases:   map[string][]*gen.Release{},
	}
}

func findLatestVersion(releases []gen.ReleaseAbbreviated) string {
	if len(releases) == 0 {
		return defaultVersion
	}

	latest := releases[0].Version
	for _, r := range releases[1:] {
		vVersion, err := version.NewVersion(r.Version)
		if err != nil {
			log.Log.Warnf("invalid version: %s", r.Version)
			continue
		}

		vlatest, err := version.NewVersion(latest)
		if err != nil {
			log.Log.Warnf("invalid version: %s", latest)
			continue
		}

		if vVersion.Compare(vlatest) >= 1 {
			latest = r.Version
		}
	}
	return latest
}

func ReleaseToAbbreviatedRelease(release *gen.Release) *gen.ReleaseAbbreviated {
	return &gen.ReleaseAbbreviated{
		Uri:       release.Uri,
		Slug:      release.Slug,
		Version:   release.Version,
		Supported: release.Supported,
		CreatedAt: release.CreatedAt,
		DeletedAt: release.DeletedAt,
		FileUri:   release.FileUri,
		FileSize:  release.FileSize,
	}
}

func (s *FilesystemBackend) GetAllReleases() ([]*gen.Release, error) {
	s.muReleases.Lock()
	defer s.muReleases.Unlock()
	result := []*gen.Release{}

	for _, v := range s.Releases {
		result = append(result, v...)
	}

	return result, nil
}

func MetadataToRelease(metadata *model.ReleaseMetadata) *gen.Release {
	var releaseMetadataInterface map[string]interface{}
	inrec, _ := json.Marshal(metadata)
	json.Unmarshal(inrec, &releaseMetadataInterface)

	return &gen.Release{
		Slug: fmt.Sprintf("%s-%s", metadata.Name, metadata.Version),
		Uri:  fmt.Sprintf("/v3/releases/%s-%s", metadata.Name, metadata.Version),
		Module: gen.ReleaseModule{
			Name: metadata.Name,
			Slug: metadata.Name,
			Uri:  fmt.Sprintf("/v3/modules/%s", metadata.Name),
			Owner: gen.ModuleOwner{
				Uri:        fmt.Sprintf("/v3/users/%s", metadata.Author),
				Slug:       metadata.Author,
				Username:   metadata.Author,
				GravatarId: "",
			},
			DeprecatedAt: nil,
		},
		Version:   metadata.Version,
		Metadata:  releaseMetadataInterface,
		Tags:      metadata.Tags,
		Supported: false,
		Pdk:       false,
	}
}

func ModuleFromRelease(release *gen.Release) *gen.Module {
	return &gen.Module{
		Uri:            fmt.Sprintf("/v3/modules/%s", release.Module.Name),
		Slug:           release.Module.Slug,
		Name:           strings.Split(release.Module.Slug, "-")[1],
		Downloads:      0,
		CreatedAt:      time.Now().String(),
		UpdatedAt:      time.Now().String(),
		DeprecatedAt:   nil,
		DeprecatedFor:  nil,
		SupersededBy:   gen.ModuleSupersededBy{},
		Supported:      release.Supported,
		Endorsement:    nil,
		ModuleGroup:    "Gorge",
		Premium:        false,
		Owner:          release.Module.Owner,
		CurrentRelease: gen.ModuleCurrentRelease(*release),
		Releases:       []gen.ReleaseAbbreviated{*ReleaseToAbbreviatedRelease(release)},
		FeedbackScore:  0,
	}
}

func (s *FilesystemBackend) AddRelease(releaseData []byte) (*gen.Release, error) {
	s.muModules.Lock()
	s.muReleases.Lock()
	defer s.muModules.Unlock()
	defer s.muReleases.Unlock()

	metadata, readme, err := ReadReleaseMetadataFromBytes(releaseData)
	if err != nil {
		return nil, err
	}

	// Validate metadata.Name to ensure it does not contain path separators or parent directory references
	if strings.Contains(metadata.Name, "/") || strings.Contains(metadata.Name, "\\") || strings.Contains(metadata.Name, "..") {
		return nil, errors.New("invalid module name")
	}

	releaseSlug := fmt.Sprintf("%s-%s", metadata.Name, metadata.Version)
	if !utils.CheckReleaseSlug(releaseSlug) {
		return nil, errors.New("invalid release slug")
	}

	// No need to re-read releases we know of
	for _, release := range s.Releases[metadata.Name] {
		if release.Slug == releaseSlug {
			return release, nil
		}
	}
	release := MetadataToRelease(metadata)

	md5Hash := md5.New()
	sha256Hash := sha256.New()

	bytesReader := bytes.NewReader(releaseData)

	_, err = io.Copy(md5Hash, bytesReader)
	if err != nil {
		return nil, err
	}

	_, err = bytesReader.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(sha256Hash, bytesReader)
	if err != nil {
		return nil, err
	}

	md5Sum := fmt.Sprintf("%x", md5Hash.Sum(nil))
	sha256Sum := fmt.Sprintf("%x", sha256Hash.Sum(nil))

	release.FileMd5 = md5Sum
	release.FileSha256 = sha256Sum
	release.FileUri = fmt.Sprintf("/v3/files/%s.tar.gz", releaseSlug)
	release.FileSize = int32(len(releaseData))
	release.Readme = readme
	release.License = metadata.License

	var module *gen.Module
	var ok bool
	if module, ok = s.Modules[metadata.Name]; !ok {
		module = ModuleFromRelease(release)
		s.Modules[metadata.Name] = module
	} else {
		module.Releases = append(module.Releases, *ReleaseToAbbreviatedRelease(release))
		if findLatestVersion(module.Releases) == release.Version {
			module.CurrentRelease = gen.ModuleCurrentRelease(*release)
		}
	}
	s.Releases[metadata.Name] = append(s.Releases[metadata.Name], release)

	releaseFile := fmt.Sprintf("%s.tar.gz", releaseSlug)
	releaseFilePath := fmt.Sprintf("%s/%s/%s", config.ModulesDir, metadata.Name, releaseFile)
	moduleDir := filepath.Join(config.ModulesDir, metadata.Name)
	if _, err := os.Stat(moduleDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(moduleDir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if _, err := os.Stat(releaseFilePath); os.IsNotExist(err) {
		// Only write if file does not exist
		err = os.WriteFile(releaseFilePath, releaseData, 0644)
		if err != nil {
			return nil, err
		}
	}

	return release, nil
}

func (s *FilesystemBackend) GetAllModules() ([]*gen.Module, error) {
	s.muModules.Lock()
	defer s.muModules.Unlock()

	result := []*gen.Module{}

	for _, v := range s.Modules {
		result = append(result, v)
	}

	return result, nil
}

func (s *FilesystemBackend) GetModuleBySlug(slug string) (*gen.Module, error) {
	s.muModules.Lock()
	defer s.muModules.Unlock()
	if module, ok := s.Modules[slug]; !ok {
		return nil, errors.New("module not found")
	} else {
		return module, nil
	}
}

func (s *FilesystemBackend) GetReleaseBySlug(slug string) (*gen.Release, error) {
	s.muReleases.Lock()
	defer s.muReleases.Unlock()
	for _, moduleReleases := range s.Releases {
		for _, release := range moduleReleases {
			if release.Slug == slug {
				return release, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

func (s *FilesystemBackend) DeleteModuleBySlug(slug string) error {
	s.muModules.Lock()
	s.muReleases.Lock()
	defer s.muModules.Unlock()
	defer s.muReleases.Unlock()

	modulePath := filepath.Join(config.ModulesDir, slug)
	err := os.RemoveAll(modulePath)
	if err != nil {
		return err
	}

	delete(s.Releases, slug)
	delete(s.Modules, slug)

	return nil
}

func (s *FilesystemBackend) DeleteReleaseBySlug(slug string) error {
	s.muModules.Lock()
	s.muReleases.Lock()
	defer s.muModules.Unlock()
	defer s.muReleases.Unlock()

	for module, releases := range s.Releases {
		newReleases := []*gen.Release{}
		for _, release := range releases {
			if release.Slug == slug {
				releasePath := filepath.Join(config.ModulesDir, release.Module.Slug, fmt.Sprintf("%s.tar.gz", slug))
				err := os.Remove(releasePath)
				if err != nil {
					return err
				}
			} else {
				newReleases = append(newReleases, release)
			}
		}
		s.Releases[module] = newReleases

		newAbbrReleases := []gen.ReleaseAbbreviated{}
		for _, abbrRelease := range s.Modules[module].Releases {
			if abbrRelease.Slug != slug {
				newAbbrReleases = append(newAbbrReleases, abbrRelease)
			}
		}
		s.Modules[module].Releases = newAbbrReleases

		if s.Modules[module].CurrentRelease.Slug == slug {
			latestReleaseVersion := findLatestVersion(s.Modules[module].Releases)
			for _, modRelease := range s.Releases[module] {
				if modRelease.Version == latestReleaseVersion {
					s.Modules[module].CurrentRelease = gen.ModuleCurrentRelease(*modRelease)
					break
				}
			}
		}
	}

	return nil
}

func (s *FilesystemBackend) LoadModules() error {
	// don't overwrite data we already have on modules and releases
	if s.Modules == nil {
		s.Modules = make(map[string]*gen.Module)
	}
	if s.Releases == nil {
		s.Releases = make(map[string][]*gen.Release)
	}

	err := filepath.Walk(s.ModulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".tar.gz") {
			return nil
		}

		log.Log.Debugf("Reading %s\n", path)
		releaseBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = s.AddRelease(releaseBytes)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func ReadReleaseMetadataFromFile(path string) (*model.ReleaseMetadata, string, error) {
	var jsonData bytes.Buffer
	var releaseMetadata model.ReleaseMetadata
	readme := new(strings.Builder)

	f, err := os.Open(path)
	if err != nil {
		return nil, readme.String(), err
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		return nil, readme.String(), err
	}

	tarReader := tar.NewReader(g)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, readme.String(), err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		switch filepath.Base(header.Name) {
		case "metadata.json":
			_, err = io.Copy(&jsonData, tarReader)
			if err != nil {
				return nil, readme.String(), err
			}

			if err := json.Unmarshal(jsonData.Bytes(), &releaseMetadata); err != nil {
				return nil, readme.String(), err
			}

		case "README.md":
			_, err = io.Copy(readme, tarReader)
			if err != nil {
				return nil, readme.String(), err
			}
		default:
			continue
		}
	}
	return &releaseMetadata, readme.String(), nil
}

func ReadReleaseMetadataFromBytes(data []byte) (*model.ReleaseMetadata, string, error) {
	if len(data) == 0 {
		return nil, "", errors.New("empty data provided")
	}

	var jsonData bytes.Buffer
	var releaseMetadata model.ReleaseMetadata
	readme := new(strings.Builder)

	f := bytes.NewReader(data)
	g, err := gzip.NewReader(f)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer g.Close()

	tarReader := tar.NewReader(g)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, readme.String(), err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		switch filepath.Base(header.Name) {
		case "metadata.json":
			_, err = io.Copy(&jsonData, tarReader)
			if err != nil {
				return nil, readme.String(), err
			}

			if err := json.Unmarshal(jsonData.Bytes(), &releaseMetadata); err != nil {
				return nil, readme.String(), err
			}

			if !utils.CheckModuleSlug(releaseMetadata.Name) {
				return nil, readme.String(), errors.New("invalid module name")
			}
		case "README.md":
			_, err = io.Copy(readme, tarReader)
			if err != nil {
				return nil, readme.String(), err
			}
		default:
			continue
		}
	}
	return &releaseMetadata, readme.String(), nil
}

func (b *FilesystemBackend) UpdateModule(module *gen.Module) error {
	// Convert module to JSON
	data, err := json.Marshal(module)
	if err != nil {
		return err
	}

	// Write to file
	filename := filepath.Join(b.ModulesDir, module.Slug+".json")
	return os.WriteFile(filename, data, 0644)
}
