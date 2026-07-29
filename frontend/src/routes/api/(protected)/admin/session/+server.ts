import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ cookies }) => {
	const adminSecret = cookies.get("admin_secret");

	if (adminSecret) {
		return json({ authenticated: true });
	}

	return json({ authenticated: false }, { status: 401 });
};
