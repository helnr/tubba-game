import useAuth from "../../../hooks/useAuth";
import { useState, useRef } from "react";
import { useNavigate } from "react-router";
import axios from "axios";
import Menu from "../../views/Menu";
import JoinGame from "../../views/JoinGame";

import "./home.css";
export function Home() {
	const [view, setView] = useState("menu");
	const navigate = useNavigate();
	const gameCode = useRef("");

	const {
		auth: { user, error },
		loading,
	} = useAuth();
	if (loading) {
		return <div className="home">Loading...</div>;
	}

	if (error) {
		return <div className="home">{error}</div>;
	}

	const onCreate = () => {
		// Make a new request to the backend to create a new game
		// get the game code and navigate to the game page with the game code

		const createGame = async () => {
			try {
				const apiURL = import.meta.env.VITE_API_URL;
				const response = await axios.post(
					`${apiURL}/game`,
					{},
					{
						headers: {
							"Content-Type": "application/json",
						},
						withCredentials: true,
					}
				);
				if (response.status === 201) {
					const gameData = await response.data.game;
					console.log(gameData);
					navigate("/game", { state: { game: gameData, user: user } });
				}
			} catch (error) {
				alert("Error creating game");
				console.log(error);
			}
		};

		createGame();
	};

	const onJoin = () => {
		setView("join");
		// console.log("Join Game!");
	};

	const onMenu = () => {
		setView("menu");
	};

	return (
		<div className="home">
			{user && (
				<div className="user-info">
					<h2>hello {user.name}</h2>
				</div>
			)}

			{view === "join" ? (
				<div className="view">
					<JoinGame onMenu={onMenu} />
				</div>
			) : view === "create" ? (
				<div>create</div>
			) : (
				<Menu onCreate={onCreate} onJoin={onJoin} />
			)}
		</div>
	);
}
