import { Route, Routes } from "react-router";
import Auth from "./components/pages/auth";
import Home from "./components/pages/home";
import "./App.css";
export default function App() {
	return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/auth" element={<Auth />} />
			<Route path="*" element={<Home />} />
		</Routes>
	);
}
