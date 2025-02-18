package types

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameStatus string

const (
	STATUS_LOADING GameStatus = "loading"
	STATUS_LOBBY GameStatus = "lobby"
	STATUS_READY GameStatus = "ready"
	STATUS_STARTED GameStatus = "started"
	STATUS_PAUSED GameStatus = "paused"
	STATUS_ENDED GameStatus = "ended"

	CARD_RED = "red"
	CARD_YELLOW = "yellow"
	CARD_BLUE = "blue"
	CARD_GREEN = "green"

	CARD_VALUES = "0123456789"
)

var (
	seededRand *rand.Rand = rand.New(rand.NewSource((time.Now().UnixNano())))
	colors [4]string = [4]string{CARD_BLUE, CARD_RED, CARD_GREEN, CARD_YELLOW}
)

type GameHandler interface {
	RouteEvent(event EventMessage, player *Player) error
	AddPlayer(player *Player)
	RemovePlayer(player *Player)
	GetGameByID(id string) (*Game, error)
}


type Player struct {
	Conn *websocket.Conn
	Manager GameHandler
	GameID string
	Egress chan EventMessage
	Ctx context.Context
	Cancel context.CancelFunc
	sync.RWMutex


	Game *Game
	User *User
	Cards *Cards `json:"cards"`
	Team string `json:"team"`
	NextPlayer *Player
	IsTurn bool `json:"is_turn"`
	IsOwner bool `json:"is_owner"`
	Connected bool
}

func NewPlayer(conn *websocket.Conn, manager GameHandler, gameID string, game *Game, user *User) *Player {
	ctx, cancel := context.WithCancel(context.Background())
	return &Player{
		Conn: conn,
		Manager: manager,
		GameID: gameID,
		Egress: make(chan EventMessage),
		Ctx: ctx,
		Cancel: cancel,


		Game: game,
		User: user,
		Cards: &Cards{},
		Team: "",
		NextPlayer: nil,
		IsTurn: false,
		IsOwner: false,
		Connected: true,
	}
}

func NewPlayerFromOldPlayer(oldPlayer *Player, conn *websocket.Conn, manager GameHandler, gameID string, game *Game, user *User) *Player {
	newPlayer := NewPlayer(conn, manager, gameID, game, user)
	newPlayer.Cards =  oldPlayer.Cards
	newPlayer.Team = oldPlayer.Team
	newPlayer.NextPlayer = oldPlayer.NextPlayer
	newPlayer.IsTurn = oldPlayer.IsTurn
	newPlayer.IsOwner = oldPlayer.IsOwner
	return newPlayer
}

func (p *Player) ReadMessages() {
	defer func() {
		if p.Game.Status == STATUS_LOBBY || p.Game.Status == STATUS_ENDED {
			p.Conn.Close()
			p.Manager.RemovePlayer(p)
		}else {
			p.Connected = false
		}
		p.Cancel()
	}()
	outerLoop:
	for {
		select {
		case <-p.Ctx.Done():
			return

		default:
			_, message, err := p.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return

			}

			var eventMessage EventMessage
			err = json.Unmarshal(message, &eventMessage)
			if err != nil {
				log.Println(err)
				continue outerLoop
			}

			err = p.Manager.RouteEvent(eventMessage, p)	
			if err != nil {
				continue outerLoop
			}

		}
	}
}

func (p *Player) WriteMessages() {
	defer func() {
		if p.Game.Status == STATUS_LOBBY {
			p.Manager.RemovePlayer(p)
		}else {
			p.Connected = false
		}
		p.Cancel()

	}()

	for {
		select {
		case message, ok := <-p.Egress:
			if !ok {
				err := p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			err = p.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				return
			}

		case <-p.Ctx.Done():
			return
		}
	}

}

func (p *Player) Boadcast(message EventMessage) {
	for _, player := range *p.Game.Players {
		if p.User.ID != player.User.ID && player.Connected {
			player.Egress <- message
		}
	}
}

func (p *Player) BoadcastGameEvent() {
	for _, player := range *p.Game.Players {
		if p.User.ID != player.User.ID && player.Connected {
			player.Egress <- *NewEventMessage(GameEvent, player.GameData())
		}
	}
}

func (p *Player) Update(message EventMessage) {
	p.Egress <- message
}

func (p *Player) UpdateError(message string) {
	p.Egress <- *NewEventMessage(ErrorEvent, message)
}

func (p *Player) UpdateGameEvent() {
	p.Egress <- *NewEventMessage(GameEvent, p.GameData())
}

func (p *Player) Info() *PlayerInfo {
	return &PlayerInfo{
		ID: p.User.ID,
		Name: p.User.Name,
		Team: p.Team,
		IsTurn: p.IsTurn,
	}
}

func (p *Player) PlayerData() *PlayerData{
	return &PlayerData{
		ID: p.User.ID,
		Name: p.User.Name,
		Team: p.Team,
		Cards: p.Cards,
		IsTurn: p.IsTurn,
		IsOwner: p.IsOwner,
	}
}

func (p *Player) GameData() *GameData{
	return &GameData{
		Status: p.Game.Status,
		CurrentCard: p.Game.CurrentCard,
		Players: p.Game.PlayersInfo(),
		MainPlayer: *p.PlayerData(),
		EndGame: p.Game.EndGame,
	}
}


type Players []*Player

func (p *Players) Get(id primitive.ObjectID) *Player {
	for i := 0; i < len(*p); i++ {
		if (*p)[i].User.ID == id {
			return (*p)[i]
		}
	}	
	return nil
}

func (p *Players) Len() int {
	return len(*p)
}

func (p *Players) Has(id primitive.ObjectID) bool {
	for _, player := range *p {
		if player.User.ID == id {
			return true
		}
	}
	return false
}	

func (p *Players) Add(player *Player) error {
	if p.Has(player.User.ID) {
		return fmt.Errorf("Player already exists")
	}
	*p = append(*p, player)
	return nil
}

func (players *Players) Remove(player *Player) error {
	for i, p := range *players {
		if p.User.ID == player.User.ID {
			*players = append((*players)[:i], (*players)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Player not found")
}

func (players *Players) BindTurns() bool {
	r1 := []*Player{}
	r2 := []*Player{}

	for _, p := range *players {
		if p.Team == "" {
			log.Println("Player has no team ")
			return false
		}
		if p.Team == "team1" {
			r1 = append(r1, p)
		}else if p.Team == "team2" {
			r2 = append(r2, p)
		}
	}


	if !(len(r1) == 2 && len(r2) == 2) {
		log.Println("Teams not found ")
		return false
	}

	// Make the first player in team2 play after the first player in team1
	tmp := r1[1]
	r1[1] = r2[0]
	r2[0] = tmp


	r1[0].NextPlayer = r1[1]
	r1[1].NextPlayer = r2[0]
	r2[0].NextPlayer = r2[1]
	r2[1].NextPlayer = r1[0]
	return true
}

func (players *Players) UnbindTurns() {
	for _, p := range *players {
		p.NextPlayer = nil
	}
}

type Game struct {
	ID string `json:"id" bson:"_id,omitempty"`
	Status GameStatus `json:"status"`
	AvailableCards *Cards `json:"total_cards"`
	PlayedCards *Cards `json:"played_cards"`
	Players *Players `json:"players"`
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
	MaxPlayers int 
	OwnerID primitive. ObjectID `json:"game_owner" bson:"gameOwner"`
	CurrentCard *Card `json:"current_card"`
	EndGame *EndGameData `json:"end_game"`
}

func NewGame() *Game {
	game := &Game{
		Status: STATUS_LOBBY,
		AvailableCards: NewCards(),
		PlayedCards: &Cards{},
		Players: &Players{},
		CreatedAt: time.Now(),
		MaxPlayers: 4,
	}

	return game
}

func (g *Game) CheckTeamsFull() bool {
	lengths := map[string]int{}
	for _, p := range *g.Players {
		if p.Team != "" { lengths[p.Team]++ }
	}

	if len(lengths) == 2 {
		for _, v := range lengths {
			if v != 2 {
				return false
			}
		}
		return true
	}

	return false
}

func (g *Game) CountTeam(team string) int {
	count := 0
	for _, p := range *g.Players {
		if p.Team == team {
			count++
		}
	}
	return count
}

func (g *Game) PlayersInfo() []PlayerInfo {
	players := []PlayerInfo{}

	for _, p := range *g.Players {
		if p.Connected == false {
			continue
		}
		players = append(players, *p.Info())
	}
	return players
}

func (g *Game) Info() *GameInfo {
	return &GameInfo{
		ID: g.ID,
		Status: g.Status,
		Players: g.PlayersInfo(),
	}
}


type Card struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type Cards []*Card

func NewCards() *Cards {
	cards := Cards{}
	for _, color := range colors {
		for _, r := range CARD_VALUES {
			value := string(r)
			card := &Card{
				Value: value,
				Color: color,
			}
			cards.Add(card)
		}
	}
	cards.Shuffle()
	return &cards
}

func (cards *Cards) Get(id int) *Card {
	return (*cards)[id]
}

func (cards *Cards) Find(card *Card) int {
	for i, c := range *cards {
		if c.Value == card.Value && c.Color == card.Color {
			return i
		}
	}
	return -1
}

func (cards *Cards) Add(card *Card) {
	(*cards) = append(*cards, card)
}

func (cards *Cards) DeleteCard(id int) {
	(*cards) = append((*cards)[:id], (*cards)[id+1:]...)
} 

func (cards *Cards) DrawCard(target *Cards) bool {
	if len(*cards) <= 0 {
		return false
	}

	card := (*cards)[0]
	(*cards) = (*cards)[1:]
	target.Add(card)
	return true
}

func (cards *Cards) Shuffle() {
	seededRand.Shuffle(len(*cards), func(i, j int) {
		(*cards)[i], (*cards)[j] = (*cards)[j], (*cards)[i]
	})
}






type PlayerInfo struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Team string `json:"team"`
	IsTurn bool `json:"is_turn"`
}

type PlayerData struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Team string `json:"team"`
	Cards *Cards `json:"cards"`
	IsTurn bool `json:"is_turn"`
	IsOwner bool `json:"is_owner"`
}

type GameInfo struct {
	ID string `json:"id" bson:"_id,omitempty"`
	Status GameStatus `json:"status"`
	Players []PlayerInfo `json:"players"`
}

type GameData struct {
	Status GameStatus `json:"status"`
	CurrentCard *Card `json:"current_card"`
	Players []PlayerInfo `json:"players"`
	MainPlayer PlayerData `json:"main_player"`
	EndGame *EndGameData `json:"end_game"`
}

type EndGameData struct {
	SenderName string `json:"sender_name"`
	TargetName string `json:"target_name"`
	TargetCards *Cards `json:"target_cards"`
	WinnerTeam string `json:"winner_team"`
}

type GameStore interface {
	GetGameByID(string) (*Game, error)
	AddGame(*Game) (string, error) 
	RemovePlayer(string, *Player) error 
	AddPlayer(string, *Player) error 
}