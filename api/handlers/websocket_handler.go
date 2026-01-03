package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler gère les connexions WebSocket pour temps réel
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]string // conn -> joueurID
}

// NewWebSocketHandler crée un nouveau handler WebSocket
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Accepter toutes les origines pour le MVP
				// En prod : vérifier l'origine Angular
				return true
			},
		},
		clients: make(map[*websocket.Conn]string),
	}
}

// RegisterRoutes enregistre les routes WebSocket
func (h *WebSocketHandler) RegisterRoutes(router *gin.RouterGroup) {
	ws := router.Group("/ws")
	{
		// Connexion WebSocket générale
		ws.GET("/connect", h.HandleWebSocket)
		
		// WebSocket spécifique aux combats
		ws.GET("/combat/:combat_id", h.HandleCombatWebSocket)
	}
}

// WebSocketMessage structure des messages WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`    // "combat_update", "notification", "error", etc.
	Data    interface{} `json:"data"`
	From    string      `json:"from,omitempty"`
	To      string      `json:"to,omitempty"`
	CombatID string     `json:"combat_id,omitempty"`
}

// HandleWebSocket gère les connexions WebSocket générales
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Erreur upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// TODO: Récupérer joueurID depuis token/session
	joueurID := c.Query("joueur_id")
	if joueurID == "" {
		joueurID = "anonymous"
	}

	h.clients[conn] = joueurID
	defer delete(h.clients, conn)

	log.Printf("Client WebSocket connecté: %s", joueurID)

	// Envoyer message de bienvenue
	welcomeMsg := WebSocketMessage{
		Type: "welcome",
		Data: map[string]interface{}{
			"message":   "Connexion WebSocket établie",
			"joueur_id": joueurID,
		},
	}
	conn.WriteJSON(welcomeMsg)

	// Écouter les messages du client
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erreur lecture WebSocket: %v", err)
			break
		}

		// Router les messages selon le type
		h.handleMessage(conn, joueurID, &msg)
	}

	log.Printf("Client WebSocket déconnecté: %s", joueurID)
}

// HandleCombatWebSocket gère les WebSocket spécifiques aux combats
func (h *WebSocketHandler) HandleCombatWebSocket(c *gin.Context) {
	combatID := c.Param("combat_id")
	if combatID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Combat ID required",
		})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Erreur upgrade WebSocket combat: %v", err)
		return
	}
	defer conn.Close()

	joueurID := c.Query("joueur_id")
	if joueurID == "" {
		joueurID = "anonymous"
	}

	log.Printf("Client connecté au combat %s: %s", combatID, joueurID)

	// Écouter les messages spécifiques au combat
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erreur lecture WebSocket combat: %v", err)
			break
		}

		msg.CombatID = combatID
		h.handleCombatMessage(conn, joueurID, &msg)
	}
}

// handleMessage route les messages WebSocket généraux
func (h *WebSocketHandler) handleMessage(conn *websocket.Conn, joueurID string, msg *WebSocketMessage) {
	switch msg.Type {
	case "ping":
		response := WebSocketMessage{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": "now",
			},
		}
		conn.WriteJSON(response)

	case "join_combat":
		// TODO: Joindre un combat
		response := WebSocketMessage{
			Type: "error",
			Data: map[string]string{
				"message": "Combat join not implemented yet",
			},
		}
		conn.WriteJSON(response)

	case "notification_read":
		// TODO: Marquer notification comme lue
		log.Printf("Notification marked as read by %s", joueurID)

	default:
		log.Printf("Message WebSocket non géré: %s", msg.Type)
	}
}

// handleCombatMessage route les messages WebSocket de combat
func (h *WebSocketHandler) handleCombatMessage(conn *websocket.Conn, joueurID string, msg *WebSocketMessage) {
	switch msg.Type {
	case "combat_action":
		// TODO: Exécuter action de combat
		response := WebSocketMessage{
			Type: "error",
			Data: map[string]string{
				"message": "Combat actions not implemented yet",
			},
		}
		conn.WriteJSON(response)

	case "get_combat_status":
		// TODO: Retourner statut du combat
		response := WebSocketMessage{
			Type: "combat_status",
			Data: map[string]interface{}{
				"combat_id": msg.CombatID,
				"status":    "not_implemented",
			},
		}
		conn.WriteJSON(response)

	default:
		log.Printf("Message combat WebSocket non géré: %s", msg.Type)
	}
}

// BroadcastToCombat diffuse un message à tous les participants d'un combat
func (h *WebSocketHandler) BroadcastToCombat(combatID string, message *WebSocketMessage) {
	message.CombatID = combatID
	
	// TODO: Filtrer les clients par combat
	for conn := range h.clients {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Erreur envoi WebSocket: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

// SendToJoueur envoie un message à un joueur spécifique
func (h *WebSocketHandler) SendToJoueur(joueurID string, message *WebSocketMessage) {
	for conn, clientID := range h.clients {
		if clientID == joueurID {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Erreur envoi WebSocket à %s: %v", joueurID, err)
				conn.Close()
				delete(h.clients, conn)
			}
			break
		}
	}
}