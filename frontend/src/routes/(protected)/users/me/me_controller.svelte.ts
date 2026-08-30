import { goto, invalidateAll } from "$app/navigation";
import { customFetch } from "$lib/api/client";
import { BaseController } from "$lib/controllers/base_controller.svelte";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import { type TwoFactorGenerateResponse } from "$lib/types/users/me/two_factor";
import { normalizeApiError } from "$lib/utils/errors";

export class MeController extends BaseController {
	#is2FAEnabled = $state(false);
	#hasPassword = $state(true);
	#twoFactorSecret = $state("");
	#qrCodeUrl = $state("");

	#password = $state("");
	#confirmPassword = $state("");
	#code = $state("");

	constructor(is2FAEnabled: boolean = false, hasPassword: boolean = true) {
		super();
		this.#is2FAEnabled = is2FAEnabled;
		this.#hasPassword = hasPassword;
	}

	get is2FAEnabled() {
		return this.#is2FAEnabled;
	}

	get twoFactorSecret() {
		return this.#twoFactorSecret;
	}

	get qrCodeUrl() {
		return this.#qrCodeUrl;
	}

	get password() {
		return this.#password;
	}

	set password(val: string) {
		this.#password = val;
		this.clearApiError();
	}

	get confirmPassword() {
		return this.#confirmPassword;
	}

	set confirmPassword(value: string) {
		this.#confirmPassword = value;
		this.clearApiError();
	}

	get code() {
		return this.#code;
	}

	set code(value: string) {
		this.#code = value.replace(/\D/g, "");
		this.clearApiError();
	}

	get isPasswordValid() {
		return this.#password.length >= 8 || (!this.#hasPassword && this.#password.length === 0);
	}

	get isCodeValid() {
		return this.#code.length === 6;
	}

	async generate2FA(): Promise<boolean> {
		if (this.#is2FAEnabled || !this.isPasswordValid || this._isLoading) return false;
		this._isLoading = true;

		const { data, error } = await customFetch<TwoFactorGenerateResponse>(
			window.fetch,
			"/api/users/me/2fa/generate",
			{ method: "POST", body: JSON.stringify({ password: this.#password }) },
			TWO_FACTOR_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		this._apiError = null;
		this.#twoFactorSecret = data.secret;
		this.#password = "";
		this.#qrCodeUrl = data.qrCode;
		return true;
	}

	async save2FA(): Promise<boolean> {
		if (this.#is2FAEnabled || !this.isCodeValid || this._isLoading) return false;
		this._isLoading = true;

		const { error } = await customFetch<null>(
			window.fetch,
			"/api/users/me/2fa/enable",
			{ method: "POST", body: JSON.stringify({ code: this.#code }) },
			TWO_FACTOR_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			return false;
		}

		this._apiError = null;
		this.#twoFactorSecret = "";
		this.#qrCodeUrl = "";
		await invalidateAll();
		return true;
	}

	async disable2FA(): Promise<boolean> {
		if (!this.#is2FAEnabled || !this.isPasswordValid || this._isLoading) return false;

		this._isLoading = true;

		const { error } = await customFetch<null>(
			window.fetch,
			"/api/users/me/2fa/disable",
			{ method: "POST", body: JSON.stringify({ password: this.#password }) },
			TWO_FACTOR_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			return false;
		}

		this._apiError = null;
		this.#password = "";
		await goto("/auth/login", { invalidateAll: true });
		return true;
	}
}
