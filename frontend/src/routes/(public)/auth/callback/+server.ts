import { redirect, type RequestHandler } from "@sveltejs/kit";
import { standardApiLimiter, twoThousandUrlSize, useMiddlewares } from "$lib/server/middlewares";

const callback: RequestHandler = async ({ url }) => {
	const accessToken = url.searchParams.get("access_token") === "true";
	const preToken = url.searchParams.get("pre_auth_token") === "true";
	const error = url.searchParams.get("error");

	if (error) {
		redirect(303, `/auth/login?error=${encodeURIComponent(error)}`);
	}

	if (preToken) {
		redirect(303, "/auth/verify-2fa");
	}

	if (accessToken) {
		redirect(303, "/");
	}

	redirect(303, "/auth/login?error=AUTH_UNEXPECTED_ERROR");
};

export const GET = useMiddlewares(twoThousandUrlSize, standardApiLimiter)(callback);
