import { json, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { normalizeApiError } from "$lib/utils/errors";
import { deleteCookie, setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { authFlowLimiter, fiveHundredQuerySize } from "$lib/server/middlewares";
import { hundredKbBodySize, useMiddlewares } from "$lib/server/middlewares";
import { TWO_FACTOR_ERRORS } from "$lib/errors/auth/verify-2fa";
import type { TwoFactorServerResponse } from "$lib/types/auth/two-factor";
import type { LoginResponse } from "$lib/types/auth/login";
import { extractDeviceHeaders } from "$lib/utils/headers";

const verify2FA: RequestHandler = async (event) => {
	const clientIp = event.getClientAddress();
	const turnstileToken = event.request.headers.get("x-cf-turnstile-response") || "";
	const deviceHeaders = extractDeviceHeaders(event.request);
	const cookieHeader = event.request.headers.get("cookie") || "";

	if (cookieHeader === "") {
		return json(normalizeApiError("MISSING_COOKIE"), {
			status: 400
		});
	}

	const { data, error, status, headers } = await customFetch<LoginResponse>(
		event.fetch,
		`${API_URL}/api/auth/verify-2fa`,
		{
			method: "POST",
			headers: {
				cookie: cookieHeader,
				"Content-Type": "application/json",
				"X-Forwarded-For": clientIp,
				"X-CF-Turnstile-Response": turnstileToken,
				...deviceHeaders
			},
			body: JSON.stringify(await event.request.json())
		},
		TWO_FACTOR_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	if (!data) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR"), { status: 500 });
	}

	deleteCookie(event.cookies, "pre_auth_token");

	//this will be removed
	if (data.accessToken) {
		setAuthCookie(event.cookies, "access_token", data.accessToken, 15);
	}

	const sanitizedResponse: TwoFactorServerResponse = {
		accessToken: Boolean(data.accessToken),
		requires2FA: data.requires2FA
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
)(verify2FA);
