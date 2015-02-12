package cli

const default_fields_template = "{{ . | toJson}}\n"

const default_list_template = "{{ range .issues }}{{ .key | append \":\" | printf \"%-12s\"}} {{ .fields.summary }}\n{{ end }}"

const default_view_template = `issue: {{ .key }}
status: {{ .fields.status.name }}
summary: {{ .fields.summary }}
project: {{ .fields.project.key }}
components: {{ range .fields.components }}{{ .name }} {{end}}
issuetype: {{ .fields.issuetype.name }}
assignee: {{ .fields.assignee.name }}
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
          
fields:
  summary: {{ .fields.summary }}
  components: # {{ range .meta.components.allowedValues }}{{.name}}, {{end}}{{ range .fields.components }}
    - name: {{ .name }}{{end}}
  assignee:
    name: {{ .fields.assignee.name }}
  reporter:
    name: {{ .fields.reporter.name }}
  # watchers
  customfield_10110: {{ range .fields.customfield_10110 }}
    - name: {{ .name }}{{end}}
  priority: # {{ range .meta.priority.allowedValues }}{{.name}}, {{end}}
    name: {{ .fields.priority.name }}
  description: |
    {{ .fields.description | indent 4 }}
`
const default_transitions_template = `{{ range .transitions }}{{color "+bh"}}{{.name | printf "%-13s" }}{{color "reset"}} -> {{.to.name}}
{{end}}`

const default_issuetypes_template = `{{ range .projects }}{{ range .issuetypes }}{{color "+bh"}}{{.name | append ":" | printf "%-13s" }}{{color "reset"}} {{.description}}
{{end}}{{end}}`
