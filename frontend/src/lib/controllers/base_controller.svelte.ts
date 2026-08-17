import type { ApiError } from "$lib/types/api";

export abstract class BaseController {
	protected _isLoading = $state(false);
	protected _apiError = $state<ApiError | null>(null);

	constructor(initialError: ApiError | null = null) {
		this._apiError = initialError;
	}

	get isLoading() {
		return this._isLoading;
	}

	get apiError() {
		return this._apiError;
	}

	protected clearApiError(targetCodes?: string | string[]) {
		if (!this._apiError) return;

		if (!targetCodes) {
			this._apiError = null;
			return;
		}

		const codes = Array.isArray(targetCodes) ? targetCodes : [targetCodes];
		if (codes.includes(this._apiError.code)) {
			this._apiError = null;
		}
	}
}
