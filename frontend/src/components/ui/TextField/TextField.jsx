import "./TextField.css";

export default function TextField(props) {
	let register = {};

	if (props.register) {
		register = props.register(
			props.name,
			props.rules ? { ...props.rules } : {}
		);
	}

	return (
		<input
			{...register}
			className="text-field"
			placeholder={props.placeholder}
			type={props.type || "text"}
		/>
	);
}
