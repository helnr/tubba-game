import Game from "./game";

export const GameState = {
	type: "game_state",
	payload: {
		status: "started",
		total_cards: 40,
		played_cards: 4,
	},
};

export const TeamState = {
	type: "team_state",
	payload: {
		team_one: [],
		team_two: [],
	},
};

export const PlayerState = {
	type: "player_state",
	payload: {
		isTurn: true,
		cards: ["2", "3", "4", "5"],
		team: "one",
	},
};

export default Game;
