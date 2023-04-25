package mocks

import "strings"

// English only as no Welsh translations exist
var enLocale = []string{
	"[FeedbackTitle]",
	"one = \"Feedback\"",
	"[FeedbackThanks]",
	"one = \"Thank you\"",
	"[FeedbackChooseType]",
	"one = \"Choose feedback type\"",
	"[FeedbackWhatOptNewService]",
	"one = \"This service\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	return []byte(strings.Join(enLocale, "\n")), nil
}
