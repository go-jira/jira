package jiradata

import ("strings")

func IsEpicTitle(s string) bool {
	return strings.HasPrefix(s, "# ")
}

func IsEpicSummary(s string) bool {
	return strings.HasPrefix(s, "**")
}

func IsTicketTitle(s string) bool {
	return strings.HasPrefix(s, "## ")
}

func IsSeparator(s string) bool {
	return s == ""
}

func IsDescriptionLine(s string) bool {
	return !strings.HasPrefix(s, "*") && !strings.HasPrefix(s, "#") && !IsSeparator(s)
}

func IsItem(s string) bool {
	return strings.HasPrefix(s, "- ")
}
