package mocks

import "strings"

// English only as no Welsh translations exist
var enLocale = []string{
	"[FeedbackTitle]",
	"one = \"Feedback\"",
	"[FeedbackThanks]",
	"one = \"Thank you\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	return []byte(strings.Join(enLocale, "\n")), nil
}
