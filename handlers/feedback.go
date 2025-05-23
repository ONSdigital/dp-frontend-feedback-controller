package handlers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	feedbackAPIModel "github.com/ONSdigital/dp-feedback-api/models"
	feedbackAPI "github.com/ONSdigital/dp-feedback-api/sdk"
	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces"
	"github.com/ONSdigital/dp-frontend-feedback-controller/mapper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/model"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	core "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/schema"
)

// FeedbackThanks loads the Feedback Thank you page
func (f *Feedback) FeedbackThanks() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		feedbackThanks(w, req, req.Referer(), f.Render, f.CacheService, lang, f.Config.SiteDomain, f.Config.EnableNewNavBar)
	})
}

func feedbackThanks(w http.ResponseWriter, req *http.Request, uri string, rend interfaces.Renderer, cacheHelperService *cacheHelper.Helper, lang, siteDomain string, enableNewNavBar bool) {
	basePage := rend.NewBasePageModel()
	p := mapper.CreateGetFeedbackThanks(req, basePage, lang, uri, siteDomain)

	if enableNewNavBar {
		ctx := req.Context()
		mappedNavContent, err := cacheHelperService.GetMappedNavigationContent(ctx, lang)
		if err == nil {
			p.NavigationContent = mappedNavContent
		}
	}

	rend.BuildPage(w, p, "feedback-thanks")
}

// GetFeedback handles the loading of a feedback page
func (f *Feedback) GetFeedback() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getFeedback(w, req, []core.ErrorItem{}, model.FeedbackForm{URL: req.Referer()}, lang, f.Render, f.CacheService, f.Config.EnableNewNavBar)
	})
}

func getFeedback(w http.ResponseWriter, req *http.Request, validationErrors []core.ErrorItem, ff model.FeedbackForm, lang string, rend interfaces.Renderer, cacheHelperService *cacheHelper.Helper, enableNewNavBar bool) {
	basePage := rend.NewBasePageModel()
	p := mapper.CreateGetFeedback(req, basePage, validationErrors, ff, lang)

	if enableNewNavBar {
		ctx := context.Background()
		mappedNavContent, err := cacheHelperService.GetMappedNavigationContent(ctx, lang)
		if err == nil {
			p.NavigationContent = mappedNavContent
		}
	}

	rend.BuildPage(w, p, "feedback")
}

// AddFeedback handles a users feedback request
func (f *Feedback) AddFeedback() http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		addFeedback(w, req, f.Render, lang, f.Config.SiteDomain, f.CacheService, f.Config)
	})
}

func addFeedback(w http.ResponseWriter, req *http.Request, rend interfaces.Renderer, lang, siteDomain string, cacheService *cacheHelper.Helper, cfg *config.Config) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		log.Error(ctx, "unable to parse request form", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var ff model.FeedbackForm
	if err := decoder.Decode(&ff, req.Form); err != nil {
		log.Error(ctx, "unable to decode request form", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validationErrors := validateForm(&ff, siteDomain)
	if len(validationErrors) > 0 {
		getFeedback(w, req, validationErrors, ff, lang, rend, cacheService, false)
		return
	}

	feedbackAPIClient := feedbackAPI.New(cfg.APIRouterURL)

	isPageUsefulVal := false
	var isGeneralFeedbackVal bool

	if ff.Type == mapper.WholeSite {
		isGeneralFeedbackVal = true
	} else {
		isGeneralFeedbackVal = false
	}

	f := &feedbackAPIModel.Feedback{
		IsPageUseful:      &isPageUsefulVal,
		IsGeneralFeedback: &isGeneralFeedbackVal,
		OnsURL:            ff.URL,
		Feedback:          ff.Description,
		Name:              ff.Name,
		EmailAddress:      ff.Email,
	}

	opts := feedbackAPI.Options{AuthToken: cfg.ServiceAuthToken}

	err := feedbackAPIClient.PostFeedback(ctx, f, opts)

	if err != nil {
		statusCode := err.Status()
		log.Error(ctx, "failed to send feedback", err, log.Data{"code": statusCode})
		return
	}

	returnTo := ff.URL

	if returnTo == mapper.WholeSite || returnTo == "" {
		returnTo = "https://www.ons.gov.uk"
	}

	redirectURL := fmt.Sprintf("/feedback/thanks?returnTo=%s", returnTo)
	http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
}

// validateForm is a helper function that validates a slice of FeedbackForm to determine if there are form validation errors
func validateForm(ff *model.FeedbackForm, siteDomain string) (validationErrors []core.ErrorItem) {
	if ff.Type == "" && ff.FormLocation != "footer" {
		validationErrors = append(validationErrors, core.ErrorItem{
			Description: core.Localisation{
				LocaleKey: "FeedbackChooseType",
				Plural:    1,
			},
			URL: "#type-error",
		})
		ff.IsTypeErr = true
	}

	ff.URL = strings.TrimSpace(ff.URL)
	if ff.Type == mapper.ASpecificPage {
		if ff.URL == "" {
			validationErrors = append(validationErrors, core.ErrorItem{
				Description: core.Localisation{
					LocaleKey: "FeedbackWhatEnterURL",
					Plural:    1,
				},
				URL: "#type-error",
			})
			ff.IsURLErr = true
		} else if !mapper.IsSiteDomainURL(ff.URL, siteDomain) {
			validationErrors = append(validationErrors, core.ErrorItem{
				Description: core.Localisation{
					LocaleKey: "FeedbackValidURL",
					Plural:    1,
				},
				URL: "#type-error",
			})
			ff.IsURLErr = true
		}
	} else if ff.Type != mapper.ASpecificPage && ff.URL != "" {
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

	if ff.Email != "" {
		if ok, err := regexp.MatchString("^[A-Za-z0-9.`!#$%&'*+-/=?^_{|}~]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}$", ff.Email); !ok || err != nil {
			validationErrors = append(validationErrors, core.ErrorItem{
				Description: core.Localisation{
					LocaleKey: "FeedbackAlertEmail",
					Plural:    1,
				},
				URL: "#email-error",
			})
			ff.IsEmailErr = true
		}
	}
	return validationErrors
}
