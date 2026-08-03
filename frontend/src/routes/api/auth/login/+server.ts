import { json } from "@sveltejs/kit";
import type { RequestHandler } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import { AUTH_ERRORS, type LoginServerResponse } from "../../../(public)/auth/login/login.svelte";
import { setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import {
	authFlowLimiter,
	fiveHundredQuerySize,
	hundredKbBodySize,
	useMiddlewares
} from "$lib/server/middlewares";

export interface LoginResponse {
	access_token?: string;
	requires_2fa: boolean;
	pre_auth_token?: string;
}

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
		return json(
			normalizeApiError("INTERNAL_SERVER_ERROR", "Resposta inválida do servidor.", AUTH_ERRORS),
			{ status: 500 }
		);
	}

	// This will be removed
	if (data.access_token) {
		setAuthCookie(event.cookies, "access_token", data.access_token, 15);
	}

	const sanitizedResponse: LoginServerResponse = {
		access_token: Boolean(data.access_token),
		requires_2fa: data.requires_2fa,
		pre_auth_token: data.pre_auth_token || ""
	};

	const response = json(sanitizedResponse);

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const POST = useMiddlewares(fiveHundredQuerySize, hundredKbBodySize, authFlowLimiter)(login);
