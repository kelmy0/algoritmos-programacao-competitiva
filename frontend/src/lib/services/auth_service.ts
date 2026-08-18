import { browser } from "$app/env";
import { AUTH_ERRORS } from "$lib/errors/auth/auth_errors";
import type { ApiError } from "$lib/types/api";
import type { RefreshServerResponse } from "$lib/types/jwt";
import { normalizeApiError } from "$lib/utils/errors";

export class AuthService {
	private static currentError: ApiError | null = null;
	private static refreshPromise: Promise<boolean> | null = null;
	private static currentExpiresAt: number | null = null;

	static getLastError(): ApiError | null {
		return this.currentError;
	}

	static setExpiresAt(expiresAt: number | null): void {
		this.currentExpiresAt = expiresAt;
	}

	static async silentRefresh(fetchImpl: typeof fetch = window.fetch): Promise<boolean> {
		if (!browser) return false;

		if (this.refreshPromise) {
			return this.refreshPromise;
		}

		this.refreshPromise = (async () => {
			this.currentError = null;

			try {
				const response = await fetchImpl("/api/auth/refresh", {
					method: "POST",
					headers: { "Content-Type": "application/json" }
				});

				if (!response.ok) {
					const errorData = await response.json().catch(() => null);
					this.currentError = normalizeApiError(errorData, "Sessão expirada.", AUTH_ERRORS);
					this.currentExpiresAt = null;
					return false;
				}

				const data: RefreshServerResponse = await response.json();

				if (!data.accessToken) {
					this.currentError = normalizeApiError("SESSION_EXPIRED");
					this.currentExpiresAt = null;
					return false;
				}

				if (data.expiresAt) {
					this.currentExpiresAt = data.expiresAt;
				}

				return true;
			} catch (error) {
				this.currentError = normalizeApiError(error, "Não foi possível renovar a sessão.");
				return false;
			} finally {
				this.refreshPromise = null;
			}
		})();

		return this.refreshPromise;
	}

	static async ensureValidSession(
		fetchImpl: typeof fetch,
		pageExpiresAt?: number | null
	): Promise<boolean> {
		if (!browser) return true;

		const expiresAt = this.currentExpiresAt ?? pageExpiresAt;
		if (!expiresAt) return true;

		const nowInMs = Date.now();
		const BUFFER_MS = 100 * 1000;

		if (expiresAt - nowInMs <= BUFFER_MS) {
			return await AuthService.silentRefresh(fetchImpl);
		}

		return true;
	}
}
