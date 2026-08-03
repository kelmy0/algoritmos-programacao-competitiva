import { goto } from "$app/navigation";
import { page } from "$app/state";
import { customFetch } from "$lib/api/client";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";
import { isValidEmail } from "../sign-up/sign_up.svelte";

export interface LoginServerResponse {
	access_token: boolean;
	requires_2fa: boolean;
	pre_auth_token?: string;
}

export const AUTH_ERRORS: Record<string, string> = {
	AUTH_INVALID_EMAIL_PASSWORD: "E-mail ou senha incorretos. Verifique seus dados.",
	USER_ALREADY_EXISTS:
		"Este e-mail já está cadastrado. Tente entrar por outro método ou use um email diferente.",
	SOCIAL_ACCOUNT_ALREADY_LINKED:
		"Este email já esta ligado a outra conta. Tente entrar por outro método ou use um email diferente. "
};

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

		if (data.requires_2fa) {
			await goto(`/auth/verify-2fa?token=${data.pre_auth_token}`);
			return;
		}

		if (data.access_token) {
			const redirectTo = page.url.searchParams.get("redirectTo");
			const isSafeRedirect =
				redirectTo && redirectTo.startsWith("/") && !redirectTo.startsWith("//");
			const targetUrl = isSafeRedirect ? redirectTo : "/";

			await goto(targetUrl, { invalidateAll: true });
		}
	}
}
