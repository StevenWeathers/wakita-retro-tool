package database

import (
	"log"
)

// CreateRetroAction adds a new action to the retrospective
func (d *Database) CreateRetrospectiveAction(RetrospectiveID string, UserID string, Content string) ([]*RetrospectiveAction, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_action VALUES (retrospective_id = $1, content = $2);`, RetrospectiveID, Content,
	); err != nil {
		log.Println(err)
	}

	actions := d.GetRetrospectiveActions(RetrospectiveID)

	return actions, nil
}

// DeleteRetrospectiveAction removes a goal from the current board by ID
func (d *Database) DeleteRetrospectiveAction(RetrospectiveID string, userID string, ActionID string) ([]*RetrospectiveAction, error) {
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
		`SELECT id, retrospective_id, content, completed FROM retrospective_action WHERE id = $1;`,
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
