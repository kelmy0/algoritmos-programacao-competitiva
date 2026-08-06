import { redirect, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { standardApiLimiter, useMiddlewares } from "$lib/server/middlewares";
import { clearAllAuthCookies } from "$lib/utils/cookies";

const logout: RequestHandler = async (event) => {
	const cookieHeader = event.request.headers.get("cookie") || "";

	if (event.locals.accessToken) {
		try {
			await event.fetch(`${API_URL}/api/auth/logout`, {
				method: "POST",
				headers: {
					Authorization: `Bearer ${event.locals.accessToken}`,
					cookie: cookieHeader
				}
			});
		} catch (err) {
			console.error("Error notifying logout to API:", err);
		}
	}

	clearAllAuthCookies(event.cookies);

	event.locals.user = null;
	event.locals.accessToken = null;

	redirect(303, "/auth/login");
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	standardApiLimiter
)(logout);
