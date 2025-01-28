import { useLocation, useNavigate, Navigate } from "react-router";
import { useState, useEffect, useCallback } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import axios from "axios";

import Player from "../../ui/Player";
import Button from "../../ui/Button";
import Table from "../../views/Table";
import Lobby from "../../views/Lobby";

import { EventMessage, EventMessageTypes } from "../../../types/game.js";

const STATUS = {
	LOADING: "loading",
	LOBBY: "lobby",
	READY: "ready",
	STARTED: "started",
};

import "./game.css";

export default function Game() {
	const navigate = useNavigate();
	const location = useLocation();
	const game = location.state?.game;
	const user = location.state?.user;
	const [gameData, setGameData] = useState({
		status: STATUS.LOADING,
		total_cards: 40,
		played_cards: 0,
		players: [],
	});
	const [playerStatus, setPlayerStatus] = useState({
		isTurn: true,
		cards: ["1", "2", "3", "4"],
		team: "one",
	});

	const { sendJsonMessage, lastMessage, readyState } = useWebSocket(
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
			onClose: () => {
				sendJsonMessage(newEvent(EventMessageTypes.LeaveEvent, {}));
			},
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

		return () => {
			window.removeEventListener("beforeunload", leave);
		};
	}, [sendJsonMessage, newEvent]);

	const handleEvent = useCallback(
		(event, payload) => {
			if (event === EventMessageTypes.GameEvent) {
				setGameData(payload);
			} else if (event === EventMessageTypes.ErrorEvent) {
				alert(payload.error);
				navigate("/");
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

	const socketStatus = {
		[ReadyState.CONNECTING]: "Connecting",
		[ReadyState.OPEN]: "Open",
		[ReadyState.CLOSING]: "Closing",
		[ReadyState.CLOSED]: "Closed",
	}[readyState];

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
						isReady={status === "ready"}
						players={[
							{
								name: "Mohammed Zangooh Gah Hah Fedora linux Zangooh Gah Hah Fedora linux",
							},
							{ name: "Ahmed" },
							{ name: "Ali" },
							{ name: "Hasan" },
						]}
					/>
				);
			case "started":
				return <Table />;

			default:
				break;
		}
		return null;
	}

	return (
		<div className="game">
			<div className="team one">
				<Player name="Player 1" />
				<Player name="Player 2" />
			</div>

			<div className="view">{setView(game?.status)}</div>

			<div className={`team ${playerStatus?.team}`}>
				<Player name={event ? event.type : ""} />
				<Player name="Player 4" mainPlayer={true} cards={playerStatus?.cards} />
			</div>

			{game?.status === "started" && (
				<Button
					value="Tubba!"
					onClick={() => setGameData({ status: "lobby" })}
				/>
			)}
			{game?.status === "lobby" && (
				<Button
					value={socketStatus}
					disabled={true}
					onClick={() =>
						sendJsonMessage(newEvent("game_event", { status: "ready" }))
					}
				/>
			)}
			{game?.status === "ready" && (
				<Button
					value="Start"
					onClick={() => setGameData({ status: "started" })}
				/>
			)}
		</div>
	);
}
