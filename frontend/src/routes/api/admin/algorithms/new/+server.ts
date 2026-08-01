import { algorithmSchema } from "$lib/schemas/algorithm";
import { json, type RequestHandler } from "@sveltejs/kit";
import { ADMIN_ALGORITHMS_ERRORS } from "../../../../(protected)/admin/algorithms/new/newAlgorithm.svelte";
import { normalizeApiError } from "$lib/utils/errors";
import { API_URL } from "$env/static/private";
import type { Algorithm } from "$lib/types/algorithm";
import { customFetch } from "$lib/api/client";
import {
	authFlowLimiter,
	fiveHundredQuerySize,
	requireAuth,
	requirePermission,
	tenMbBodySize,
	useMiddlewares
} from "$lib/server/middlewares";

const createAlgorithm: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError(
			"MISSING_ADMIN_COOKIE",
			"Falta a senha das rotas admin.",
			ADMIN_ALGORITHMS_ERRORS
		);
		return json(normalizedError, { status: 401 });
	}

	const body = await event.request.json().catch(() => null);
	const result = algorithmSchema.safeParse(body);

	if (!result.success) {
		const errorCode = result.error.issues[0].message;
		const normalizedError = normalizeApiError(
			errorCode,
			"Formato de dados inválido.",
			ADMIN_ALGORITHMS_ERRORS
		);
		return json(normalizedError, { status: 400 });
	}

	const clientIp = event.getClientAddress();

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(
		event.fetch,
		`${API_URL}/api/admin/algorithms`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": clientIp,
				"X-Admin-Secret": adminSecret,
				Authorization: `Bearer ${event.locals.accessToken}`
			},
			body: JSON.stringify(result.data)
		},
		ADMIN_ALGORITHMS_ERRORS
	);

	if (apiError) {
		return json(apiError, { status });
	}

	return json({ algorithm: data?.data });
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	tenMbBodySize,
	authFlowLimiter,
	requireAuth,
	requirePermission("create:algorithms")
)(createAlgorithm);
