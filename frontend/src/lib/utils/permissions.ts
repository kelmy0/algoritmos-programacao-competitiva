import { error } from "@sveltejs/kit";
import { normalizeApiError } from "$lib/utils/errors";

export interface UserPermissionsContext {
	permissions?: string[];
	isEmployee?: boolean;
}

export function checkAdminAccess(
	user: UserPermissionsContext | null | undefined,
	requiredPermission?: string
) {
	if (!user?.isEmployee) {
		error(404, normalizeApiError("PAGE_NOT_FOUND"));
	}

	if (requiredPermission && !user.permissions?.includes(requiredPermission)) {
		error(404, normalizeApiError("PAGE_NOT_FOUND"));
	}
}
