package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bparsons094/go-server-base/database"
	"github.com/bparsons094/go-server-base/models"
	"github.com/bparsons094/go-server-base/utils"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type UserConnections struct {
	Connections []*websocket.Conn
}

type WebSocketService struct {
	ConnectedUsers         sync.Map
	ConnectionLastResponse sync.Map
	ReceiveNotify          chan interface{}
}

type WebSocketMessage struct {
	Type       string       `json:"type"`
	Payload    interface{}  `json:"payload"`
	Authorized bool         `json:"authorized"`
	User       *models.User `json:"user"`
}

func NewWebSocketService() *WebSocketService {
	receiveChannel := make(chan interface{})
	service := &WebSocketService{
		ConnectedUsers: sync.Map{},
		ReceiveNotify:  receiveChannel,
	}

	go service.startKeepAlive()

	return service
}

func (s *WebSocketService) HandleWebSocketConnection(c *websocket.Conn) {
	user, err := s.AuthUser(c)
	if err != nil {
		s.sendMessage(c, WebSocketMessage{
			Type:       "connection",
			Payload:    err.Error(),
			Authorized: false,
		})
		c.Close()
		return
	}

	s.onConnect(c, user)
	defer s.onDisconnect(c, user)

	// Handle incoming WebSocket connections here
	for {
		messageType, msg, err := c.ReadMessage()

		if messageType != websocket.TextMessage {
			log.Println("Message Type is not TextMessage", messageType)
			break
		}

		if err != nil {
			log.Println("Read Message Error: ", err)
			break
		}

		var message WebSocketMessage
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Error Unmarshalling Message: ", err)
			continue
		}

		// Handle Different Message Types
		switch message.Type {
		case "connection":
			// s.handleConnection(message, c)
		case "message":
			// s.handleMessage(message, c)
		case "pong":
			s.ConnectionLastResponse.Store(c, time.Now().UTC())
		}
	}
}

func (s *WebSocketService) AuthUser(c *websocket.Conn) (models.User, error) {
	token := c.Query("token")

	if token == "" {
		return models.User{}, errors.New("Missing token")
	}

	sub, err := utils.ValidateToken(token)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	DB := database.GetDatabase()
	if err := DB.Where("id = ?", sub).First(&user).Error; err != nil {
		return models.User{}, errors.New("User not found")
	}

	return user, nil
}

func (s *WebSocketService) onConnect(c *websocket.Conn, user models.User) {

	userConnInterface, _ := s.ConnectedUsers.LoadOrStore(user.ID, &UserConnections{})

	userConnections := userConnInterface.(*UserConnections)
	userConnections.Connections = append(userConnections.Connections, c)

	// Store the connection's last response time
	s.ConnectionLastResponse.Store(c, time.Now().UTC())

	// Send welcome message to the connected user
	welcomeMessage := WebSocketMessage{
		Type:       "connection",
		Payload:    "Welcome to The WebSocket!",
		Authorized: true,
		User:       &user,
	}
	s.sendMessage(c, welcomeMessage)

	// Broadcast new user connection message
	broadcastMessage := WebSocketMessage{
		Type:       "connection",
		Payload:    "A new user has connected!",
		Authorized: true,
		User:       &user,
	}
	s.broadcast(broadcastMessage)
}

func (s *WebSocketService) broadcast(data WebSocketMessage) {
	s.ConnectedUsers.Range(func(key, value interface{}) bool {
		userConnections := value.(*UserConnections)
		for _, conn := range userConnections.Connections {
			s.sendMessage(conn, data)
		}
		return true
	})
}

func (s *WebSocketService) sendMessage(c *websocket.Conn, data WebSocketMessage) {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Println("Error Marshalling Message: ", err)
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("Error Sending Message: ", err)
	}
}

func (s *WebSocketService) onDisconnect(c *websocket.Conn, user models.User) {
	userConnInterface, ok := s.ConnectedUsers.Load(user.ID)
	if !ok {
		log.Printf("No connections found for user: %v", user.ID)
		return
	}

	userConnections := userConnInterface.(*UserConnections)

	// Synchronize access to the user's connections
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	for i, conn := range userConnections.Connections {
		if conn == c {
			userConnections.Connections[i] = userConnections.Connections[len(userConnections.Connections)-1]
			userConnections.Connections = userConnections.Connections[:len(userConnections.Connections)-1]
			break
		}
	}

	err := c.Close()
	if err != nil {
		log.Printf("Error closing connection for user %v: %v", user.ID, err)
	}
}

func (s *WebSocketService) startKeepAlive() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.handleKeepAlive()
		}
	}
}

func (s *WebSocketService) handleKeepAlive() {
	// Loop through all connected users and send a ping message to each connection every 15 seconds.
	s.ConnectedUsers.Range(func(key, value interface{}) bool {
		userConnections := value.(*sync.Map)
		userConnections.Range(func(subKey, subValue interface{}) bool {
			userConns := subValue.(*UserConnections)
			for _, conn := range userConns.Connections {

				lastResponseTime, ok := s.ConnectionLastResponse.Load(conn)
				if !ok {
					log.Println("Error getting last response time for connection")
					continue
				}

				if time.Since(lastResponseTime.(time.Time)) > time.Minute {
					log.Println("Connection has not responded in over a minute, closing connection for user", subKey)
					s.onDisconnect(conn, models.User{ID: subKey.(uuid.UUID)})
					continue
				}

				connMemoryAddress := fmt.Sprintf("%p", conn)
				s.sendMessage(conn, WebSocketMessage{
					Type:       "ping",
					Payload:    connMemoryAddress,
					Authorized: true,
				})
			}

			return true
		})

		return true
	})
}
