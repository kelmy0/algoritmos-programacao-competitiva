import { redirect, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { authFlowLimiter, fiveHundredQuerySize, useMiddlewares } from "$lib/server/middlewares";

const redirectSocial: RequestHandler = async (event) => {
	const { provider } = event.params;

	if (provider != "google" && provider != "github") {
		redirect(303, "/auth/login");
	}

	redirect(303, `${API_URL}/api/auth/${provider}`);
};

export const GET = useMiddlewares(fiveHundredQuerySize, authFlowLimiter)(redirectSocial);
