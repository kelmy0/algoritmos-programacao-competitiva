import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { requireAuth, useMiddlewares } from "$lib/server/middlewares";
import type { RequestHandler } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import { disable2FASchema, generate2FASchema } from "$lib/schemas/me";
import { normalizeApiError } from "$lib/utils/errors";
import { extractDeviceHeaders } from "$lib/utils/headers";
import { clearAllAuthCookies } from "$lib/utils/cookies";

const disable2FA: RequestHandler = async (event) => {
	const cookieHeader = event.request.headers.get("cookie") || "";

	if (cookieHeader === "") {
		return json(normalizeApiError("MISSING_COOKIE"), {
			status: 400
		});
	}

	const body = await event.request.json().catch(() => null);
	const result = disable2FASchema.safeParse(body);

	if (!result.success) {
		return json(normalizeApiError("INVALID_REQUEST_BODY"), { status: 400 });
	}

	const deviceHeaders = extractDeviceHeaders(event.request);

	const { error, status } = await customFetch<null>(
		event.fetch,
		`${API_URL}/api/users/me/2fa/disable`,
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

	if (error) {
		return json(error, { status });
	}

	clearAllAuthCookies(event.cookies);

	event.locals.user = null;
	event.locals.accessToken = null;

	return new Response(null, { status: 204 });
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	requireAuth,
	authFlowLimiter
)(disable2FA);
