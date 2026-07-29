import type { LayoutLoad } from "./$types";
import { checkAdminAccess } from "$lib/utils/permissions";
import { customFetch } from "$lib/api/client";

export const ssr = false;

const ROUTE_PERMISSIONS: Record<string, string> = {
	"/admin/algorithms/new": "create:algorithms",
	"/admin/algorithms/edit": "create:algorithms",
	"/admin/algorithms/trash": "create:algorithms"
	//'/admin/algorithms/edit': 'update:algorithms',
	//"/admin/algorithms/trash": "delete:algorithms"
};

export const load: LayoutLoad = async ({ url, parent, fetch: svelteFetch }) => {
	const { user } = await parent();
	const requiredPermission = Object.entries(ROUTE_PERMISSIONS).find(([path]) =>
		url.pathname.startsWith(path)
	)?.[1];

	checkAdminAccess(user, requiredPermission);

	const { data } = await customFetch<{ authenticated: boolean }>(
		svelteFetch,
		"/api/admin/session",
		{ method: "GET" }
	);

	return {
		user,
		hasAdminSecret: Boolean(data?.authenticated)
	};
};
