package model

type SupportedOS struct {
	Name     string   `json:"operatingsystem"`
	Releases []string `json:"operatingsystemrelease,omitempty"`
}

type ModuleDependency struct {
	Name               string `json:"name"`
	VersionRequirement string `json:"version_requirement,omitempty"`
}

type ModuleRequirement ModuleDependency

type ReleaseMetadata struct {
	Name                   string              `json:"name"`
	Version                string              `json:"version"`
	Author                 string              `json:"author"`
	License                string              `json:"license"`
	Summary                string              `json:"summary"`
	Source                 string              `json:"source"`
	Dependencies           []ModuleDependency  `json:"dependencies"`
	Requirements           []ModuleRequirement `json:"requirements,omitempty"`
	ProjectUrl             string              `json:"project_url,omitempty"`
	IssuesUrl              string              `json:"issues_url,omitempty"`
	OperatingsystemSupport []SupportedOS       `json:"operatingsystem_support,omitempty"`
	Tags                   []string            `json:"tags,omitempty"`
}
