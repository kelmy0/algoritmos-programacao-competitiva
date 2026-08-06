import { goto } from "$app/navigation";
import { page } from "$app/state";
import { customFetch } from "$lib/api/client";
import { AUTH_ERRORS } from "$lib/errors/auth/auth_errors";
import type { ApiError } from "$lib/types/api";
import type { LoginServerResponse } from "$lib/types/auth/login";
import { normalizeApiError } from "$lib/utils/errors";
import { isValidEmail } from "$lib/utils/sanitize";

export class LoginController {
	email = $state("");
	password = $state("");
	turnstileToken = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	showPassword = $state(false);

	turnstileComponent: { reset: () => void } | null = null;

	constructor(initialError: ApiError | null = null) {
		this.apiError = initialError;
	}

	touched = $state({
		email: false,
		password: false
	});

	get isEmailValid() {
		return isValidEmail(this.email);
	}

	get isPasswordValid() {
		return this.password.length >= 8;
	}

	onInput() {
		if (this.apiError) {
			this.apiError = null;
		}
	}

	togglePassword() {
		this.showPassword = !this.showPassword;
	}

	onTurnstileSuccess(token: string) {
		this.turnstileToken = token;
	}

	onTurnstileExpire() {
		this.turnstileToken = "";
	}

	async login(event: SubmitEvent) {
		event.preventDefault();
		this.touched.email = true;
		this.touched.password = true;

		if (!this.isEmailValid || !this.isPasswordValid || this.isLoading) return;

		if (!this.turnstileToken) {
			this.apiError = normalizeApiError("CAPTCHA_REQUIRED");
			return;
		}

		this.isLoading = true;
		this.apiError = null;

		const { data, error } = await customFetch<LoginServerResponse>(
			window.fetch,
			"/api/auth/login",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this.turnstileToken
				},
				body: JSON.stringify({
					email: this.email,
					password: this.password
				})
			},
			AUTH_ERRORS
		);

		this.isLoading = false;
		if (error) {
			this.apiError = error;
			this.turnstileToken = "";
			this.turnstileComponent?.reset();
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Falha ao processar resposta do servidor.",
				AUTH_ERRORS
			);
			return;
		}

		if (data.requires2FA) {
			await goto(`/auth/verify-2fa?token=${data.preAuthToken}`);
			return;
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
	}
}
