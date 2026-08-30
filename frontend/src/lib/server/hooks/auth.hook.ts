import type { Handle } from "@sveltejs/kit";
import { API_URL, JWT_ACCESS_PUBLIC_KEY } from "$env/static/private";
import { setAuthCookie, deleteCookie, setIdleCookie, syncServerCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import type { AccessJwtPayload, RefreshResponse } from "$lib/types/jwt";
import { validateJWT } from "$lib/utils/jwt";
import { clearAllAuthCookies } from "$lib/utils/cookies";
import { extractDeviceHeaders } from "$lib/utils/headers";

let serverRefreshPromise: Promise<{
	data: RefreshResponse | null;
	headers: Headers | null;
	error: unknown;
}> | null = null;

export const handleAuth: Handle = async ({ event, resolve }) => {
	event.locals.user = null;
	event.locals.accessToken = null;

	const isRefreshRoute =
		event.url.pathname === "/api/auth/refresh" && event.request.method === "POST";

	const accessToken = event.cookies.get("access_token");
	const refreshToken = event.cookies.get("refresh_token");

	let isTokenValid = false;

	if (accessToken) {
		try {
			const { claims, valid } = await validateJWT<AccessJwtPayload>(
				JWT_ACCESS_PUBLIC_KEY,
				accessToken
			);

			if (!valid || !claims) {
				throw new Error("INVALID_ACCESS_TOKEN");
			}

			const nowInSeconds = Math.floor(Date.now() / 1000);
			const BUFFER_SECONDS = 60;

			if (claims.exp && claims.exp - nowInSeconds > BUFFER_SECONDS) {
				isTokenValid = true;

				event.locals.user = {
					id: claims.sub,
					name: claims.name,
					username: claims.username,
					email: claims.email,
					permissions: claims.permissions || [],
					isEmployee: claims.isEmployee,
					is2FAEnabled: claims.is2FAEnabled || false,
					hasPassword: claims.hasPassword ?? false
				};
				event.locals.accessToken = accessToken;
			}
		} catch {
			isTokenValid = false;
			deleteCookie(event.cookies, "access_token");
		}
	}

	let rotatedCookies: string[] = [];

	if (!isTokenValid && refreshToken && !serverRefreshPromise && !isRefreshRoute) {
		const cookieHeader = event.request.headers.get("cookie") || "";
		const clientIp = event.getClientAddress();
		const deviceHeaders = extractDeviceHeaders(event.request);

		serverRefreshPromise = (async () => {
			try {
				return await customFetch<RefreshResponse>(event.fetch, `${API_URL}/api/auth/refresh`, {
					method: "POST",
					headers: {
						cookie: cookieHeader,
						"X-Forwarded-For": clientIp,
						"X-Real-Ip": clientIp,
						...deviceHeaders
					}
				});
			} finally {
				serverRefreshPromise = null;
			}
		})();

		const { data, error: apiError, headers } = await serverRefreshPromise;

		if (!apiError && data?.accessToken) {
			const { claims, valid } = await validateJWT<AccessJwtPayload>(
				JWT_ACCESS_PUBLIC_KEY,
				data.accessToken
			);

			if (valid && claims) {
				event.locals.user = {
					id: claims.sub,
					name: claims.name,
					username: claims.username,
					email: claims.email,
					permissions: claims.permissions || [],
					isEmployee: claims.isEmployee,
					is2FAEnabled: claims.is2FAEnabled || false,
					hasPassword: claims.hasPassword ?? false
				};
				event.locals.accessToken = data.accessToken;

				setAuthCookie(event.cookies, "access_token", data.accessToken, 15);

				if (headers) {
					rotatedCookies = headers.getSetCookie();
					for (const cookieString of rotatedCookies) {
						const [nameValue] = cookieString.split(";");
						const eqIdx = nameValue.indexOf("=");

						if (eqIdx !== -1) {
							const name = nameValue.substring(0, eqIdx).trim();
							const value = nameValue.substring(eqIdx + 1).trim();

							if (cookieString.includes("Max-Age=0") || !value) {
								deleteCookie(event.cookies, name);
							} else {
								syncServerCookie(event.cookies, name, value);
							}
						}
					}
				}
			} else {
				clearAllAuthCookies(event.cookies);
			}
		} else {
			clearAllAuthCookies(event.cookies);
		}
	}

	const response = await resolve(event);
	for (const cookieString of rotatedCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};
