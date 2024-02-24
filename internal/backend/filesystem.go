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

	"github.com/dadav/gorge/internal/config"
	"github.com/dadav/gorge/internal/model"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"golang.org/x/mod/semver"
)

type FilesystemBackend struct {
	muModules  sync.RWMutex
	Modules    map[string]*gen.Module
	ModulesDir string
	muReleases sync.RWMutex
	Releases   map[string][]*gen.Release
}

var _ Backend = (*FilesystemBackend)(nil)

func NewFilesystemBackend(path string) *FilesystemBackend {
	return &FilesystemBackend{
		ModulesDir: path,
	}
}

func findLatestVersion(releases []gen.ReleaseAbbreviated) string {
	latest := "0.0.0"
	for i, r := range releases {
		if i == 0 {
			latest = r.Version
			continue
		}

		if semver.Compare(r.Version, latest) >= 1 {
			latest = r.Version
		}
	}
	return latest
}

func currentReleaseToAbbreviatedRelease(release *gen.ModuleCurrentRelease) *gen.ReleaseAbbreviated {
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

func (s *FilesystemBackend) GetAllReleases() []*gen.Release {
	s.muReleases.Lock()
	defer s.muReleases.Unlock()
	result := []*gen.Release{}

	for _, v := range s.Releases {
		result = append(result, v...)
	}

	return result
}

func (s *FilesystemBackend) AddRelease(name, version string, data []byte) error {
	releaseFile := fmt.Sprintf("%s-%s.tar.gz", name, version)
	cacheFileDir := filepath.Join(config.ModulesDir, name)
	if _, err := os.Stat(cacheFileDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(cacheFileDir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	cacheFile := filepath.Join(cacheFileDir, releaseFile)
	_, err := os.Stat(cacheFile)
	if err == nil {
		return nil
	}

	err = os.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *FilesystemBackend) GetAllModules() []*gen.Module {
	s.muModules.Lock()
	defer s.muModules.Unlock()

	result := []*gen.Module{}

	for _, v := range s.Modules {
		result = append(result, v)
	}

	return result
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
	return nil, errors.New("release not found")
}

func (s *FilesystemBackend) LoadModules() error {
	s.muModules.Lock()
	s.muReleases.Lock()
	defer s.muModules.Unlock()
	defer s.muReleases.Unlock()

	s.Modules = make(map[string]*gen.Module)
	s.Releases = make(map[string][]*gen.Release)

	err := filepath.Walk(s.ModulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".tar.gz") {
			return nil
		}

		releaseMetadata, releaseReadme, err := ReadReleaseMetadata(path)
		if err != nil {
			return err
		}

		moduleSlug := releaseMetadata.Name
		moduleName := strings.TrimPrefix(releaseMetadata.Name, fmt.Sprintf("%s-", releaseMetadata.Author))
		releaseSlug := fmt.Sprintf("%s-%s", releaseMetadata.Name, releaseMetadata.Version)
		releasePath := fmt.Sprintf("/v3/files/%s.tar.gz", releaseSlug)

		var releaseMetadataInterface map[string]interface{}
		inrec, _ := json.Marshal(releaseMetadata)
		json.Unmarshal(inrec, &releaseMetadataInterface)

		md5Hash := md5.New()
		sha256Hash := sha256.New()

		releaseFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer releaseFile.Close()

		_, err = io.Copy(md5Hash, releaseFile)
		if err != nil {
			return err
		}

		_, err = releaseFile.Seek(0, 0)
		if err != nil {
			return err
		}

		_, err = io.Copy(sha256Hash, releaseFile)
		if err != nil {
			return err
		}

		md5Sum := fmt.Sprintf("%x", md5Hash.Sum(nil))
		sha256Sum := fmt.Sprintf("%x", sha256Hash.Sum(nil))
		owner := gen.ModuleOwner{
			Uri:      fmt.Sprintf("/v3/users/%s", releaseMetadata.Author),
			Slug:     releaseMetadata.Author,
			Username: releaseMetadata.Author,
		}

		newRelease := gen.Release{
			Uri:  fmt.Sprintf("/%s/releases/%s", config.ApiVersion, releaseSlug),
			Slug: releaseMetadata.Name,
			Module: gen.ReleaseModule{
				Uri:   fmt.Sprintf("/v3/modules/%s", moduleSlug),
				Slug:  moduleSlug,
				Name:  moduleName,
				Owner: owner,
			},
			Version:    releaseMetadata.Version,
			Metadata:   releaseMetadataInterface,
			Tags:       releaseMetadata.Tags,
			FileUri:    releasePath,
			FileSize:   int32(info.Size()),
			FileMd5:    md5Sum,
			FileSha256: sha256Sum,
			Readme:     releaseReadme,
			License:    releaseMetadata.License,
		}

		currentRelease := gen.ModuleCurrentRelease(newRelease)

		if module, ok := s.Modules[moduleSlug]; !ok {
			newModule := gen.Module{
				Uri:            fmt.Sprintf("/%s/modules/%s", config.ApiVersion, releaseMetadata.Name),
				Slug:           releaseMetadata.Name,
				Name:           moduleName,
				Owner:          owner,
				CurrentRelease: gen.ModuleCurrentRelease(newRelease),
				Releases:       []gen.ReleaseAbbreviated{*currentReleaseToAbbreviatedRelease(&currentRelease)},
			}

			s.Modules[moduleSlug] = &newModule
			s.Releases[moduleSlug] = []*gen.Release{&newRelease}
		} else {
			s.Releases[moduleSlug] = append(s.Releases[moduleSlug], &newRelease)
			module.Releases = append(module.Releases, *currentReleaseToAbbreviatedRelease(&currentRelease))
			if findLatestVersion(module.Releases) == currentRelease.Version {
				module.CurrentRelease = currentRelease
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ReadReleaseMetadata(path string) (*model.ReleaseMetadata, string, error) {
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
