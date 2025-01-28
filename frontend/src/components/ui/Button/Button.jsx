import "./Button.css";
import PropTypes from "prop-types";
export default function Button(props) {
	return (
		<button
			onClick={props.onClick}
			className={`main-btn ${props.disabled ? "disabled" : ""}`}>
			{props.value}
		</button>
	);
}

Button.propTypes = {
	value: PropTypes.string,
	onClick: PropTypes.func,
	disabled: PropTypes.bool,
};
