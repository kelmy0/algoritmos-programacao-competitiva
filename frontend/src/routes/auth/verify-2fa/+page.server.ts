import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export async function load({ locals }: Parameters<PageServerLoad>[0]) {
	if (locals.user) {
		redirect(303, '/');
	}

	return {};
}
