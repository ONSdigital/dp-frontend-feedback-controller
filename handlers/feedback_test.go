package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-frontend-feedback-controller/email/emailtest"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces/interfacestest"
	"github.com/ONSdigital/dp-frontend-models/model/feedback"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_getFeedback(t *testing.T) {
	Convey("Given a request without a query string", t, func() {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		url := "whatever"
		errorType := ""
		description := ""
		name := ""
		email := ""
		lang := "en"
		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		Convey("When getFeedback is called", func() {
			getFeedback(w, req, url, errorType, description, name, email, lang, mockRenderer)
			Convey("Then a 200 request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost?service=dev", nil)
		w := httptest.NewRecorder()
		url := "whatever"
		errorType := ""
		description := ""
		name := ""
		email := ""
		lang := "en"
		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		Convey("When getFeedback is called", func() {
			getFeedback(w, req, url, errorType, description, name, email, lang, mockRenderer)
			Convey("Then the page model is sent to the renderer", func() {
				var expectedPage feedback.Page
				expectedPage.Language = "en"
				expectedPage.Metadata.Title = "Feedback"
				expectedPage.PreviousURL = url
				expectedPage.Metadata.Description = url
				expectedPage.ServiceDescription = "ONS developer website"
				expectedPage.Type = "feedback"

				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})
			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func Test_addFeedback(t *testing.T) {
	Convey("Given a valid request", t, func() {
		req := httptest.NewRequest("GET", "http://localhost?description=whatever", nil)
		w := httptest.NewRecorder()
		isPositive := false
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

		Convey("When addFeedback is called", func() {
			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to, lang)
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
		req := httptest.NewRequest("GET", "http://localhost?description=whatever", nil)
		w := httptest.NewRecorder()
		isPositive := false
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

		Convey("When addFeedback is called", func() {
			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to, lang)
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
		isPositive := false
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

		Convey("When addFeedback is called", func() {
			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to, lang)
			Convey("Then the renderer is not called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 0)
			})

			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 0)
			})

			Convey("Then a 500 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given a request for feedback with an empty description value", t, func() {
		req := httptest.NewRequest("POST", "http://localhost?service=dev", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		isPositive := false
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
		Convey("When addFeedback is called", func() {
			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to, lang)
			Convey("Then the renderer is called to render the feedback page", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})
			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 0)
			})
			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
	Convey("Given a request for feedback with an invalid email address", t, func() {
		body := strings.NewReader("email=hello&description=hfjkshk")
		req := httptest.NewRequest("POST", "http://localhost?service=dev", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		isPositive := false
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

		Convey("When addFeedback is called", func() {
			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to, lang)
			Convey("Then the renderer is called to render the feedback page", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})
			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 0)
			})
			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func Test_feedbackThanks(t *testing.T) {

	Convey("Given a valid request", t, func() {

		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()
		url := "www.test.com"
		errorType := ""

		mockRenderer := &interfacestest.RendererMock{
			BuildPageFunc: func(w io.Writer, pageModel interface{}, templateName string) {},
			NewBasePageModelFunc: func() coreModel.Page {
				return coreModel.Page{}
			},
		}

		Convey("When feedbackThanks is called", func() {

			feedbackThanks(w, req, url, errorType, mockRenderer)

			Convey("Then the renderer is called", func() {
				So(len(mockRenderer.BuildPageCalls()), ShouldEqual, 1)
			})

			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
