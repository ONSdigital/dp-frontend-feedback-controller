package mapper

import (
	"net/http"
	"net/http/httptest"
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
			wholeSite := "https://ons.gov.uk"
			sut := CreateGetFeedbackThanks(req, bp, lang, url, wholeSite)

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

		Convey("When the return to parameter is set", func() {
			req := httptest.NewRequest(http.MethodGet, "/?returnTo=Whole%20site", nil)
			bp := core.Page{}
			url := "https://localhost/a/page/somewhere"
			lang := "en"
			wholeSite := "https://ons.gov.uk"
			sut := CreateGetFeedbackThanks(req, bp, lang, url, wholeSite)

			Convey("Then it sets the returnTo property", func() {
				So(sut.ReturnTo, ShouldEqual, wholeSite)
			})
		})
	})
}
