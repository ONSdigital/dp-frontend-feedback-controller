package mapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-frontend-feedback-controller/mocks"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	core "github.com/ONSdigital/dp-renderer/v2/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateGetFeedback(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page request", t, func() {
		Convey("When the parameters area valid", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			bp := core.Page{}
			validationErr := []core.ErrorItem{}
			ff := model.FeedbackForm{}
			lang := "en"
			sut := CreateGetFeedback(req, bp, validationErr, ff, lang)

			Convey("Then it sets the page metadata", func() {
				So(sut.Type, ShouldEqual, "feedback")
				So(sut.Metadata.Title, ShouldEqual, "Feedback")
				So(sut.Metadata.Description, ShouldEqual, ff.URL)
				So(sut.Language, ShouldEqual, lang)
			})

			Convey("Then it maps the expected radio inputs", func() {
				So(sut.TypeRadios, ShouldNotBeEmpty)
				So(sut.TypeRadios.Radios, ShouldHaveLength, 2)
			})

			Convey("Then it maps the contact text inputs", func() {
				So(sut.Contact, ShouldNotBeEmpty)
				So(sut.Contact, ShouldHaveLength, 2)
			})

			Convey("Then it maps the description textarea input", func() {
				So(sut.DescriptionField, ShouldNotBeEmpty)
			})

			Convey("Then it maps the previous url field", func() {
				So(sut.PreviousURL, ShouldEqual, ff.URL)
			})
		})

		Convey("When a valid service parameter is passed", func() {
			req := httptest.NewRequest(http.MethodGet, "/?service=cmd", nil)
			bp := core.Page{}
			validationErr := []core.ErrorItem{}
			ff := model.FeedbackForm{}
			lang := "en"
			sut := CreateGetFeedback(req, bp, validationErr, ff, lang)

			Convey("Then it maps the additional radio input", func() {
				So(sut.TypeRadios, ShouldNotBeEmpty)
				So(sut.TypeRadios.Radios, ShouldHaveLength, 3)
			})
		})

		Convey("When validation errors are passed", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			bp := core.Page{}
			lang := "en"
			validationErr := []core.ErrorItem{
				{
					Description: core.Localisation{
						Text: "Error one",
					},
					Language: lang,
					ID:       "error-one",
					URL:      "#error-one",
				},
			}
			ff := model.FeedbackForm{}
			sut := CreateGetFeedback(req, bp, validationErr, ff, lang)

			Convey("Then it maps the error panel", func() {
				So(sut.Error.Title, ShouldNotBeEmpty)
			})

			Convey("The error items are mapped", func() {
				So(sut.Error.ErrorItems, ShouldResemble, validationErr)
			})

			ff.IsURLErr = true
			sut = CreateGetFeedback(req, bp, validationErr, ff, lang)
			Convey("Then it changes the radio validation field description", func() {
				So(sut.TypeRadios.ValidationErr.ErrorItem.Description.Text, ShouldEqual, "Enter URL or name of the page")
			})
		})
	})
}

func TestCreateGetFeedbackThanks(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page request", t, func() {
		Convey("When the parameters area valid", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			bp := core.Page{}
			url := "https://localhost/a/page/somewhere"
			lang := "en"
			wholeSiteURL := "https://ons.gov.uk"
			sut := CreateGetFeedbackThanks(req, bp, lang, url, wholeSiteURL)

			Convey("Then it sets the page metadata", func() {
				So(sut.Metadata.Title, ShouldEqual, "Thank you")
				So(sut.Page.Type, ShouldEqual, "feedback")
				So(sut.Language, ShouldEqual, lang)
			})

			Convey("Then it sets the previousUrl property", func() {
				So(sut.PreviousURL, ShouldEqual, url)
			})

			Convey("Then it sets the returnTo property", func() {
				So(sut.ReturnTo, ShouldEqual, url)
			})
		})

		var (
			referrer     = "https://any.localhost/a/page/somewhere"
			lang         = "en"
			wholeSiteURL = "https://cy.localhost"
			encWholeSite = url.QueryEscape(WholeSite)
			bp           = core.Page{}
		)

		Convey("When the returnTo parameter is set to whole-site and whole-site is explicit", func() {
			req := httptest.NewRequest(http.MethodGet, "/?returnTo="+encWholeSite, nil)
			sut := CreateGetFeedbackThanks(req, bp, lang, referrer, wholeSiteURL)

			Convey("Then it sets the returnTo property to the whole-site", func() {
				So(sut.ReturnTo, ShouldEqual, wholeSiteURL)
			})
		})
		Convey("When the returnTo parameter is set to whole-site but whole-site is not explicit", func() {
			req := httptest.NewRequest(http.MethodGet, "/?returnTo="+encWholeSite, nil)
			sut := CreateGetFeedbackThanks(req, bp, lang, referrer, "")

			Convey("Then it sets the returnTo property to the default whole-site", func() {
				So(sut.ReturnTo, ShouldEqual, "https://www.ons.gov.uk")
			})
		})
	})
}

func TestURLFunctions(t *testing.T) {
	Convey("Given the IsSiteDomainURL functions", t, func() {

		type testSiteDomainStruct struct {
			name      string
			pageURL   string
			siteURL   string
			isAllowed bool
		}
		tests := []testSiteDomainStruct{
			{
				"sub-domain off an explicit site domain",
				"https://anything.ons.gov.uk:443/ook",
				"ons.gov.uk",
				true,
			},
			{
				"non-site domain URL is not recognised for explicit site domain",
				"https://anything.example.com",
				"ons.gov.uk",
				false,
			},
			{
				"non-URL is not recognised for explicit site domain",
				"blah",
				"ons.gov.uk",
				false,
			},
			{
				"URL of the config's site domain is recognised",
				"https://localhost",
				"",
				true,
			},
			{
				"sub-domain/host of the config's site domain is recognised",
				"anything.localhost",
				"",
				true,
			},
			{
				"non-site domain URL is not recognised for config's site domain",
				"https://not-site-domain.example.com",
				"",
				false,
			},
			{
				"non-URL is not recognised for config's site domain",
				"blah",
				"",
				false,
			},
		}
		// test IsSiteDomainURL
		for _, check := range tests {
			Convey("When "+check.name, func() {
				allowedStr := fmt.Sprint(check.isAllowed)
				Convey("Then "+check.name+" is "+allowedStr, func() {
					isAllowedURL := IsSiteDomainURL(check.pageURL, check.siteURL)
					So(isAllowedURL, ShouldEqual, check.isAllowed)
				})
			})
		}

		// test NormaliseURL
		Convey("When given a normal URL", func() {
			origURL := "http://is.url"
			normalURL := NormaliseURL(origURL)
			Convey("Then it returns as itself", func() {
				So(normalURL, ShouldEqual, origURL)
			})
		})
		Convey("When given a user-typed URL", func() {
			origURL := "is.url/path"
			normalURL := NormaliseURL(origURL)
			Convey("Then it is normalised with `https://`", func() {
				So(normalURL, ShouldEqual, `https://`+origURL)
			})
		})
	})
}
