import { redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { API_URL } from "$env/static/private";

export const GET: RequestHandler = async (event) => {
	const { provider } = event.params;

	if (provider != "google" && provider != "github") {
		redirect(303, "/auth/login");
	}

	redirect(303, `${API_URL}/api/auth/${provider}`);
};
