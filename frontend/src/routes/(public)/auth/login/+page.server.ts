import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { normalizeApiError } from "$lib/utils/errors";
import { AUTH_ERRORS } from "$lib/errors/auth/auth_errors";
import { deleteAuthCookie } from "$lib/server/cookies";

export const load: PageServerLoad = async ({ url, locals, cookies }) => {
	const errorCode = url.searchParams.get("error");

	if (errorCode === "TOKEN_NO_LONGER_VALID") {
		deleteAuthCookie(cookies, "access_token");
		deleteAuthCookie(cookies, "refresh_token");
		locals.user = null;
		locals.accessToken = null;
	}

	if (locals.user && !errorCode) {
		const redirectTo = url.searchParams.get("redirectTo");
		const isSafeRedirect = redirectTo && redirectTo.startsWith("/") && !redirectTo.startsWith("//");

		if (isSafeRedirect) {
			redirect(303, redirectTo);
		} else {
			redirect(303, "/");
		}
	}

	const initialError = errorCode
		? normalizeApiError(errorCode, "Erro ao realizar autenticação.", AUTH_ERRORS)
		: null;

	return {
		initialError
	};
};
