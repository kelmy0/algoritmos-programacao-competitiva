import { JWT_ACCESS_PUBLIC_KEY } from "$env/static/private";
import type { AccessJwtPayload } from "$lib/types/jwt";
import { validateJWT } from "$lib/utils/jwt";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
	let expiresAt: number | null = null;

	if (locals.accessToken) {
		const { claims, valid } = await validateJWT<AccessJwtPayload>(
			JWT_ACCESS_PUBLIC_KEY,
			locals.accessToken
		);
		if (!valid || !claims?.exp) {
			expiresAt = null;
		} else {
			expiresAt = claims.exp * 1000;
		}
	}

	return {
		user: locals.user,
		accessToken: Boolean(locals.accessToken),
		expiresAt
	};
};
