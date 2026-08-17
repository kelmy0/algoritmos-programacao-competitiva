import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";
import { BaseController } from "./base_controller.svelte";

export abstract class BaseAuthController extends BaseController {
	protected _turnstileToken = $state("");
	protected _turnstileComponent = $state<{ reset: () => void } | null>(null);

	constructor(initialError: ApiError | null = null) {
		super(initialError);
	}

	get turnstileToken() {
		return this._turnstileToken;
	}

	setTurnstileComponent(component: { reset: () => void } | null) {
		this._turnstileComponent = component;
	}

	onTurnstileSuccess(token: string) {
		this._turnstileToken = token;
	}

	onTurnstileExpire() {
		this._turnstileToken = "";
	}

	protected resetTurnstile() {
		this._turnstileToken = "";
		this._turnstileComponent?.reset();
	}

	protected validateTurnstile(): boolean {
		if (!this._turnstileToken) {
			this._apiError = normalizeApiError("CAPTCHA_REQUIRED");
			return false;
		}
		return true;
	}
}
