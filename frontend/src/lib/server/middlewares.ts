import { normalizeApiError } from "$lib/utils/errors";
import { error, type RequestEvent, type RequestHandler } from "@sveltejs/kit";

export type Middleware = (event: RequestEvent) => Promise<Response | void> | Response | void;

export function useMiddlewares(...middlewares: Middleware[]) {
	return (handler: RequestHandler): RequestHandler => {
		return async (event) => {
			for (const middleware of middlewares) {
				const result = await middleware(event);

				if (result instanceof Response) {
					return result;
				}
			}

			return handler(event);
		};
	};
}

export const requireAuth: Middleware = (event) => {
	if (!event.locals.user) {
		error(401, normalizeApiError("UNAUTHORIZED", "Não autenticado."));
	}
};

export function requirePermission(permission?: string): Middleware {
	return (event) => {
		const user = event.locals.user;

		if (!user?.is_employee) {
			error(404, normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada."));
		}

		if (permission && !user.permissions?.includes(permission)) {
			error(404, normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada."));
		}
	};
}
