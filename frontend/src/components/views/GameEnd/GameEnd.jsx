import PropTypes from "prop-types";
import { useNavigate } from "react-router";
import short from "short-uuid";

import Button from "../../ui/Button";
import Card from "../../ui/Card";

import classes from "./GameEnd.module.css";

const WinState = {
	WIN: "win",
	LOSE: "lose",
};

export default function GameEnd(props) {
	const navigate = useNavigate();

	const winState = props.winState;
	const winnerTeam = props.winnerTeam;
	const lostTeam = props.lostTeam;
	const endGame = props.endGame;

	console.log(endGame);

	const winPlayers = winnerTeam.map((player) => {
		return (
			<div key={player.id} className={classes["player"]}>
				<p>{player.name}</p>
			</div>
		);
	});

	const lostPlayers = lostTeam.map((player) => {
		return (
			<div key={player.id} className={classes["player"]}>
				<p>{player.name}</p>
			</div>
		);
	});

	const gameDataMarkup = (
		<div className={classes["game-data"]}>
			<p>{`${endGame.sender_name} revealed ${endGame.target_name} cards!`}</p>
			<div className={classes["cards"]}>
				{endGame.target_cards.map((card) => {
					return (
						<Card
							style={{ fontSize: "8px" }}
							key={short.generate()}
							cardInfo={card}
						/>
					);
				})}
			</div>
		</div>
	);

	const homeButton = <Button value="Home" onClick={() => navigate("/")} />;

	const teamMarkup =
		winState === WinState.WIN ? (
			<div className={`${classes["team"]} ${classes["win"]}`}>
				<h1>You Won !</h1>
				{winPlayers}
				{homeButton}
			</div>
		) : (
			<div className={`${classes["team"]} ${classes["lose"]}`}>
				<h1>You Lost :(</h1>
				{lostPlayers}
				{homeButton}
			</div>
		);

	return (
		<div className={`${classes["game-end"]}`}>
			{gameDataMarkup}
			{teamMarkup}
		</div>
	);
}

GameEnd.propTypes = {
	winState: PropTypes.string,
	winnerTeam: PropTypes.arrayOf(
		PropTypes.shape({
			name: PropTypes.string,
			team: PropTypes.string,
			id: PropTypes.string,
		})
	),
	lostTeam: PropTypes.arrayOf(
		PropTypes.shape({
			name: PropTypes.string,
			team: PropTypes.string,
			id: PropTypes.string,
		})
	),
	endGame: PropTypes.shape({
		sender_name: PropTypes.string,
		target_name: PropTypes.string,
		target_cards: PropTypes.arrayOf(
			PropTypes.shape({
				value: PropTypes.string,
				color: PropTypes.string,
			})
		),
	}),
};
