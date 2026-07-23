import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export async function load({ locals, url }: Parameters<PageServerLoad>[0]) {
	if (locals.user) {
		redirect(303, '/');
	}

	const token = url.searchParams.get('token');

	if (!token) {
		redirect(303, '/auth/login?error=MISSING_PRE_TOKEN');
	}

	return {};
}
