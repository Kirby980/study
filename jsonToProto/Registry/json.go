package Registry

import "regexp"

func RemoveJSONComment(jsonStr string) string {
	reMultiline := regexp.MustCompile(`(?s)/\*.*?\*/`)
	cleaned := reMultiline.ReplaceAllString(jsonStr, "")

	// 单行注释处理（排除合法URL）
	reSingleLine := regexp.MustCompile(`(?m)(^|[^:])//.*$`)
	cleaned = reSingleLine.ReplaceAllString(cleaned, "$1")

	cleaned = regexp.MustCompile(`,(\s*[}\]])`).ReplaceAllString(cleaned, "$1")
	return cleaned
}
