import { error } from '@sveltejs/kit';
import { normalizeApiError } from '$lib/utils/errors';

export interface UserPermissionsContext {
    permissions?: string[];
    is_employee?: boolean;
}

export function checkAdminAccess(user: UserPermissionsContext | null | undefined, requiredPermission?: string) {
    if (!user?.is_employee) {
        error(404, normalizeApiError('PAGE_NOT_FOUND', 'Página não encontrada.'));
    }

    if (requiredPermission && !user.permissions?.includes(requiredPermission)) {
        error(404, normalizeApiError('PAGE_NOT_FOUND', 'Página não encontrada.'));
    }
}