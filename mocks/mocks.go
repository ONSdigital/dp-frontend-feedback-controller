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
	"[FeedbackWhatEnterURL]",
	"one = \"Enter URL or name of the page\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) { //nolint:all // app does not use welsh
	return []byte(strings.Join(enLocale, "\n")), nil
}
