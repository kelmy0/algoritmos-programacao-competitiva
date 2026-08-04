import { json, type RequestHandler } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import { deleteAuthCookie, setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { useMiddlewares } from "$lib/server/middlewares";
import type { JwtPayload, RefreshResponse } from "$lib/types/jwt";
import { jwtDecode } from "jwt-decode";

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
		deleteAuthCookie(event.cookies, "access_token");
		return json(error, { status });
	}

	if (!data || !data.accessToken) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR"), {
			status: 500
		});
	}

	const decoded = jwtDecode<JwtPayload>(data.accessToken);
	setAuthCookie(event.cookies, "access_token", data.accessToken, 15);

	const response = json({
		accessToken: true,
		expiresAt: decoded.exp ? decoded.exp * 1000 : undefined
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
