package board

func (battle Battle) areaNames() []string {
	nameMap := make(map[string]struct{})

	for _, result := range battle.Results {
		if result.DefenderArea != "" {
			nameMap[result.DefenderArea] = struct{}{}
		} else if result.Move.To != "" {
			nameMap[result.Move.To] = struct{}{}
		}
	}

	names := make([]string, 0)
	for name := range nameMap {
		names = append(names, name)
	}

	return names
}

// Returns whether the battle was between two moves moving against each other.
func (battle Battle) isBorderConflict() bool {
	return len(battle.Results) == 2 &&
		(battle.Results[0].Move.To == battle.Results[1].Move.From) &&
		(battle.Results[1].Move.To == battle.Results[0].Move.From)
}

// Parses the results of the battle and finds the winners and losers.
//
// In case of a battle against an unconquered area or a danger zone,
// only one player is returned in one of the lists.
//
// In case of a battle between players, multiple winners are returned
// in the case of a tie.
func (battle Battle) parseResults() (winners []Player, losers []Player) {
	// Checks if the battle was against an unconquered area.
	if len(battle.Results) == 1 {
		result := battle.Results[0]

		// Checks that order meets the requirement for crossing the danger zone.
		if battle.DangerZone != "" {
			if result.Total >= RequirementDangerZone {
				return []Player{result.Move.Player}, make([]Player, 0)
			} else {
				return make([]Player, 0), []Player{result.Move.Player}
			}
		}

		// Checks that order meets the requirement for conquering the neutral area.
		if result.Total >= RequirementConquer {
			return []Player{result.Move.Player}, make([]Player, 0)
		} else {
			return make([]Player, 0), []Player{result.Move.Player}
		}
	}

	// Finds the highest result among the players in the battle.
	highestResult := 0
	for _, result := range battle.Results {
		if result.Total > highestResult {
			highestResult = result.Total
		}
	}

	// Sorts combatants based on whether they exceeded the highest result.
	winners = make([]Player, 0)
	losers = make([]Player, 0)
	for _, result := range battle.Results {
		if result.Total >= highestResult {
			winners = append(winners, result.Move.Player)
		} else {
			losers = append(losers, result.Move.Player)
		}
	}

	return winners, losers
}

// Checks if the given player is contained in the given list of players.
func containsPlayer(players []Player, player Player) bool {
	for _, otherPlayer := range players {
		if otherPlayer == player {
			return true
		}
	}

	return false
}
