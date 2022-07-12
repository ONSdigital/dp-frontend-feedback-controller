package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces"
	"github.com/ONSdigital/dp-frontend-models/model"
	"github.com/ONSdigital/dp-frontend-models/model/feedback"
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
func FeedbackThanks(renderer interfaces.Renderer) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		feedbackThanks(w, req, renderer)
	})
}

func feedbackThanks(w http.ResponseWriter, req *http.Request, renderer interfaces.Renderer) {
	var p model.Page
	ctx := req.Context()

	p.Metadata.Title = "Thank you"
	returnTo := req.URL.Query().Get("returnTo")

	if returnTo == "Whole site" || returnTo == "" {
		returnTo = "https://www.ons.gov.uk"
	}
	p.Metadata.Description = returnTo

	b, err := json.Marshal(p)
	if err != nil {
		log.Error(ctx, "unable to marshal page data", err, log.Data{"setting-response-status": http.StatusInternalServerError})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templateHTML, err := renderer.Do("feedback-thanks", b)
	if err != nil {
		log.Error(ctx, "failed to render feedback-thanks template", err, log.Data{"setting-response-status": http.StatusInternalServerError})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(templateHTML)
}

// GetFeedback handles the loading of a feedback page
func GetFeedback(renderer interfaces.Renderer) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getFeedback(w, req, req.Referer(), "", "", "", "", lang, renderer)
	})
}

func getFeedback(w http.ResponseWriter, req *http.Request, url, errorType, description, name, email, lang string, renderer interfaces.Renderer) {
	var p feedback.Page

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
	p.Email = email
	p.PreviousURL = url

	b, err := json.Marshal(p)
	if err != nil {
		log.Error(req.Context(), "unable to marshal feedback page data", err, log.Data{"setting-response-status": http.StatusInternalServerError})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templateHTML, err := renderer.Do("feedback", b)
	if err != nil {
		log.Error(req.Context(), "failed to render feedback template", err, log.Data{"setting-response-status": http.StatusInternalServerError})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(templateHTML)
}

// AddFeedback handles a users feedback request and sends a message to slack
func AddFeedback(to, from string, isPositive bool, renderer interfaces.Renderer, emailSender email.Sender) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		addFeedback(w, req, isPositive, renderer, emailSender, from, to, lang)
	})
}

func addFeedback(w http.ResponseWriter, req *http.Request, isPositive bool, renderer interfaces.Renderer, emailSender email.Sender, from, to, lang string) {
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
		getFeedback(w, req, f.URL, "description", f.Description, f.Name, f.Email, lang, renderer)
		return
	}

	if len(f.Email) > 0 && !isPositive {
		if ok, err := regexp.MatchString(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$`, f.Email); !ok || err != nil {
			getFeedback(w, req, f.URL, "email", f.Description, f.Name, f.Email, lang, renderer)
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

	redirectURL := "/feedback/thanks?returnTo=" + f.URL
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
