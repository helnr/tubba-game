import "./Button.css";
export default function Button(props) {
	return (
		<button onClick={props.onClick} className="main-btn">
			{props.value}
		</button>
	);
}
