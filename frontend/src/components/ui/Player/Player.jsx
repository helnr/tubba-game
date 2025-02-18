import Card from "../Card";

import "./Player.css";
import PropTypes from "prop-types";
export default function Player(props) {
	// const colors = {
	// 	red: "var(--red)",
	// 	blue: "var(--blue)",
	// 	green: "var(--dark-green)",
	// 	yellow: "var(--yellow)",
	// };
	const player = props.player;
	const sendCardEvent = props.sendCardEvent;
	const selected = props.selected || false;
	const onClick = props.onClick;
	const className = props?.className || "";

	const cards = player?.cards?.map((card, index) => {
		return <Card key={index} cardInfo={card} sendEventFunc={sendCardEvent} />;
	});

	return (
		<div
			onClick={onClick}
			className={`${className} player ${props.mainPlayer ? "main-player" : ""}`}
			style={{
				...props.style,
				outline: selected
					? "4px solid var(--yellow)"
					: player?.is_turn
					? "2px solid var(--foreground-color)"
					: "none",
			}}>
			{props.mainPlayer && cards.length != 0 ? (
				<div className="player-cards">{cards}</div>
			) : (
				<span className="player-name">{player?.name}</span>
			)}
		</div>
	);
}

Player.propTypes = {
	style: PropTypes.object,
	mainPlayer: PropTypes.bool,
	player: PropTypes.shape({
		name: PropTypes.string,
		cards: PropTypes.arrayOf(
			PropTypes.shape({
				value: PropTypes.string,
				color: PropTypes.string,
			})
		),
		is_turn: PropTypes.bool,
	}),
	sendCardEvent: PropTypes.func,
	selected: PropTypes.bool,
	onClick: PropTypes.func,
	className: PropTypes.string,
};
