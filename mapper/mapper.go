package mapper

import (
	"html"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dis-design-system-go/helper"
	core "github.com/ONSdigital/dis-design-system-go/model"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
)

const (
	WholeSite     = "The whole website"
	ASpecificPage = "A specific page"
)

var cfg *config.Config

// CreateGetFeedback returns a mapped feedback page to the feedback model
func CreateGetFeedback(req *http.Request, basePage core.Page, validationErrors []core.ErrorItem, ff model.FeedbackForm, lang string) model.Feedback {
	p := model.Feedback{
		Page: basePage,
	}
	p.Breadcrumb = []core.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
	}

	var services = make(map[string]string)
	services["cmd"] = "customising data by applying filters"
	services["dev"] = "ONS developer"
	services["search"] = "search"
	serviceDescription := services[req.URL.Query().Get("service")]

	p.Language = lang
	p.Type = "feedback"
	p.URI = req.URL.Path
	p.Metadata.Title = helper.Localise("FeedbackTitle", lang, 1)
	p.Metadata.Description = ff.URL

	if len(p.Metadata.Description) > 50 {
		p.Metadata.Description = p.Metadata.Description[len(p.Metadata.Description)-50 : len(p.Metadata.Description)]
	}

	if len(validationErrors) > 0 {
		p.Page.Error = core.Error{
			Title:      p.Metadata.Title,
			ErrorItems: validationErrors,
			Language:   lang,
		}
	}

	radioErrDetail := helper.Localise("FeedbackChooseType", lang, 1)
	if ff.IsURLErr {
		radioErrDetail = helper.Localise("FeedbackWhatEnterURL", lang, 1)
	}
	p.TypeRadios = core.RadioFieldset{
		Legend: core.Localisation{
			LocaleKey: "FeedbackTitleWhat",
			Plural:    1,
		},
		Radios: []core.Radio{
			{
				Input: core.Input{
					ID:        "whole-site",
					IsChecked: ff.Type == WholeSite,
					Label: core.Localisation{
						LocaleKey: "FeedbackWholeWebsite",
						Plural:    1,
					},
					Name:  "type",
					Value: WholeSite,
				},
			},
			{
				Input: core.Input{
					ID:        "specific-page",
					IsChecked: ff.Type == ASpecificPage || (ff.URL != "" && serviceDescription == ""),
					Label: core.Localisation{
						LocaleKey: "FeedbackASpecificPage",
						Plural:    1,
					},
					Name:  "type",
					Value: ASpecificPage,
				},
				OtherInput: core.Input{
					Autocomplete: "url",
					ID:           "page-url-field",
					Name:         "url",
					Value: func() string {
						if serviceDescription != "" {
							return ""
						}
						return ff.URL
					}(),
					Label: core.Localisation{
						LocaleKey: "FeedbackWhatEnterURL",
						Plural:    1,
					},
					Type: core.Url,
				},
			},
		},
		ValidationErr: core.ValidationErr{
			HasValidationErr: ff.IsTypeErr || ff.IsURLErr,
			ErrorItem: core.ErrorItem{
				Description: core.Localisation{
					Text: radioErrDetail,
				},
				ID: "type-error",
			},
		},
	}

	if serviceDescription != "" {
		p.TypeRadios.Radios = append(
			p.TypeRadios.Radios[:1],
			core.Radio{
				Input: core.Input{
					ID:        "new-service",
					IsChecked: true,
					Label: core.Localisation{
						Text: helper.Localise("FeedbackWhatOptNewService", lang, 1, serviceDescription),
					},
					Name:  "type",
					Value: "The new service",
				},
			},
			p.TypeRadios.Radios[1])
	}

	p.Contact = []core.TextField{
		{
			Input: core.Input{
				Autocomplete: "name",
				ID:           "name-field",
				Name:         "name",
				Value:        ff.Name,
				Label: core.Localisation{
					LocaleKey: "FeedbackTitleName",
					Plural:    1,
				},
			},
		},
		{
			Input: core.Input{
				Autocomplete: "email",
				ID:           "email-field",
				Name:         "email",
				Label: core.Localisation{
					LocaleKey: "FeedbackTitleEmail",
					Plural:    1,
				},
				Type:  core.Email,
				Value: ff.Email,
			},
			ValidationErr: core.ValidationErr{
				HasValidationErr: ff.IsEmailErr,
				ErrorItem: core.ErrorItem{
					Description: core.Localisation{
						LocaleKey: "FeedbackAlertEmail",
						Plural:    1,
					},
					ID: "email-error",
				},
			},
		},
	}

	p.DescriptionField = core.TextareaField{
		Input: core.Input{
			Autocomplete: "off",
			Description: core.Localisation{
				LocaleKey: "FeedbackHintEntry",
				Plural:    1,
			},
			ID: "description-field",
			Label: core.Localisation{
				LocaleKey: "FeedbackTitleEntry",
				Plural:    1,
			},
			Language: lang,
			Name:     "description",
			Value:    ff.Description,
		},
		ValidationErr: core.ValidationErr{
			HasValidationErr: ff.IsDescriptionErr,
			ErrorItem: core.ErrorItem{
				Description: core.Localisation{
					LocaleKey: "FeedbackAlertEntry",
					Plural:    1,
				},
				ID: "feedback-error",
			},
		},
	}

	p.PreviousURL = ff.URL

	return p
}

func CreateGetFeedbackThanks(req *http.Request, basePage core.Page, lang, referrer, wholeSiteURL string) model.Feedback {
	if wholeSiteURL == "" {
		wholeSiteURL = "https://www.ons.gov.uk"
	}
	if referrer == "" {
		referrer = wholeSiteURL
	}

	p := model.Feedback{
		Page:        basePage,
		PreviousURL: referrer,
	}
	p.Language = lang
	p.Type = "feedback"
	p.URI = req.URL.Path
	p.Metadata.Title = helper.Localise("FeedbackThanks", lang, 1)

	// returnTo is rendered on page so needs XSS protection
	returnTo := html.EscapeString(req.URL.Query().Get("returnTo"))
	if returnTo == WholeSite {
		returnTo = wholeSiteURL
	} else if returnTo == "" {
		returnTo = referrer
	} else if IsSiteDomainURL(returnTo, "") {
		returnTo = NormaliseURL(returnTo)
	} else {
		returnTo = referrer
	}

	p.ReturnTo = returnTo

	return p
}

// IsSiteDomainURL is true when urlString is a URL and its host ends with `.`+siteDomain (when siteDomain is blank, or uses config.SiteDomain)
func IsSiteDomainURL(urlString, siteDomain string) bool {
	if urlString == "" {
		return false
	}
	urlString = NormaliseURL(urlString)
	urlObject, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	if siteDomain == "" {
		if cfg == nil {
			if cfg, err = config.Get(); err != nil {
				return false
			}
		}
		siteDomain = cfg.SiteDomain
	}
	hostName := urlObject.Hostname()
	if hostName != siteDomain && !strings.HasSuffix(hostName, "."+siteDomain) {
		return false
	}
	return true
}

// NormaliseURL when a string is a URL without a scheme (e.g. `host.name/path`), add it (`https://`)
func NormaliseURL(urlString string) string {
	if strings.HasPrefix(urlString, "http") {
		return urlString
	}
	return "https://" + urlString
}
