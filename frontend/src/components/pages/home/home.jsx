import useAuth from "../../../hooks/useAuth";
import Button from "../../ui/Button";
import { useDispatch } from "react-redux";
import { logout } from "../../../services/redux/slices/authSlice";

import "./home.css";
export function Home() {
	const dispatch = useDispatch();

	const {
		auth: { user, error },
		loading,
	} = useAuth();
	if (loading) {
		return <div className="home">Loading...</div>;
	}

	if (error) {
		return <div className="home">{error}</div>;
	}

	return (
		<div className="home">
			{user && (
				<div className="user-info">
					<h2>hello {user.name}</h2>
				</div>
			)}
			<Button value="create game" />
			<Button value="join game" />
			<Button
				onClick={() => {
					dispatch(logout());
				}}
				value="logout"
			/>
		</div>
	);
}
