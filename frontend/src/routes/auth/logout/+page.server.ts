import { redirect, type Actions, type RequestEvent } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';

async function handleLogout({ fetch, locals, cookies }: RequestEvent) {
	if (locals.accessToken) {
		try {
			await fetch(`${PUBLIC_API_URL}/api/auth/logout`, {
				method: 'POST',
				headers: { Authorization: `Bearer ${locals.accessToken}` }
			});
		} catch (err) {
			console.error('Erro ao deslogar na API:', err);
		}
	}

	cookies.delete('access_token', { path: '/' });
	cookies.delete('refresh_token', { path: '/' });
	locals.user = null;
	locals.accessToken = null;

	redirect(303, '/auth/login');
}

export const actions: Actions = {
	default: handleLogout
};
