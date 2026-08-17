import { invalidateAll } from "$app/navigation";
import { customFetch } from "$lib/api/client";
import { BaseController } from "$lib/controllers/base_controller.svelte";
import { ADMIN_PASSWORD_ERRORS } from "$lib/errors/admin/password";
import { normalizeApiError } from "$lib/utils/errors";

export class AdminController extends BaseController {
	#password = $state("");

	get password() {
		return this.#password;
	}

	set password(value: string) {
		this.#password = value;
		this.clearApiError();
	}

	get isPasswordValid() {
		return this.#password.length >= 8;
	}

	async sendPassword(): Promise<boolean> {
		if (!this.isPasswordValid || this._isLoading) return false;

		this._isLoading = true;
		this._apiError = null;

		const { data, error } = await customFetch<{ correct: boolean }>(
			window.fetch,
			"/api/admin",
			{
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ password: this.#password })
			},
			ADMIN_PASSWORD_ERRORS
		);

		this._isLoading = false;
		this.#password = "";

		if (error) {
			this._apiError = error;
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		if (!data.correct) {
			this._apiError = normalizeApiError("INCORRECT_ADMIN_PASSWORD", "", ADMIN_PASSWORD_ERRORS);
			return false;
		}

		await invalidateAll();
		return true;
	}
}
