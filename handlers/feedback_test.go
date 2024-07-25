package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cacheClient "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/client"
	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/email/emailtest"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces/interfacestest"
	"github.com/ONSdigital/dp-frontend-feedback-controller/mocks"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	topicModel "github.com/ONSdigital/dp-topic-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

const siteDomain = "ons.gov.uk"

func Test_getFeedback(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		ff := model.FeedbackForm{}
		ff.URL = "whatever"
		lang := "en"
		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When getFeedback is called", func() {
			getFeedback(w, req, []coreModel.ErrorItem{}, ff, lang, mockRenderer, mockNagivationCache, false)
			Convey("Then a 200 request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func Test_addFeedback(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid request", t, func() {
		body := strings.NewReader("description=testing1234&type=test")
		req := httptest.NewRequest("POST", "http://localhost", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		from := ""
		to := ""
		lang := "en"

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		mockSender := &emailtest.SenderMock{
			SendFunc: func(from string, to []string, msg []byte) error {
				return nil
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When addFeedback is called", func() {
			addFeedback(w, req, mockRenderer, mockSender, from, to, lang, siteDomain, mockNagivationCache)
			Convey("Then the renderer is not called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 0)
			})
			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 1)
			})
			Convey("Then a 301 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusMovedPermanently)
			})
		})
	})

	Convey("Given an error returned from the sender", t, func() {
		body := strings.NewReader("description=testing1234&type=test")
		req := httptest.NewRequest("POST", "http://localhost", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		from := ""
		to := ""
		lang := "en"

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		mockSender := &emailtest.SenderMock{
			SendFunc: func(from string, to []string, msg []byte) error {
				return errors.New("email is broken")
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When addFeedback is called", func() {
			addFeedback(w, req, mockRenderer, mockSender, from, to, lang, siteDomain, mockNagivationCache)
			Convey("Then the renderer is not called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 0)
			})

			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 1)
			})

			Convey("Then a 500 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given a request with invalid form data", t, func() {
		req := httptest.NewRequest("POST", "http://localhost?!@£$@$£%£$%^^&^&*", nil)
		w := httptest.NewRecorder()
		from := ""
		to := ""
		lang := "en"

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		mockSender := &emailtest.SenderMock{
			SendFunc: func(from string, to []string, msg []byte) error {
				return nil
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When addFeedback is called", func() {
			addFeedback(w, req, mockRenderer, mockSender, from, to, lang, siteDomain, mockNagivationCache)
			Convey("Then the renderer is not called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 0)
			})

			Convey("Then the email sender is not called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 0)
			})

			Convey("Then a 400 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given a request for feedback with an empty description value", t, func() {
		body := strings.NewReader("description=")
		req := httptest.NewRequest("POST", "http://localhost?service=dev", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		from := ""
		to := ""
		lang := "en"

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		mockSender := &emailtest.SenderMock{
			SendFunc: func(from string, to []string, msg []byte) error {
				return nil
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When addFeedback is called", func() {
			addFeedback(w, req, mockRenderer, mockSender, from, to, lang, siteDomain, mockNagivationCache)
			Convey("Then the renderer is called to render the feedback page", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})
			Convey("Then the email sender is not called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 0)
			})
			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func Test_feedbackThanks(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	lang := "en"
	config.Get() // need to seed config
	Convey("Given a valid request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		url := "www.test.com"

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When feedbackThanks is called", func() {
			feedbackThanks(w, req, url, mockRenderer, mockNagivationCache, lang, siteDomain, false)
			Convey("Then the renderer is called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})
			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})

	Convey("Given a reflective XSS request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost?returnTo=<script>alert(1)</script>", nil)
		w := httptest.NewRecorder()
		url := "https://www.referrer-test.com"
		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}
		mockNagivationCache := &cacheHelper.Helper{
			Clienter: &cacheClient.ClienterMock{
				AddNavigationCacheFunc: func(ctx context.Context, updateInterval *time.Duration) error {
					return nil
				},
				CloseFunc: func() {
				},
				GetNavigationDataFunc: func(ctx context.Context, lang string) (*topicModel.Navigation, error) {
					return &topicModel.Navigation{}, nil
				},
				StartBackgroundUpdateFunc: func(ctx context.Context, errorChannel chan error) {
				},
			}}
		Convey("When feedbackThanks is called", func() {
			feedbackThanks(w, req, url, mockRenderer, mockNagivationCache, lang, siteDomain, false)
			Convey("Then the handler sanitises the request text to the referrer", func() {
				dataSentToRender := mockRenderer.BuildPageCalls()[0].PageModel.(model.Feedback)
				returnToUrl := dataSentToRender.ReturnTo
				So(returnToUrl, ShouldEqual, url)
			})
		})
	})
}

func TestValidateForm(t *testing.T) {
	Convey("Given the validateForm function", t, func() {
		testCases := []struct {
			givenDescription    string
			given               *model.FeedbackForm
			expectedDescription string
			expected            []coreModel.ErrorItem
		}{
			{
				givenDescription: "the form is valid",
				given: &model.FeedbackForm{
					Type:        "Whole site",
					Description: "Some text",
				},
				expectedDescription: "no validation errors are returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "the form does not have a type selected",
				given: &model.FeedbackForm{
					Description: "Some text",
				},
				expectedDescription: "a type validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackChooseType",
							Plural:    1,
						},
						URL: "#type-error",
					},
				},
			},
			{
				givenDescription: "the a specific page/url type is chosen but the child input field is empty",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					Description: "Some text",
				},
				expectedDescription: "a page/url validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackWhatEnterURL",
							Plural:    1,
						},
						URL: "#type-error",
					},
				},
			},
			{
				givenDescription: "the a specific page/url type is chosen and the url is invalid",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					Description: "Some text",
					URL:         "not a url",
				},
				expectedDescription: "a url validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackValidURL",
							Plural:    1,
						},
						URL: "#type-error",
					},
				},
			},
			{
				givenDescription: "the a specific page/url type is chosen and the url is valid but not allowed",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					Description: "Some text",
					URL:         "https://not-site-domain.com",
				},
				expectedDescription: "a url validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackValidURL",
							Plural:    1,
						},
						URL: "#type-error",
					},
				}},
			{
				givenDescription: "the a specific page/url type is chosen and the url is valid without a path",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					Description: "Some text",
					URL:         "https://cy.ons.gov.uk",
				},
				expectedDescription: "no validation errors are returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "the a specific page/url type is chosen and the url is valid and has a path",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					Description: "Some text",
					URL:         "https://cy.ons.gov.uk/path",
				},
				expectedDescription: "no validation errors are returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "the a whole site type is chosen but the child input for a specific page is not empty",
				given: &model.FeedbackForm{
					Type:        "Whole site",
					Description: "Some text",
					URL:         "http://somewhere.com",
				},
				expectedDescription: "no validation error is returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "the form does not have a type selected but is located on the footer",
				given: &model.FeedbackForm{
					FormLocation: "footer",
					Description:  "Some text",
				},
				expectedDescription: "no validation error is returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "the form does not have any feedback",
				given: &model.FeedbackForm{
					Type: "Whole site",
				},
				expectedDescription: "a description validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackAlertEntry",
							Plural:    1,
						},
						URL: "#feedback-error",
					},
				},
			},
			{
				givenDescription: "the feedback provided is whitespace",
				given: &model.FeedbackForm{
					Type:        "Whole site",
					Description: " ",
				},
				expectedDescription: "a description validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackAlertEntry",
							Plural:    1,
						},
						URL: "#feedback-error",
					},
				},
			},
			{
				givenDescription: "the email field has an invalid email address",
				given: &model.FeedbackForm{
					Type:        "Whole site",
					Description: "A description",
					Email:       "a.string",
				},
				expectedDescription: "an email validation error is returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackAlertEmail",
							Plural:    1,
						},
						URL: "#email-error",
					},
				},
			},
			{
				givenDescription: "the email field has a valid email address",
				given: &model.FeedbackForm{
					Type:        "Whole site",
					Description: "A description",
					Email:       "hello@world.com",
				},
				expectedDescription: "no validation errors are returned",
				expected:            []coreModel.ErrorItem(nil),
			},
			{
				givenDescription: "multiple form validation errors",
				given: &model.FeedbackForm{
					Type:        "A specific page",
					URL:         "",
					Description: "",
					Email:       "not an email address",
				},
				expectedDescription: "validation errors are returned",
				expected: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackWhatEnterURL",
							Plural:    1,
						},
						URL: "#type-error",
					},
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackAlertEntry",
							Plural:    1,
						},
						URL: "#feedback-error",
					},
					{
						Description: coreModel.Localisation{
							LocaleKey: "FeedbackAlertEmail",
							Plural:    1,
						},
						URL: "#email-error",
					},
				},
			},
		}
		for _, t := range testCases {
			Convey(fmt.Sprintf("When %s", t.givenDescription), func() {
				Convey(fmt.Sprintf("Then %s", t.expectedDescription), func() {
					So(validateForm(t.given, siteDomain), ShouldResemble, t.expected)
				})
			})
		}
	})
}
