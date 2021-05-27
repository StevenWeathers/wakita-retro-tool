package database

import (
	"database/sql"
	"encoding/json"
	"log"
)

// CreateRetroItem adds a new item to the retrospective
func (d *Database) CreateRetrospectiveItem(RetrospectiveID string, UserID string, Type string, Content string) ([]*RetrospectiveItem, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_item VALUES (retrospective_id = $1, type = $2, content = $3);`, RetrospectiveID, Type, Content,
	); err != nil {
		log.Println(err)
	}

	items := d.GetRetrospectiveItems(RetrospectiveID)

	return items, nil
}

// DeleteRetrospectiveItem removes a goal from the current board by ID
func (d *Database) DeleteRetrospectiveItem(RetrospectiveID string, userID string, ItemID string) ([]*RetrospectiveItem, error) {
	if _, err := d.db.Exec(
		`DELETE FROM retrospective_item WHERE id = $1;`, ItemID); err != nil {
		log.Println(err)
	}

	items := d.GetRetrospectiveItems(RetrospectiveID)

	return items, nil
}

// GetRetrospectiveItems retrieves retrospective items from the DB
func (d *Database) GetRetrospectiveItems(RetrospectiveID string) []*RetrospectiveItem {
	var items = make([]*RetrospectiveItem, 0)

	itemRows, itemsErr := d.db.Query(
		`SELECT id, retrospective_id, user_id, parent_id, content, votes, type FROM retrospective_item WHERE id = $1;`,
		RetrospectiveID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var votes string
			var parentId sql.NullString
			var ri = &RetrospectiveItem{
				ID:              "",
				RetrospectiveID: "",
				UserID:          "",
				ParentID:        "",
				Content:         "",
				Type:            "",
				Votes:           make([]string, 0),
			}
			if err := itemRows.Scan(&ri.ID, &ri.RetrospectiveID, &ri.UserID, &parentId, &ri.Content, &votes, &ri.Type); err != nil {
				log.Println(err)
			} else {
				ri.ParentID = parentId.String
				jsonErr := json.Unmarshal([]byte(votes), &ri.Votes)
				if jsonErr != nil {
					log.Println(jsonErr)
				}
				items = append(items, ri)
			}
		}
	}

	return items
}
