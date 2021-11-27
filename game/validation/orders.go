package validation

import (
	"errors"
	"immerse-ntnu/hermannia/server/game"
)

func ValidateOrder(order game.Order, season game.Season) error {
	if order.Player.Color != order.From.Control {
		return errors.New("must control area that is ordered")
	}

	switch season {
	case game.Winter:
		return validateWinterOrder(order)
	default:
		return validateNonWinterOrder(order)
	}
}

func validateNonWinterOrder(order game.Order) error {
	if order.UnitBuild != "" {
		return errors.New("units can only be built in winter")
	}

	switch {
	case order.Type == game.Move || order.Type == game.Support:
		return validateMoveOrSupport(order)
	case order.Type == game.Besiege || order.Type == game.Transport:
		return validateBesiegeOrTransport(order)
	}

	return errors.New("invalid order type")
}

func validateMoveOrSupport(order game.Order) error {
	if order.To == nil {
		return errors.New("moves and supports must have destination")
	}

	if _, ok := order.From.Neighbors[order.To.Name]; !ok {
		return errors.New("destination not adjacent to origin")
	}

	if order.From.Unit.Type == game.Ship {
		if !(order.To.Sea || order.To.IsCoast()) {
			return errors.New("ship order destination must be sea or coast")
		}
	} else {
		if order.To.Sea {
			return errors.New("only ships can order to seas")
		}
	}

	switch order.Type {
	case game.Move:
		return validateMove(order)
	case game.Support:
		return validateSupport(order)
	}

	return errors.New("invalid order type")
}

func validateMove(order game.Order) error {
	if order.From.Unit == nil || order.From.Unit.Color != order.Player.Color {
		secondHorseMove := false

		for _, firstOrder := range order.From.IncomingMoves {
			if firstOrder.From.Unit == order.From.Unit {
				secondHorseMove = true
				break
			}
		}

		if !secondHorseMove {
			return errors.New("must have unit in origin area")
		}
	}

	return nil
}

func validateSupport(order game.Order) error {
	return nil
}

func validateBesiegeOrTransport(order game.Order) error {
	if order.To != nil {
		return errors.New("besiege or transport orders cannot have destination")
	}

	switch order.Type {
	case game.Besiege:
		return validateBesiege(order)
	case game.Transport:
		return validateTransport(order)
	}

	return errors.New("invalid order type")
}

func validateBesiege(order game.Order) error {
	if !order.From.Castle {
		return errors.New("besieged area must have castle")
	}

	if order.From.Control != game.Uncontrolled {
		return errors.New("besieged area cannot already be controlled")
	}

	if order.From.Unit.Type == game.Ship {
		return errors.New("ships cannot besiege")
	}

	return nil
}

func validateTransport(order game.Order) error {
	if order.From.Unit.Type != game.Ship {
		return errors.New("only ships can transport")
	}

	return nil
}

func validateWinterOrder(order game.Order) error {
	switch order.Type {
	case game.Move:
		return validateWinterMove(order)
	case game.Build:
		return validateBuild(order)
	}

	return errors.New("invalid order type")
}

func validateWinterMove(order game.Order) error {
	if order.To.Control != order.Player.Color {
		return errors.New("must control destination area in winter move")
	}

	if order.From.Unit.Type == game.Ship && !order.To.IsCoast() {
		return errors.New("ship winter move destination must be coast")
	}

	if order.UnitBuild != "" {
		return errors.New("cannot build unit with move order")
	}

	return nil
}

func validateBuild(order game.Order) error {
	if order.From.Unit != nil {
		return errors.New("cannot build in area already occupied")
	}

	switch order.UnitBuild {
	case game.Ship:
		if !order.From.IsCoast() {
			return errors.New("ships can only be built on coast")
		}
	case game.Footman:
		break
	case game.Horse:
		break
	case game.Catapult:
		break
	default:
		return errors.New("invalid unit type")
	}

	return nil
}