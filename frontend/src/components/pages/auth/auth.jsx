import "./auth.css";
import Login from "../../views/Login";
import Register from "../../views/Register";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import useAuth from "../../../hooks/useAuth";

export function Auth() {
	const {
		auth: { loggedIn, error },
		loading,
	} = useAuth();
	const [view, setView] = useState("login");
	const navigate = useNavigate();

	useEffect(() => {
		if (loggedIn) {
			navigate("/");
		}
	}, [loggedIn, navigate]);

	if (loading) {
		return <div className="auth">Loading...</div>;
	}

	return (
		<div className="auth">
			<ul>
				<li key="login">
					<button
						className={view === "login" ? "auth-btn active" : "auth-btn"}
						onClick={() => setView("login")}>
						LOGIN
					</button>
				</li>
				<li key="register">
					<button
						className={view === "register" ? "auth-btn active" : "auth-btn"}
						onClick={() => setView("register")}>
						REGISTER
					</button>
				</li>
			</ul>
			<div className="view">
				{view === "login" ? <Login error={error} /> : <Register />}
			</div>
		</div>
	);
}
