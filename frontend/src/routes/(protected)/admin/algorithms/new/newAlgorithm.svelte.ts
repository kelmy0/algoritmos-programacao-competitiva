import { customFetch } from "$lib/api/client";
import { ADMIN_ALGORITHMS_ERRORS } from "$lib/errors/admin/algorithms";
import type { NewAlgorithmResponse } from "$lib/types/algorithm";
import type { AlgorithmEditor } from "$lib/states/editor.svelte";
import { normalizeApiError } from "$lib/utils/errors";
import { BaseEditorController } from "$lib/controllers/base_editor_controller.svelte";

export class NewAlgorithmController extends BaseEditorController {
	async save(editor: AlgorithmEditor): Promise<boolean> {
		if (this._isLoading) return false;

		const payload = editor.getPayload();
		if (!payload) return false;

		this._isLoading = true;
		this._apiError = null;

		const { data, error } = await customFetch<NewAlgorithmResponse>(
			window.fetch,
			"/api/admin/algorithms/new",
			{
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(payload)
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			this._isSuccess = false;

			editor.setApiError(error);
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		this._isSuccess = true;
		this._link = `/admin/algorithms/my-algorithms/${data.slug}-${data.publicId}`;

		return true;
	}
}
