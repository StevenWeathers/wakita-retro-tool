package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// handleAppStats gets the applications stats
func (s *server) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := s.database.GetAppStats()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ActiveRetroUserCount := 0
		for _, s := range h.arenas {
			ActiveRetroUserCount = ActiveRetroUserCount + len(s)
		}

		AppStats.ActiveRetroCount = len(h.arenas)
		AppStats.ActiveRetroUserCount = ActiveRetroUserCount

		s.respondWithJSON(w, http.StatusOK, AppStats)
	}
}

// handleGetRegisteredUsers gets a list of registered users
func (s *server) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Users := s.database.GetRegisteredUsers(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Users)
	}
}

// handleUserCreate registers a user as a registered user
func (s *server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		UserName, UserEmail, UserPassword, accountErr := ValidateUserAccount(
			keyVal["userName"].(string),
			keyVal["userEmail"].(string),
			keyVal["userPassword1"].(string),
			keyVal["userPassword2"].(string),
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUser, VerifyID, err := s.database.CreateUserRegistered(UserName, UserEmail, UserPassword, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendWelcome(UserName, UserEmail, VerifyID)

		s.respondWithJSON(w, http.StatusOK, newUser)
	}
}

// handleUserPromote handles promoting a user to ADMIN by ID
func (s *server) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		err := s.database.PromoteUser(keyVal["userId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDemote handles demoting a user to REGISTERED by ID
func (s *server) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		err := s.database.DemoteUser(keyVal["userId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanRetrospectives handles cleaning up old retrospectives (ADMIN Manaually Triggered)
func (s *server) handleCleanRetrospectives() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_retrospectives_days_old")

		err := s.database.CleanRetrospectives(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
func (s *server) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_guests_days_old")

		err := s.database.CleanGuests(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleGetOrganizations gets a list of organizations
func (s *server) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := s.database.OrganizationList(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetTeams gets a list of teams
func (s *server) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.TeamList(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetAPIKeys gets a list of APIKeys
func (s *server) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.GetAPIKeys(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}
