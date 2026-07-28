import { jwtDecode } from "jwt-decode";
import type { LayoutServerLoad } from "./$types";

interface JwtPayload {
	exp?: number;
}

export const load: LayoutServerLoad = async ({ locals }) => {
	let expiresAt: number | null = null;

	if (locals.accessToken) {
		try {
			const decoded = jwtDecode<JwtPayload>(locals.accessToken);
			if (decoded.exp) {
				expiresAt = decoded.exp * 1000;
			}
		} catch {
			expiresAt = null;
		}
	}

	return {
		user: locals.user,
		accessToken: Boolean(locals.accessToken),
		expiresAt
	};
};
