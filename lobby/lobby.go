package lobby

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Global list of game lobbies.
var lobbies = make(map[string]*Lobby)

// A collection of players for a game.
type Lobby struct {
	ID   string
	Game Game

	Lock sync.RWMutex   // Used to synchronize the adding/removal of players.
	Wait sync.WaitGroup // Used to wait for the lobby to fill up with players.

	// Maps player IDs (unique to the lobby) to their socket connections for sending and receiving.
	Players map[string]*Player
}

// A player connected to a game lobby.
type Player struct {
	Socket *websocket.Conn
	Active bool // Whether the connection is initialized/not timed out.

	Lock *sync.RWMutex // Used to synchronize reading and setting the Active field.
}

// Represents a game instance. Used by lobbies to enable different types of games.
type Game interface {
	// Takes a player identifier string (unique to this game instance, format depends on the game),
	// and returns a receiver to handle messages from the player,
	// or an error if adding the player failed.
	AddPlayer(playerID string) (interface {
		HandleMessage(msgType string, msg []byte)
	}, error)

	// Returns the range of possible player IDs for this game.
	PlayerIDs() []string

	// Starts the game.
	Start()
}

// Signature for functions that construct a game instance.
// Takes the lobby to which players can connect,
// and an untyped options parameter that can be parsed by the game instance for use in setup.
type GameConstructor func(lobby *Lobby, options any) (Game, error)

// Creates and registers a new lobby with the given ID,
// and uses the given constructor to construct its game instance.
// Returns error if lobby ID is already taken, or if game construction failed.
func New(id string, gameConstructor GameConstructor) (*Lobby, error) {
	if id == "" {
		return nil, errors.New("lobby name cannot be blank")
	}

	lobby := Lobby{
		ID: id,
	}

	game, err := gameConstructor(&lobby, nil)
	if err != nil {
		return nil, err
	}

	lobby.Game = game
	playerIDs := game.PlayerIDs()
	lobby.AddPlayerSlots(playerIDs)

	err = RegisterLobby(&lobby)
	if err != nil {
		return nil, err
	}

	return &lobby, nil
}

// Takes the given list of player IDs and adds slots for each of them in the lobby.
// Adds the length of the given IDs to the lobby's wait group, so it can be used to wait for the lobby to fill up.
func (lobby *Lobby) AddPlayerSlots(playerIDs []string) {
	lobby.Players = make(map[string]*Player, len(playerIDs))
	for _, playerID := range playerIDs {
		lobby.Players[playerID] = &Player{
			Lock: new(sync.RWMutex),
		}
	}
	lobby.Wait.Add(len(playerIDs))
}

// Registers a lobby in the global list of lobbies.
// Returns error if lobby with same ID already exists.
func RegisterLobby(lobby *Lobby) error {
	if _, ok := lobbies[lobby.ID]; ok {
		return errors.New("lobby with ID \"" + lobby.ID + "\" already exists")
	}

	lobbies[lobby.ID] = lobby
	return nil
}

// Returns the player in the lobby corresponding to the given player ID,
// or false if none is found.
func (lobby *Lobby) GetPlayer(playerID string) (interface{ Send(msg any) error }, bool) {
	lobby.Lock.RLock()
	defer lobby.Lock.RUnlock()
	player, ok := lobby.Players[playerID]
	return player, ok
}

// Sets the player in the lobby corresponding to the given player ID.
// Returns an error if no matching player is found.
func (lobby *Lobby) setPlayer(playerID string, player Player) error {
	lobby.Lock.Lock()
	defer lobby.Lock.Unlock()

	if _, ok := lobby.Players[playerID]; !ok {
		return errors.New("invalid player ID")
	}

	lobby.Players[playerID] = &player
	return nil
}

// Returns a player's Active flag in a thread-safe manner.
func (player *Player) isActive() bool {
	player.Lock.RLock()
	defer player.Lock.RUnlock()
	return player.Active
}

// Sets a player's Active flag in a thread-safe manner.
func (player *Player) setActive(active bool) {
	player.Lock.Lock()
	defer player.Lock.Unlock()
	player.Active = active
}

// Marshals the given message to JSON and sends it over the player's socket connection.
// Returns an error if the player is inactive, or if the marshaling/sending failed.
func (player *Player) Send(msg any) error {
	if !player.isActive() {
		return errors.New("cannot send to inactive player")
	}

	err := player.Socket.WriteJSON(msg)
	return err
}

// Marshals the given message to JSON and sends it to all connected players.
// Returns an error if it failed to marshal or send to at least one of the players.
func (lobby *Lobby) Send(msg any) error {
	marshaledMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, player := range lobby.Players {
		err = player.Socket.WriteMessage(websocket.TextMessage, marshaledMsg)
	}

	return err
}

// Returns the current connected players in a lobby, and the max number of potential players.
func (lobby *Lobby) PlayerCount() (current int, max int) {
	for _, player := range lobby.Players {
		if player.isActive() {
			current++
		}
	}

	max = len(lobby.Players)

	return current, max
}

// Returns a map of player IDs to whether they are taken (true if taken).
func (lobby *Lobby) AvailablePlayerIDs() map[string]bool {
	available := make(map[string]bool)

	for id, player := range lobby.Players {
		if player.isActive() {
			available[id] = true
		} else {
			available[id] = false
		}
	}

	return available
}

// Removes a lobby from the lobby map and closes its connections.
func (lobby *Lobby) Close() error {
	for id, player := range lobby.Players {
		player.Socket.Close()
		player.setActive(false)
		lobby.setPlayer(id, Player{})
	}
	delete(lobbies, lobby.ID)

	return nil
}

// If there is only 1 lobby on the server, returns that,
// otherwise returns lobby corresponding to lobby parameter in request.
// Returns error on absent lobby parameter or lobby not found.
func findLobby(req *http.Request) (*Lobby, error) {
	if len(lobbies) == 1 {
		for _, lobby := range lobbies {
			return lobby, nil
		}
	}

	params, ok := checkParams(req, "lobby")
	if !ok {
		return nil, errors.New("lacking lobby query parameter")
	}

	lobbyID := params.Get("lobby")
	lobby, ok := lobbies[lobbyID]
	if !ok {
		return nil, errors.New("no lobby found with provided id")
	}

	return lobby, nil
}
