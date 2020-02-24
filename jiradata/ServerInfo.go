package jiradata

type ServerInfo struct {
	BaseURL        string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
	BuildDate      string `json:"buildDate,omitempty" yaml:"buildDate,omitempty"`
	BuildNumber    int    `json:"buildNumber,omitempty" yaml:"buildNumber,omitempty"`
	DeploymentType string `json:"deploymentType,omitempty" yaml:"deploymentType,omitempty"`
	SCMInfo        string `json:"scmInfo,omitempty" yaml:"scmInfo,omitempty"`
	ServerTime     string `json:"serverTime,omitempty" yaml:"serverTime,omitempty"`
	ServerTitle    string `json:"serverTitle,omitempty" yaml:"serverTitle,omitempty"`
	Version        string `json:"version,omitempty" yaml:"version,omitempty"`
	VersionNumbers []int  `json:"versionNumbers,omitempty" yaml:"versionNumbers,omitempty"`
}
