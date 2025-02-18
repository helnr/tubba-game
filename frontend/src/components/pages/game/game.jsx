import { useLocation, useNavigate, Navigate } from "react-router";
import { useState, useEffect, useCallback } from "react";
import useWebSocket from "react-use-websocket";

import Player from "../../ui/Player";
import Button from "../../ui/Button";
import Table from "../../views/Table";
import Lobby from "../../views/Lobby";
import GameEnd from "../../views/GameEnd";

import {
	EventMessage,
	EventMessageTypes,
	STATUS,
} from "../../../types/game.js";

import classes from "./game.module.css";

export default function Game() {
	const navigate = useNavigate();
	const location = useLocation();
	const game = location.state?.game;
	const user = location.state?.user;
	const [selectedPlayer, setSelectedPlayer] = useState(null);
	const [gameData, setGameData] = useState({
		status: game?.status,
		current_card: null,
		players: [],
		main_player: {
			cards: [],
			id: user.id,
			is_owner: false,
			is_turn: false,
			name: user.name,
			team: "",
		},
		end_game: {
			winner_team: "",
			sender_name: "",
			target_name: "",
			target_cards: [],
		},
	});
	let team1, team2, teamPlayer, opponents;

	if (gameData.main_player.team != "") {
		team1 = gameData.players.filter((player) => player.team === "team1");
		team2 = gameData.players.filter((player) => player.team === "team2");

		teamPlayer = (gameData.main_player.team === "team1" ? team1 : team2).filter(
			(player) => player.id !== gameData.main_player.id
		)[0];
		opponents = gameData.main_player.team === "team1" ? team2 : team1;
	} else {
		teamPlayer = null;
		opponents = [];
	}

	const { sendJsonMessage, lastMessage } = useWebSocket(
		`ws://localhost:8080/game/join/${game.id}`,
		{
			share: true,
			shouldReconnect: () => {
				return true;
			},
			reconnectAttempts: 3,
			onError: (error) => {
				console.log(error);
			},
			onOpen: () => {
				sendJsonMessage(newEvent(EventMessageTypes.JoinEvent, {}));
				console.log("Connected");
			},
			onClose: () => {},
		}
	);

	const newEvent = useCallback((event, payload) => {
		return new EventMessage(event, payload);
	}, []);

	useEffect(() => {
		const leave = () => {
			sendJsonMessage(newEvent(EventMessageTypes.LeaveEvent, {}));
		};
		window.addEventListener("beforeunload", leave);
		window.addEventListener("popstate", leave);

		return () => {
			window.removeEventListener("beforeunload", leave);
			window.removeEventListener("popstate", leave);
		};
	}, [sendJsonMessage, newEvent]);

	console.log(gameData);

	function sendCardEvent(card) {
		sendJsonMessage(
			newEvent(EventMessageTypes.PlayedCardEvent, { card: card })
		);
	}

	function sendTubbahEvent() {
		if (selectedPlayer) {
			sendJsonMessage(
				newEvent(EventMessageTypes.TubbaEvent, { id: selectedPlayer.id })
			);
		} else {
			alert("Please select a player first");
		}
	}

	function changeTeam(team) {
		if (team === gameData.main_player.team) {
			sendJsonMessage(
				newEvent(EventMessageTypes.ChangeTeamEvent, { team: "" })
			);
			return;
		}

		sendJsonMessage(
			newEvent(EventMessageTypes.ChangeTeamEvent, { team: team })
		);
	}

	const handleEvent = useCallback(
		(event, payload) => {
			if (event === EventMessageTypes.GameEvent) {
				console.log(payload);
				setGameData(payload);
			} else if (event === EventMessageTypes.ErrorEvent) {
				alert(payload.error);
				if (payload.type === "navigate") {
					navigate("/");
				}
			} else if (event === EventMessageTypes.JoinEvent) {
				console.log(payload);
			}
		},
		[navigate]
	);

	useEffect(() => {
		if (lastMessage !== null) {
			const data = JSON.parse(lastMessage.data);
			const event = Object.assign(new EventMessage(), data);
			handleEvent(event.type, event.payload);
		}
	}, [lastMessage, handleEvent]);

	// const socketStatus = {
	// 	[ReadyState.CONNECTING]: "Connecting",
	// 	[ReadyState.OPEN]: "Open",
	// 	[ReadyState.CLOSING]: "Closing",
	// 	[ReadyState.CLOSED]: "Closed",
	// }[readyState];

	if (!game) {
		return <Navigate to="/" />;
	}

	function setView(status) {
		switch (status) {
			case "loading":
				return <div>Loading...</div>;
			case "lobby":
			case "ready":
				return (
					<Lobby
						onReady={() => {
							console.log("Ready");
						}}
						gameCode={game.id}
						status={gameData.status}
						players={gameData.players}
						mainPlayer={gameData.main_player}
						teamPlayer={teamPlayer}
						opponents={opponents}
						playerID={user.id}
						changeTeam={changeTeam}
					/>
				);
			case "started":
				return (
					<Table
						tubbaEventFunc={sendTubbahEvent}
						currentCard={gameData.current_card}
					/>
				);

			default:
				break;
		}
		return null;
	}

	function selectPlayer(player) {
		return () => {
			if (gameData.main_player.team != "" && player) {
				if (player.id !== selectedPlayer?.id) {
					setSelectedPlayer(player);
				} else {
					setSelectedPlayer(null);
				}
			}
		};
	}

	return gameData.status !== STATUS.ENDED ? (
		<div className={`${classes["game"]} ${classes[gameData.status]}`}>
			<div
				className={`${classes["team"]}${
					gameData.main_player.team === "team1"
						? " " + classes["two"]
						: gameData.main_player.team === "team2"
						? " " + classes["one"]
						: ""
				}`}>
				<Player
					style={{ cursor: "pointer" }}
					player={opponents[0]}
					className={classes["player"]}
					onClick={selectPlayer(opponents[0])}
					selected={
						selectedPlayer && selectedPlayer.id === opponents[0].id
							? true
							: false
					}
				/>
				<Player
					style={{ cursor: "pointer" }}
					player={opponents[1]}
					className={classes["player"]}
					onClick={selectPlayer(opponents[1])}
					selected={
						selectedPlayer && selectedPlayer.id === opponents[1].id
							? true
							: false
					}
				/>
			</div>

			<div className={classes["view"]}>{setView(gameData.status)}</div>

			<div
				className={`${classes["team"]} ${classes["main-team"]}${
					gameData.main_player.team === "team1"
						? " " + classes["one"]
						: gameData.main_player.team === "team2"
						? " " + classes["two"]
						: ""
				}`}>
				<Player
					style={{ cursor: "pointer" }}
					player={teamPlayer}
					className={classes["player"]}
					onClick={selectPlayer(teamPlayer)}
					selected={
						selectedPlayer && selectedPlayer.id === teamPlayer.id ? true : false
					}
				/>
				<Player
					className={classes["player"]}
					mainPlayer={true}
					player={gameData.main_player}
					sendCardEvent={sendCardEvent}
				/>
			</div>

			{/* {gameData.status === "started" && <Button value="Tubba!" />} */}
			{gameData.status === "ready" && (
				<Button
					onClick={() => {
						sendJsonMessage(newEvent(EventMessageTypes.StartEvent, {}));
					}}
					value="Start"
				/>
			)}
		</div>
	) : (
		<GameEnd
			winnerTeam={gameData.end_game.winner_team === "team1" ? team1 : team2}
			lostTeam={gameData.end_game.winner_team === "team1" ? team2 : team1}
			winState={
				gameData.end_game.winner_team === gameData.main_player.team
					? "win"
					: "lose"
			}
			endGame={gameData.end_game}
		/>
	);
}
