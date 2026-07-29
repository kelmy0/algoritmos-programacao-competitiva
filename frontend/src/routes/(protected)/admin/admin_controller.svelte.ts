import { invalidateAll } from "$app/navigation";
import { customFetch } from "$lib/api/client";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";

export const ADMIN_PASSWORD_ERRORS = {
	INCORRECT_ADMIN_PASSWORD: "Senha incorreta, tente novamente!",
	MISSING_ADMIN_PASSWORD: "A senha é obrigatória!"
};

interface AdminPasswordResponse {
	correct: boolean;
}

export class AdminController {
	password = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	showPassword = $state(false);

	touched = $state({
		password: false
	});

	get isPasswordValid() {
		return this.password.length >= 8;
	}

	onInput() {
		this.apiError = null;
	}

	togglePassword() {
		this.showPassword = !this.showPassword;
	}

	async sendPassword(e: SubmitEvent) {
		e.preventDefault();

		this.touched.password = true;

		if (!this.isPasswordValid) {
			return;
		}

		this.isLoading = true;
		this.apiError = null;

		const { data, error } = await customFetch<AdminPasswordResponse>(
			window.fetch,
			"/api/admin",
			{
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ password: this.password })
			},
			ADMIN_PASSWORD_ERRORS
		);

		if (error) {
			this.apiError = error;
			this.isLoading = false;
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Falha ao processar resposta do servidor.",
				ADMIN_PASSWORD_ERRORS
			);
			this.isLoading = false;
			return;
		}

		if (!data.correct) {
			this.apiError = normalizeApiError(
				"INCORRECT_ADMIN_PASSWORD",
				"Senha incorreta.",
				ADMIN_PASSWORD_ERRORS
			);
			this.isLoading = false;
			return;
		}

		try {
			await invalidateAll();
		} finally {
			this.password = "";
			this.touched.password = false;
			this.isLoading = false;
		}
	}
}
