package database

import (
	"database/sql"
	"encoding/json"
	"log"
)

// CreateRetrospectiveItemWorked adds a worked item to the retrospective
func (d *Database) CreateRetrospectiveItemWorked(RetrospectiveID string, UserID string, Content string) ([]*RetrospectiveItem, error) {
	var Type string = "worked"
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_item
		(retrospective_id, type, content, user_id)
		VALUES ($1,$2, $3, $4);`,
		RetrospectiveID, Type, Content, UserID,
	); err != nil {
		log.Println(err)
	}

	items, _, _ := d.GetRetrospectiveItems(RetrospectiveID)

	return items, nil
}

// CreateRetrospectiveItemImprove adds a improve item to the retrospective
func (d *Database) CreateRetrospectiveItemImprove(RetrospectiveID string, UserID string, Content string) ([]*RetrospectiveItem, error) {
	var Type string = "improve"
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_item
		(retrospective_id, type, content, user_id)
		VALUES ($1,$2, $3, $4);`,
		RetrospectiveID, Type, Content, UserID,
	); err != nil {
		log.Println(err)
	}

	_, items, _ := d.GetRetrospectiveItems(RetrospectiveID)

	return items, nil
}

// CreateRetrospectiveItemQuestion adds a question item to the retrospective
func (d *Database) CreateRetrospectiveItemQuestion(RetrospectiveID string, UserID string, Content string) ([]*RetrospectiveItem, error) {
	var Type string = "question"
	if _, err := d.db.Exec(
		`INSERT INTO retrospective_item
		(retrospective_id, type, content, user_id)
		VALUES ($1,$2, $3, $4);`,
		RetrospectiveID, Type, Content, UserID,
	); err != nil {
		log.Println(err)
	}

	_, _, items := d.GetRetrospectiveItems(RetrospectiveID)

	return items, nil
}

// DeleteRetrospectiveItem removes a goal from the current board by ID
func (d *Database) DeleteRetrospectiveItem(RetrospectiveID string, userID string, ItemID string) (WorkedItems []*RetrospectiveItem, ImproveItems []*RetrospectiveItem, QuestionItems []*RetrospectiveItem, DeleteError error) {
	if _, err := d.db.Exec(
		`DELETE FROM retrospective_item WHERE id = $1;`, ItemID); err != nil {
		log.Println(err)
	}

	workedItems, improveItems, questionItems := d.GetRetrospectiveItems(RetrospectiveID)

	return workedItems, improveItems, questionItems, nil
}

// GetRetrospectiveItems retrieves retrospective items from the DB
func (d *Database) GetRetrospectiveItems(RetrospectiveID string) (Worked []*RetrospectiveItem, Improve []*RetrospectiveItem, Question []*RetrospectiveItem) {
	var itemsWorked = make([]*RetrospectiveItem, 0)
	var itemsImprove = make([]*RetrospectiveItem, 0)
	var itemsQuestion = make([]*RetrospectiveItem, 0)

	itemRows, itemsErr := d.db.Query(
		`SELECT id, retrospective_id, user_id, parent_id, content, votes, type FROM retrospective_item WHERE retrospective_id = $1;`,
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
				if ri.Type == "worked" {
					itemsWorked = append(itemsWorked, ri)
				}
				if ri.Type == "improve" {
					itemsImprove = append(itemsImprove, ri)
				}
				if ri.Type == "question" {
					itemsQuestion = append(itemsQuestion, ri)
				}
			}
		}
	} else {
		log.Println(itemsErr)
	}

	return itemsWorked, itemsImprove, itemsQuestion
}
