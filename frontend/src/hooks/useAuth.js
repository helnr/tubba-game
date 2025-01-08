import { useSelector } from "react-redux";
import { useNavigate } from "react-router";
import { useEffect } from "react";
import { useState } from "react";

const useAuth = () => {
	const auth = useSelector((state) => state.auth);
	const [loading, setLoading] = useState(true);
	const navigate = useNavigate();

	useEffect(() => {
		if (!auth.loggedIn) {
			navigate("/auth");
		}

		setLoading(false);
	}, [auth, navigate]);

	return { auth, loading };
};

export default useAuth;
