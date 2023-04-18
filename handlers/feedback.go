package handlers

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/dp-renderer/helper"
	core "github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/schema"
)

// FeedbackForm represents a user's feedback
type FeedbackForm struct {
	FormLocation     string `schema:"feedback-form-type"`
	Type             string `schema:"type"`
	IsTypeErr        bool   `schema:"is_type_err"`
	URI              string `schema:":uri"`
	URL              string `schema:"url"`
	IsURLErr         bool   `schema:"is_url_err"`
	Description      string `schema:"description"`
	IsDescriptionErr bool   `schema:"is_description_err"`
	Name             string `schema:"name"`
	Email            string `schema:"email"`
	IsEmailErr       bool   `schema:"is_email_err"`
}

// FeedbackThanks loads the Feedback Thank you page
func (f *Feedback) FeedbackThanks() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		feedbackThanks(w, req, req.Referer(), f.Render, f.CacheService, lang)
	})
}

func feedbackThanks(w http.ResponseWriter, req *http.Request, url string, rend interfaces.Renderer, cacheHelperService *cacheHelper.Helper, lang string) {
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
	if cfg.EnableNewNavBar {
		mappedNavContent, err := cacheHelperService.GetMappedNavigationContent(ctx, lang)
		if err == nil {
			p.NavigationContent = mappedNavContent
		}
	}
	p.Type = "feedback"
	p.Metadata.Title = helper.Localise("FeedbackThanks", lang, 1)
	p.PreviousURL = url

	// returnTo is rendered on page so needs XSS protection
	returnTo := html.EscapeString(req.URL.Query().Get("returnTo"))
	if returnTo == "Whole site" {
		returnTo = wholeSite
	} else if returnTo == "" {
		returnTo = url
	}

	p.ReturnTo = returnTo

	rend.BuildPage(w, p, "feedback-thanks")
}

// GetFeedback handles the loading of a feedback page
func (f *Feedback) GetFeedback() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getFeedback(w, req, []core.ErrorItem{}, FeedbackForm{URL: req.Referer()}, lang, f.Render, f.CacheService)
	})
}

func getFeedback(w http.ResponseWriter, req *http.Request, validationErrors []core.ErrorItem, ff FeedbackForm, lang string, rend interfaces.Renderer, cacheHelperService *cacheHelper.Helper) {
	basePage := rend.NewBasePageModel()
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

	p.ServiceDescription = services[req.URL.Query().Get("service")]

	p.Language = lang
	p.Type = "feedback"
	p.Metadata.Title = helper.Localise("FeedbackTitle", lang, 1)
	p.Metadata.Description = ff.URL
	ctx := context.Background()
	cfg, err := config.Get()
	if err != nil {
		log.Warn(ctx, "Unable to retrieve configuration", log.FormatErrors([]error{err}))
	}
	if cfg.EnableNewNavBar {
		mappedNavContent, err := cacheHelperService.GetMappedNavigationContent(ctx, lang)
		if err == nil {
			p.NavigationContent = mappedNavContent
		}
	}

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

	p.Radio = ff.Type
	p.IsRadioErr = ff.IsTypeErr
	p.Feedback = ff.Description
	p.IsFeedbackErr = ff.IsDescriptionErr
	p.Name = ff.Name
	p.Email = ff.Email
	p.IsEmailErr = ff.IsEmailErr
	p.PreviousURL = ff.URL
	p.IsURLErr = ff.IsURLErr

	rend.BuildPage(w, p, "feedback")
}

// AddFeedback handles a users feedback request and sends a message to slack
func (f *Feedback) AddFeedback() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		addFeedback(w, req, f.Render, f.EmailSender, f.Config.FeedbackFrom, f.Config.FeedbackTo, lang, f.CacheService)
	})
}

func addFeedback(w http.ResponseWriter, req *http.Request, rend interfaces.Renderer, emailSender email.Sender, from, to, lang string, cacheService *cacheHelper.Helper) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		log.Error(ctx, "unable to parse request form", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var ff FeedbackForm
	if err := decoder.Decode(&ff, req.Form); err != nil {
		log.Error(ctx, "unable to decode request form", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validationErrors := validateForm(&ff)

	if len(validationErrors) > 0 {
		getFeedback(w, req, validationErrors, ff, lang, rend, cacheService)
		return
	}

	if ff.URL == "" {
		ff.URL = "Whole site"
	}

	if err := emailSender.Send(
		from,
		[]string{to},
		generateFeedbackMessage(ff, from, to),
	); err != nil {
		log.Error(ctx, "failed to send message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	returnTo := ff.URL

	if returnTo == "Whole site" || returnTo == "" {
		returnTo = "https://www.ons.gov.uk"
	}

	redirectURL := fmt.Sprintf("/feedback/thanks?returnTo=%s", returnTo)
	http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
}

// validateForm is a helper function that validates a slice of FeedbackForm to determine if there are form validation errors
func validateForm(ff *FeedbackForm) (validationErrors []core.ErrorItem) {
	if ff.Type == "" && ff.FormLocation != "footer" {
		validationErrors = append(validationErrors, core.ErrorItem{
			Description: core.Localisation{
				LocaleKey: "FeedbackChooseType",
				Plural:    1,
			},
			URL: "#radio-error",
		})
		ff.IsTypeErr = true
	}

	ff.URL = strings.TrimSpace(ff.URL)
	if ff.Type == "A specific page" && ff.URL == "" {
		validationErrors = append(validationErrors, core.ErrorItem{
			Description: core.Localisation{
				LocaleKey: "FeedbackWhatEnterURL",
				Plural:    1,
			},
			URL: "#radio-error",
		})
		ff.IsURLErr = true
	}

	if ff.Type != "A specific page" && ff.URL != "" {
		ff.URL = ""
	}

	ff.Description = strings.TrimSpace(ff.Description)
	if ff.Description == "" {
		validationErrors = append(validationErrors, core.ErrorItem{
			Description: core.Localisation{
				LocaleKey: "FeedbackAlertEntry",
				Plural:    1,
			},
			URL: "#feedback-error",
		})
		ff.IsDescriptionErr = true
	}

	if len(ff.Email) > 0 {
		if ok, err := regexp.MatchString(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$`, ff.Email); !ok || err != nil {
			validationErrors = append(validationErrors, core.ErrorItem{
				Description: core.Localisation{
					LocaleKey: "FeedbackAlertEmail",
					Plural:    1,
				},
				URL: "#email-error",
			})
		}
		ff.IsEmailErr = true
	}
	return validationErrors
}

func generateFeedbackMessage(f FeedbackForm, from, to string) []byte {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("From: %s\n", from))
	b.WriteString(fmt.Sprintf("To: %s\n", to))
	b.WriteString("Subject: Feedback received\n\n")

	if len(f.Type) > 0 {
		b.WriteString(fmt.Sprintf("Feedback Type: %s\n", f.Type))
	}

	b.WriteString(fmt.Sprintf("Page URL: %s\n", f.URL))
	b.WriteString(fmt.Sprintf("Description: %s\n", f.Description))

	if len(f.Name) > 0 {
		b.WriteString(fmt.Sprintf("Name: %s\n", f.Name))
	}

	if len(f.Email) > 0 {
		b.WriteString(fmt.Sprintf("Email address: %s\n", f.Email))
	}

	return b.Bytes()
}
