import type { HandleServerError } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";
import { sequence } from "@sveltejs/kit/hooks";
import { handleAuth } from "$lib/server/hooks/auth.hook";
import { handleAdminSession } from "$lib/server/hooks/admin.hook";

export const handle = sequence(handleAuth, handleAdminSession);

export const handleError: HandleServerError = ({ error, event, status }) => {
	if (status === 404) {
		return normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada.");
	}

	const apiError = normalizeApiError(error, "Ocorreu um erro interno no servidor.");
	console.error(`[Server Error ${event.url.pathname}]:`, apiError);

	return normalizeApiError("INTERNAL_ERROR", "Ocorreu um erro inesperado no servidor.");
};
