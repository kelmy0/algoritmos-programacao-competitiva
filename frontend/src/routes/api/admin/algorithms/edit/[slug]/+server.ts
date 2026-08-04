import { normalizeApiError } from "$lib/utils/errors";
import { json, type RequestHandler } from "@sveltejs/kit";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import type { Algorithm } from "$lib/types/algorithm";
import { ADMIN_ALGORITHMS_ERRORS } from "$lib/errors/admin/algorithms";
import { algorithmSchema } from "$lib/schemas/algorithm";
import { authFlowLimiter, thousandQuerySize } from "$lib/server/middlewares";
import { requirePermission, standardApiLimiter } from "$lib/server/middlewares";
import { tenMbBodySize, useMiddlewares, requireAuth } from "$lib/server/middlewares";
import { ADMIN_PASSWORD_ERRORS } from "$lib/errors/admin/password";

const myAlgorithm: RequestHandler = async (event) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (!adminSecret) {
		const normalizedError = normalizeApiError("MISSING_ADMIN_COOKIE", "", ADMIN_PASSWORD_ERRORS);
		return json(normalizedError, { status: 401 });
	}

	const { slug } = event.params;
	const clientIp = event.getClientAddress();

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(
		event.fetch,
		`${API_URL}/api/admin/algorithms/${slug}`,
		{
			method: "GET",
			headers: {
				"X-Forwarded-For": clientIp,
				"X-Admin-Secret": adminSecret,
				Authorization: `Bearer ${event.locals.accessToken}`
			}
		}
	);

	if (apiError) {
		return json(apiError, { status });
	}

	if (!data?.data) {
		return json({});
	}

	return json(data);
};

const editAlgorithm: RequestHandler = async (event) => {
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

	const { slug } = event.params;
	const clientIp = event.getClientAddress();

	const {
		data,
		error: apiError,
		status
	} = await customFetch<{ data: Algorithm }>(
		event.fetch,
		`${API_URL}/api/admin/algorithms/${slug}`,
		{
			method: "PUT",
			headers: {
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

export const GET = useMiddlewares(
	thousandQuerySize,
	standardApiLimiter,
	requireAuth,
	requirePermission("create:algorithms")
)(myAlgorithm);

export const PUT = useMiddlewares(
	thousandQuerySize,
	tenMbBodySize,
	authFlowLimiter,
	requireAuth,
	requirePermission("create:algorithms")
)(editAlgorithm);
