import { json, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { authFlowLimiter, fiveHundredQuerySize } from "$lib/server/middlewares";
import { hundredKbBodySize, useMiddlewares } from "$lib/server/middlewares";
import { SIGN_UP_ERRORS } from "$lib/errors/auth/sign-up";
import type { SignUpResponse, SignUpServerResponse } from "$lib/types/auth/sign-up";
import { extractDeviceHeaders } from "$lib/utils/headers";

const signUp: RequestHandler = async (event) => {
	const body = await event.request.json().catch(() => ({}));
	const turnstileToken = event.request.headers.get("x-cf-turnstile-response") || "";
	const clientIp = event.getClientAddress();
	const deviceHeaders = extractDeviceHeaders(event.request);

	const { data, error, status, headers } = await customFetch<SignUpResponse>(
		event.fetch,
		`${API_URL}/api/auth/sign-up`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": clientIp,
				"X-CF-Turnstile-Response": turnstileToken,
				...deviceHeaders
			},
			body: JSON.stringify(body)
		},
		SIGN_UP_ERRORS
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

	const sanitizedResponse: SignUpServerResponse = {
		success: data.success,
		autoLogin: data.autoLogin
	};

	const response = json(sanitizedResponse);

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
)(signUp);
