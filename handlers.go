package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

var ActiveAlerts []interface{}

type contextKey string

var (
	contextKeyUserID         contextKey = "userId"
	apiKeyHeaderName         string     = "X-API-Key"
	contextKeyOrgRole        contextKey = "orgRole"
	contextKeyDepartmentRole contextKey = "departmentRole"
	contextKeyTeamRole       contextKey = "teamRole"
)

type userAccount struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

type userPassword struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// ValidateUserAccount makes sure user name, email, and password are valid before creating the account
func ValidateUserAccount(name string, email string, pwd1 string, pwd2 string) (UserName string, UserEmail string, UserPassword string, validateErr error) {
	v := validator.New()
	a := userAccount{
		Name:      name,
		Email:     email,
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return name, email, pwd1, err
}

// ValidateUserPassword makes sure user password is valid before updating the password
func ValidateUserPassword(pwd1 string, pwd2 string) (UserPassword string, validateErr error) {
	v := validator.New()
	a := userPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return pwd1, err
}

// respondWithJSON takes a payload and writes the response
func (s *server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getJSONRequestBody gets a JSON request body broken into a key/value map
func (s *server) getJSONRequestBody(r *http.Request, w http.ResponseWriter) map[string]interface{} {
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &keyVal) // check for errors

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return keyVal
}

// createUserCookie creates the users cookie
func (s *server) createUserCookie(w http.ResponseWriter, isRegistered bool, UserID string) {
	var cookiedays = 365 // 356 days
	if isRegistered == true {
		cookiedays = 30 // 30 days
	}

	encoded, err := s.cookie.Encode(s.config.SecureCookieName, UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	cookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    encoded,
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   86400 * cookiedays,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearUserCookies wipes the frontend and backend cookies
// used in the event of bad cookie reads
func (s *server) clearUserCookies(w http.ResponseWriter) {
	feCookie := &http.Cookie{
		Name:   s.config.FrontendCookieName,
		Value:  "",
		Path:   s.config.PathPrefix + "/",
		MaxAge: -1,
	}
	beCookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    "",
		Path:     s.config.PathPrefix + "/",
		Domain:   s.config.AppDomain,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
}

// validateUserCookie returns the userID from secure cookies or errors if failures getting it
func (s *server) validateUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var userID string

	if cookie, err := r.Cookie(s.config.SecureCookieName); err == nil {
		var value string
		if err = s.cookie.Decode(s.config.SecureCookieName, cookie.Value, &value); err == nil {
			userID = value
		} else {
			log.Println("error in reading user cookie : " + err.Error() + "\n")
			s.clearUserCookies(w)
			return "", errors.New("invalid user cookies")
		}
	} else {
		log.Println("error in reading user cookie : " + err.Error() + "\n")
		s.clearUserCookies(w)
		return "", errors.New("invalid user cookies")
	}

	return userID, nil
}

/*
	Handlers
*/

// handleIndex parses the index html file, injecting any relevant data
func (s *server) handleIndex() http.HandlerFunc {
	type AppConfig struct {
		AvatarService                string
		ToastTimeout                 int
		AllowGuests                  bool
		AllowRegistration            bool
		DefaultLocale                string
		AuthMethod                   string
		AppVersion                   string
		CookieName                   string
		PathPrefix                   string
		APIEnabled                   bool
		CleanupGuestsDaysOld         int
		CleanupRetrospectivesDaysOld int
		ShowActiveCountries          bool
	}
	type UIConfig struct {
		AnalyticsEnabled bool
		AnalyticsID      string
		AppConfig        AppConfig
		ActiveAlerts     []interface{}
	}

	// get the html template from dist, have it ready for requests
	tmplContent, ioErr := fs.ReadFile(f, "dist/index.html")
	if ioErr != nil {
		log.Println("Error opening index template")
		log.Fatal(ioErr)
	}

	tmplString := string(tmplContent)
	tmpl, tmplErr := template.New("index").Parse(tmplString)
	if tmplErr != nil {
		log.Println("Error parsing index template")
		log.Fatal(tmplErr)
	}

	appConfig := AppConfig{
		AvatarService:                viper.GetString("config.avatar_service"),
		ToastTimeout:                 viper.GetInt("config.toast_timeout"),
		AllowGuests:                  viper.GetBool("config.allow_guests"),
		AllowRegistration:            viper.GetBool("config.allow_registration") && viper.GetString("auth.method") == "normal",
		DefaultLocale:                viper.GetString("config.default_locale"),
		AuthMethod:                   viper.GetString("auth.method"),
		APIEnabled:                   viper.GetBool("config.allow_external_api"),
		AppVersion:                   s.config.Version,
		CookieName:                   s.config.FrontendCookieName,
		PathPrefix:                   s.config.PathPrefix,
		CleanupGuestsDaysOld:         viper.GetInt("config.cleanup_guests_days_old"),
		CleanupRetrospectivesDaysOld: viper.GetInt("config.cleanup_retrospectives_days_old"),
		ShowActiveCountries:          viper.GetBool("config.show_active_countries"),
	}

	ActiveAlerts = s.database.GetActiveAlerts()

	data := UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
		AppConfig:        appConfig,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data.ActiveAlerts = ActiveAlerts // get latest alerts from memory

		tmpl.Execute(w, data)
	}
}

/*
	Retrospective Handlers
*/

// handleRetrospectiveCreate handles creating a retrospective (arena)
func (s *server) handleRetrospectiveCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)

		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			log.Println("error in reading request body: " + bodyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var keyVal struct {
			RetrospectiveName string `json:"retrospectiveName"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newRetrospective, err := s.database.CreateRetrospective(userID, keyVal.RetrospectiveName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if retrospective created with team association
		TeamID, ok := vars["teamId"]
		if ok {
			OrgRole := r.Context().Value(contextKeyOrgRole)
			DepartmentRole := r.Context().Value(contextKeyDepartmentRole)
			TeamRole := r.Context().Value(contextKeyTeamRole).(string)
			var isAdmin bool
			if DepartmentRole != nil && DepartmentRole.(string) == "ADMIN" {
				isAdmin = true
			}
			if OrgRole != nil && OrgRole.(string) == "ADMIN" {
				isAdmin = true
			}

			if isAdmin == true || TeamRole != "" {
				err := s.database.TeamAddRetrospective(TeamID, newRetrospective.RetrospectiveID)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		s.respondWithJSON(w, http.StatusOK, newRetrospective)
	}
}

// handleRetrospectiveGet looks up retrospective or returns notfound status
func (s *server) handleRetrospectiveGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RetrospectiveID := vars["id"]

		retrospective, err := s.database.GetRetrospective(RetrospectiveID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		s.respondWithJSON(w, http.StatusOK, retrospective)
	}
}

// handleRetrospectivesGet looks up retrospectives associated with userID
func (s *server) handleRetrospectivesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserID).(string)

		retrospectives, err := s.database.GetRetrospectivesByUser(userID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		s.respondWithJSON(w, http.StatusOK, retrospectives)
	}
}
