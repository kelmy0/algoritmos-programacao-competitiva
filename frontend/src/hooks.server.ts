import type { Handle, HandleServerError } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { jwtDecode } from 'jwt-decode';
import { normalizeApiError } from '$lib/utils/errors';
import { setAuthCookie } from '$lib/server/cookies';
import { deleteAuthCookie } from '$lib/server/cookies';

interface JwtPayload {
	sub: string;
	username: string;
	email: string;
	permissions: string[];
	isEmployee: boolean;
	exp?: number;
}

export async function handle({ event, resolve }: Parameters<Handle>[0]) {
	event.locals.user = null;
	event.locals.accessToken = null;

	const accessToken = event.cookies.get('access_token');
	const refreshToken = event.cookies.get('refresh_token');

	let isTokenValid = false;

	if (accessToken) {
		try {
			const decoded = jwtDecode<JwtPayload>(accessToken);
			const nowInSeconds = Math.floor(Date.now() / 1000);
			const BUFFER_SECONDS = 90;

			if (decoded.exp && decoded.exp - nowInSeconds > BUFFER_SECONDS) {
				isTokenValid = true;

				event.locals.user = {
					id: decoded.sub,
					username: decoded.username,
					email: decoded.email,
					permissions: decoded.permissions || [],
					is_employee: decoded.isEmployee
				};
				event.locals.accessToken = accessToken;
			}
		} catch {
			isTokenValid = false;
		}
	}

	if (isTokenValid && event.locals.user) {
		return await resolve(event);
	}

	if (refreshToken) {
		const cookieHeader = event.request.headers.get('cookie') || '';
		const clientIp = event.getClientAddress();

		try {
			const refreshRes = await event.fetch(`${PUBLIC_API_URL}/api/auth/refresh`, {
				method: 'POST',
				headers: {
					cookie: cookieHeader,
					'x-forwarded-for': clientIp,
					'x-real-ip': clientIp
				}
			});

			if (refreshRes.ok) {
				const { access_token } = await refreshRes.json();
				event.locals.accessToken = access_token;

				const decoded = jwtDecode<JwtPayload>(access_token);
				event.locals.user = {
					id: decoded.sub,
					username: decoded.username,
					email: decoded.email,
					permissions: decoded.permissions || [],
					is_employee: decoded.isEmployee
				};

				setAuthCookie(event.cookies, 'access_token', access_token, 15);
			} else {
				deleteAuthCookie(event.cookies, 'access_token');
				deleteAuthCookie(event.cookies, 'refresh_token');
			}
		} catch (err) {
			deleteAuthCookie(event.cookies, 'access_token');
			deleteAuthCookie(event.cookies, 'refresh_token');
		}
	}

	return await resolve(event);
}

export const handleError: HandleServerError = ({ error, event }) => {
	const apiError = normalizeApiError(error, 'Ocorreu um erro interno no servidor.');

	console.error(`[Server Error ${event.url.pathname}]:`, apiError);

	return {
		message: apiError.message,
		code: apiError.code
	};
};
