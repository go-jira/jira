package cli

var all_templates = map[string]string{
	"debug":          default_debug_template,
	"fields":         default_debug_template,
	"editmeta":       default_debug_template,
	"transmeta":      default_debug_template,
	"createmeta":     default_debug_template,
	"issuelinktypes": default_debug_template,
	"list":           default_list_template,
	"view":           default_view_template,
	"edit":           default_edit_template,
	"transitions":    default_transitions_template,
	"issuetypes":     default_issuetypes_template,
	"create":         default_create_template,
	"comment":        default_comment_template,
}

const default_debug_template = "{{ . | toJson}}\n"

const default_list_template = "{{ range .issues }}{{ .key | append \":\" | printf \"%-12s\"}} {{ .fields.summary }}\n{{ end }}"

const default_view_template = `issue: {{ .key }}
status: {{ .fields.status.name }}
summary: {{ .fields.summary }}
project: {{ .fields.project.key }}
components: {{ range .fields.components }}{{ .name }} {{end}}
issuetype: {{ .fields.issuetype.name }}
assignee: {{ if .fields.assignee }}{{ .fields.assignee }}{{end}}
reporter: {{ .fields.reporter.name }}
watchers: {{ range .fields.customfield_10110 }}{{ .name }} {{end}}
blockers: {{ range .fields.issuelinks }}{{if .outwardIssue}}{{ .outwardIssue.key }}[{{.outwardIssue.fields.status.name}}]{{end}}{{end}}
depends: {{ range .fields.issuelinks }}{{if .inwardIssue}}{{ .inwardIssue.key }}[{{.inwardIssue.fields.status.name}}]{{end}}{{end}}
priority: {{ .fields.priority.name }}
description: |
  {{ .fields.description | indent 2 }}

comments:
{{ range .fields.comment.comments }}  - | # {{.author.name}} at {{.created}}
    {{ .body | indent 4}}
{{end}}
`
const default_edit_template = `update:
  comment:
    - add: 
        body: |
          {{ or .overrides.comment ""}}
fields:
  summary: {{ or .overrides.summary .fields.summary }}
  components: # {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{if .overrides.components }}{{ range (split "," .overrides.components)}}
    - name: {{.}}{{end}}{{else}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}{{end}}
  assignee:
    name: {{ if .overrides.assignee }}{{.overrides.assignee}}{{else}}{{if .fields.assignee }}{{ .fields.assignee.name }}{{end}}{{end}}
  reporter:
    name: {{ or .overrides.reporter .fields.reporter.name }}
  # watchers
  customfield_10110: {{ range .fields.customfield_10110 }}
    - name: {{ .name }}{{end}}{{if .overrides.watcher}}
    - name: {{ .overrides.watcher}}{{end}}
  priority: # {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority .fields.priority.name }}
  description: |
    {{ or .overrides.description (or .fields.description "") | indent 4 }}
`
const default_transitions_template = `{{ range .transitions }}{{.id }}: {{.name}}
{{end}}`

const default_issuetypes_template = `{{ range .projects }}{{ range .issuetypes }}{{color "+bh"}}{{.name | append ":" | printf "%-13s" }}{{color "reset"}} {{.description}}
{{end}}{{end}}`

const default_create_template = `fields:
  project:
    key: {{ .overrides.project }}
  issuetype:
    name: {{ .overrides.issuetype }}
  summary: {{ or .overrides.summary "" }}
  priority: # {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority "" }}
  components: # {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{ range split "," (or .overrides.components "")}}
    - name: {{ . }}{{end}}
  description: |
    {{ or .overrides.description "" | indent 4 }}
  assignee:
    name: {{ or .overrides.assignee .overrides.user}}
  reporter:
    name: {{ or .overrides.reporter .overrides.user }}
  # watchers
  customfield_10110:
    - name:
`

const default_comment_template = `body: |
  
`
