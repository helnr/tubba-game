import Form from "../../ui/Form";
import Button from "../../ui/Button";
import TextField from "../../ui/TextField";
import { useForm } from "react-hook-form";
import { useDispatch } from "react-redux";
import z from "zod";
import axios from "axios";
import { zodResolver } from "@hookform/resolvers/zod";
import { setLoggedIn } from "../../../services/redux/slices/authSlice";

import "./Register.css";

const registerSchema = z.object({
	name: z.string(),
	email: z.string().email({ message: "email is not valid" }),
	password: z
		.string()
		.min(8, { message: "password must be at least 8 characters" }),
});

export default function Register() {
	const dispatch = useDispatch();

	const {
		register,
		handleSubmit,
		setError,
		formState: { errors, isSubmitting },
	} = useForm({
		resolver: zodResolver(registerSchema),
	});

	const onSubmit = async (data) => {
		try {
			const apiURL = import.meta.env.VITE_API_URL;
			const response = await axios.post(`${apiURL}/auth/register`, data, {
				headers: {
					"Content-Type": "application/json",
				},
				withCredentials: true,
			});

			if (response.status === 201) {
				const data = await response.data;
				if (data.status === "success") {
					dispatch(setLoggedIn(data.data));
				}
			}
		} catch (error) {
			setError("root", { message: error.response.data.error });
		}
	};

	return (
		<Form onSubmit={handleSubmit(onSubmit)}>
			<TextField name="name" register={register} placeholder="Name" />
			<TextField
				name="email"
				register={register}
				type="email"
				placeholder="Email"
			/>
			{errors.email && <p>{errors.email.message}</p>}
			<TextField
				name="password"
				register={register}
				type="password"
				placeholder="Password"
			/>
			{errors.password && <p>{errors.password.message}</p>}
			<Button value="Register" />
			{isSubmitting && <p>Registreing...</p>}
			{errors.root && <p>Error: {errors.root.message}</p>}
		</Form>
	);
}
