import {
	authFlowLimiter,
	fiveHundredQuerySize,
	requireAuth,
	requirePermission,
	useMiddlewares
} from "$lib/server/middlewares";
import { json, type RequestHandler } from "@sveltejs/kit";

const checkAdmin: RequestHandler = async ({ cookies }) => {
	const adminSecret = cookies.get("admin_secret");

	if (adminSecret) {
		return json({ authenticated: true });
	}

	return json({ authenticated: false }, { status: 401 });
};

export const GET = useMiddlewares(
	fiveHundredQuerySize,
	authFlowLimiter,
	requireAuth,
	requirePermission()
)(checkAdmin);
