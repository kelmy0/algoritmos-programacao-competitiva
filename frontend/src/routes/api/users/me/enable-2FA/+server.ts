import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { requireAuth, useMiddlewares } from "$lib/server/middlewares";
import type { RequestHandler } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import { enable2FASchema } from "$lib/schemas/me";
import { deleteCookie } from "$lib/server/cookies";

const generate2FA: RequestHandler = async (event) => {
	const body = await event.request.json().catch(() => null);
	const result = enable2FASchema.safeParse(body);

	const { error, status } = await customFetch<null>(
		event.fetch,
		`${API_URL}/api/users/me/2fa/enable`,
		{
			method: "POST",
			headers: {
				Authorization: `Bearer ${event.locals.accessToken}`
			},
			body: JSON.stringify(result.data)
		},
		TWO_FACTOR_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	deleteCookie(event.cookies, "access_token");

	return new Response(null, { status: 204 });
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	requireAuth,
	authFlowLimiter
)(generate2FA);
