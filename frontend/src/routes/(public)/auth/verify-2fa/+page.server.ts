import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url, locals }) => {
	if (locals.user) {
		redirect(303, "/");
	}

	return {};
};
