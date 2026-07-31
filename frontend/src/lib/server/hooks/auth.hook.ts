import type { Handle } from "@sveltejs/kit";
import { jwtDecode } from "jwt-decode";
import { API_URL } from "$env/static/private";
import { setAuthCookie, deleteAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";

interface JwtPayload {
	sub: string;
	username: string;
	email: string;
	permissions: string[];
	is_employee: boolean;
	exp?: number;
}

interface RefreshResponse {
	access_token: string;
}

export const handleAuth: Handle = async ({ event, resolve }) => {
	event.locals.user = null;
	event.locals.accessToken = null;

	const accessToken = event.cookies.get("access_token");
	const refreshToken = event.cookies.get("refresh_token");

	let isTokenValid = false;

	if (accessToken) {
		try {
			const decoded = jwtDecode<JwtPayload>(accessToken);
			const nowInSeconds = Math.floor(Date.now() / 1000);
			const BUFFER_SECONDS = 90;

			if (decoded.exp && decoded.exp - nowInSeconds > BUFFER_SECONDS) {
				isTokenValid = true;

				event.locals.user = {
					id: decoded.sub,
					username: decoded.username,
					email: decoded.email,
					permissions: decoded.permissions || [],
					is_employee: decoded.is_employee
				};
				event.locals.accessToken = accessToken;
			}
		} catch {
			isTokenValid = false;
		}
	}

	if (!isTokenValid && refreshToken) {
		const cookieHeader = event.request.headers.get("cookie") || "";
		const clientIp = event.getClientAddress();

		const { data, error: apiError } = await customFetch<RefreshResponse>(
			event.fetch,
			`${API_URL}/api/auth/refresh`,
			{
				method: "POST",
				headers: {
					cookie: cookieHeader,
					"x-forwarded-for": clientIp,
					"x-real-ip": clientIp
				}
			}
		);

		if (!apiError && data?.access_token) {
			event.locals.accessToken = data.access_token;

			const decoded = jwtDecode<JwtPayload>(data.access_token);
			event.locals.user = {
				id: decoded.sub,
				username: decoded.username,
				email: decoded.email,
				permissions: decoded.permissions || [],
				is_employee: decoded.is_employee
			};

			setAuthCookie(event.cookies, "access_token", data.access_token, 15);
		} else {
			deleteAuthCookie(event.cookies, "access_token");
			deleteAuthCookie(event.cookies, "refresh_token");
		}
	}

	return await resolve(event);
};
