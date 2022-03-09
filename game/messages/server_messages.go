package messages

// Messages from server to client.
const (
	MessageAskSupport         = "askSupport"
	MessageOrdersReceived     = "ordersReceived"
	MessageOrdersConfirmation = "ordersConfirmation"
)

// Message sent from server when asking a supporting player who to support in an embattled area.
type AskSupport struct {
	Base
	From     string   `json:"from"`
	To       string   `json:"to"`
	Battlers []string `json:"battlers"` // List of possible players to support in the battle.
}

// Message sent from server to all clients when valid orders are received from all players.
type OrdersReceived struct {
	Base
	Orders map[string][]Order `json:"orders"` // Maps a player's ID to their submitted orders.
}

// Message sent from server to all clients when valid orders are received from a player.
// Used to show who the server is waiting for.
type OrdersConfirmation struct {
	Base
	Player string `json:"player"` // The player who submitted orders.
}

// Message sent from server to all clients when the results for a round are calculated.
type RoundResult struct {
	Base
	// Maps a player's ID to their submitted orders, now with added Status field.
	Orders map[string][]OrderWithStatus `json:"orders"`
	// Maps area names to the chronological list of battles that took place in that area.
	Battles map[string][]Battle `json:"battles"`
}

// Order message with added Status field for showing the server's calculated result.
type OrderWithStatus struct {
	Order
	Status string `json:"status"`
}

type Battle struct {
	Results    []Result
	DangerZone string
}

type Result struct {
	Total        int
	Parts        []Modifier
	Move         Order
	DefenderArea string
}

type Modifier struct {
	Type        string
	Value       int
	SupportFrom string
}
