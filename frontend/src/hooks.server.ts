import type { Handle } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { jwtDecode } from 'jwt-decode';

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
	let isTokenValid = false;

	if (accessToken) {
		try {
			const decoded = jwtDecode<JwtPayload>(accessToken);

			if (decoded.exp && decoded.exp * 1000 > Date.now()) {
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

	const cookieHeader = event.request.headers.get('cookie');

	if (cookieHeader) {
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

				event.cookies.set('access_token', access_token, {
					path: '/',
					httpOnly: true,
					sameSite: 'lax',
					maxAge: 60 * 15
				});
			}
		} catch (err) {
			console.error('Error renewing session:', err);
		}
	}

	return await resolve(event);
}
