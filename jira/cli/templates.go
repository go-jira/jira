package cli

const default_fields_template = "{{ . | toJson}}\n"

const default_list_template = "{{ range .issues }}{{ .key | append \":\" | printf \"%-12s\"}} {{ .fields.summary }}\n{{ end }}"

const default_view_template = `issue: {{ .key }}
summary: {{ .fields.summary }}
project: {{ .fields.project.key }}
components: {{ range .fields.components }}{{ .name }} {{end}}
issuetype: {{ .fields.issuetype.name }}
assignee: {{ .fields.assignee.name }}
reporter: {{ .fields.reporter.name }}
watchers: {{ range .fields.customfield_10110 }}{{ .name }} {{end}}
blockers: {{ range .fields.issuelinks }}{{if .outwardIssue}}{{ .outwardIssue.key }}{{end}}{{end}}
depends: {{ range .fields.issuelinks }}{{if .inwardIssue}}{{ .inwardIssue.key }}{{end}}{{end}}
priority: {{ .fields.priority.name }}
description: |
  {{ .fields.description | indent 2 }}

comments:
{{ range .fields.comment.comments }}  - | # {{.author.name}} at {{.created}}
    {{ .body | indent 4}}{{end}}
`
