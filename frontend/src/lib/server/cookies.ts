import { ENV } from '$env/static/private';
import type { Cookies } from '@sveltejs/kit';

const DEFAULT_COOKIE_OPTIONS = {
	path: '/',
	httpOnly: true,
	sameSite: 'lax' as const,
	secure: ENV !== 'development'
};

export function setAuthCookie(cookies: Cookies, name: string, value: string, minutes: number) {
	cookies.set(name, value, {
		...DEFAULT_COOKIE_OPTIONS,
		maxAge: 60 * minutes
	});
}

export function deleteAuthCookie(cookies: Cookies, name: string) {
	cookies.delete(name, {
		path: DEFAULT_COOKIE_OPTIONS.path
	});
}
