import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url, locals }) => {
	if (locals.user) {
		redirect(303, "/");
	}

	const token = url.searchParams.get("token");

	if (!token) {
		redirect(303, "/auth/login?error=MISSING_PRE_TOKEN");
	}

	return {};
};
