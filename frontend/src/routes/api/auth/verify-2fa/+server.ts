import { json, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import {
	TWO_FACTOR_ERRORS,
	type TwoFactorServerResponse
} from "../../../(public)/auth/verify-2fa/two_factor_verify.svelte";
import { customFetch } from "$lib/api/client";
import type { LoginResponse } from "../login/+server";
import { rateLimit, useMiddlewares } from "$lib/server/middlewares";

const verify2FA: RequestHandler = async (event) => {
	const clientIp = event.getClientAddress();
	const { data, error, status, headers } = await customFetch<LoginResponse>(
		event.fetch,
		`${API_URL}/api/auth/verify-2fa`,
		{
			method: "POST",
			headers: { "Content-Type": "application/json", "X-Forwarded-For": clientIp },
			body: JSON.stringify(await event.request.json())
		},
		TWO_FACTOR_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	if (!data) {
		return json(
			normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Resposta inválida do servidor.",
				TWO_FACTOR_ERRORS
			),
			{ status: 500 }
		);
	}

	//this will be removed
	if (data.access_token) {
		setAuthCookie(event.cookies, "access_token", data.access_token, 15);
	}

	const sanitizedResponse: TwoFactorServerResponse = {
		access_token: Boolean(data.access_token),
		requires_2fa: data.requires_2fa
	};

	const response = json(sanitizedResponse);

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const GET = useMiddlewares(rateLimit({ capacity: 5, fillRate: 0.1 }))(verify2FA);
