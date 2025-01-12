import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";

const initialState = {
	loggedIn: false,
	user: null,
	error: null,
};

const authSlice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		setLoggedIn: (state, action) => {
			state.loggedIn = true;
			state.user = action.payload;
			state.error = null;
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(login.fulfilled, (state, action) => {
				state.loggedIn = true;
				state.user = action.payload;
				state.error = null;
			})
			.addCase(login.rejected, (state, action) => {
				state.loggedIn = false;
				state.user = null;
				state.error = action.payload || "Login failed";
			})
			.addCase(logout.fulfilled, (state) => {
				state.loggedIn = false;
				state.user = null;
			});
	},
});

const login = createAsyncThunk(
	"auth/login",
	async (user, { rejectWithValue }) => {
		try {
			const apiURL = import.meta.env.VITE_API_URL;
			const response = await axios.post(`${apiURL}/auth/login`, user, {
				headers: {
					"Content-Type": "application/json",
				},
				withCredentials: true,
			});

			if (response.status !== 200) {
				throw new Error(response.data.error);
			}
			const data = await response.data;

			return data.data;
		} catch (error) {
			return rejectWithValue(error.response.data.error);
		}
	}
);

const logout = createAsyncThunk("auth/logout", async () => {
	return null;
});

export { login, logout };

export const { setLoggedIn } = authSlice.actions;

export default authSlice.reducer;
