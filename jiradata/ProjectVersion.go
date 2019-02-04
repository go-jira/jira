package jiradata

type ProjectVersion struct {
	Self            string `json:"self,omitempty" yaml:"self,omitempty"`
	ID              string `json:"id,omitempty" yaml:"id,omitempty"`
	Description     string `json:"description,omitempty" yaml:"description,omitempty"`
	Name            string `json:"name,omitempty" yaml:"name,omitempty"`
	Archived        bool   `json:"archived,omitempty" yaml:archived,omitempty"`
	Released        bool   `json:"released,omitempty" yaml:released,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty" yaml:"releaseDate,omitempty"`
	UserReleaseDate string `json:"userReleaseDate,omitempty" yaml:"userReleaseDate,omitempty"`
	ProjectID       int    `json:"projectId,omitempty" yaml:"projectId,omitempty"`
}
