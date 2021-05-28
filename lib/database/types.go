package database

import (
	"database/sql"
	"time"
)

// Config holds all the configuration for the db
type Config struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslmode  string
}

// Database contains all the methods to interact with DB
type Database struct {
	config *Config
	db     *sql.DB
}

// Color is a color legend
type Color struct {
	Color  string `json:"color"`
	Legend string `json:"legend"`
}

// RetrospectiveUser aka user
type RetrospectiveUser struct {
	UserID   string `json:"id"`
	UserName string `json:"name"`
	Active   bool   `json:"active"`
}

// Retrospective A story mapping board
type Retrospective struct {
	RetrospectiveID   string                 `json:"id" db:"id"`
	OwnerID           string                 `json:"ownerId" db:"ownder_id"`
	RetrospectiveName string                 `json:"name" db:"name"`
	Users             []*RetrospectiveUser   `json:"users"`
	Items             []*RetrospectiveItem   `json:"items"`
	Actions           []*RetrospectiveAction `json:"actions"`
	Phase             int                    `json:"phase" db:"phase"`
}

// RetrospectiveItem can be a pro (went well/worked), con (needs improvement), or a question
type RetrospectiveItem struct {
	ID              string   `json:"id" db:"id"`
	RetrospectiveID string   `json:"retrospectiveId" db:"retrospective_id"`
	UserID          string   `json:"userId" db:"user_id"`
	ParentID        string   `json:"parentId" db:"parent_id"`
	Content         string   `json:"content" db:"content"`
	Type            string   `json:"type" db:"type"`
	Votes           []string `json:"votes" db:"votes"`
}

// RetrospectiveAction is an action the team can take based on retro feedback
type RetrospectiveAction struct {
	ID              string `json:"id" db:"id"`
	RetrospectiveID string `json:"retrospectiveId" db:"retrospective_id"`
	Content         string `json:"content" db:"content"`
	Completed       bool   `json:"completed" db:"completed"`
}

// User aka user
type User struct {
	UserID     string `json:"id"`
	UserName   string `json:"name"`
	UserEmail  string `json:"email"`
	UserAvatar string `json:"avatar"`
	UserType   string `json:"type"`
	Verified   bool   `json:"verified"`
	Country    string `json:"country"`
	Company    string `json:"company"`
	JobTitle   string `json:"jobTitle"`
}

// APIKey structure
type APIKey struct {
	ID          string    `json:"id"`
	Prefix      string    `json:"prefix"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

// ApplicationStats includes user, retrospective counts
type ApplicationStats struct {
	RegisteredCount      int `json:"registeredUserCount"`
	UnregisteredCount    int `json:"unregisteredUserCount"`
	RetrospectiveCount   int `json:"retrospectiveCount"`
	OrganizationCount    int `json:"organizationCount"`
	DepartmentCount      int `json:"departmentCount"`
	TeamCount            int `json:"teamCount"`
	APIKeyCount          int `json:"apikeyCount"`
	ActiveRetroCount     int `json:"activeRetroCount"`
	ActiveRetroUserCount int `json:"activeRetroUserCount"`
}

// Organization can be a company
type Organization struct {
	OrganizationID string `json:"id"`
	Name           string `json:"name"`
	CreatedDate    string `json:"createdDate"`
	UpdatedDate    string `json:"updatedDate"`
}

type OrganizationUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Department struct {
	DepartmentID string `json:"id"`
	Name         string `json:"name"`
	CreatedDate  string `json:"createdDate"`
	UpdatedDate  string `json:"updatedDate"`
}

type DepartmentUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Team struct {
	TeamID      string `json:"id"`
	Name        string `json:"name"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type TeamUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Alert struct {
	AlertID        string `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Type           string `json:"type" db:"type"`
	Content        string `json:"content" db:"content"`
	Active         bool   `json:"active" db:"active"`
	AllowDismiss   bool   `json:"allowDismiss" db:"allow_dismiss"`
	RegisteredOnly bool   `json:"registeredOnly" db:"registered_only"`
	CreatedDate    string `json:"createdDate" db:"created_date"`
	UpdatedDate    string `json:"updatedDate" db:"updated_date"`
}
