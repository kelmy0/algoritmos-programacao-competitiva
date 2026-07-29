import { json, redirect, type Handle, type HandleServerError } from "@sveltejs/kit";
import { PUBLIC_API_URL } from "$env/static/public";
import { jwtDecode } from "jwt-decode";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie, setIdleCookie } from "$lib/server/cookies";
import { deleteAuthCookie } from "$lib/server/cookies";
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

export const handle: Handle = async ({ event, resolve }) => {
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
			`${PUBLIC_API_URL}/api/auth/refresh`,
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
	const user = event.locals.user;
	const isProtectedRoute = event.route.id?.includes("(protected)");
	const isApiRoute = event.url.pathname.startsWith("/api");

	if (isProtectedRoute && !user) {
		if (isApiRoute) {
			const normalizedError = normalizeApiError(
				"INVALID_ACCESS_TOKEN",
				"Seu token de acesso é inválido ou expirou."
			);
			return json(normalizedError, { status: 401 });
		}

		const redirectTo = event.url.pathname + event.url.search;
		redirect(303, `/auth/login?redirectTo=${encodeURIComponent(redirectTo)}`);
	}

	const adminSecret = event.cookies.get("admin_secret");

	if (adminSecret) {
		setIdleCookie(event.cookies, "admin_secret", adminSecret, 60);
	}

	return await resolve(event);
};

export const handleError: HandleServerError = ({ error, event, status }) => {
	if (status === 404) {
		return normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada.");
	}

	const apiError = normalizeApiError(error, "Ocorreu um erro interno no servidor.");
	console.error(`[Server Error ${event.url.pathname}]:`, apiError);

	return normalizeApiError("INTERNAL_ERROR", "Ocorreu um erro inesperado no servidor.");
};
