import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
	plugins: [react()],
	server: {
		proxy: {
			"/api/game/join": {
				target: "ws://localhost:8080",
				changeOrigin: true,
				secure: false,
				ws: true,
				rewrite: (path) => path.replace(/^\/api\/game\/join/, "/api/game/join"),
			},
			"/api": {
				target: "http://localhost:8080",
				changeOrigin: true,
				secure: false,
				rewrite: (path) => path.replace(/^\/api/, "/api"),
			},
		},
	},
});
