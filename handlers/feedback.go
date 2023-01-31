package handlers

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/schema"
)

// Feedback represents a user's feedback
type Feedback struct {
	Type             string `schema:"type"`
	URI              string `schema:":uri"`
	URL              string `schema:"url"`
	Description      string `schema:"description"`
	Name             string `schema:"name"`
	Email            string `schema:"email"`
	FeedbackFormType string `schema:"feedback-form-type"`
}

// FeedbackThanks loads the Feedback Thank you page
func FeedbackThanks(rend interfaces.Renderer) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		feedbackThanks(w, req, req.Referer(), "", rend)
	})
}

func feedbackThanks(w http.ResponseWriter, req *http.Request, url, errorType string, rend interfaces.Renderer) {
	ctx := req.Context()
	basePage := rend.NewBasePageModel()
	p := model.Feedback{
		Page: basePage,
	}

	var wholeSite string
	cfg, err := config.Get()
	if err != nil {
		log.Warn(ctx, "Unable to retrieve configuration", log.FormatErrors([]error{err}))
	} else {
		wholeSite = cfg.SiteDomain
	}

	p.Type = "feedback"
	p.Metadata.Title = "Thank you"
	p.ErrorType = errorType
	p.PreviousURL = url

	// returnTo is redered on page so needs XSS protection
	returnTo := html.EscapeString(req.URL.Query().Get("returnTo"))
	if returnTo == "Whole site" {
		returnTo = wholeSite
	} else if returnTo == "" {
		returnTo = url
	}

	p.Metadata.Description = returnTo

	rend.BuildPage(w, p, "feedback-thanks")
}

// GetFeedback handles the loading of a feedback page
func GetFeedback(rend interfaces.Renderer) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getFeedback(w, req, req.Referer(), "", "", "", "", lang, rend)
	})
}

func getFeedback(w http.ResponseWriter, req *http.Request, url, errorType, description, name, userEmail, lang string, rend interfaces.Renderer) {
	basePage := rend.NewBasePageModel()
	p := model.Feedback{
		Page: basePage,
	}

	var services = make(map[string]string)
	services["cmd"] = "Customising data by applying filters"
	services["dev"] = "ONS developer website"

	p.ServiceDescription = services[req.URL.Query().Get("service")]

	p.Language = lang
	p.Type = "feedback"
	p.Metadata.Title = "Feedback"
	p.Metadata.Description = url

	if len(p.Metadata.Description) > 50 {
		p.Metadata.Description = p.Metadata.Description[len(p.Metadata.Description)-50 : len(p.Metadata.Description)]
	}

	p.ErrorType = errorType
	p.Feedback = description
	p.Name = name
	p.Email = userEmail
	p.PreviousURL = url

	rend.BuildPage(w, p, "feedback")
}

// AddFeedback handles a users feedback request and sends a message to slack
func AddFeedback(to, from string, isPositive bool, rend interfaces.Renderer, emailSender email.Sender) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		addFeedback(w, req, isPositive, rend, emailSender, from, to, lang)
	})
}

func addFeedback(w http.ResponseWriter, req *http.Request, isPositive bool, rend interfaces.Renderer, emailSender email.Sender, from, to, lang string) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		log.Error(ctx, "unable to parse request form", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f Feedback
	if err := decoder.Decode(&f, req.Form); err != nil {
		log.Error(ctx, "unable to decode request form", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if f.Description == "" && !isPositive {
		getFeedback(w, req, f.URL, "description", f.Description, f.Name, f.Email, lang, rend)
		return
	}

	if len(f.Email) > 0 && !isPositive {
		if ok, err := regexp.MatchString(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$`, f.Email); !ok || err != nil {
			getFeedback(w, req, f.URL, "email", f.Description, f.Name, f.Email, lang, rend)
			return
		}
	}

	if f.URL == "" {
		f.URL = "Whole site"
	}

	if err := emailSender.Send(
		from,
		[]string{to},
		generateFeedbackMessage(f, from, to, isPositive),
	); err != nil {
		log.Error(ctx, "failed to send message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	returnTo := f.URL

	if returnTo == "Whole site" || returnTo == "" {
		returnTo = "https://www.ons.gov.uk"
	}

	redirectURL := "/feedback/thanks?returnTo=" + returnTo
	http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
}

func generateFeedbackMessage(f Feedback, from, to string, isPositive bool) []byte {
	var description string
	if isPositive {
		description = "Positive feedback received"
	} else {
		description = f.Description
	}

	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("From: %s\n", from))
	b.WriteString(fmt.Sprintf("To: %s\n", to))
	b.WriteString(fmt.Sprintf("Subject: Feedback received\n\n"))

	if len(f.Type) > 0 {
		b.WriteString(fmt.Sprintf("Feedback Type: %s\n", f.Type))
	}

	b.WriteString(fmt.Sprintf("Page URL: %s\n", f.URL))
	b.WriteString(fmt.Sprintf("Description: %s\n", description))

	if len(f.Name) > 0 {
		b.WriteString(fmt.Sprintf("Name: %s\n", f.Name))
	}

	if len(f.Email) > 0 {
		b.WriteString(fmt.Sprintf("Email address: %s\n", f.Email))
	}

	return b.Bytes()
}
