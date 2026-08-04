import { redirect, type RequestHandler } from "@sveltejs/kit";
import { API_URL } from "$env/static/private";
import { deleteAuthCookie } from "$lib/server/cookies";
import { fiveHundredQuerySize, hundredKbBodySize } from "$lib/server/middlewares";
import { standardApiLimiter, useMiddlewares } from "$lib/server/middlewares";

const logout: RequestHandler = async (event) => {
	if (event.locals.accessToken) {
		try {
			await event.fetch(`${API_URL}/api/auth/logout`, {
				method: "POST",
				headers: {
					Authorization: `Bearer ${event.locals.accessToken}`
				}
			});
		} catch (err) {
			console.error("Error notifying logout to API:", err);
		}
	}

	deleteAuthCookie(event.cookies, "access_token");
	deleteAuthCookie(event.cookies, "refresh_token");
	deleteAuthCookie(event.cookies, "admin_secret");

	event.locals.user = null;
	event.locals.accessToken = null;

	redirect(303, "/auth/login");
};

export const POST = useMiddlewares(
	fiveHundredQuerySize,
	hundredKbBodySize,
	standardApiLimiter
)(logout);
