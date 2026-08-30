import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { requireAuth, useMiddlewares } from "$lib/server/middlewares";
import type { RequestHandler } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import { enable2FASchema } from "$lib/schemas/me";
import { extractDeviceHeaders } from "$lib/utils/headers";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import { clearAllAuthCookies } from "$lib/utils/cookies";

interface TwoFactorEnableResponse {
	accessToken: string;
}

const enable2FA: RequestHandler = async (event) => {
	const cookieHeader = event.request.headers.get("cookie") || "";

	if (cookieHeader === "") {
		return json(normalizeApiError("MISSING_COOKIE"), {
			status: 400
		});
	}

	const body = await event.request.json().catch(() => null);
	const result = enable2FASchema.safeParse(body);

	if (!result.success) {
		return json(normalizeApiError("INVALID_REQUEST_BODY"), { status: 400 });
	}

	const deviceHeaders = extractDeviceHeaders(event.request);

	const { data, error, status, headers } = await customFetch<TwoFactorEnableResponse>(
		event.fetch,
		`${API_URL}/api/users/me/2fa/enable`,
		{
			method: "POST",
			headers: {
				cookie: cookieHeader,
				Authorization: `Bearer ${event.locals.accessToken}`,
				...deviceHeaders
			},
			body: JSON.stringify(result.data)
		},
		TWO_FACTOR_ERRORS
	);

	if (error || !data?.accessToken) {
		if ((status === 401 && error?.code !== "2FA_INVALID_CODE") || status === 403) {
			clearAllAuthCookies(event.cookies);
		}
		return json(error || normalizeApiError("INTERNAL_SERVER_ERROR"), { status: status || 500 });
	}

	setAuthCookie(event.cookies, "access_token", data.accessToken, 15);

	const response = new Response(null, { status: 204 });

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	requireAuth,
	authFlowLimiter
)(enable2FA);
