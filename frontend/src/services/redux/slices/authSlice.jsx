import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";

const initialState = {
	loggedIn: false,
	user: null,
};

const authSlice = createSlice({
	name: "auth",
	initialState,
	extraReducers: (builder) => {
		builder
			.addCase(login.fulfilled, (state, action) => {
				state.loggedIn = true;
				state.user = action.payload;
			})
			.addCase(logout.fulfilled, (state) => {
				state.loggedIn = false;
				state.user = null;
			});
	},
});

const login = createAsyncThunk("auth/login", async (user) => {
	return user;
});

const logout = createAsyncThunk("auth/logout", async () => {
	return null;
});

export { login, logout };

export default authSlice.reducer;
