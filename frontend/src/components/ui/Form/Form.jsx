import "./Form.css";

export default function Form(props) {
	return (
		<form onSubmit={props.onSubmit} className="form-view">
			{props.children}
		</form>
	);
}
