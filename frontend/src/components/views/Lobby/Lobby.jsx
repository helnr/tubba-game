import Button from "../../ui/Button";

import "./Lobby.css";
import PropTypes from "prop-types";
export default function Lobby(props) {
	const players = props.players?.map((player, index) => {
		return <p key={index}>{player.name}</p>;
	});
	return (
		<div className="lobby">
			<div className="lobby-players">{players}</div>
			{props.isReady ? (
				<p className="message">Ready!</p>
			) : (
				<p className="message">Waiting for players...</p>
			)}
		</div>
	);
}

Lobby.propTypes = {
	players: PropTypes.arrayOf(
		PropTypes.shape({
			name: PropTypes.string,
		})
	),

	isReady: PropTypes.bool,
	onReady: PropTypes.func,
};
