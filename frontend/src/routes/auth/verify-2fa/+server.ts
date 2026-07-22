import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';

export async function POST({ fetch, request, cookies }: Parameters<RequestHandler>[0]) {
	try {
		const body = await request.json();

		const apiRes = await fetch(`${PUBLIC_API_URL}/api/auth/verify-2fa`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});

		if (!apiRes.ok) {
			const errorData = await apiRes.json();
			return json(errorData, { status: apiRes.status });
		}

		const setCookieHeader = apiRes.headers.get('set-cookie');
		const data = await apiRes.json();

		if (data.access_token) {
			cookies.set('access_token', data.access_token, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				maxAge: 60 * 15
			});
		}

		const response = json(data);
		if (setCookieHeader) {
			response.headers.append('set-cookie', setCookieHeader);
		}

		return response;
	} catch (err) {
		return json({ code: 'NETWORK_ERROR', message: 'Erro na rede' }, { status: 500 });
	}
}
