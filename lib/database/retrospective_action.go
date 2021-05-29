package database

import (
	"errors"
	"log"
)

// CreateRetroAction adds a new action to the retrospective
func (d *Database) CreateRetrospectiveAction(RetrospectiveID string, UserID string, Content string) ([]*RetrospectiveAction, error) {
	err := d.ConfirmOwner(RetrospectiveID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`INSERT INTO retrospective_action (retrospective_id, content) VALUES ($1, $2);`, RetrospectiveID, Content,
	); err != nil {
		log.Println(err)
	}

	actions := d.GetRetrospectiveActions(RetrospectiveID)

	return actions, nil
}

// UpdatedRetrospectiveAction updates an actions status
func (d *Database) UpdatedRetrospectiveAction(RetrospectiveID string, userID string, ActionID string, Completed bool) (Actions []*RetrospectiveAction, DeleteError error) {
	err := d.ConfirmOwner(RetrospectiveID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`UPDATE retrospective_action SET completed = $2, updated_date = NOW() WHERE id = $1;`, ActionID, Completed); err != nil {
		log.Println(err)
	}

	actions := d.GetRetrospectiveActions(RetrospectiveID)

	return actions, nil
}

// DeleteRetrospectiveAction removes a goal from the current board by ID
func (d *Database) DeleteRetrospectiveAction(RetrospectiveID string, userID string, ActionID string) ([]*RetrospectiveAction, error) {
	err := d.ConfirmOwner(RetrospectiveID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`DELETE FROM retrospective_action WHERE id = $1;`, ActionID); err != nil {
		log.Println(err)
	}

	actions := d.GetRetrospectiveActions(RetrospectiveID)

	return actions, nil
}

// GetRetrospectiveActions retrieves retrospective actions from the DB
func (d *Database) GetRetrospectiveActions(RetrospectiveID string) []*RetrospectiveAction {
	var actions = make([]*RetrospectiveAction, 0)

	actionRows, actionsErr := d.db.Query(
		`SELECT id, retrospective_id, content, completed FROM retrospective_action WHERE retrospective_id = $1;`,
		RetrospectiveID,
	)
	if actionsErr == nil {
		defer actionRows.Close()
		for actionRows.Next() {
			var ri = &RetrospectiveAction{
				ID:              "",
				RetrospectiveID: "",
				Content:         "",
				Completed:       false,
			}
			if err := actionRows.Scan(&ri.ID, &ri.RetrospectiveID, &ri.Content, &ri.Completed); err != nil {
				log.Println(err)
			} else {
				actions = append(actions, ri)
			}
		}
	}

	return actions
}
