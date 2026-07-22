import { PUBLIC_API_URL } from '$env/static/public';
import type { ApiError } from '$lib/types/api';
import { json, type RequestHandler } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

interface JwtPayload {
	sub: string;
	username: string;
	email: string;
	permissions: string[];
	isEmployee: boolean;
}

export async function POST({ fetch, request, cookies }: Parameters<RequestHandler>[0]) {
	const cookieHeader = request.headers.get('cookie') || '';

	try {
		const refreshRes = await fetch(`${PUBLIC_API_URL}/api/auth/refresh`, {
			method: 'POST',
			headers: { cookie: cookieHeader }
		});

		if (!refreshRes.ok) {
			cookies.delete('access_token', { path: '/' });

			const errorData: ApiError = await refreshRes.json().catch(() => ({
				code: 'AUTH_SESSION_EXPIRED',
				message: 'A sessão expirou.'
			}));

			return json(errorData, { status: refreshRes.status });
		}

		const { access_token } = await refreshRes.json();
		const decoded = jwtDecode<JwtPayload>(access_token);
		cookies.set('access_token', access_token, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 60 * 15
		});

		return json({
			accessToken: access_token,
			user: {
				id: decoded.sub,
				username: decoded.username,
				email: decoded.email,
				permissions: decoded.permissions || [],
				is_employee: decoded.isEmployee
			}
		});
	} catch (error) {
		const internalError: ApiError = {
			code: 'AUTH_UNEXPECTED_ERROR',
			message: 'Ocorreu um erro inesperado ao renovar a sessão.'
		};

		return json(internalError, { status: 500 });
	}
}
