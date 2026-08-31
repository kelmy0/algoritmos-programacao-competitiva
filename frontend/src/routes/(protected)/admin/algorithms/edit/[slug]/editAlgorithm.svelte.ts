import { customFetch } from "$lib/api/client";
import { ADMIN_ALGORITHMS_ERRORS } from "$lib/errors/admin/algorithms";
import type { Algorithm, EditAlgorithmResponse } from "$lib/types/algorithm";
import type { AlgorithmEditor } from "$lib/states/editor.svelte";
import { normalizeApiError } from "$lib/utils/errors";
import { BaseEditorController } from "$lib/controllers/base_editor_controller.svelte";

export type ActionType = "save" | "delete" | "restore";

export class EditAlgorithmController extends BaseEditorController {
	#lastAction = $state<ActionType>("save");
	#publicId = $state("");
	#slug = $state("");

	constructor(getAlgorithm?: () => Algorithm | null) {
		super();
		if (getAlgorithm) {
			$effect.pre(() => {
				const algo = getAlgorithm();
				if (algo) {
					if (!this.#publicId) this.#publicId = algo.publicId;
					if (!this.#slug) this.#slug = algo.slug;
				}
			});
		}
	}

	get lastAction() {
		return this.#lastAction;
	}

	get publicId() {
		return this.#publicId;
	}

	get slug() {
		return this.#slug;
	}

	async save(editor: AlgorithmEditor): Promise<boolean> {
		if (this._isLoading) return false;

		const payload = editor.getPayload();
		if (!payload) return false;

		this.#lastAction = "save";
		if (!this.#validateIdentifier()) return false;

		this._isLoading = true;
		this._apiError = null;

		const { data, error } = await customFetch<EditAlgorithmResponse>(
			window.fetch,
			`/api/admin/algorithms/edit/${this.#slug}-${this.#publicId}`,
			{
				method: "PUT",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(payload)
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		this._isLoading = false;

		if (error || !data) {
			const finalError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			this._apiError = finalError;
			this._isSuccess = false;

			editor.setApiError(finalError);
			return false;
		}

		this._isSuccess = true;
		this.#slug = data.slug;
		this._link = `/admin/algorithms/my-algorithms/${this.slug}-${this.publicId}`;

		return true;
	}

	async delete(): Promise<boolean> {
		if (this._isLoading) return false;

		this.#lastAction = "delete";
		if (!this.#validateIdentifier()) return false;

		this._isLoading = true;

		const { error, status } = await customFetch<null>(
			window.fetch,
			`/api/admin/algorithms/delete/${this.slug}-${this.publicId}`,
			{ method: "DELETE" },
			ADMIN_ALGORITHMS_ERRORS
		);

		this._isLoading = false;

		return this.#handleActionResult(error, status === 204);
	}

	async restore(): Promise<boolean> {
		if (this._isLoading) return false;

		this.#lastAction = "restore";
		if (!this.#validateIdentifier()) return false;

		this._isLoading = true;

		const { error, status } = await customFetch<null>(
			window.fetch,
			`/api/admin/algorithms/restore/${this.slug}-${this.publicId}`,
			{ method: "PATCH" },
			ADMIN_ALGORITHMS_ERRORS
		);

		this._isLoading = false;

		return this.#handleActionResult(error, status === 204);
	}

	#validateIdentifier(): boolean {
		if (!this.publicId || !this.slug) {
			this._apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"",
				ADMIN_ALGORITHMS_ERRORS
			);
			return false;
		}
		return true;
	}

	#handleActionResult(error: any, isSuccessStatus: boolean): boolean {
		if (error || !isSuccessStatus) {
			this._isSuccess = false;
			this._apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		this._isSuccess = true;
		this._apiError = null;
		this._link = `/admin/algorithms/my-algorithms/${this.slug}-${this.publicId}`;
		return true;
	}
}
