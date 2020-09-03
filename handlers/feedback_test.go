package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_getFeedback(t *testing.T) {

	Convey("Given a request without a query string", t, func() {

		Convey("When getFeedback is called", func() {
			req := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			url := "whatever"
			errorType := ""
			purpose := ""
			description := ""
			name := ""
			email := ""

			mockRenderer := MockRenderer{}

			getFeedback(w, req, url, errorType, purpose, description, name, email, mockRenderer)

			Convey("Then a 400 bad request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given a valid request", t, func() {

		Convey("When getFeedback is called", func() {
			req := httptest.NewRequest("GET", "http://localhost?service=dev", nil)
			w := httptest.NewRecorder()

			url := "whatever"
			errorType := ""
			purpose := ""
			description := ""
			name := ""
			email := ""

			mockRenderer := MockRenderer{}

			getFeedback(w, req, url, errorType, purpose, description, name, email, mockRenderer)

			Convey("Then a 200 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})

}

type MockRenderer struct {
}

func (r MockRenderer) Do(path string, b []byte) ([]byte, error) {
	return nil, nil
}
