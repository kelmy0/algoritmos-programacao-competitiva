import { error } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { normalizeApiError } from '$lib/utils/errors';

export const ssr = false;

const ROUTE_PERMISSIONS: Record<string, string> = {
	'/admin/algorithms/new': 'create:algorithms',
	'/admin/algorithms/edit': 'update:algorithms',
	'/admin/algorithms/trash': 'delete:algorithms'
};

export const load: LayoutLoad = async ({ url, parent }) => {
	const { user } = await parent();
	const requiredPermission = Object.entries(ROUTE_PERMISSIONS).find(([path]) =>
		url.pathname.startsWith(path)
	)?.[1];

	if (requiredPermission && !user?.permissions?.includes(requiredPermission)) {
		error(404, normalizeApiError('PAGE_NOT_FOUND', 'Página não encontrada.'));
	}

	return { user };
};
