import { goto } from "$app/navigation";
import type { ApiError } from "$lib/types/api";
import { page } from "$app/state";
import { normalizeApiError } from "$lib/utils/errors";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/auth/verify-2fa";
import type { TwoFactorServerResponse } from "$lib/types/auth/two-factor";

interface TwoFactorRequest {
	pre_auth_token: string;
	code: string;
}

export class TwoFactorController {
	token = "";
	code = $state("");
	turnstileToken = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);

	turnstileComponent: { reset: () => void } | null = null;

	touched = $state({
		code: false
	});

	get isCodeValid() {
		return this.code.length === 6;
	}

	onInput(event: Event) {
		const input = event.target as HTMLInputElement;

		this.code = input.value.replace(/\D/g, "");

		if (this.apiError) {
			this.apiError = null;
		}

		if (this.code.length === 6 && !this.isLoading) {
			const form = input.closest("form");
			if (form) {
				form.requestSubmit();
			}
		}
	}

	getToken() {
		const token = page.url.searchParams.get("token");

		if (!token) {
			goto("/auth/login?error=MISSING_PRE_TOKEN");
			return;
		}

		this.token = token;
	}

	onTurnstileSuccess(token: string) {
		this.turnstileToken = token;
	}

	onTurnstileExpire() {
		this.turnstileToken = "";
	}

	async sendCode(e: SubmitEvent) {
		e.preventDefault();

		this.touched.code = true;

		if (!this.isCodeValid || this.isLoading) {
			return;
		}

		if (!this.turnstileToken) {
			this.apiError = normalizeApiError("CAPTCHA_REQUIRED");
			return;
		}

		this.isLoading = true;
		this.apiError = null;

		const bodyRequest: TwoFactorRequest = {
			pre_auth_token: this.token,
			code: this.code
		};

		const { data, error } = await customFetch<TwoFactorServerResponse>(
			window.fetch,
			"/api/auth/verify-2fa",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this.turnstileToken
				},
				body: JSON.stringify(bodyRequest)
			},
			TWO_FACTOR_ERRORS
		);

		this.isLoading = false;
		if (error) {
			this.apiError = error;
			this.turnstileToken = "";
			this.turnstileComponent?.reset();
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return;
		}

		if (data.requires2FA) {
			await goto(`/auth/login?error=AUTH_UNEXPECTED_ERROR`);
			return;
		}

		await goto("/", { invalidateAll: true });
	}
}
