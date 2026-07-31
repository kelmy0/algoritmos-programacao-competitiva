import { redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { API_URL } from "$env/static/private";
import { deleteAuthCookie } from "$lib/server/cookies";

export const POST: RequestHandler = async ({ fetch: svelteFetch, locals, cookies }) => {
	if (locals.accessToken) {
		try {
			await svelteFetch(`${API_URL}/api/auth/logout`, {
				method: "POST",
				headers: {
					Authorization: `Bearer ${locals.accessToken}`
				}
			});
		} catch (err) {
			console.error("Erro ao notificar logout na API Go:", err);
		}
	}

	deleteAuthCookie(cookies, "access_token");
	deleteAuthCookie(cookies, "refresh_token");
	deleteAuthCookie(cookies, "admin_secret");

	locals.user = null;
	locals.accessToken = null;

	redirect(303, "/auth/login");
};
