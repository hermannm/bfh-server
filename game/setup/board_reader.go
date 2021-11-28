package setup

import (
	"encoding/json"
	"fmt"
	"immerse-ntnu/hermannia/server/game"
	"io/ioutil"
	"os"
)

type board struct {
	Areas     []area     `json:"areas"`
	Neighbors []neighbor `json:"neighbors"`
}

type area struct {
	Name   string           `json:"name"`
	Sea    bool             `json:"sea"`
	Forest bool             `json:"forest"`
	Castle bool             `json:"castle"`
	Home   game.PlayerColor `json:"home"`
}

type neighbor struct {
	Area1      string `json:"area1"`
	Area2      string `json:"area2"`
	River      bool   `json:"river"`
	Cliffs     bool   `json:"cliffs"`
	DangerZone string `json:"dangerZone"`
}

func ReadBoard(players int) (game.Board, error) {
	file, err := os.Open(fmt.Sprintf("./game/setup/board_%dplayers.json", players))
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonBoard board

	err = json.Unmarshal(bytes, &jsonBoard)
	if err != nil {
		return nil, err
	}

	board := make(game.Board)

	for _, jsonArea := range jsonBoard.Areas {
		area := game.BoardArea{
			Name:             jsonArea.Name,
			Control:          jsonArea.Home,
			Home:             jsonArea.Home,
			Sea:              jsonArea.Sea,
			Forest:           jsonArea.Forest,
			Castle:           jsonArea.Castle,
			Neighbors:        make([]game.Neighbor, 0),
			IncomingMoves:    make([]*game.Order, 0),
			IncomingSupports: make([]*game.Order, 0),
			Combats:          make([]game.Combat, 0),
		}

		board[area.Name] = &area
	}

	for _, jsonNeighbor := range jsonBoard.Neighbors {
		area1, ok1 := board[jsonNeighbor.Area1]
		area2, ok2 := board[jsonNeighbor.Area2]

		if !ok1 || !ok2 {
			return nil, fmt.Errorf(
				"error in board config: neighbor relation %s <-> %s",
				jsonNeighbor.Area1,
				jsonNeighbor.Area2,
			)
		}

		area1.Neighbors = append(area1.Neighbors, game.Neighbor{
			Area:       area2,
			River:      jsonNeighbor.River,
			Cliffs:     jsonNeighbor.Cliffs,
			DangerZone: jsonNeighbor.DangerZone,
		})

		area2.Neighbors = append(area2.Neighbors, game.Neighbor{
			Area:       area1,
			River:      jsonNeighbor.River,
			Cliffs:     jsonNeighbor.Cliffs,
			DangerZone: jsonNeighbor.DangerZone,
		})
	}

	return board, nil
}