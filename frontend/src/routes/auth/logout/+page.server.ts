import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export async function load({ fetch, locals, cookies }: Parameters<PageServerLoad>[0]) {
	if (locals.accessToken) {
		try {
			await fetch(`${PUBLIC_API_URL}/api/auth/logout`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${locals.accessToken}`
				}
			});
		} catch (err) {
			console.error('Erro ao comunicar logout com o backend:', err);
		}
	}

	cookies.delete('access_token', { path: '/' });
	cookies.delete('refresh_token', { path: '/' });
	locals.user = null;
	locals.accessToken = null;
	throw redirect(303, '/auth/login');
}
