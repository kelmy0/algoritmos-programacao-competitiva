import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, url }) => {
	if (!locals.user) {
		const redirectTo = url.pathname + url.search;
		redirect(303, `/auth/login?redirectTo=${encodeURIComponent(redirectTo)}`);
	}

	return {};
};
