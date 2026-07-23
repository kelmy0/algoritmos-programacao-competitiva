import { json, type RequestHandler } from '@sveltejs/kit';
import { jwtDecode, type JwtPayload as BaseJwtPayload } from 'jwt-decode';
import { PUBLIC_API_URL } from '$env/static/public';
import { normalizeApiError } from '$lib/utils/errors';
import { deleteAuthCookie, setAuthCookie } from '$lib/server/cookies';

interface JwtPayload extends BaseJwtPayload {
	username: string;
	email: string;
	permissions?: string[];
	isEmployee: boolean;
}

export const POST: RequestHandler = async ({ fetch, request, cookies }) => {
	const cookieHeader = request.headers.get('cookie') || '';

	try {
		const refreshRes = await fetch(`${PUBLIC_API_URL}/api/auth/refresh`, {
			method: 'POST',
			headers: { cookie: cookieHeader }
		});

		if (!refreshRes.ok) {
			cookies.delete('access_token', { path: '/' });

			const rawError = await refreshRes.json().catch(() => null);
			const normalizedError = normalizeApiError(
				rawError,
				'Sua sessão expirou. Faça login novamente.'
			);

			return json(normalizedError, { status: refreshRes.status });
		}

		const { access_token } = await refreshRes.json();
		const decoded = jwtDecode<JwtPayload>(access_token);

		setAuthCookie(cookies, 'access_token', access_token, 15);

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
	} catch (err) {
		deleteAuthCookie(cookies, 'access_token');

		const normalizedError = normalizeApiError(
			err,
			'Ocorreu um erro inesperado ao renovar a sessão.'
		);

		return json(normalizedError, { status: 500 });
	}
};
