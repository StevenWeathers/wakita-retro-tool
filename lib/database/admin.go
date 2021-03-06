package database

import (
	"errors"
	"log"
	"strings"
)

// ConfirmAdmin confirms whether the user is infact a ADMIN
func (d *Database) ConfirmAdmin(AdminID string) error {
	var userType string
	e := d.db.QueryRow("SELECT coalesce(type, '') FROM users WHERE id = $1;", AdminID).Scan(&userType)
	if e != nil {
		log.Println(e)
		return errors.New("could not find users type")
	}

	if userType != "ADMIN" {
		return errors.New(("user is not an admin"))
	}

	return nil
}

// GetAppStats gets counts of users (registered and unregistered), and retrospectives
func (d *Database) GetAppStats() (*ApplicationStats, error) {
	var Appstats ApplicationStats

	statsErr := d.db.QueryRow(`
		SELECT
			unregistered_user_count,
			registered_user_count,
			retrospective_count,
			organization_count,
			department_count,
			team_count,
			apikey_count
		FROM get_app_stats();
		`,
	).Scan(
		&Appstats.UnregisteredCount,
		&Appstats.RegisteredCount,
		&Appstats.RetrospectiveCount,
		&Appstats.OrganizationCount,
		&Appstats.DepartmentCount,
		&Appstats.TeamCount,
		&Appstats.APIKeyCount,
	)
	if statsErr != nil {
		log.Println("Unable to get application stats: ", statsErr)
		return nil, statsErr
	}

	return &Appstats, nil
}

// PromoteUser promotes a user to ADMIN type
func (d *Database) PromoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call promote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to promote user to ADMIN")
	}

	return nil
}

// DemoteUser demotes a user to REGISTERED type
func (d *Database) DemoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call demote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to demote user to REGISTERED")
	}

	return nil
}

// CleanRetrospectives deletes retrospectives older than X days
func (d *Database) CleanRetrospectives(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_retrospectives($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean retrospectives")
	}

	return nil
}

// CleanGuests deletes guest users older than X days
func (d *Database) CleanGuests(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_guest_users($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean Guest Warriors")
	}

	return nil
}

// OrganizationList gets a list of organizations
func (d *Database) OrganizationList(Limit int, Offset int) []*Organization {
	var organizations = make([]*Organization, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM organization_list($1, $2);`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org Organization

			if err := rows.Scan(
				&org.OrganizationID,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		log.Println(err)
	}

	return organizations
}

// TeamList gets a list of teams
func (d *Database) TeamList(Limit int, Offset int) []*Team {
	var teams = make([]*Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM team_list($1, $2);`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team Team

			if err := rows.Scan(
				&team.TeamID,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		log.Println(err)
	}

	return teams
}

// GetAPIKeys gets a list of api keys
func (d *Database) GetAPIKeys(Limit int, Offset int) []*APIKey {
	var APIKeys = make([]*APIKey, 0)
	rows, err := d.db.Query(
		`SELECT apk.id, apk.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM api_keys apk
		LEFT JOIN users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT $1
		OFFSET $2;`,
		Limit,
		Offset,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserID,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.ID = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys
}
