package main

import (
	"fmt"

	"hermannm.dev/bfh-server/game/board"
	"hermannm.dev/bfh-server/game/boardsetup"
)

func printBoard(brd board.Board, areas map[string]board.Unit, neighbors bool) {
	for _, area := range brd.Areas {
		if _, ok := areas[area.Name]; !ok {
			continue
		}

		areaString := area.Name
		if area.Control != "" {
			areaString += fmt.Sprintf(" (control: %s)", area.Control)
		}
		if area.Home != "" {
			areaString += fmt.Sprintf(" (home: %s)", area.Home)
		}
		if area.Sea {
			areaString += " (sea) "
		}
		if area.Forest {
			areaString += " (forest)"
		}
		if area.Castle {
			areaString += " (castle)"
		}
		fmt.Println(areaString)

		if !area.IsEmpty() {
			fmt.Println("Unit:", area.Unit.Player, area.Unit.Type)
		}

		if neighbors {
			fmt.Println("\nNEIGHBORS:")
			for _, neighbor := range area.Neighbors {
				neighborString := neighbor.Name
				if neighbor.AcrossWater {
					neighborString += " (across water)"
				}
				if neighbor.Cliffs {
					neighborString += " (cliffs)"
				}
				if neighbor.DangerZone != "" {
					neighborString += fmt.Sprintf(" (danger: %s)", neighbor.DangerZone)
				}
				fmt.Println(neighborString)
			}
		}

		fmt.Print("\n-----------------------------------\n\n")
	}
}

func adjustBoard(brd board.Board, areas map[string]board.Unit) {
	for areaName, unit := range areas {
		if unit.Type != "" {
			area := brd.Areas[areaName]
			area.Unit = unit
			if !area.Sea {
				area.Control = unit.Player
			}
			brd.Areas[areaName] = area
		}
	}
}

func printResolvePrint(brd board.Board, areas map[string]board.Unit, round board.Round) {
	fmt.Print("---BEFORE---\n\n")
	printBoard(brd, areas, false)

	brd.Resolve(round)

	fmt.Print("---AFTER---\n\n")
	printBoard(brd, areas, false)
}

func main() {
	brd, err := boardsetup.ReadBoard("hermannia_5players")

	if err != nil {
		fmt.Println(err.Error())
	}

	testTransportWithDangerZone(brd)
}

func testTransportWithDangerZone(brd board.Board) {
	areas := map[string]board.Unit{
		"Winde": {
			Type:   board.UnitFootman,
			Player: "green",
		},
		"Mare Gond": {
			Type:   board.UnitShip,
			Player: "green",
		},
		"Mare Elle": {
			Type:   board.UnitShip,
			Player: "green",
		},
		"Mare Ovond": {
			Type:   board.UnitShip,
			Player: "green",
		},
		"Mare Duna": {
			Type:   board.UnitShip,
			Player: "red",
		},
		"Mare Furie": {
			Type:   board.UnitShip,
			Player: "red",
		},
		"Fond": {},
	}

	adjustBoard(brd, areas)

	round := board.Round{
		Season: board.SeasonSpring,
		FirstOrders: []board.Order{
			{
				Type:   board.OrderMove,
				Player: "green",
				From:   "Winde",
				To:     "Fond",
				Unit:   brd.Areas["Winde"].Unit,
			},
			{
				Type:   board.OrderTransport,
				Player: "green",
				From:   "Mare Gond",
				Unit:   brd.Areas["Mare Gond"].Unit,
			},
			{
				Type:   board.OrderTransport,
				Player: "green",
				From:   "Mare Elle",
				Unit:   brd.Areas["Mare Elle"].Unit,
			},
			{
				Type:   board.OrderTransport,
				Player: "green",
				From:   "Mare Ovond",
				Unit:   brd.Areas["Mare Ovond"].Unit,
			},
			{
				Type:   board.OrderMove,
				Player: "red",
				From:   "Mare Duna",
				To:     "Mare Gond",
				Unit:   brd.Areas["Mare Gond"].Unit,
			},
			{
				Type:   board.OrderMove,
				Player: "red",
				From:   "Mare Furie",
				To:     "Mare Elle",
				Unit:   brd.Areas["Mare Furie"].Unit,
			},
		},
	}

	printResolvePrint(brd, areas, round)
}

func testTransportBattle(brd board.Board) {
	areas := map[string]board.Unit{
		"Worp": {
			Type:   board.UnitFootman,
			Player: "green",
		},
		"Mare Gond": {
			Type:   board.UnitShip,
			Player: "green",
		},
		"Mare Elle": {
			Type:   board.UnitShip,
			Player: "red",
		},
		"Zona": {},
	}

	adjustBoard(brd, areas)

	round := board.Round{
		Season: board.SeasonSpring,
		FirstOrders: []board.Order{
			{
				Type:   board.OrderMove,
				Player: "green",
				From:   "Worp",
				To:     "Zona",
				Unit:   brd.Areas["Worp"].Unit,
			},
			{
				Type:   board.OrderTransport,
				Player: "green",
				From:   "Mare Gond",
				Unit:   brd.Areas["Mare Gond"].Unit,
			},
			{
				Type:   board.OrderMove,
				Player: "red",
				From:   "Mare Elle",
				To:     "Mare Gond",
				Unit:   brd.Areas["Mare Gond"].Unit,
			},
		},
	}

	printResolvePrint(brd, areas, round)
}
