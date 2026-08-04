import { customFetch } from "$lib/api/client";
import { ADMIN_ALGORITHMS_ERRORS } from "$lib/errors/admin/algorithms";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";
import type { Algorithm } from "$lib/types/algorithm";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";

export class NewAlgorithmController {
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	isSuccess = $state(false);
	link = $state("");

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

	async submit(content: AlgorithmPayload) {
		if (this.isLoading) {
			return;
		}

		this.isLoading = true;

		const { data, error } = await customFetch<{ algorithm: Algorithm }>(
			window.fetch,
			"/api/admin/algorithms/new",
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
			this.apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return;
		}

		this.isSuccess = true;
		this.apiError = null;
		this.link = `/admin/algorithms/my-algorithms/${data.algorithm.slug}-${data.algorithm.publicId}`;
		scrollToAndFocus(this.alertDiv);

		this.isLoading = false;
	}
}
