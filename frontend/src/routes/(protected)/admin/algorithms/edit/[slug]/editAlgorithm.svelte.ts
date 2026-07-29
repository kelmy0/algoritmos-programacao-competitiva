import { customFetch } from "$lib/api/client";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";
import type { Algorithm } from "$lib/types/algorithm";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";
import { ADMIN_ALGORITHMS_ERRORS } from "../../new/newAlgorithm.svelte";

export class EditAlgorithmController {
	isLoading = $state(false);
	isDeleting = $state(false);
	apiError = $state<ApiError | null>(null);
	isSuccess = $state(false);
	link = $state("");
	publicId = $state("");
	slug = $state("");
	isDeleteModalOpen = $state(false);
	isSaved = $state(true);

	alertDiv = $state<HTMLDivElement | null>(null);
	actionErrorLabel = $state("salvar");

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

	openDeleteModal() {
		this.isDeleteModalOpen = true;
	}

	closeDeleteModal() {
		this.isDeleteModalOpen = false;
	}

	clearApiError(codes: string[]) {
		if (this.apiError && codes.includes(this.apiError.code)) {
			this.apiError = null;
		}
	}

	async editAlgorithm(content: AlgorithmPayload) {
		if (this.isDeleting || this.isLoading) {
			return;
		}

		this.actionErrorLabel = "salvar";

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

		scrollToAndFocus(this.alertDiv);
		this.isLoading = false;

		if (error || !data) {
			this.apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			this.isSuccess = false;
			return;
		}

		this.isSaved = true;
		this.isSuccess = true;
		this.apiError = null;
		this.link = `/admin/algorithms/edit/${data.algorithm.Slug}-${data.algorithm.PublicId}`;
	}

	async handleDelete() {
		if (this.isLoading || this.isDeleting) {
			return;
		}

		if (!this.publicId || !this.slug) {
			this.apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			return;
		}

		this.actionErrorLabel = "deletar";
		this.isDeleting = true;

		const { data, error } = await customFetch<{ deleted: boolean }>(
			window.fetch,
			`/api/admin/algorithms/delete/${this.slug}-${this.publicId}`,
			{
				method: "DELETE",
				headers: {
					"Content-Type": "application/json"
				}
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		this.isDeleteModalOpen = false;
		this.isDeleting = false;
		scrollToAndFocus(this.alertDiv);

		if (error || !data || !data.deleted) {
			this.isSuccess = false;
			this.apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			return;
		}

		this.isSaved = false;
		this.isSuccess = true;
		this.link = `/admin/algorithms/trash/${this.slug}-${this.publicId}`;
	}
}
