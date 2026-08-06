import { deleteCookie } from "$lib/server/cookies";
import type { Handle } from "@sveltejs/kit";

export function clearAllAuthCookies(cookies: Parameters<Handle>[0]["event"]["cookies"]) {
	deleteCookie(cookies, "access_token");
	deleteCookie(cookies, "refresh_token");
	deleteCookie(cookies, "admin_secret");
}
