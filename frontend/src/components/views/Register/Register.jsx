import Form from "../../ui/Form";
import Button from "../../ui/Button";
import TextField from "../../ui/TextField";
import { useForm } from "react-hook-form";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import "./Register.css";

const registerSchema = z.object({
	name: z.string(),
	email: z.string().email({ message: "email is not valid" }),
	password: z
		.string()
		.min(8, { message: "password must be at least 8 characters" }),
});

export default function Register() {
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
			console.log(data);
			await new Promise((r) => setTimeout(r, 2000));
			throw new Error("Something went wrong");
		} catch (error) {
			setError("root", { message: error.message });
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
			{isSubmitting && <p>Sending Request...</p>}
			{errors.root && <p>{errors.root.message}</p>}
		</Form>
	);
}
