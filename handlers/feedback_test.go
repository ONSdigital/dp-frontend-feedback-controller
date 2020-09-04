package handlers

import (
	"errors"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces/interfacestest"
	"net/http"
	"net/http/httptest"
	"testing"

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
