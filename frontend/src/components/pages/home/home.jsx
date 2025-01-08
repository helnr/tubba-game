import useAuth from "../../../hooks/useAuth";
import Button from "../../ui/Button";
import "./home.css";
export function Home() {
	const { loading } = useAuth();
	if (loading) {
		return <div className="home">Loading...</div>;
	}

	return (
		<div className="home">
			<Button value="create game" />
			<Button value="join game" />
		</div>
	);
}
