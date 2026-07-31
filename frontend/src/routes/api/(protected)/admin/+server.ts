import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";
import { ADMIN_PASSWORD_ERRORS } from "../../../(protected)/admin/admin_controller.svelte";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import { setIdleCookie } from "$lib/server/cookies";

interface AdminPasswordResponse {
	message: string;
	status: string;
}

export const POST: RequestHandler = async (event) => {
	const { password } = await event.request.json();

	if (!password || password.length < 8) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_PASSWORD",
			"Falta a senha das rotas administradoras.",
			ADMIN_PASSWORD_ERRORS
		);
		return json(normalizedError, { status: 401 });
	}

	const { data, error, status } = await customFetch<AdminPasswordResponse>(
		fetch,
		`${API_URL}/api/admin/ping`,
		{
			method: "GET",
			headers: {
				"X-Admin-Secret": password,
				Authorization: `Bearer ${event.locals.accessToken}`
			}
		},
		ADMIN_PASSWORD_ERRORS
	);

	if (error) {
		if (error.code === "PAGE_NOT_FOUND") {
			return json({ correct: false });
		}

		return json(error, { status });
	}

	if (!data) {
		return json(
			normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Resposta inválida do servidor.",
				ADMIN_PASSWORD_ERRORS
			),
			{ status: 500 }
		);
	}

	if (data.message !== "pong") {
		return json({ correct: false });
	}

	setIdleCookie(event.cookies, "admin_secret", password, 15);

	return json({ correct: true });
};
