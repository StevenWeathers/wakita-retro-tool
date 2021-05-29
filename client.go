package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// SocketEvent is the event structure used for socket messages
type SocketEvent struct {
	EventType  string `json:"type"`
	EventValue string `json:"value"`
	EventUser  string `json:"userId"`
}

// CreateSocketEvent makes a SocketEvent struct and turns it into json []byte
func CreateSocketEvent(EventType string, EventValue string, EventUser string) []byte {
	newEvent := &SocketEvent{
		EventType:  EventType,
		EventValue: EventValue,
		EventUser:  EventUser,
	}

	event, _ := json.Marshal(newEvent)

	return event
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump(srv *server) {
	var forceClosed bool
	c := s.conn
	defer func() {
		RetrospectiveID := s.arena
		UserID := s.userID

		Users := srv.database.RetreatUser(RetrospectiveID, UserID)
		updatedUsers, _ := json.Marshal(Users)

		retreatEvent := CreateSocketEvent("user_retreated", string(updatedUsers), UserID)
		m := message{retreatEvent, RetrospectiveID}
		h.broadcast <- m

		h.unregister <- s
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(writeWait)); err != nil {
				log.Printf("abandon error: %v", err)
			}
		}
		if err := c.ws.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var badEvent bool
		var targetedEvent bool
		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors
		userID := s.userID
		retrospectiveID := s.arena

		switch keyVal["type"] {
		case "create_item_worked":
			var rs struct {
				Content string `json:"content"`
				Phase   int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, err := srv.database.CreateRetrospectiveItemWorked(retrospectiveID, userID, rs.Content)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_worked_updated", string(updatedItems), "")
		case "create_item_improve":
			var rs struct {
				Content string `json:"content"`
				Phase   int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, err := srv.database.CreateRetrospectiveItemImprove(retrospectiveID, userID, rs.Content)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_improve_updated", string(updatedItems), "")
		case "create_item_question":
			var rs struct {
				Content string `json:"content"`
				Phase   int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, err := srv.database.CreateRetrospectiveItemQuestion(retrospectiveID, userID, rs.Content)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_question_updated", string(updatedItems), "")
		case "nest_item_worked":
			var rs struct {
				ItemID   string `json:"id"`
				ParentID string `json:"parentId"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, _, _, err := srv.database.NestRetrospectiveItem(retrospectiveID, userID, rs.ItemID, rs.ParentID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_worked_updated", string(updatedItems), "")
		case "unnest_item_worked":
			var rs struct {
				ItemID string `json:"id"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, _, _, err := srv.database.UnNestRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_worked_updated", string(updatedItems), "")
		case "nest_item_improve":
			var rs struct {
				ItemID   string `json:"id"`
				ParentID string `json:"parentId"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, items, _, err := srv.database.NestRetrospectiveItem(retrospectiveID, userID, rs.ItemID, rs.ParentID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_improve_updated", string(updatedItems), "")
		case "unnest_item_improve":
			var rs struct {
				ItemID string `json:"id"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, items, _, err := srv.database.UnNestRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_improve_updated", string(updatedItems), "")
		case "nest_item_question":
			var rs struct {
				ItemID   string `json:"id"`
				ParentID string `json:"parentId"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, _, items, err := srv.database.NestRetrospectiveItem(retrospectiveID, userID, rs.ItemID, rs.ParentID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_question_updated", string(updatedItems), "")
		case "unnest_item_question":
			var rs struct {
				ItemID string `json:"id"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, _, items, err := srv.database.UnNestRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_question_updated", string(updatedItems), "")
		case "delete_item_worked":
			var rs struct {
				ItemID string `json:"id"`
				Phase  int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			items, _, _, err := srv.database.DeleteRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_worked_updated", string(updatedItems), "")
		case "delete_item_improve":
			var rs struct {
				ItemID string `json:"id"`
				Phase  int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, items, _, err := srv.database.DeleteRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_improve_updated", string(updatedItems), "")
		case "delete_item_question":
			var rs struct {
				ItemID string `json:"id"`
				Phase  int    `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			_, _, items, err := srv.database.DeleteRetrospectiveItem(retrospectiveID, userID, rs.ItemID)
			if err != nil {
				badEvent = true
				break
			}
			if rs.Phase == 1 {
				targetedEvent = true
				items = srv.database.FilterItemsByUser(userID, items)
			}

			updatedItems, _ := json.Marshal(items)
			msg = CreateSocketEvent("item_question_updated", string(updatedItems), "")
		case "create_action":
			var rs struct {
				Content string `json:"content"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			actions, err := srv.database.CreateRetrospectiveAction(retrospectiveID, userID, rs.Content)
			if err != nil {
				badEvent = true
				break
			}

			updatedActions, _ := json.Marshal(actions)
			msg = CreateSocketEvent("action_updated", string(updatedActions), "")
		case "update_action":
			var rs struct {
				ActionID  string `json:"id"`
				Completed bool   `json:"completed"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			actions, err := srv.database.UpdatedRetrospectiveAction(retrospectiveID, userID, rs.ActionID, rs.Completed)
			if err != nil {
				badEvent = true
				break
			}

			updatedActions, _ := json.Marshal(actions)
			msg = CreateSocketEvent("action_updated", string(updatedActions), "")
		case "delete_action":
			var rs struct {
				ActionID string `json:"id"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			actions, err := srv.database.DeleteRetrospectiveAction(retrospectiveID, userID, rs.ActionID)
			if err != nil {
				badEvent = true
				break
			}

			updatedActions, _ := json.Marshal(actions)
			msg = CreateSocketEvent("action_updated", string(updatedActions), "")
		case "advance_phase":
			var rs struct {
				Phase int `json:"phase"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &rs)

			err := srv.database.RetrospectiveAdvancePhase(retrospectiveID, userID, rs.Phase)
			if err != nil {
				badEvent = true
				break
			}

			msg = CreateSocketEvent("phase_updated", strconv.Itoa(rs.Phase), "")
		case "promote_owner":
			retrospective, err := srv.database.SetRetrospectiveOwner(retrospectiveID, userID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}

			updatedRetrospective, _ := json.Marshal(retrospective)
			msg = CreateSocketEvent("retrospective_updated", string(updatedRetrospective), "")
		case "concede_retrospective":
			err := srv.database.DeleteRetrospective(retrospectiveID, userID)
			if err != nil {
				badEvent = true
				break
			}
			msg = CreateSocketEvent("retrospective_conceded", "", "")
		case "abandon_retrospective":
			_, err := srv.database.AbandonRetrospective(retrospectiveID, userID)
			if err != nil {
				badEvent = true
				break
			}
			badEvent = true // don't want this event to cause write panic
			forceClosed = true
		default:
		}

		if !badEvent && !targetedEvent {
			m := message{msg, s.arena}
			h.broadcast <- m
		}

		if targetedEvent {
			c.write(websocket.TextMessage, msg)
		}

		if forceClosed {
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (s *server) serveWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		retrospectiveID := vars["id"]

		// upgrade to WebSocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// make sure user cookies are valid
		userID, cookieErr := s.validateUserCookie(w, r)
		if cookieErr != nil {
			cm := websocket.FormatCloseMessage(4001, "unauthorized")
			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("unauthorized close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}

		// make sure retrospective is legit
		b, retrospectiveErr := s.database.GetRetrospective(retrospectiveID)
		if retrospectiveErr != nil {
			cm := websocket.FormatCloseMessage(4004, "retrospective not found")
			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("not found close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}
		if b.Phase == 1 {
			b.WorkedItems = s.database.FilterItemsByUser(userID, b.WorkedItems)
			b.ImproveItems = s.database.FilterItemsByUser(userID, b.ImproveItems)
			b.QuestionItems = s.database.FilterItemsByUser(userID, b.QuestionItems)
		}
		retrospective, _ := json.Marshal(b)

		// make sure user exists
		_, userErr := s.database.GetRetrospectiveUser(retrospectiveID, userID)

		if userErr != nil {
			log.Println("error finding user : " + userErr.Error() + "\n")
			cm := websocket.FormatCloseMessage(4003, "duplicate session")

			if fmt.Sprint(userErr) == "User Not found" {
				s.clearUserCookies(w)
				cm = websocket.FormatCloseMessage(4001, "unauthorized")
			}

			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("unauthorized close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}

		c := &connection{send: make(chan []byte, 256), ws: ws}
		ss := subscription{c, retrospectiveID, userID}
		h.register <- ss

		Users, _ := s.database.AddUserToRetrospective(ss.arena, userID)
		updatedUsers, _ := json.Marshal(Users)

		initEvent := CreateSocketEvent("init", string(retrospective), userID)
		_ = c.write(websocket.TextMessage, initEvent)

		joinedEvent := CreateSocketEvent("user_joined", string(updatedUsers), userID)
		m := message{joinedEvent, ss.arena}
		h.broadcast <- m

		go ss.writePump()
		go ss.readPump(s)
	}
}
