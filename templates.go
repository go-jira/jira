package jira

var allTemplates = map[string]string{
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

const defaultIssuetypesTemplate = `{{ range .projects }}{{ range .issuetypes }}{{color "+bh"}}{{.name | append ":" | printf "%-13s" }}{{color "reset"}} {{.description}}
{{end}}{{end}}`

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
update:
  comment:
    - add: 
        body: |~
          {{ or .overrides.comment "" | indent 10 }}
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
