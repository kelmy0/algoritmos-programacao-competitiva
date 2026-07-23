import { json, type RequestHandler } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { normalizeApiError } from '$lib/utils/errors';
import { setAuthCookie } from '$lib/server/cookies';
import { TWO_FACTOR_ERRORS } from '../../../auth/verify-2fa/two_factor_verify.svelte';

export const POST: RequestHandler = async ({ fetch, request, cookies }) => {
	try {
		const body = await request.json();

		const apiRes = await fetch(`${PUBLIC_API_URL}/api/auth/verify-2fa`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});

		if (!apiRes.ok) {
			const rawError = await apiRes.json().catch(() => null);

			const normalizedError = normalizeApiError(
				rawError,
				'Código de autenticação inválido ou expirado.',
				TWO_FACTOR_ERRORS
			);

			return json(normalizedError, { status: apiRes.status });
		}

		const data = await apiRes.json();
		const setCookieHeader = apiRes.headers.get('set-cookie');

		if (data.access_token) {
			setAuthCookie(cookies, 'access_token', data.access_token, 15);
		}

		const response = json(data);

		if (setCookieHeader) {
			response.headers.append('set-cookie', setCookieHeader);
		}

		return response;
	} catch (err) {
		const normalizedError = normalizeApiError(
			err,
			'Não foi possível conectar ao servidor. Verifique sua conexão.',
			TWO_FACTOR_ERRORS
		);

		return json(normalizedError, { status: 500 });
	}
};
