import { goto } from "$app/navigation";
import { page } from "$app/state";
import { customFetch } from "$lib/api/client";
import { BaseAuthController } from "$lib/controllers/base_auth_controller.svelte";
import { AUTH_ERRORS } from "$lib/errors/auth/auth_errors";
import type { ApiError } from "$lib/types/api";
import type { LoginServerResponse } from "$lib/types/auth/login";
import { normalizeApiError } from "$lib/utils/errors";
import { isValidEmail } from "$lib/utils/sanitize";

export class LoginController extends BaseAuthController {
	#email = $state("");
	#password = $state("");

	constructor(initialError: ApiError | null = null) {
		super(initialError);
	}

	get email() {
		return this.#email;
	}

	set email(value: string) {
		this.#email = value;
		this.clearApiError();
	}

	get password() {
		return this.#password;
	}

	set password(value: string) {
		this.#password = value;
		this.clearApiError();
	}

	get isEmailValid() {
		return isValidEmail(this.#email);
	}

	get isPasswordValid() {
		return this.#password.length >= 8;
	}

	async login(): Promise<boolean> {
		if (!this.isEmailValid || !this.isPasswordValid || this._isLoading || !this.validateTurnstile())
			return false;

		this._isLoading = true;
		this._apiError = null;

		const { data, error } = await customFetch<LoginServerResponse>(
			window.fetch,
			"/api/auth/login",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this._turnstileToken
				},
				body: JSON.stringify({
					email: this.#email,
					password: this.#password
				})
			},
			AUTH_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			this.resetTurnstile();
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		if (data.requires2FA) {
			await goto(`/auth/verify-2fa`);
			return true;
		}

		if (data.accessToken) {
			const redirectTo = page.url.searchParams.get("redirectTo");

			let targetUrl = "/";

			if (redirectTo) {
				try {
					const parsedUrl = new URL(redirectTo, page.url.origin);

					if (parsedUrl.origin === page.url.origin) {
						targetUrl = parsedUrl.pathname + parsedUrl.search + parsedUrl.hash;
					}
				} catch {
					targetUrl = "/";
				}
			}

			await goto(targetUrl, { invalidateAll: true });
		}

		return true;
	}
}
