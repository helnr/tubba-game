import Card from "../../ui/Card";
import CircleButton from "../../ui/CircleButton";
import PropsTypes from "prop-types";

import "./Table.css";
export default function Table(props) {
	const current_card = props?.currentCard;
	const sendTubbahEvent = props?.tubbaEventFunc;

	return (
		<div className="table">
			<CircleButton onClick={sendTubbahEvent} width={50}></CircleButton>
			{current_card ? (
				<>
					<Card className="table-card" cardInfo={current_card} />
				</>
			) : (
				<p className="info">No card...</p>
			)}
		</div>
	);
}

Table.propTypes = {
	currentCard: PropsTypes.shape({
		value: PropsTypes.string,
		color: PropsTypes.string,
	}),
	tubbaEventFunc: PropsTypes.func,
};
