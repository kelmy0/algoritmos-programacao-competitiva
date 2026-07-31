import { json } from "@sveltejs/kit";
import type { RequestHandler } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import {
	AUTH_ERRORS,
	type LoginServerResponse
} from "../../../../(public)/auth/login/login.svelte";
import { setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";

export interface LoginResponse {
	access_token?: string;
	requires_2fa: boolean;
	pre_auth_token?: string;
}

export const POST: RequestHandler = async ({ fetch: svelteFetch, request, cookies }) => {
	const { data, error, status, headers } = await customFetch<LoginResponse>(
		svelteFetch,
		`${API_URL}/api/auth/login`,
		{
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(await request.json())
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
		setAuthCookie(cookies, "access_token", data.access_token, 15);
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
