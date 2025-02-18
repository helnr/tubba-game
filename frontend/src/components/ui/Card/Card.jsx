import "./Card.css";
import PropTypes from "prop-types";

export default function Card(props) {
	const card = props.cardInfo;
	const sendEventFunc = props.sendEventFunc;

	const cl = () => {
		if (sendEventFunc) {
			sendEventFunc(card);
		}
	};
	return (
		<button
			onClick={cl}
			className={`card ${props.className ? props.className : ""}`}
			style={props.style}>
			<div className="inner">
				<span>{card?.value}</span>
				<span>{card?.value}</span>
				<span>{card?.value}</span>
			</div>
		</button>
	);
}

Card.propTypes = {
	style: PropTypes.object,
	className: PropTypes.string,
	cardInfo: PropTypes.shape({
		value: PropTypes.string,
		color: PropTypes.string,
	}),
	sendEventFunc: PropTypes.func,
};
