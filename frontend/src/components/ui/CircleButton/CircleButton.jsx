import PropsType from "prop-types";

import classes from "./CircleButton.module.css";
export default function CircleButton(props) {
	const onClick = props.onClick;
	const width = props.width || 40;

	return (
		<button
			onClick={onClick}
			style={{ width: `${width}px` }}
			className={`${classes["circle-button"]} circle-button`}>
			{props.children}
		</button>
	);
}

CircleButton.propTypes = {
	onClick: PropsType.func,
	width: PropsType.number,
	children: PropsType.node,
};
