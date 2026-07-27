import { json, type RequestHandler } from "@sveltejs/kit";
import { PUBLIC_API_URL } from "$env/static/public";
import { normalizeApiError } from "$lib/utils/errors";
import { setAuthCookie } from "$lib/server/cookies";
import {
	SIGN_UP_ERRORS,
	type SignUpServerResponse
} from "../../../../(public)/auth/sign-up/sign_up.svelte";
import { customFetch } from "$lib/api/client";

interface SignUpResponse {
	access_token?: string;
	success: boolean;
	auto_login: boolean;
}

export const POST: RequestHandler = async ({ fetch: svelteFetch, request, cookies }) => {
	const body = await request.json().catch(() => ({}));

	const { data, error, status, headers } = await customFetch<SignUpResponse>(
		svelteFetch,
		`${PUBLIC_API_URL}/api/auth/sign-up`,
		{
			method: "POST",
			headers: { "Content-Type": "application/json" },
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
		setAuthCookie(cookies, "access_token", data.access_token, 15);
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
