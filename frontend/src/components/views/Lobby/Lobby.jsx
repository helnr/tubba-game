import { STATUS } from "../../../types/game";

import classes from "./Lobby.module.css";
import PropTypes from "prop-types";
export default function Lobby(props) {
	let players = [];
	if (props.players) {
		const mainPlayer = props.mainPlayer;
		if (mainPlayer) {
			players.push(
				<div
					className={`${classes["main-player"]} ${classes[mainPlayer.team]}`}
					key={mainPlayer.id}>
					<p>{mainPlayer.name}</p>
					<button
						onClick={() => props.changeTeam("team1")}
						className={classes["team1"]}>
						Team 1
					</button>
					<button
						onClick={() => props.changeTeam("team2")}
						className={classes["team2"]}>
						Team 2
					</button>
				</div>
			);
		}

		const top = [];
		const bottom = [];
		const unknown = [];

		for (let i = 0; i < props.players.length; i++) {
			const player = props.players[i];
			const playerP = (
				<p className={classes[player.team]} key={player.id}>
					{player.name}
				</p>
			);

			if (player.id === mainPlayer.id) {
				continue;
			} else if (player.team === "") {
				unknown.push(playerP);
			} else if (player.team === mainPlayer.team) {
				top.push(playerP);
			} else {
				bottom.push(playerP);
			}
		}

		// players = [...players, ...top, ...bottom, ...unknown];
		players = [...unknown, ...bottom, ...top, ...players];
	}

	function copyGameCode(e) {
		e.preventDefault();
		navigator.clipboard.writeText(props.gameCode).then(() => {
			e.target.textContent = `${props.gameCode} Copied!`;
			window.setTimeout(() => {
				e.target.textContent = props.gameCode;
			}, 3000);
		});
	}

	return (
		<div className={classes.lobby}>
			<div className={classes["lobby-players"]}>{players}</div>
			{props.status === STATUS.LOBBY && props.gameCode && (
				<button onClick={copyGameCode} className={classes["game-code"]}>
					{props.gameCode}
				</button>
			)}
			{props.status === STATUS.LOBBY && (
				<p className={classes["message"]}>Waiting for players...</p>
			)}
		</div>
	);
}

Lobby.propTypes = {
	players: PropTypes.arrayOf(
		PropTypes.shape({
			name: PropTypes.string,
			team: PropTypes.string,
			id: PropTypes.string,
		})
	),
	mainPlayer: PropTypes.shape({
		name: PropTypes.string,
		team: PropTypes.string,
		id: PropTypes.string,
	}),
	teamPlayer: PropTypes.shape({
		name: PropTypes.string,
	}),
	opponents: PropTypes.arrayOf(
		PropTypes.shape({
			name: PropTypes.string,
		})
	),

	gameCode: PropTypes.string,
	status: PropTypes.string,
	onReady: PropTypes.func,
	playerID: PropTypes.string,
	changeTeam: PropTypes.func,
};
