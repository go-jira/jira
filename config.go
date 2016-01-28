package jira

type StoredQuery struct {
	Description string `yaml:"description"`
	JQL         string `yaml:"jql"`
	Template    string `yaml:"template"`
}

type Options struct {
	// general options
	Browse   bool   `env:"JIRA_BROWSE"      yaml:"browse"`
	Endpoint string `env:"JIRA_ENDPOINT"    yaml:"endpoint"`
	Insecure bool   `env:"JIRA_INSECURE"    yaml:"insecure"`
	Template string `env:"JIRA_TEMPLATE"    yaml:"template"`
	User     string `env:"JIRA_USER"        yaml:"user"`

	// query options
	Assignee    string `env:"JIRA_ASSIGNEE"    yaml:"assignee"`
	Component   string `env:"JIRA_COMPONENT"   yaml:"component"`
	QueryFields string `env:"JIRA_QUERYFIELDS" yaml:"queryfields"`
	IssueType   string `env:"JIRA_ISSUETYPE"   yaml:"issuetype"`
	MaxResults  uint64 `env:"JIRA_MAX_RESULTS" yaml:"max_results"` // aka "limit"
	Project     string `env:"JIRA_PROJECT"     yaml:"project"`
	Query       string `env:"JIRA_QUERY"       yaml:"query"`
	Reporter    string `env:"JIRA_REPORTER"    yaml:"reporter"`
	Sort        string `env:"JIRA_SORT"        yaml:"sort"`
	Watcher     string `env:"JIRA_WATCHER"     yaml:"watcher"`

	// edit options
	Comment string `env:"JIRA_COMMENT"     yaml:"comment"`

	// create options
	// IssueType
	// Comment

	// command options
	Directory string `env:"JIRA_DIRECTORY"   yaml:"directory"`

	// misc options
	Edit       *bool  `env:"JIRA_EDIT"        yaml:"edit,omitempty"`
	NoEdit     bool   `env:"JIRA_NOEDIT"      yaml:"noedit"`
	DryRun     bool   `env:"JIRA_DRYRUN"      yaml:"dryrun"`
	Quiet      bool   `env:"JIRA_QUIET"       yaml:"quiet"`
	Command    string `env:"JIRA_COMMAND"     yaml:"command"`
	Resolution string `env:"JIRA_RESOLUTION"  yaml:"resolution"`
	Editor     string `env:"JIRA_EDITOR"      yaml:"editor"`
	SaveFile   string `env:"JIRA_SAVEFILE"    yaml:"saveFile"`
	Method     string `env:"JIRA_METHOD"      yaml:"method"`

	// stored queries
	Queries map[string]StoredQuery `yaml:"queries"`
	
	// overrides
	Overrides map[string]interface{} `yaml:"overrides"`
}
