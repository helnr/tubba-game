package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/helnr/tubba-game/backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type MainPlayerGameData struct {
	GameInfo *types.GameInfo `json:"game_info"`
}





type GameEventPayload struct {
	Status types.GameStatus `json:"status"`
}

type JoinEventPayload struct {
	UserID string `json:"user_id"`
	Username string `json:"username"`
}

type ChangeTeamPayload struct {
	Team string `json:"team"`
}

type TubbaEventPayload struct {
	ID string `json:"id"`
}

type EventHandler func(message *types.EventMessage, player *types.Player) error

func HandleGameEvent(message *types.EventMessage, player *types.Player) error {
	fmt.Printf("Players in game %d: ", player.Game.Players.Len())
	for _, player := range *player.Game.Players {
		fmt.Printf("%s ", player.User.Name)
	}
	fmt.Printf("\n")

	var payload GameEventPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return err
	}

	log.Println("Game event: ", payload.Status)
	player.Boadcast(*message)
	return nil
}

func HandleReadyEvent(message *types.EventMessage, player *types.Player) error {
	player.Game.Status = types.STATUS_READY
	log.Printf("Ready %p\n", player)
	return nil
}

func HandleLeaveEvent(message *types.EventMessage, player *types.Player) error {

	player.Lock()
	defer player.Unlock()
	player.Cancel()
	player.Conn.Close()
	if player.Game.Status == types.STATUS_LOBBY {
		player.Manager.RemovePlayer(player)
	}else {
		player.Connected = false
		if player.Game.Status == types.STATUS_READY {
			player.Game.Status = types.STATUS_LOBBY
		}
	}

	player.BoadcastGameEvent()

	return nil
}

func HandleJoinEvent(message *types.EventMessage, player *types.Player) error {
	player.UpdateGameEvent()
	player.BoadcastGameEvent()
	return nil
}

func HandleCardPlayedEvent(message *types.EventMessage, player *types.Player) error {
	dataObj := struct{Card *types.Card `json:"card"`}{
		Card: &types.Card{
			
		},
	}
	err := json.Unmarshal(message.Payload, &dataObj)
	if err != nil {
		return err
	}

	player.Lock()
	defer player.Unlock()

	if !player.IsTurn {
		for _, p := range *player.Game.Players {
			if p.IsTurn {
				message := *types.NewEventMessage(types.ErrorEvent, types.ErrorEventPayload{Error: "Not your turn", Type: "message"})
				player.Update(message)
				return fmt.Errorf("Not your turn")
			}
		}
	}

	player.IsTurn = false
	player.NextPlayer.IsTurn = true

	card_idx := player.Cards.Find(dataObj.Card)
	if card_idx == -1 {
		message := *types.NewEventMessage(types.ErrorEvent, types.ErrorEventPayload{Error: "Card not found", Type: "message"})
		player.Update(message)
		return fmt.Errorf("Card not found")
	}

	player.Game.PlayedCards.Add((*player.Cards)[card_idx])
	player.Cards.DeleteCard(card_idx)
	player.Game.CurrentCard = dataObj.Card
	ok := player.Game.AvailableCards.DrawCard(player.Cards)
	if !ok {
		player.Game.CurrentCard = nil
		player.Game.PlayedCards.Shuffle()
		tmp := player.Game.AvailableCards
		player.Game.AvailableCards = player.Game.PlayedCards
		player.Game.PlayedCards = tmp
		player.Game.AvailableCards.DrawCard(player.Cards)
	}

	player.UpdateGameEvent()
	player.BoadcastGameEvent()
	return nil
}

func HandleStartEvent(message *types.EventMessage, player *types.Player) error {
	player.Lock()
	player.Unlock()
	player.Game.Status = types.STATUS_STARTED
	for _, p := range *player.Game.Players {
		if p.IsOwner {
			p.IsTurn = true
		}
		for i := 0; i < 4; i++ {
			player.Game.AvailableCards.DrawCard(p.Cards)
		}
	}
	player.UpdateGameEvent()
	player.BoadcastGameEvent()
	return nil
}

func HandleChangeTeamEvent(message *types.EventMessage, player *types.Player) error {
	var payload ChangeTeamPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return err
	}

	if payload.Team != "" {
		if player.Game.CountTeam(payload.Team) >= 2 {
			return fmt.Errorf("Team is full")
		}
	}


	player.Lock()
	player.Team = payload.Team

	if player.Game.CheckTeamsFull() {
		if status := player.Game.Players.BindTurns(); status {
			player.Game.Status = types.STATUS_READY
		}
	}else {
		player.Game.Status = types.STATUS_LOBBY
		player.Game.Players.UnbindTurns()
	}

	player.UpdateGameEvent()
	player.BoadcastGameEvent()
	player.Unlock()

	return nil
}

func HandleTubbasEvent(message *types.EventMessage, player *types.Player) error {
	var payload TubbaEventPayload
	err := json.Unmarshal(message.Payload, &payload)
	if err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		return err
	}



	player.Lock()
	defer player.Unlock()

	p := player.Game.Players.Get(id)

	if p == nil {
		return fmt.Errorf("Player not found")
	}

	cardsLength := len(*p.Cards)
	var win bool = true

	for i := 1; i < cardsLength; i++ {
		if (*p.Cards)[i].Value != (*p.Cards)[i-1].Value {
			win = false	
			break
		}
	}

	endGameData := &types.EndGameData{
		SenderName: player.User.Name,
		TargetName: p.User.Name,
		TargetCards: p.Cards,
	}

	player.Game.Status = types.STATUS_ENDED
	player.Game.EndGame = endGameData
	if win {
		endGameData.WinnerTeam = player.Team
	}else {
		if player.Team == p.Team {
			endGameData.WinnerTeam = player.NextPlayer.Team
		}else {
			endGameData.WinnerTeam = p.Team
		}
	}
	

	player.UpdateGameEvent()
	player.BoadcastGameEvent()
	return nil
}

func SendGameEvent(player *types.Player) error {
	game, err := player.Manager.GetGameByID(player.GameID)
	if err != nil {
		return err
	}
	payload := GameEventPayload{
		Status: game.Status,
	}
	raw_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	player.Boadcast(types.EventMessage{
		Type: types.GameEvent,
		Payload: raw_payload,
	})

	return nil
}

func SendJoinEvent(player *types.Player) error {
	payload := JoinEventPayload{
		UserID: player.User.ID.Hex(),
		Username: player.User.Name,
	}
	raw_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	player.Boadcast(types.EventMessage{
		Type: types.JoinEvent,
		Payload: raw_payload,
	})

	return nil
}

