---
number: {{.issue.Number}}
title: "{{.issue.Title}}"
state: "{{.issue.State}}"
tags:
{{range .issue.Labels}}  - "{{.Name}}"
{{end}}
date: "{{.issue.CreatedAt}}"
lastmod: "{{.issue.UpdatedAt}}"
milestone: "{{if .issue.Milestone}}{{.issue.Milestone.Title}}{{end}}"
---

{{.issue.Body}}

