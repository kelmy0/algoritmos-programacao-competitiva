import { goto, invalidateAll } from "$app/navigation";
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

	isChangePasswordModalOpen = $state(false);

	password = $state("");
	confirmPassword = $state("");
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

	// MODALS
	open2FAModal() {
		this.is2FAModalOpen = true;
	}

	close2FAModal() {
		this.is2FAModalOpen = false;
	}

	openChangePasswordModal() {
		this.isChangePasswordModalOpen = true;
	}

	closeChangePasswordModal() {
		this.isChangePasswordModalOpen = false;
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
			"/api/users/me/2fa/generate",
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
		this.password = "";
		this.touched.password = false;
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
			"/api/users/me/2fa/enable",
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
		this.touched.code = false;
		this.close2FAModal();
		await invalidateAll();
	}

	async disable2FA(e: SubmitEvent) {
		e.preventDefault();
		this.touched.password = true;

		if (!this.is2FAEnabled || !this.isPasswordValid) return;

		this.isLoading = true;

		const bodyrequest = { password: this.password };

		const { error } = await customFetch<null>(
			window.fetch,
			"/api/users/me/2fa/disable",
			{ method: "POST", body: JSON.stringify(bodyrequest) },
			TWO_FACTOR_ERRORS
		);

		this.isLoading = false;

		if (error) {
			this.apiError = error;
			return;
		}

		this.apiError = null;
		this.password = "";
		this.touched.password = false;
		this.close2FAModal();
		await goto("/auth/login", { invalidateAll: true });
	}
}
