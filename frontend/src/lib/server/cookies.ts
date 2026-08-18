import { ENV } from "$env/static/private";
import type { Cookies } from "@sveltejs/kit";

const DEFAULT_COOKIE_OPTIONS = {
	path: "/",
	httpOnly: true,
	secure: ENV !== "development",
	partitioned: ENV !== "development"
};

export function setAuthCookie(cookies: Cookies, name: string, value: string, minutes: number) {
	cookies.set(name, value, {
		...DEFAULT_COOKIE_OPTIONS,
		sameSite: "lax" as const,
		maxAge: 60 * minutes
	});
}

export function setSessionCookie(cookies: Cookies, name: string, value: string) {
	cookies.set(name, value, {
		...DEFAULT_COOKIE_OPTIONS,
		sameSite: "strict" as const
	});
}

export function setIdleCookie(cookies: Cookies, name: string, value: string, minutes: number = 15) {
	cookies.set(name, value, {
		...DEFAULT_COOKIE_OPTIONS,
		sameSite: "strict",
		maxAge: 60 * minutes
	});
}

export function deleteCookie(cookies: Cookies, name: string) {
	cookies.delete(name, {
		path: DEFAULT_COOKIE_OPTIONS.path
	});
}

export function syncServerCookie(cookies: Cookies, name: string, value: string) {
	cookies.set(name, value, {
		path: DEFAULT_COOKIE_OPTIONS.path
	});
}
