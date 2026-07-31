import { normalizeApiError } from "$lib/utils/errors";
import { json, redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, url }) => {
	if (!locals.user) {
		const isApiRoute = url.pathname.startsWith("/api");

		if (isApiRoute) {
			const normalizedError = normalizeApiError(
				"INVALID_ACCESS_TOKEN",
				"Seu token de acesso é inválido ou expirou."
			);
			return json(normalizedError, { status: 401 });
		}

		const redirectTo = url.pathname + url.search;
		redirect(303, `/auth/login?redirectTo=${encodeURIComponent(redirectTo)}`);
	}

	return {
		user: locals.user
	};
};
