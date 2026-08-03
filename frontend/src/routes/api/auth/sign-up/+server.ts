import { json, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import {
	SIGN_UP_ERRORS,
	type SignUpServerResponse
} from "../../../(public)/auth/sign-up/sign_up.svelte";
import { customFetch } from "$lib/api/client";
import {
	authFlowLimiter,
	fiveHundredQuerySize,
	hundredKbBodySize,
	useMiddlewares
} from "$lib/server/middlewares";

interface SignUpResponse {
	access_token?: string;
	success: boolean;
	auto_login: boolean;
}

const signUp: RequestHandler = async (event) => {
	const body = await event.request.json().catch(() => ({}));
	const turnstileToken = event.request.headers.get("x-cf-turnstile-response") || "";
	const clientIp = event.getClientAddress();

	const { data, error, status, headers } = await customFetch<SignUpResponse>(
		event.fetch,
		`${API_URL}/api/auth/sign-up`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": clientIp,
				"X-CF-Turnstile-Response": turnstileToken
			},
			body: JSON.stringify(body)
		},
		SIGN_UP_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	if (!data) {
		return json(
			normalizeApiError("INTERNAL_SERVER_ERROR", "Resposta inválida do servidor.", SIGN_UP_ERRORS),
			{ status: 500 }
		);
	}

	// This will be removed
	if (data.access_token) {
		setAuthCookie(event.cookies, "access_token", data.access_token, 15);
	}

	const sanitizedResponse: SignUpServerResponse = {
		success: data.success,
		autoLogin: data.auto_login
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
