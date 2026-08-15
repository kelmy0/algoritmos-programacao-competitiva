import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { requireAuth, useMiddlewares } from "$lib/server/middlewares";
import type { RequestHandler } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import { deleteCookie } from "$lib/server/cookies";
import { generate2FASchema } from "$lib/schemas/me";

const disable2FA: RequestHandler = async (event) => {
	const body = await event.request.json().catch(() => null);
	const result = generate2FASchema.safeParse(body);

	const { error, status } = await customFetch<null>(
		event.fetch,
		`${API_URL}/api/users/me/2fa/disable`,
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

	return new Response(null, { status: 204 });
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	requireAuth,
	authFlowLimiter
)(disable2FA);
