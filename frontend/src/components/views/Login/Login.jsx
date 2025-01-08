import Form from "../../ui/Form";
import Button from "../../ui/Button";
import TextField from "../../ui/TextField";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import store from "../../../services/redux/store";
import { login } from "../../../services/redux/slices/authSlice";

const loginSchema = z.object({
	email: z.string().email({ message: "email is not valid" }),
	password: z
		.string()
		.min(8, { message: "password must be at least 8 characters" }),
});

import "./Login.css";

export default function Login() {
	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
		setError,
	} = useForm({
		resolver: zodResolver(loginSchema),
	});

	const onSubmit = async (data) => {
		try {
			await new Promise((r) => setTimeout(r, 2000));
			store.dispatch(login(data));
		} catch (error) {
			setError("root", { message: error.message });
		}
	};

	return (
		<Form onSubmit={handleSubmit(onSubmit)}>
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
			<Button value="Login" />
			{isSubmitting && <p>Welcome Back...</p>}
			{errors.root && <p>{errors.root.message}</p>}
		</Form>
	);
}
