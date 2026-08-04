import { json } from "@sveltejs/kit";
import type { RequestHandler } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import { authFlowLimiter, fiveHundredQuerySize } from "$lib/server/middlewares";
import { hundredKbBodySize, useMiddlewares } from "$lib/server/middlewares";
import { AUTH_ERRORS } from "$lib/errors/auth/auth_errors";
import type { LoginResponse, LoginServerResponse } from "$lib/types/auth/login";

const login: RequestHandler = async (event) => {
	const clientIp = event.getClientAddress();
	const turnstileToken = event.request.headers.get("x-cf-turnstile-response") || "";

	const { data, error, status, headers } = await customFetch<LoginResponse>(
		event.fetch,
		`${API_URL}/api/auth/login`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": clientIp,
				"X-CF-Turnstile-Response": turnstileToken
			},
			body: JSON.stringify(await event.request.json())
		},
		AUTH_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	if (!data) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR"), { status: 500 });
	}

	// This will be removed
	if (data.accessToken) {
		setAuthCookie(event.cookies, "access_token", data.accessToken, 15);
	}

	const sanitizedResponse: LoginServerResponse = {
		accessToken: Boolean(data.accessToken),
		requires2FA: data.requires2FA,
		preAuthToken: data.preAuthToken || ""
	};

	const response = json(sanitizedResponse);

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const POST = useMiddlewares(fiveHundredQuerySize, hundredKbBodySize, authFlowLimiter)(login);
