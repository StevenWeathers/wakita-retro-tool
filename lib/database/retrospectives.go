package database

import (
	"errors"
	"log"
)

//CreateRetrospective adds a new retrospective to the db
func (d *Database) CreateRetrospective(OwnerID string, RetrospectiveName string) (*Retrospective, error) {
	var b = &Retrospective{
		RetrospectiveID:   "",
		OwnerID:           OwnerID,
		RetrospectiveName: RetrospectiveName,
		Users:             make([]*RetrospectiveUser, 0),
	}

	e := d.db.QueryRow(
		`SELECT * FROM create_retrospective($1, $2);`,
		OwnerID,
		RetrospectiveName,
	).Scan(&b.RetrospectiveID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Error Creating Retrospective")
	}

	return b, nil
}

// GetRetrospective gets a retrospective by ID
func (d *Database) GetRetrospective(RetrospectiveID string) (*Retrospective, error) {
	var b = &Retrospective{
		RetrospectiveID:   RetrospectiveID,
		OwnerID:           "",
		RetrospectiveName: "",
		Users:             make([]*RetrospectiveUser, 0),
	}

	// get retrospective
	e := d.db.QueryRow(
		`SELECT
			id, name, owner_id
		FROM retrospective WHERE id = $1`,
		RetrospectiveID,
	).Scan(
		&b.RetrospectiveID,
		&b.RetrospectiveName,
		&b.OwnerID,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Not found")
	}

	b.Users = d.GetRetrospectiveUsers(RetrospectiveID)

	return b, nil
}

// GetRetrospectivesByUser gets a list of retrospectives by UserID
func (d *Database) GetRetrospectivesByUser(UserID string) ([]*Retrospective, error) {
	var retrospectives = make([]*Retrospective, 0)
	retrospectiveRows, retrospectivesErr := d.db.Query(`
		SELECT * FROM get_retrospectives_by_user($1);
	`, UserID)
	if retrospectivesErr != nil {
		return nil, errors.New("Not found")
	}

	defer retrospectiveRows.Close()
	for retrospectiveRows.Next() {
		var b = &Retrospective{
			RetrospectiveID:   "",
			OwnerID:           "",
			RetrospectiveName: "",
			Users:             make([]*RetrospectiveUser, 0),
		}
		if err := retrospectiveRows.Scan(
			&b.RetrospectiveID,
			&b.RetrospectiveName,
			&b.OwnerID,
		); err != nil {
			log.Println(err)
		} else {
			retrospectives = append(retrospectives, b)
		}
	}

	return retrospectives, nil
}

// ConfirmOwner confirms the user is infact owner of the retrospective
func (d *Database) ConfirmOwner(RetrospectiveID string, userID string) error {
	var ownerID string
	e := d.db.QueryRow("SELECT owner_id FROM retrospective WHERE id = $1", RetrospectiveID).Scan(&ownerID)
	if e != nil {
		log.Println(e)
		return errors.New("Retrospective Not found")
	}

	if ownerID != userID {
		return errors.New("Not Owner")
	}

	return nil
}

// GetRetrospectiveUser gets a user from db by ID and checks retrospective active status
func (d *Database) GetRetrospectiveUser(RetrospectiveID string, UserID string) (*RetrospectiveUser, error) {
	var active bool
	var w RetrospectiveUser

	e := d.db.QueryRow(
		`SELECT * FROM get_retrospective_user($1, $2);`,
		RetrospectiveID,
		UserID,
	).Scan(
		&w.UserID,
		&w.UserName,
		&active,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("User Not found")
	}

	if active {
		return nil, errors.New("User Already Active in Retrospective")
	}

	return &w, nil
}

// GetRetrospectiveUsers retrieves the users for a given retrospective from db
func (d *Database) GetRetrospectiveUsers(RetrospectiveID string) []*RetrospectiveUser {
	var users = make([]*RetrospectiveUser, 0)
	rows, err := d.db.Query(
		`SELECT * FROM get_retrospective_users($1);`,
		RetrospectiveID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w RetrospectiveUser
			if err := rows.Scan(&w.UserID, &w.UserName, &w.Active); err != nil {
				log.Println(err)
			} else {
				users = append(users, &w)
			}
		}
	}

	return users
}

// AddUserToRetrospective adds a user by ID to the retrospective by ID
func (d *Database) AddUserToRetrospective(RetrospectiveID string, UserID string) ([]*RetrospectiveUser, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_user (retrospective_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (retrospective_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		RetrospectiveID,
		UserID,
	); err != nil {
		log.Println(err)
	}

	users := d.GetRetrospectiveUsers(RetrospectiveID)

	return users, nil
}

// RetreatUser removes a user from the current retrospective by ID
func (d *Database) RetreatUser(RetrospectiveID string, UserID string) []*RetrospectiveUser {
	if _, err := d.db.Exec(
		`UPDATE retrospective_user SET active = false WHERE retrospective_id = $1 AND user_id = $2`, RetrospectiveID, UserID); err != nil {
		log.Println(err)
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		log.Println(err)
	}

	users := d.GetRetrospectiveUsers(RetrospectiveID)

	return users
}

// AbandonRetrospective removes a user from the current retrospective by ID and sets abandoned true
func (d *Database) AbandonRetrospective(RetrospectiveID string, UserID string) ([]*RetrospectiveUser, error) {
	if _, err := d.db.Exec(
		`UPDATE retrospective_user SET active = false, abandoned = true WHERE retrospective_id = $1 AND user_id = $2`, RetrospectiveID, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	users := d.GetRetrospectiveUsers(RetrospectiveID)

	return users, nil
}

// SetRetrospectiveOwner sets the ownerId for the retrospective
func (d *Database) SetRetrospectiveOwner(RetrospectiveID string, userID string, OwnerID string) (*Retrospective, error) {
	err := d.ConfirmOwner(RetrospectiveID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call set_retrospective_owner($1, $2);`, RetrospectiveID, OwnerID); err != nil {
		log.Println(err)
	}

	retrospective, err := d.GetRetrospective(RetrospectiveID)
	if err != nil {
		return nil, errors.New("Unable to promote owner")
	}

	return retrospective, nil
}

// DeleteRetrospective removes all retrospective associations and the retrospective itself from DB by RetrospectiveID
func (d *Database) DeleteRetrospective(RetrospectiveID string, userID string) error {
	err := d.ConfirmOwner(RetrospectiveID, userID)
	if err != nil {
		return errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call delete_retrospective($1);`, RetrospectiveID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
