import "./Card.css";
import PropTypes from "prop-types";

export default function Card(props) {
	const cardSymbol = props.symbol;

	const cl = () => {
		console.log(cardSymbol);
	};
	return (
		<button onClick={cl} className="card" style={props.style}>
			<div className="inner">
				<span>{cardSymbol}</span>
				<span>{cardSymbol}</span>
				<span>{cardSymbol}</span>
			</div>
		</button>
	);
}

Card.propTypes = {
	style: PropTypes.object,
	symbol: PropTypes.string,
};
