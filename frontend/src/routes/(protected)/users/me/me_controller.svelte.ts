import { customFetch } from "$lib/api/client";
import { TWO_FACTOR_ERRORS } from "$lib/errors/users/me/two_factor";
import type { ApiError } from "$lib/types/api";
import { type TwoFactorGenerateResponse } from "$lib/types/users/me/two_factor";
import { normalizeApiError } from "$lib/utils/errors";

export class MeController {
	apiError = $state<ApiError | null>(null);
	isLoading = $state(false);

	is2FAModalOpen = $state(false);
	is2FAEnabled = $state(false);
	twoFactorSecret = $state("");
	qrCodeUrl = $state("");

	password = $state("");
	showPassword = $state(false);

	code = $state("");

	touched = $state({
		password: false,
		code: false
	});

	get isPasswordValid() {
		return this.password.length >= 8;
	}

	get isCodeValid() {
		return this.code.length === 6;
	}

	open2FAModal() {
		this.is2FAModalOpen = true;
	}

	close2FAModal() {
		this.is2FAModalOpen = false;
	}

	onInput() {
		if (this.apiError) {
			this.apiError = null;
		}
	}

	on2FAInput(e: Event) {
		this.onInput();
		const input = e.target as HTMLInputElement;

		this.code = input.value.replace(/\D/g, "");
	}

	togglePassword() {
		this.showPassword = !this.showPassword;
	}

	async generate2FA(e: SubmitEvent) {
		e.preventDefault();
		this.touched.password = true;

		if (this.is2FAEnabled || !this.isPasswordValid) return;
		this.isLoading = true;

		const bodyrequest = { password: this.password };

		const { data, error } = await customFetch<TwoFactorGenerateResponse>(
			window.fetch,
			"/api/users/me/generate-2FA",
			{ method: "POST", body: JSON.stringify(bodyrequest) },
			TWO_FACTOR_ERRORS
		);

		this.isLoading = false;

		if (error) {
			this.apiError = error;
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return;
		}

		this.apiError = null;
		this.twoFactorSecret = data.secret;
		this.qrCodeUrl = data.qrCode;
	}

	async save2FA(e: SubmitEvent) {
		e.preventDefault();
		this.touched.code = true;

		if (this.is2FAEnabled || !this.isCodeValid) return;
		this.isLoading = true;

		const bodyrequest = { code: this.code };

		const { error } = await customFetch<null>(
			window.fetch,
			"/api/users/me/enable-2FA",
			{ method: "POST", body: JSON.stringify(bodyrequest) },
			TWO_FACTOR_ERRORS
		);

		this.isLoading = false;

		if (error) {
			this.apiError = error;
			return;
		}

		this.apiError = null;
		this.twoFactorSecret = "";
		this.qrCodeUrl = "";
	}
}
