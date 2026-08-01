import { redirect, type RequestHandler } from "@sveltejs/kit";
import { setAuthCookie } from "$lib/server/cookies";
import { standardApiLimiter, twoThousandUrlSize, useMiddlewares } from "$lib/server/middlewares";

const callback: RequestHandler = async ({ url, cookies }) => {
	const accessToken = url.searchParams.get("access_token");
	const preToken = url.searchParams.get("pre_token");
	const error = url.searchParams.get("error");

	if (error) {
		redirect(303, `/auth/login?error=${encodeURIComponent(error)}`);
	}

	if (preToken) {
		redirect(303, `/auth/verify-2fa?token=${encodeURIComponent(preToken)}`);
	}

	if (accessToken) {
		setAuthCookie(cookies, "access_token", accessToken, 15);
		redirect(303, "/");
	}

	redirect(303, "/auth/login?error=AUTH_UNEXPECTED_ERROR");
};

export const GET = useMiddlewares(twoThousandUrlSize, standardApiLimiter)(callback);
