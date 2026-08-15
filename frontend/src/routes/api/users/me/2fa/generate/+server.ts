import { API_URL } from "$env/static/private";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { authFlowLimiter, fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { requireAuth, useMiddlewares } from "$lib/server/middlewares";
import { type TwoFactorGenerateResponse } from "$lib/types/users/me/two_factor";
import { normalizeApiError } from "$lib/utils/errors";
import type { RequestHandler } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import QRCode from "qrcode";
import { generate2FASchema } from "$lib/schemas/me";

const generate2FA: RequestHandler = async (event) => {
	const body = await event.request.json().catch(() => null);
	const result = generate2FASchema.safeParse(body);

	if (!result.success) {
		return json(normalizeApiError("INVALID_REQUEST_BODY"), { status: 400 });
	}

	const { data, error, status } = await customFetch<TwoFactorGenerateResponse>(
		event.fetch,
		`${API_URL}/api/users/me/2fa/generate`,
		{
			method: "POST",
			headers: {
				Authorization: `Bearer ${event.locals.accessToken}`
			},
			body: JSON.stringify(result.data)
		},
		TWO_FACTOR_ERRORS
	);

	if (error) {
		return json(error, { status });
	}

	if (!data?.secret || !data.qrCode) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR"), { status: 500 });
	}

	const qrCodeBase64 = await QRCode.toDataURL(data.qrCode, {
		margin: 2,
		width: 200,
		color: {
			dark: "#000000",
			light: "#ffffff"
		}
	});

	const response: TwoFactorGenerateResponse = {
		qrCode: qrCodeBase64,
		secret: data.secret
	};

	return json(response, { status: 200 });
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	requireAuth,
	authFlowLimiter
)(generate2FA);
