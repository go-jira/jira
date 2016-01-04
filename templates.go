package jira

var all_templates = map[string]string{
	"debug":          default_debug_template,
	"fields":         default_debug_template,
	"editmeta":       default_debug_template,
	"transmeta":      default_debug_template,
	"createmeta":     default_debug_template,
	"issuelinktypes": default_debug_template,
	"list":           default_list_template,
	"table":          default_table_template,
	"view":           default_view_template,
	"edit":           default_edit_template,
	"transitions":    default_transitions_template,
	"issuetypes":     default_issuetypes_template,
	"create":         default_create_template,
	"comment":        default_comment_template,
	"transition":     default_transition_template,
	"request":        default_debug_template,
}

const default_debug_template = "{{ . | toJson}}\n"

const default_list_template = "{{ range .issues }}{{ .key | append \":\" | printf \"%-12s\"}} {{ .fields.summary }}\n{{ end }}"

const default_table_template = `+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
| {{ "Issue" | printf "%-14s" }} | {{ "Summary" | printf "%-55s" }} | {{ "Priority" | printf "%-12s" }} | {{ "Status" | printf "%-12s" }} | {{ "Age" | printf "%-10s" }} | {{ "Reporter" | printf "%-12s" }} | {{ "Assignee" | printf "%-12s" }} |
+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
{{ range .issues }}| {{ .key | printf "%-14s"}} | {{ .fields.summary | abbrev 55 | printf "%-55s" }} | {{.fields.priority.name | printf "%-12s" }} | {{.fields.status.name | printf "%-12s" }} | {{.fields.created | age | printf "%-10s" }} | {{if .fields.reporter}}{{ .fields.reporter.name | printf "%-12s"}}{{else}}<unassigned>{{end}} | {{if .fields.assignee }}{{.fields.assignee.name | printf "%-12s" }}{{else}}<unassigned>{{end}} |
{{ end }}+{{ "-" | rep 16 }}+{{ "-" | rep 57 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+{{ "-" | rep 12 }}+{{ "-" | rep 14 }}+{{ "-" | rep 14 }}+
`

const default_view_template = `issue: {{ .key }}
created: {{ .fields.created }}
status: {{ .fields.status.name }}
summary: {{ .fields.summary }}
project: {{ .fields.project.key }}
components: {{ range .fields.components }}{{ .name }} {{end}}
issuetype: {{ .fields.issuetype.name }}
assignee: {{ if .fields.assignee }}{{ .fields.assignee.name }}{{end}}
reporter: {{ if .fields.reporter }}{{ .fields.reporter.name }}{{end}}
watchers: {{ range .fields.customfield_10110 }}{{ .name }} {{end}}
blockers: {{ range .fields.issuelinks }}{{if .outwardIssue}}{{ .outwardIssue.key }}[{{.outwardIssue.fields.status.name}}]{{end}}{{end}}
depends: {{ range .fields.issuelinks }}{{if .inwardIssue}}{{ .inwardIssue.key }}[{{.inwardIssue.fields.status.name}}]{{end}}{{end}}
priority: {{ .fields.priority.name }}
description: |
  {{ or .fields.description "" | indent 2 }}

comments:
{{ range .fields.comment.comments }}  - | # {{.author.name}} at {{.created}}
    {{ or .body "" | indent 4}}
{{end}}
`
const default_edit_template = `# issue: {{ .key }}
update:
  comment:
    - add: 
        body: |~
          {{ or .overrides.comment "" | indent 10 }}
fields:
  summary: {{ or .overrides.summary .fields.summary }}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{if .overrides.components }}{{ range (split "," .overrides.components)}}
    - name: {{.}}{{end}}{{else}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}{{end}}
  assignee:
    name: {{ if .overrides.assignee }}{{.overrides.assignee}}{{else}}{{if .fields.assignee }}{{ .fields.assignee.name }}{{end}}{{end}}
#	reporter:
#	  name: {{ if .overrides.reporter }}{{ .overrides.reporter }}{{else if .fields.reporter}}{{ .fields.reporter.name }}{{end}}
  # watchers
  customfield_10110: {{ range .fields.customfield_10110 }}
    - name: {{ .name }}{{end}}{{if .overrides.watcher}}
    - name: {{ .overrides.watcher}}{{end}}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority .fields.priority.name }}
  description: |~
    {{ or .overrides.description (or .fields.description "") | indent 4 }}
# comments:
# {{ range .fields.comment.comments }}  - | # {{.author.name}} at {{.created}}
#     {{ or .body "" | indent 4 | comment}}
# {{end}}
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
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority "unassigned" }}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{ range split "," (or .overrides.components "")}}
    - name: {{ . }}{{end}}
  description: |~
    {{ or .overrides.description "" | indent 4 }}
  assignee:
    name: {{ or .overrides.assignee "" }}
  reporter:
    name: {{ or .overrides.reporter .overrides.user }}
  # watchers
  customfield_10110: {{ range split "," (or .overrides.watchers "")}}
    - name: {{.}}{{end}}
    - name:
`

const default_comment_template = `body: |~
  {{ or .overrides.comment "" | indent 2 }}
`

const default_transition_template = `update:
  comment:
    - add: 
        body: |~
          {{ or .overrides.comment "" | indent 10 }}
fields:{{if .meta.fields.assignee}}
  assignee:
    name: {{if .overrides.assignee}}{{.overrides.assignee}}{{else}}{{if .fields.assignee}}{{.fields.assignee.name}}{{end}}{{end}}{{end}}{{if .meta.fields.components}}
  components: # Values: {{ range .meta.fields.components.allowedValues }}{{.name}}, {{end}}{{if .overrides.components }}{{ range (split "," .overrides.components)}}
    - name: {{.}}{{end}}{{else}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}{{end}}{{end}}{{if .meta.fields.description}}
  description: {{or .overrides.description .fields.description }}{{end}}{{if .meta.fields.fixVersions}}{{if .meta.fields.fixVersions.allowedValues}}
  fixVersions: # Values: {{ range .meta.fields.fixVersions.allowedValues }}{{.name}}, {{end}}{{if .overrides.fixVersions}}{{ range (split "," .overrides.fixVersions)}}
    - name: {{.}}{{end}}{{else}}{{range .fields.fixVersions}}
    - name: {{.}}{{end}}{{end}}{{end}}{{end}}{{if .meta.fields.issuetype}}
  issuetype: # Values: {{ range .meta.fields.issuetype.allowedValues }}{{.name}}, {{end}}
    name: {{if .overrides.issuetype}}{{.overrides.issuetype}}{{else}}{{if .fields.issuetype}}{{.fields.issuetype.name}}{{end}}{{end}}{{end}}{{if .meta.fields.labels}}
  labels: {{range .fields.labels}}
    - {{.}}{{end}}{{if .overrides.labels}}{{range (split "," .overrides.labels)}}
    - {{.}}{{end}}{{end}}{{end}}{{if .meta.fields.priority}}
  priority: # Values: {{ range .meta.fields.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ or .overrides.priority "unassigned" }}{{end}}{{if .meta.fields.reporter}}
  reporter:
    name: {{if .overrides.reporter}}{{.overrides.reporter}}{{else}}{{if .fields.reporter}}{{.fields.reporter.name}}{{end}}{{end}}{{end}}{{if .meta.fields.resolution}}
  resolution: # Values: {{ range .meta.fields.resolution.allowedValues }}{{.name}}, {{end}}
    name: {{if .overrides.resolution}}{{.overrides.resolution}}{{else if .fields.resolution}}{{.fields.resolution.name}}{{else}}Fixed{{end}}{{end}}{{if .meta.fields.summary}}
  summary: {{or .overrides.summary .fields.summary}}{{end}}{{if .meta.fields.versions.allowedValues}}
  versions: # Values: {{ range .meta.fields.versions.allowedValues }}{{.name}}, {{end}}{{if .overrides.versions}}{{ range (split "," .overrides.versions)}}
    - name: {{.}}{{end}}{{else}}{{range .fields.versions}}
    - name: {{.}}{{end}}{{end}}{{end}}
transition:
  id: {{ .transition.id }}
  name: {{ .transition.name }}
`
