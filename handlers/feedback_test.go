package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-frontend-feedback-controller/email/emailtest"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces/interfacestest"
	"github.com/ONSdigital/dp-frontend-models/model/feedback"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_getFeedback(t *testing.T) {

	Convey("Given a request without a query string", t, func() {

		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		url := "whatever"
		errorType := ""
		purpose := ""
		description := ""
		name := ""
		email := ""

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, nil
			},
		}

		Convey("When getFeedback is called", func() {

			getFeedback(w, req, url, errorType, purpose, description, name, email, mockRenderer)

			Convey("Then a 400 bad request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given a valid request", t, func() {

		req := httptest.NewRequest("GET", "http://localhost?service=dev", nil)
		w := httptest.NewRecorder()

		url := "whatever"
		errorType := ""
		purpose := "the purpose"
		description := ""
		name := ""
		email := ""

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, nil
			},
		}

		Convey("When getFeedback is called", func() {

			getFeedback(w, req, url, errorType, purpose, description, name, email, mockRenderer)

			Convey("Then the expected JSON is sent to the renderer", func() {

				var expectedPage feedback.Page
				expectedPage.Purpose = purpose
				expectedPage.Metadata.Title = "Feedback"
				expectedPage.PreviousURL = url
				expectedPage.Metadata.Description = url
				expectedPage.ServiceDescription = "ONS developer website"
				expectedJSON, _ := json.Marshal(expectedPage)

				actualJSON := string(mockRenderer.DoCalls()[0].B)

				So(mockRenderer.DoCalls()[0].Path, ShouldEqual, "feedback")
				So(actualJSON, ShouldEqual, string(expectedJSON))
			})

			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

		})
	})

	Convey("Given an error returned from the renderer", t, func() {

		req := httptest.NewRequest("GET", "http://localhost?service=dev", nil)
		w := httptest.NewRecorder()

		url := "whatever"
		errorType := ""
		purpose := ""
		description := ""
		name := ""
		email := ""

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, errors.New("renderer is broken")
			},
		}

		Convey("When getFeedback is called", func() {

			getFeedback(w, req, url, errorType, purpose, description, name, email, mockRenderer)

			Convey("Then a 500 internal server error response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
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

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, nil
			},
		}

		mockSender := &emailtest.SenderMock{
			SendFunc: func(from string, to []string, msg []byte) error {
				return nil
			},
		}

		Convey("When addFeedback is called", func() {

			addFeedback(w, req, isPositive, mockRenderer, mockSender, from, to)

			Convey("Then the renderer is not called", func() {
				So(len(mockRenderer.DoCalls()), ShouldEqual, 0)
			})

			Convey("Then the email sender is called", func() {
				So(len(mockSender.SendCalls()), ShouldEqual, 1)
			})

			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusMovedPermanently)
			})
		})
	})

}

func Test_feedbackThanks(t *testing.T) {

	Convey("Given a valid request", t, func() {

		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, nil
			},
		}

		Convey("When feedbackThanks is called", func() {

			feedbackThanks(w, req, mockRenderer)

			Convey("Then the renderer is called", func() {
				So(len(mockRenderer.DoCalls()), ShouldEqual, 1)
			})

			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})

	Convey("Given an error returned from the renderer", t, func() {

		req := httptest.NewRequest("GET", "http://localhost", nil)
		w := httptest.NewRecorder()

		mockRenderer := &interfacestest.RendererMock{
			DoFunc: func(path string, b []byte) ([]byte, error) {
				return nil, errors.New("renderer is broken")
			},
		}

		Convey("When feedbackThanks is called", func() {

			feedbackThanks(w, req, mockRenderer)

			Convey("Then a 500 internal server error response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
