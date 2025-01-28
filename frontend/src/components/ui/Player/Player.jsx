import Card from "../Card";

import "./Player.css";
import PropTypes from "prop-types";
export default function Player(props) {
	const cards = props.cards?.map((cardSymbol, index) => {
		return <Card key={index} symbol={cardSymbol} />;
	});

	return (
		<div className="player" style={props.style}>
			{props.mainPlayer ? (
				<div className="player-cards">{cards}</div>
			) : (
				<span className="player-name">{props.name}</span>
			)}
		</div>
	);
}

Player.propTypes = {
	style: PropTypes.object,
	mainPlayer: PropTypes.bool,
	name: PropTypes.string,
	cards: PropTypes.array,
};
