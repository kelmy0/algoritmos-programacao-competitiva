import { goto, invalidateAll } from "$app/navigation";
import type { ApiError } from "$lib/types/api";
import { page } from "$app/state";
import { normalizeApiError } from "$lib/utils/errors";
import { customFetch } from "$lib/api/client";

interface TwoFactorRequest {
	pre_auth_token: string;
	code: string;
}

export interface TwoFactorServerResponse {
	access_token: boolean;
	requires_2fa: boolean;
}

export const TWO_FACTOR_ERRORS: Record<string, string> = {
	INVALID_SESSION_DATA: "Está faltando o id do usuário no token. Faça login novamente!"
};

export class TwoFactorController {
	token = "";
	code = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);

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

	async sendCode(e: SubmitEvent) {
		e.preventDefault();

		this.touched.code = true;

		if (!this.isCodeValid) {
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
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(bodyRequest)
			},
			TWO_FACTOR_ERRORS
		);

		this.isLoading = false;
		if (error) {
			this.apiError = error;
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Falha ao processar resposta do servidor.",
				TWO_FACTOR_ERRORS
			);
			return;
		}

		if (data.requires_2fa) {
			await goto(`/auth/login?error=AUTH_UNEXPECTED_ERROR`);
			return;
		}

		await goto("/", { invalidateAll: true });
	}
}
