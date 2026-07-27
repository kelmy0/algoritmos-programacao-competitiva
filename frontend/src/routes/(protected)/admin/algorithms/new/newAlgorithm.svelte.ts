import { customFetch } from "$lib/api/client";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";
import type { Algorithm } from "$lib/types/algorithm";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";

export const ADMIN_ALGORITHMS_ERRORS: Record<string, string> = {
	ALGORITHM_INVALID_NAME: "Nome do algoritmo não é válido.",
	ALGORITHM_INVALID_CATEGORY: "Categoria do algoritmo não é válida.",
	ALGORITHM_INVALID_CONTENT: "Conteudo no markdown não é válido.",
	ALGORITHM_GENERATE_PUBLIC_ID_FAILED: "Falha interna ao gerar id público, tente novamente!",
	ALGORITHM_POST_FAILED: "Falha interna ao postar algoritmo, tente novamente!"
};

export class NewAlgorithmController {
	password = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	isSuccess = $state(false);
	link = $state("");

	passwordInput = $state<HTMLInputElement | null>(null);
	alertDiv = $state<HTMLDivElement | null>(null);

	touched = $state({
		password: false
	});

	hasNameError = $derived(this.apiError?.code === "ALGORITHM_INVALID_NAME");
	hasCategoryError = $derived(this.apiError?.code === "ALGORITHM_INVALID_CATEGORY");
	hasContentError = $derived(this.apiError?.code === "ALGORITHM_INVALID_CONTENT");
	hasPasswordError = $derived(
		this.apiError?.code === "INVALID_PASSWORD" || (this.touched.password && !this.isPasswordValid)
	);

	get isPasswordValid() {
		return this.password.length >= 8;
	}

	onPasswordInput() {
		this.clearApiError(["USER_PASSWORDS_DONT_MATCH", "USER_PASSWORD_NOT_VALID"]);
	}

	onPasswordBlur() {
		this.touched.password = true;
	}

	onNameInput() {
		this.clearApiError(["ALGORITHM_INVALID_NAME"]);
	}

	onCategoryInput() {
		this.clearApiError(["ALGORITHM_INVALID_CATEGORY"]);
	}

	onContentInput() {
		this.clearApiError(["ALGORITHM_INVALID_CONTENT"]);
	}

	clearApiError(codes: string[]) {
		if (this.apiError && codes.includes(this.apiError.code)) {
			this.apiError = null;
		}
	}

	async submit(content: AlgorithmPayload) {
		this.isLoading = true;
		const { data, error } = await customFetch<{ algorithm: Algorithm }>(
			window.fetch,
			"/api/admin/algorithms/new",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"x-admin-secret": this.password
				},
				body: JSON.stringify(content)
			},
			ADMIN_ALGORITHMS_ERRORS // Suas mensagens amigáveis locais
		);

		if (error) {
			this.apiError = error;
			this.isSuccess = false;
			await scrollToAndFocus(this.alertDiv);
			this.isLoading = false;
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Erro inesperado.",
				ADMIN_ALGORITHMS_ERRORS
			);
			return;
		}

		this.isSuccess = true;
		this.apiError = null;
		this.link = `/algorithms/${data.algorithm.Slug}-${data.algorithm.PublicId}`;
		await scrollToAndFocus(this.alertDiv);

		this.isLoading = false;
	}
}
