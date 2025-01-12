import { Route, Routes } from "react-router";
import Auth from "./components/pages/auth";
import Home from "./components/pages/home";
import "./App.css";
import { useEffect, useLayoutEffect, useState } from "react";
import axios from "axios";
import { useDispatch } from "react-redux";
import { setLoggedIn } from "./services/redux/slices/authSlice";

export default function App() {
	const dispatch = useDispatch();
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		const fetchUser = async () => {
			try {
				const apiURL = import.meta.env.VITE_API_URL;
				const response = await axios.get(`${apiURL}/user/me`, {
					withCredentials: true,
				});
				if (response.status === 200) {
					const data = await response.data;
					if (data.status === "success") {
						dispatch(setLoggedIn(data.data));
						setLoading(false);
					}
				}
			} catch (error) {
				console.log(error);
				setLoading(false);
			}
		};

		fetchUser();
	}, [dispatch, setLoading]);

	return loading ? null : (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/auth" element={<Auth />} />
			<Route path="*" element={<Home />} />
		</Routes>
	);
}
