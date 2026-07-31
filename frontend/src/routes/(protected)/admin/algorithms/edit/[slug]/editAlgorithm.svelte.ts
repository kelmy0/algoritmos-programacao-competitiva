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
	lastAction = $state<"save" | "delete" | "restore">("save");

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

	async handleSubmit(content: AlgorithmPayload): Promise<boolean> {
		if (this.isDeleting || this.isLoading) {
			return false;
		}
		this.lastAction = "save";

		if (!this.publicId) {
			this.apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			scrollToAndFocus(this.alertDiv);
			return false;
		}

		this.isLoading = true;
		const { data, error } = await customFetch<{ algorithm: Algorithm }>(
			window.fetch,
			`/api/admin/algorithms/edit/${this.slug}-${this.publicId}`,
			{
				method: "PUT",
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
			return false;
		}

		this.isSuccess = true;
		this.apiError = null;
		this.slug = data.algorithm.slug;
		this.publicId = data.algorithm.publicId;
		this.link = `/admin/algorithms/my-algorithms/${this.slug}-${this.publicId}`;
		return true;
	}

	async handleDelete(): Promise<boolean> {
		if (this.isLoading || this.isDeleting) {
			return false;
		}

		this.lastAction = "delete";

		if (!this.publicId || !this.slug) {
			this.apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			this.isDeleteModalOpen = false;
			scrollToAndFocus(this.alertDiv);
			return false;
		}

		this.isDeleting = true;

		const { error, status } = await customFetch<null>(
			window.fetch,
			`/api/admin/algorithms/delete/${this.slug}-${this.publicId}`,
			{
				method: "DELETE"
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		this.isDeleteModalOpen = false;
		this.isDeleting = false;
		scrollToAndFocus(this.alertDiv);

		if (error || status !== 204) {
			this.isSuccess = false;
			this.apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		this.isSuccess = true;
		this.apiError = null;
		this.link = `/admin/algorithms/my-algorithms/${this.slug}-${this.publicId}`;
		return true;
	}

	async handleRestore(): Promise<boolean> {
		if (this.isLoading || this.isDeleting) {
			return false;
		}

		this.lastAction = "restore";

		if (!this.publicId || !this.slug) {
			this.apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			scrollToAndFocus(this.alertDiv);
			return false;
		}

		this.isLoading = true;

		const { error, status } = await customFetch<null>(
			window.fetch,
			`/api/admin/algorithms/restore/${this.slug}-${this.publicId}`,
			{
				method: "PATCH"
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		this.isLoading = false;
		scrollToAndFocus(this.alertDiv);

		if (error || status !== 204) {
			this.isSuccess = false;
			this.apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		this.isSuccess = true;
		this.apiError = null;
		this.link = `/admin/algorithms/my-algorithms/${this.slug}-${this.publicId}`;
		return true;
	}
}
