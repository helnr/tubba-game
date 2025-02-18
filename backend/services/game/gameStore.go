package game

import (
	"fmt"
	"time"

	"github.com/helnr/tubba-game/backend/types"
	"github.com/helnr/tubba-game/backend/utils"
)


type Games map[string]*types.Game

func NewGamesStore() *Games { 
	gameStore := &Games{}
	ticker := time.NewTicker(1200 * time.Second)

	go ClearUnactiveGames(gameStore, ticker)

	return gameStore 
}

func (g Games) GetNumberOfPlayers(code string, inLobby bool) (int, error) {
	if game, ok := g[code]; ok {
		if !inLobby {
			return game.Players.Len(), nil
		}
		return game.Players.Len(), nil
	}
	return 0, fmt.Errorf("Game not found")
}

func (g Games) GetGameByID(id string) (*types.Game, error) {
	if game, ok := g[id]; ok {
		return game, nil
	}
	return nil, fmt.Errorf("Game not found") 
}

func (g Games) AddGame(game *types.Game) (string, error) {
	code := utils.NewCode()
	attempts := 10
	for attempts > 0 {
		if _, ok := g[code]; !ok {
			g[code] = game
			return code, nil
		}
		code = utils.NewCode()
		attempts--
	}
	return "", fmt.Errorf("Too many attempts")
}

func (g Games) RemoveGame(game *types.Game) error {
	for _, p := range *game.Players {
		p.Cancel()
	}
	delete(g, game.ID)
	return nil
}

func (g Games) AddPlayer(code string, player *types.Player) error {
	if game, ok := g[code]; ok {
		return game.Players.Add(player)
	}
	return fmt.Errorf("Game not found")
}

func (g Games) RemovePlayer(code string, player *types.Player) error {
	if game, ok := g[code]; ok {
		return game.Players.Remove(player)
	}
	return fmt.Errorf("Game not found")
}

func ClearUnactiveGames(games *Games, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			for _, game := range *games {
				if game.Status == types.STATUS_ENDED {
					games.RemoveGame(game)
				}else if game.CreatedAt.Add(1440 * time.Minute).Before(time.Now()) {
					games.RemoveGame(game)
				}
			}
		}
	}
}


