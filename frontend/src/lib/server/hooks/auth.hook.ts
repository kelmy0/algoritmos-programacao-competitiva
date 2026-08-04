import type { Handle } from "@sveltejs/kit";
import { jwtDecode } from "jwt-decode";
import { API_URL } from "$env/static/private";
import { setAuthCookie, deleteAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import type { JwtPayload, RefreshResponse } from "$lib/types/jwt";

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
					isEmployee: decoded.isEmployee
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

		if (!apiError && data?.accessToken) {
			event.locals.accessToken = data.accessToken;

			const decoded = jwtDecode<JwtPayload>(data.accessToken);
			event.locals.user = {
				id: decoded.sub,
				username: decoded.username,
				email: decoded.email,
				permissions: decoded.permissions || [],
				isEmployee: decoded.isEmployee
			};

			setAuthCookie(event.cookies, "access_token", data.accessToken, 15);
		} else {
			deleteAuthCookie(event.cookies, "access_token");
			deleteAuthCookie(event.cookies, "refresh_token");
		}
	}

	return await resolve(event);
};
