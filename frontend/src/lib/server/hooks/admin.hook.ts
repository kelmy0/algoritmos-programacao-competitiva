import type { Handle } from "@sveltejs/kit";
import { setIdleCookie } from "$lib/server/cookies";

export const handleAdminSession: Handle = async ({ event, resolve }) => {
	const adminSecret = event.cookies.get("admin_secret");

	if (adminSecret) {
		setIdleCookie(event.cookies, "admin_secret", adminSecret, 15);
	}

	return await resolve(event);
};
