import { goto } from "$app/navigation";
import { normalizeApiError } from "$lib/utils/errors";
import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/auth/verify-2fa";
import type { TwoFactorServerResponse } from "$lib/types/auth/two-factor";
import { BaseAuthController } from "$lib/controllers/base_auth_controller.svelte";

interface TwoFactorRequest {
	code: string;
}

export class TwoFactorController extends BaseAuthController {
	#code = $state("");

	get code() {
		return this.#code;
	}

	set code(value: string) {
		this.#code = value.replace(/\D/g, "");
		this.clearApiError();
	}

	get isCodeValid() {
		return this.#code.length === 6;
	}

	async sendCode(): Promise<boolean> {
		if (!this.isCodeValid || this._isLoading || !this.validateTurnstile()) return false;

		this._isLoading = true;
		this._apiError = null;

		const bodyRequest: TwoFactorRequest = {
			code: this.#code
		};

		const { data, error } = await customFetch<TwoFactorServerResponse>(
			window.fetch,
			"/api/auth/verify-2fa",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this._turnstileToken
				},
				body: JSON.stringify(bodyRequest)
			},
			TWO_FACTOR_ERRORS
		);

		this._isLoading = false;

		if (error) {
			if (error.code === "MISSING_COOKIE") {
				await goto("/auth/login?error=MISSING_COOKIE", { invalidateAll: true });
				return false;
			}

			this._apiError = error;
			this.resetTurnstile();
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		if (data.requires2FA) {
			await goto(`/auth/login?error=AUTH_UNEXPECTED_ERROR`);
			return false;
		}

		await goto("/", { invalidateAll: true });
		return true;
	}
}
