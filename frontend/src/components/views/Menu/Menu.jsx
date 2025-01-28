import Button from "../../ui/Button";
import { useDispatch } from "react-redux";
import { logout } from "../../../services/redux/slices/authSlice";

import "./Menu.css";
export default function Menu(props) {
	const dispatch = useDispatch();
	return (
		<>
			<Button onClick={props.onCreate} value="create game" />
			<Button onClick={props.onJoin} value="join game" />
			<Button
				onClick={() => {
					dispatch(logout());
				}}
				value="logout"
			/>
		</>
	);
}
