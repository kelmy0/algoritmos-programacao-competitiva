import { customFetch } from "$lib/api/client";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";
import type { Algorithm } from "$lib/types/algorithm";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";
import { ADMIN_ALGORITHMS_ERRORS } from "../../new/newAlgorithm.svelte";

export class EditAlgorithmController {
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	isSuccess = $state(false);
	link = $state("");
	publicId = $state("");

	alertDiv = $state<HTMLDivElement | null>(null);

	touched = $state({
		password: false
	});

	hasNameError = $derived(this.apiError?.code === "ALGORITHM_INVALID_NAME");
	hasCategoryError = $derived(this.apiError?.code === "ALGORITHM_INVALID_CATEGORY");
	hasContentError = $derived(this.apiError?.code === "ALGORITHM_INVALID_CONTENT");

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

	async editAlgorithm(content: AlgorithmPayload) {
		if (!this.publicId) {
			this.apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			return;
		}

		this.isLoading = true;
		const { data, error } = await customFetch<{ algorithm: Algorithm }>(
			window.fetch,
			`/api/admin/algorithms/edit/${this.publicId}`,
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify(content)
			},
			ADMIN_ALGORITHMS_ERRORS
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
		this.link = `/admin/algorithms/edit/${data.algorithm.Slug}-${data.algorithm.PublicId}`;
		await scrollToAndFocus(this.alertDiv);

		this.isLoading = false;
	}
}
