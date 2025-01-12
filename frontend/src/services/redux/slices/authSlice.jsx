import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";

const initialState = {
	loggedIn: false,
	user: null,
	status: null,
};

const authSlice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		setLoggedIn: (state, action) => {
			state.loggedIn = true;
			state.user = action.payload;
			state.status = null;
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(login.fulfilled, (state, action) => {
				state.loggedIn = true;
				state.user = action.payload;
				state.status = null;
			})
			.addCase(login.pending, (state) => {
				state.status = "Logging in...";
			})
			.addCase(login.rejected, (state, action) => {
				state.loggedIn = false;
				state.user = null;
				state.status = action.payload || "Login failed";
			})
			.addCase(logout.fulfilled, (state) => {
				state.loggedIn = false;
				state.user = null;
			})
			.addCase(logout.rejected, (state, action) => {
				state.status = action.payload || "Logout failed";
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

const logout = createAsyncThunk(
	"auth/logout",
	async (data, { rejectWithValue }) => {
		try {
			const apiURL = import.meta.env.VITE_API_URL;
			const response = await axios.post(
				`${apiURL}/user/logout`,
				{},
				{
					withCredentials: true,
				}
			);
			if (response.status !== 200) {
				throw new Error(response.data.error);
			}
			return response.data.data;
		} catch (error) {
			return rejectWithValue(error.response.data.error);
		}
	}
);

export { login, logout };

export const { setLoggedIn } = authSlice.actions;

export default authSlice.reducer;
