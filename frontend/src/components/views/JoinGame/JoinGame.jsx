import Form from "../../ui/Form";
import TextField from "../../ui/TextField";
import Button from "../../ui/Button";

import axios from "axios";
import z, { set } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import "./JoinGame.css";
export default function JoinGame(props) {
	const navigate = useNavigate();
	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
		setError,
	} = useForm({
		resolver: zodResolver(
			z.object({
				gameCode: z.string().regex(/^[a-zA-Z0-9]{6,40}$/, "invalid game code"),
			})
		),
	});

	const onSubmit = async (data) => {
		try {
			const apiURL = import.meta.env.VITE_API_URL;
			const response = await axios.get(`${apiURL}/game/${data.gameCode}`, {
				withCredentials: true,
			});

			if (response.status === 200) {
				const gameData = await response.data.game;
				console.log(gameData);
				navigate("/game", { state: { game: gameData } });
			}
		} catch (error) {
			const message = error?.response?.data?.error || "error";
			setError("root", { message });
		}
	};

	return (
		<Form>
			<TextField name="gameCode" placeholder="Game Code" register={register} />
			{errors.gameCode && <p>{errors.gameCode.message}</p>}
			<Button value="join" onClick={handleSubmit(onSubmit)} />
			<Button onClick={props.onMenu} value="menu" />
			{isSubmitting && <p>Joining Game...</p>}
			{errors.root && <p>{errors.root.message}</p>}
		</Form>
	);
}
