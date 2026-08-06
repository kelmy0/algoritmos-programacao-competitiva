import { json, type RequestHandler } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL, JWT_ACCESS_PUBLIC_KEY } from "$env/static/private";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { useMiddlewares } from "$lib/server/middlewares";
import type { AccessJwtPayload, RefreshResponse } from "$lib/types/jwt";
import { validateJWT } from "$lib/utils/jwt";
import { clearAllAuthCookies } from "$lib/utils/cookies";

const refreshToken: RequestHandler = async (event) => {
	const cookieHeader = event.request.headers.get("cookie") || "";
	const clientIp = event.getClientAddress();

	if (cookieHeader === "") {
		return json(normalizeApiError("MISSING_COOKIE"), {
			status: 400
		});
	}

	const { data, error, status, headers } = await customFetch<RefreshResponse>(
		event.fetch,
		`${API_URL}/api/auth/refresh`,
		{
			method: "POST",
			headers: { cookie: cookieHeader, "X-Forwarded-For": clientIp }
		}
	);

	if (error) {
		clearAllAuthCookies(event.cookies);
		return json(error, { status });
	}

	if (!data || !data.accessToken) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR"), {
			status: 500
		});
	}

	const { claims, valid } = await validateJWT<AccessJwtPayload>(
		JWT_ACCESS_PUBLIC_KEY,
		data.accessToken
	);

	if (!valid || !claims) {
		clearAllAuthCookies(event.cookies);
		return json(normalizeApiError("INVALID_ACCESS_TOKEN"), { status: 401 });
	}

	setAuthCookie(event.cookies, "access_token", data.accessToken, 15);

	const response = json({
		accessToken: true,
		expiresAt: claims.exp ? claims.exp * 1000 : null
	});

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	authFlowLimiter
)(refreshToken);
