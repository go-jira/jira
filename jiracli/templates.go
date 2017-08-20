package jiracli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	yaml "gopkg.in/coryb/yaml.v2"

	"github.com/mgutz/ansi"
)

func findTemplate(name string) ([]byte, error) {
	if file, err := findClosestParentPath(filepath.Join(".jira.d", "templates", name)); err == nil {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func (jc *JiraCli) getTemplate(name string) (string, error) {
	if _, err := os.Stat(name); err == nil {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	b, err := findTemplate(name)
	if err != nil {
		return "", err
	} else if b != nil {
		return string(b), nil
	}
	if s, ok := allTemplates[name]; ok {
		return s, nil
	}
	return "", fmt.Errorf("No Template found for %q", name)
}

func (jc *JiraCli) tmpTemplate(templateName string, data interface{}) (string, error) {
	tmpFile, err := jc.tmpYml(templateName)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	return tmpFile.Name(), jc.runTemplate(templateName, data, tmpFile)
}

func (jc *JiraCli) runTemplate(templateName string, data interface{}, out io.Writer) error {

	templateContent, err := jc.getTemplate(templateName)
	if err != nil {
		return err
	}

	if out == nil {
		out = os.Stdout
	}

	funcs := map[string]interface{}{
		"toJson": func(content interface{}) (string, error) {
			bytes, err := json.MarshalIndent(content, "", "    ")
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
		"append": func(more string, content interface{}) (string, error) {
			switch value := content.(type) {
			case string:
				return string(append([]byte(content.(string)), []byte(more)...)), nil
			case []byte:
				return string(append(content.([]byte), []byte(more)...)), nil
			default:
				return "", fmt.Errorf("Unknown type: %s", value)
			}
		},
		"indent": func(spaces int, content string) string {
			indent := make([]rune, spaces+1, spaces+1)
			indent[0] = '\n'
			for i := 1; i < spaces+1; i++ {
				indent[i] = ' '
			}

			lineSeps := []rune{'\n', '\u0085', '\u2028', '\u2029'}
			for _, sep := range lineSeps {
				indent[0] = sep
				content = strings.Replace(content, string(sep), string(indent), -1)
			}
			return content

		},
		"comment": func(content string) string {
			lineSeps := []rune{'\n', '\u0085', '\u2028', '\u2029'}
			for _, sep := range lineSeps {
				content = strings.Replace(content, string(sep), string([]rune{sep, '#', ' '}), -1)
			}
			return content
		},
		"color": func(color string) string {
			return ansi.ColorCode(color)
		},
		"split": func(sep string, content string) []string {
			return strings.Split(content, sep)
		},
		"join": func(sep string, content []interface{}) string {
			vals := make([]string, len(content))
			for i, v := range content {
				vals[i] = v.(string)
			}
			return strings.Join(vals, sep)
		},
		"abbrev": func(max int, content string) string {
			if len(content) > max {
				var buffer bytes.Buffer
				buffer.WriteString(content[:max-3])
				buffer.WriteString("...")
				return buffer.String()
			}
			return content
		},
		"rep": func(count int, content string) string {
			var buffer bytes.Buffer
			for i := 0; i < count; i++ {
				buffer.WriteString(content)
			}
			return buffer.String()
		},
		"age": func(content string) (string, error) {
			return fuzzyAge(content)
		},
		"dateFormat": func(format string, content string) (string, error) {
			return dateFormat(format, content)
		},
	}

	// HACK HACK HACK: convert data formats to json for backwards compatibilty with templates
	var rawData interface{}
	if jsonData, err := json.Marshal(data); err != nil {
		return err
	} else {
		defer func(mapType, iface reflect.Type) {
			yaml.DefaultMapType = mapType
			yaml.IfaceType = iface
		}(yaml.DefaultMapType, yaml.IfaceType)

		yaml.DefaultMapType = reflect.TypeOf(map[string]interface{}{})
		yaml.IfaceType = yaml.DefaultMapType.Elem()

		if err := yaml.Unmarshal(jsonData, &rawData); err != nil {
			return err
		}
	}
	// rawData, err = yamlFixup(rawData)
	// if err != nil {
	// 	return err
	// }

	tmpl, err := template.New("template").Funcs(funcs).Parse(templateContent)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(out, rawData); err != nil {
		return err
	}
	return nil
}

var allTemplates = map[string]string{
	"component-add":  defaultComponentAddTemplate,
	"debug":          defaultDebugTemplate,
	"fields":         defaultDebugTemplate,
	"editmeta":       defaultDebugTemplate,
	"transmeta":      defaultDebugTemplate,
	"createmeta":     defaultDebugTemplate,
	"issuelinktypes": defaultDebugTemplate,
	"list":           defaultListTemplate,
	"table":          defaultTableTemplate,
	"view":           defaultViewTemplate,
	"edit":           defaultEditTemplate,
	"transitions":    defaultTransitionsTemplate,
	"components":     defaultComponentsTemplate,
	"issuetypes":     defaultIssuetypesTemplate,
	"create":         defaultCreateTemplate,
	"subtask":        defaultSubtaskTemplate,
	"comment":        defaultCommentTemplate,
	"transition":     defaultTransitionTemplate,
	"request":        defaultDebugTemplate,
	"worklog":        defaultWorklogTemplate,
	"worklogs":       defaultWorklogsTemplate,
}

const defaultDebugTemplate = "{{ . | toJson}}\n"

const defaultListTemplate = "{{ range .issues }}{{ .key | append \":\" | printf \"%-12s\"}} {{ .fields.summary }}\n{{ end }}"

const defaultTableTemplate = `+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
| {{ "Issue" | printf "%-14s" }} | {{ "Summary" | printf "%-55s" }} | {{ "Priority" | printf "%-12s" }} | {{ "Status" | printf "%-12s" }} | {{ "Age" | printf "%-10s" }} | {{ "Reporter" | printf "%-12s" }} | {{ "Assignee" | printf "%-12s" }} |
+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
{{ range .issues -}}
  | {{ .key | printf "%-14s"}} | {{ .fields.summary | abbrev 55 | printf "%-55s" }} | {{.fields.priority.name | printf "%-12s" }} | {{.fields.status.name | printf "%-12s" }} | {{.fields.created | age | printf "%-10s" }} | {{if .fields.reporter}}{{ .fields.reporter.name | printf "%-12s"}}{{else}}<unassigned>{{end}} | {{if .fields.assignee }}{{.fields.assignee.name | printf "%-12s" }}{{else}}<unassigned> {{end}} |
{{ end -}}
+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
`

const defaultViewTemplate = `{{/* view template */ -}}
issue: {{ .key }}
{{if .fields.created -}}
created: {{ .fields.created | age }} ago
{{end -}}
{{if .fields.status -}}
status: {{ .fields.status.name }}
{{end -}}
summary: {{ .fields.summary }}
project: {{ .fields.project.key }}
{{if .fields.components -}}
components: {{ range .fields.components }}{{ .name }} {{end}}
{{end -}}
{{if .fields.issuetype -}}
issuetype: {{ .fields.issuetype.name }}
{{end -}}
{{if .fields.assignee -}}
assignee: {{ .fields.assignee.name }}
{{end -}}
reporter: {{ if .fields.reporter }}{{ .fields.reporter.name }}{{end}}
{{if .fields.customfield_10110 -}}
watchers: {{ range .fields.customfield_10110 }}{{ .name }} {{end}}
{{end -}}
{{if .fields.issuelinks -}}
blockers: {{ range .fields.issuelinks }}{{if .outwardIssue}}{{ .outwardIssue.key }}[{{.outwardIssue.fields.status.name}}]{{end}}{{end}}
depends: {{ range .fields.issuelinks }}{{if .inwardIssue}}{{ .inwardIssue.key }}[{{.inwardIssue.fields.status.name}}]{{end}}{{end}}
{{end -}}
{{if .fields.priority -}}
priority: {{ .fields.priority.name }}
{{end -}}
{{if .fields.votes -}}
votes: {{ .fields.votes.votes}}
{{end -}}
{{if .fields.labels -}}
labels: {{ join ", " .fields.labels }}
{{end -}}
description: |
  {{ or .fields.description "" | indent 2 }}
{{if .fields.comment.comments}}
comments:
{{ range .fields.comment.comments }}  - | # {{.author.name}}, {{.created | age}} ago
    {{ or .body "" | indent 4}}
{{end}}
{{end -}}
`
const defaultEditTemplate = `{{/* edit template */ -}}
# issue: {{ .key }}
update:
  comment:
    - add: 
        body: |~
          {{ or .overrides.comment "" | indent 10 }}
fields:
  summary: {{ or .overrides.summary .fields.summary }}
{{- if and .meta.fields.components .meta.fields.components.allowedValues }}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{if .overrides.components }}{{ range (split "," .overrides.components)}}
    - name: {{.}}{{end}}{{else}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}{{end}}{{end}}
{{- if .meta.fields.assignee}}
  assignee:
    name: {{ if .overrides.assignee }}{{.overrides.assignee}}{{else}}{{if .fields.assignee }}{{ .fields.assignee.name }}{{end}}{{end}}{{end}}
{{- if .meta.fields.reporter}}
  reporter:
    name: {{ if .overrides.reporter }}{{ .overrides.reporter }}{{else if .fields.reporter}}{{ .fields.reporter.name }}{{end}}{{end}}
{{- if .meta.fields.customfield_10110}}
  # watchers
  customfield_10110: {{ range .fields.customfield_10110 }}
    - name: {{ .name }}{{end}}{{if .overrides.watcher}}
    - name: {{ .overrides.watcher}}{{end}}{{end}}
{{- if .meta.fields.priority }}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority .fields.priority.name }}{{end}}
  description: |~
    {{ or .overrides.description (or .fields.description "") | indent 4 }}
# comments:
# {{ range .fields.comment.comments }}  - | # {{.author.name}}, {{.created | age}} ago
#     {{ or .body "" | indent 4 | comment}}
# {{end}}
`
const defaultTransitionsTemplate = `{{ range .transitions }}{{.id }}: {{.name}}
{{end}}`

const defaultComponentsTemplate = `{{ range . }}{{.id }}: {{.name}}
{{end}}`

const defaultComponentAddTemplate = `{{/* compoinent add template */ -}}
project: {{or .project ""}}
name: {{or .name ""}}
description: {{or .description ""}}
leadUserName: {{or .leadUserName ""}}
`

const defaultIssuetypesTemplate = `{{/* issuetypes template */ -}}
{{ range .issuetypes }}{{color "+bh"}}{{.name | append ":" | printf "%-13s" }}{{color "reset"}} {{.description}}
{{end}}`

const defaultCreateTemplate = `{{/* create template */ -}}
fields:
  project:
    key: {{ or .overrides.project "" }}
  issuetype:
    name: {{ or .overrides.issuetype "" }}
  summary: {{ or .overrides.summary "" }}{{if .meta.fields.priority.allowedValues}}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority ""}}{{end}}{{if .meta.fields.components.allowedValues}}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{ range split "," (or .overrides.components "")}}
    - name: {{ . }}{{end}}{{end}}
  description: |~
    {{ or .overrides.description "" | indent 4 }}{{if .meta.fields.assignee}}
  assignee:
    name: {{ or .overrides.assignee "" }}{{end}}{{if .meta.fields.reporter}}
  reporter:
    name: {{ or .overrides.reporter .overrides.user }}{{end}}{{if .meta.fields.customfield_10110}}
  # watchers
  customfield_10110: {{ range split "," (or .overrides.watchers "")}}
    - name: {{.}}{{end}}
    - name:{{end}}`

const defaultSubtaskTemplate = `{{/* create subtask template */ -}}
fields:
  project:
    key: {{ .parent.fields.project.key }}
  summary: {{ or .overrides.summary "" }}{{if .meta.fields.priority.allowedValues}}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority ""}}{{end}}{{if .meta.fields.components.allowedValues}}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{ range split "," (or .overrides.components "")}}
    - name: {{ . }}{{end}}{{end}}
  description: |~
    {{ or .overrides.description "" | indent 4 }}{{if .meta.fields.assignee}}
  assignee:
    name: {{ or .overrides.assignee "" }}{{end}}{{if .meta.fields.reporter}}
  reporter:
    name: {{ or .overrides.reporter .overrides.user }}{{end}}{{if .meta.fields.customfield_10110}}
  # watchers
  customfield_10110: {{ range split "," (or .overrides.watchers "")}}
    - name: {{.}}{{end}}
    - name:{{end}}
  issuetype:
    name: Sub-task
  parent:
    key: {{ .parent.key }}`

const defaultCommentTemplate = `body: |~
  {{ or .overrides.comment "" | indent 2 }}
`

const defaultTransitionTemplate = `{{/* transition template */ -}}
{{- if .meta.fields.comment }}
update:
  comment:
    - add: 
        body: |~
          {{ or .overrides.comment "" | indent 10 }}
{{- end -}}
fields:
{{- if .meta.fields.assignee}}
  assignee:
    name: {{if .overrides.assignee}}{{.overrides.assignee}}{{else}}{{if .fields.assignee}}{{.fields.assignee.name}}{{end}}{{end}}
{{- end -}}
{{if .meta.fields.components}}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{if .overrides.components }}{{ range (split "," .overrides.components)}}
    - name: {{.}}{{end}}{{else}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}{{end}}
{{- end -}}
{{if .meta.fields.description}}
  description: {{or .overrides.description .fields.description }}
{{- end -}}
{{if .meta.fields.fixVersions -}}
  {{if .meta.fields.fixVersions.allowedValues}}
  fixVersions: # Values: {{ range .meta.fields.fixVersions.allowedValues }}{{.name}}, {{end}}{{if .overrides.fixVersions}}{{ range (split "," .overrides.fixVersions)}}
    - name: {{.name}}{{end}}{{else}}{{range .fields.fixVersions}}
    - name: {{.name}}{{end}}{{end}}
  {{- end -}}
{{- end -}}
{{if .meta.fields.issuetype}}
  issuetype: # Values: {{ range .meta.fields.issuetype.allowedValues }}{{.name}}, {{end}}
    name: {{if .overrides.issuetype}}{{.overrides.issuetype}}{{else}}{{if .fields.issuetype}}{{.fields.issuetype.name}}{{end}}{{end}}
{{- end -}}
{{if .meta.fields.labels}}
  labels: {{range .fields.labels}}
    - {{.}}{{end}}{{if .overrides.labels}}{{range (split "," .overrides.labels)}}
    - {{.}}{{end}}{{end}}
{{- end -}}
{{if .meta.fields.priority}}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority "unassigned" }}
{{- end -}}
{{if .meta.fields.reporter}}
  reporter:
    name: {{if .overrides.reporter}}{{.overrides.reporter}}{{else}}{{if .fields.reporter}}{{.fields.reporter.name}}{{end}}{{end}}
{{- end -}}
{{if .meta.fields.resolution}}
  resolution: # Values: {{ range .meta.fields.resolution.allowedValues }}{{.name}}, {{end}}
    name: {{if .overrides.resolution}}{{.overrides.resolution}}{{else if .fields.resolution}}{{.fields.resolution.name}}{{else}}{{or .overrides.defaultResolution "Fixed"}}{{end}}
{{- end -}}
{{if .meta.fields.summary}}
  summary: {{or .overrides.summary .fields.summary}}
{{- end -}}
{{if .meta.fields.versions.allowedValues}}
  versions: # Values: {{ range .meta.fields.versions.allowedValues }}{{.name}}, {{end}}{{if .overrides.versions}}{{ range (split "," .overrides.versions)}}
    - name: {{.}}{{end}}{{else}}{{range .fields.versions}}
    - name: {{.}}{{end}}{{end}}
{{- end}}
transition:
  id: {{ .transition.id }}
  name: {{ .transition.name }}
`

const defaultWorklogTemplate = `{{/* worklog template */ -}}
# issue: {{ .issue }}
comment: |~
  {{ or .comment "" }}
timeSpent: {{ or .timeSpent "" }}
started:
`

const defaultWorklogsTemplate = `{{/* worklogs template */ -}}
{{ range .worklogs }}- # {{.author.name}}, {{.created | age}} ago
  comment: {{ or .comment "" }}
  timeSpent: {{ .timeSpent }}

{{end}}`
