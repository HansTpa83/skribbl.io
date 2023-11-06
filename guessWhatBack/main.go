package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Room struct {
	ID            string
	Name          string
	GuessWord     string
	CurrentDrawer string
	drawers       []string
	Messages      []MessageRoom
	Points        map[string]int
	canvasPoints  []CanvasUpdate
	usernames     map[string]string
	conns         map[string]*websocket.Conn
}

type Server struct {
	rooms   map[string]*Room
	clients map[string]*websocket.Conn
}
type CanvasUpdate struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
type MessageRoom struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
type ChatResponse struct {
	Type string `json:"type"`
	Data struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	} `json:"data"`
}

func NewServer() *Server {
	return &Server{
		rooms:   make(map[string]*Room),
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *Server) NewRoom(name string) *Room {
	roomID := generateRoomID()
	room := &Room{
		ID:        roomID,
		Name:      name,
		conns:     make(map[string]*websocket.Conn),
		usernames: make(map[string]string),
		Points:    make(map[string]int),
		Messages:  make([]MessageRoom, 0),
	}
	s.rooms[roomID] = room
	return room
}

func (s *Server) handleWS(ws *websocket.Conn) {
	connID := uuid.New().String()

	fmt.Println("New connection from client", ws.RemoteAddr())

	s.clients[connID] = ws

	for {
		messageJSON := make([]byte, 1024)
		n, err := ws.Read(messageJSON)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of connection")
				break
			}
			fmt.Println("Error reading message:", err)
			continue
		}

		messageJSON = messageJSON[:n]

		var message Message
		err = json.Unmarshal(messageJSON, &message)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}

		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error encoding message:", err)
			continue
		}

		fmt.Println("Message JSON:", string(jsonBytes))

		switch message.Type {
		case "room_create":
			roomName := message.Data.(map[string]interface{})["room"].(string)
			username := message.Data.(map[string]interface{})["username"].(string)
			fmt.Println("Room : ", roomName, " | username : ", username)

			room := s.NewRoom(roomName)

			room.conns[connID] = ws
			room.usernames[connID] = username

			room.Points[username] = 0

			fmt.Println("Username set for client", ws.RemoteAddr(), "in room", room.ID)

			response := Message{
				Type: "room_created",
				Data: map[string]interface{}{
					"roomId":   room.ID,
					"admin":    true,
					"username": username,
					"roomName": roomName,
				},
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}

			if _, err := ws.Write(responseJSON); err != nil {
				fmt.Println("Write error:", err)
			}

			s.sendRoomInfo(room.ID, connID)
			s.sendUsername(room.ID, connID)
		case "room_join":
			roomID := message.Data.(map[string]interface{})["room"].(string)
			room, ok := s.rooms[roomID]
			if !ok {
				room = s.NewRoom(roomID)
			}

			username := message.Data.(map[string]interface{})["username"].(string)
			room.conns[connID] = ws
			room.usernames[connID] = username
			room.Points[username] = 0

			fmt.Println("Username set for client", ws.RemoteAddr(), "in room", roomID)

			response := Message{
				Type: "room_join",
				Data: map[string]interface{}{
					"roomId":       room.ID,
					"username":     username,
					"roomName":     room.Name,
					"roomMessages": room.Messages,
				},
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}

			if _, err := ws.Write(responseJSON); err != nil {
				fmt.Println("Write error:", err)
			}
			s.sendRoomInfo(room.ID, connID)
			s.sendUsername(room.ID, connID)
		case "chat":
			roomID := message.Data.(map[string]interface{})["room"].(string)
			room, ok := s.rooms[roomID]
			if !ok {
				fmt.Println("Room not found:", roomID)
				continue
			}

			username, ok := room.usernames[connID]
			if !ok {
				fmt.Println("Username not found for client", ws.RemoteAddr())
				continue
			}

			guessWord := room.GuessWord
			content := message.Data.(map[string]interface{})["content"].(string)

			newMessage := MessageRoom{
				Username: username,
				Content:  content,
			}

			room.Messages = append(room.Messages, newMessage)

			if content == guessWord {
				room.Points[username] += 1
				nextDrawerID := ""
				for userID := range room.usernames {
					// Vérifier si l'utilisateur n'est pas déjà dans la liste des "drawers"
					if !contains(room.drawers, userID) {
						nextDrawerID = userID
						room.drawers = append(room.drawers, nextDrawerID)
						break
					}
				}

				// Vérifier si un nouveau "drawer" a été sélectionné
				if nextDrawerID != "" {
					nextDrawerUsername := room.usernames[nextDrawerID]
					room.GuessWord = getRandomWord()
					// Envoyer les informations du nouveau "drawer" aux clients
					responseDrawer := Message{
						Type: "game_started",
						Data: map[string]interface{}{
							"room":      roomID,
							"drawer":    nextDrawerUsername,
							"guessWord": room.GuessWord,
						},
					}

					responseJSONDrawer, err := json.Marshal(responseDrawer)
					if err != nil {
						fmt.Println("Error encoding response:", err)
						continue
					}

					if _, err := room.conns[nextDrawerID].Write(responseJSONDrawer); err != nil {
						fmt.Println("Write error:", err)
					}

					responseOthers := Message{
						Type: "game_started",
						Data: map[string]interface{}{
							"room":   roomID,
							"drawer": nextDrawerUsername,
						},
					}

					responseJSONOthers, err := json.Marshal(responseOthers)
					if err != nil {
						fmt.Println("Error encoding response:", err)
						continue
					}

					for connID, ws := range room.conns {
						if connID != nextDrawerID {
							if _, err := ws.Write(responseJSONOthers); err != nil {
								fmt.Println("Write error:", err)
							}
						}
					}
				}

				// Réponse pour indiquer que l'utilisateur a trouvé le mot
				response := ChatResponse{
					Type: "message_received",
					Data: struct {
						Username string `json:"username"`
						Content  string `json:"content"`
					}{
						Username: "SERVER",
						Content:  "Mot trouvé par " + username,
					},
				}

				responseJSON, err := json.Marshal(response)
				if err != nil {
					fmt.Println("Error encoding response:", err)
					continue
				}

				room.broadcast(responseJSON)

				if nextDrawerID == "" {
					room.canvasPoints = make([]CanvasUpdate, 0)
					room.GuessWord = ""
					room.CurrentDrawer = ""
					room.drawers = []string{}

					// Envoie du classement à la room
					classement := calculateClassement(room.Points)
					responseClassement := Message{
						Type: "end_game",
						Data: map[string]interface{}{
							"room":          roomID,
							"canvas":        room.canvasPoints,
							"GuessWord":     room.GuessWord,
							"classement":    classement[0],
							"CurrentDrawer": room.CurrentDrawer,
						},
					}

					fmt.Println("responseClassement :", responseClassement)

					responseJSONClassement, err := json.Marshal(responseClassement)
					if err != nil {
						fmt.Println("Error encoding response:", err)
						continue
					}

					room.broadcast(responseJSONClassement)

					room.Points = make(map[string]int)
				}

			} else {
				// Réponse pour indiquer que l'utilisateur a envoyé un message
				response := ChatResponse{
					Type: "message_received",
					Data: struct {
						Username string `json:"username"`
						Content  string `json:"content"`
					}{
						Username: username,
						Content:  content,
					},
				}

				responseJSON, err := json.Marshal(response)
				if err != nil {
					fmt.Println("Error encoding response:", err)
					continue
				}

				room.broadcast(responseJSON)
			}

		case "canvas_update":
			roomID := message.Data.(map[string]interface{})["room"].(string)
			room, ok := s.rooms[roomID]
			if !ok {
				fmt.Println("Room not found:", roomID)
				continue
			}

			pointsData, ok := message.Data.(map[string]interface{})["points"]
			if !ok || pointsData == nil {
				fmt.Println("Invalid canvas update message: points data is nil")
				continue
			}

			// Accéder aux coordonnées des points individuellement
			updateData, ok := pointsData.(map[string]interface{})
			if !ok {
				fmt.Println("Invalid canvas update message:", pointsData)
				continue
			}

			start, ok := updateData["start"].(map[string]interface{})
			if !ok {
				fmt.Println("Invalid start coordinates:", updateData["start"])
				continue
			}

			end, ok := updateData["end"].(map[string]interface{})
			if !ok {
				fmt.Println("Invalid end coordinates:", updateData["end"])
				continue
			}

			// Construire les structures Point à partir des données
			startPoint := Point{
				X: int(start["x"].(float64)),
				Y: int(start["y"].(float64)),
			}
			endPoint := Point{
				X: int(end["x"].(float64)),
				Y: int(end["y"].(float64)),
			}

			// // Ajouter les coordonnées des points au tableau dans la salle
			room.canvasPoints = append(room.canvasPoints, CanvasUpdate{Start: startPoint, End: endPoint})

			response := Message{
				Type: "canvas_update",
				Data: map[string]interface{}{
					"points": room.canvasPoints,
					"room":   roomID,
				},
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}
			room.broadcast(responseJSON)
		case "reset_canvas":
			roomID := message.Data.(map[string]interface{})["room"].(string)
			room, ok := s.rooms[roomID]
			if !ok {
				fmt.Println("Room not found:", roomID)
				continue
			}

			room.Points = make(map[string]int)

			response := Message{
				Type: "reset_canvas",
				Data: map[string]interface{}{
					"points": room.Points,
					"room":   roomID,
				},
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}
			room.broadcast(responseJSON)
		case "start_game":
			fmt.Println("start_game !!!")

			roomID := message.Data.(map[string]interface{})["room"].(string)
			room, ok := s.rooms[roomID]
			if !ok {
				fmt.Println("Room not found:", roomID)
				continue
			}

			// Set the guess word
			room.GuessWord = getRandomWord()

			// Get the list of user IDs in the room
			userIDs := make([]string, 0, len(room.usernames))
			for userID := range room.usernames {
				userIDs = append(userIDs, userID)
			}

			// Select the drawer as the first user in the room
			drawerID := userIDs[0]
			drawerUsername := room.usernames[drawerID]

			room.drawers = append(room.drawers, drawerID)

			responseDrawer := Message{
				Type: "game_started",
				Data: map[string]interface{}{
					"room":      roomID,
					"drawer":    drawerUsername,
					"guessWord": room.GuessWord,
				},
			}

			responseJSONDrawer, err := json.Marshal(responseDrawer)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}

			if _, err := room.conns[drawerID].Write(responseJSONDrawer); err != nil {
				fmt.Println("Write error:", err)
			}

			responseOthers := Message{
				Type: "game_started",
				Data: map[string]interface{}{
					"room":   roomID,
					"drawer": drawerUsername,
				},
			}

			responseJSONOthers, err := json.Marshal(responseOthers)
			if err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}

			for connID, ws := range room.conns {
				if connID != drawerID {
					if _, err := ws.Write(responseJSONOthers); err != nil {
						fmt.Println("Write error:", err)
					}
				}
			}

		default:
			fmt.Println("Unknown message type:", message.Type)
		}
	}

	err := ws.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
	}

	// Supprimer la connexion WebSocket de la liste des clients
	delete(s.clients, connID)
}

func (r *Room) broadcast(b []byte) {
	fmt.Println("Broadcasting message to room", r.ID)

	if len(r.conns) == 0 {
		fmt.Println("No connections found in room", r.ID)
		return
	}

	for connID, ws := range r.conns {
		go func(connID string, ws *websocket.Conn) {
			message := fmt.Sprintf("%s", string(b))

			if _, err := ws.Write([]byte(message)); err != nil {
				fmt.Println("Write error:", err)
			}
		}(connID, ws)
	}
}

func generateRoomID() string {
	id := uuid.New().String()
	return id
}
func contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}
func getRandomWord() string {
	words := [5]string{"voiture", "carre", "triangle", "avion", "moto"}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(words))
	return words[index]
}
func calculateClassement(points map[string]int) []string {
	// Création d'une liste de joueurs
	players := make([]string, 0, len(points))
	for player := range points {
		players = append(players, player)
	}

	// Trie des joueurs en fonction de leurs points
	sort.Slice(players, func(i, j int) bool {
		return points[players[i]] > points[players[j]]
	})

	return players
}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.handleWS))

	// Nouvelle route pour obtenir la liste des connexions WebSocket
	http.HandleFunc("/sockets", func(w http.ResponseWriter, r *http.Request) {
		sockets := server.GetSockets()

		for _, socket := range sockets {
			fmt.Fprintln(w, socket.RemoteAddr().String())
		}
	})

	go func() {
		err := http.ListenAndServe("localhost:8000", nil)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	fmt.Println("Server started on localhost:8000")

	select {}
}

func (s *Server) GetSockets() []*websocket.Conn {
	var sockets []*websocket.Conn

	for _, conn := range s.clients {
		sockets = append(sockets, conn)
	}

	return sockets
}

func (s *Server) sendRoomInfo(roomID string, connId string) {
	room, ok := s.rooms[roomID]
	if !ok {
		fmt.Println("Room not found:", roomID)
		return
	}

	usernames := make([]string, 0, len(room.usernames))
	for _, username := range room.usernames {
		usernames = append(usernames, username)
	}

	message := struct {
		Type     string        `json:"type"`
		Users    []string      `json:"users"`
		RoomID   string        `json:"roomId"`
		Messages []MessageRoom `json:"message"`
	}{
		Type:     "room_info",
		Users:    usernames,
		RoomID:   roomID,
		Messages: room.Messages,
	}

	responseJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		return
	}

	for _, conn := range room.conns {
		if _, err := conn.Write(responseJSON); err != nil {
			fmt.Println("Write error:", err)
		}
	}
}

func (s *Server) sendUsername(roomID string, connId string) {
	room, ok := s.rooms[roomID]
	if !ok {
		fmt.Println("Room not found:", roomID)
		return
	}

	username, ok := room.usernames[connId]
	if !ok {
		fmt.Println("Username not found for client")
		return
	}

	message := struct {
		Type     string `json:"type"`
		Username string `json:"username"`
	}{
		Type:     "username",
		Username: username,
	}

	responseJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		return
	}

	conn, ok := room.conns[connId]
	if !ok {
		fmt.Println("Connection not found for client:", connId)
		return
	}

	if _, err := conn.Write(responseJSON); err != nil {
		fmt.Println("Write error:", err)
	}
}
