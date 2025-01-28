package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameStatus string

const (
	STATUS_LOADING GameStatus = "loading"
	STATUS_LOBBY GameStatus = "lobby"
	STATUS_READY GameStatus = "ready"
	STATUS_STARTED GameStatus = "started"
)



type Player struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Cards []string `json:"cards"`
	Team string `json:"team"`
	IsTurn bool `json:"is_turn"`
}

func PlayerFromUser(user *User) *Player {
	return &Player{
		ID: user.ID,
		Name: user.Name,
		IsTurn: false,
	}
}

func (p *Player) Info() *PlayerInfo {
	return &PlayerInfo{
		ID: p.ID,
		Name: p.Name,
		Team: p.Team,
	}
}


type Players []Player

func (p Players) Len() int {
	return len(p)
}

func (p Players) Has(id primitive.ObjectID) bool {
	for _, player := range p {
		if player.ID == id {
			return true
		}
	}
	return false
}	

func (p Players) Add(player *Player) error {
	if p.Has(player.ID) {
		return fmt.Errorf("Player already exists")
	}
	p = append(p, *player)
	return nil
}

type Game struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status GameStatus `json:"status"`
	TotalCards int `json:"total_cards"`
	PlayedCards int `json:"played_cards"`
	Players Players `json:"players"`
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
}

func (g *Game) Info() *GameInfo {
	players := make([]PlayerInfo, len(g.Players))

	for i, p := range g.Players {
		players[i] = *p.Info()
	}


	return &GameInfo{
		ID: g.ID,
		Status: g.Status,
		Players: players,
	}
}

type PlayerInfo struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Team string `json:"team"`
}

type GameInfo struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status GameStatus `json:"status"`
	Players []PlayerInfo `json:"players"`
}

type GameStore interface {
	GetGameByID(id primitive.ObjectID) (*Game, error)
	SaveGame(game *Game) error
	UpdateGame(game *Game) error
}