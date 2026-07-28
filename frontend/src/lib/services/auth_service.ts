import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";
import { AUTH_ERRORS } from "../../routes/(public)/auth/login/login.svelte";

declare global {
	interface Window {
		__activeRefreshPromise?: Promise<boolean>;
		__refreshTimeoutId?: ReturnType<typeof setTimeout>;
	}
}

interface RefreshResponse {
	accessToken: boolean;
	expiresAt?: number;
}

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
		if (typeof window === "undefined") return false;

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

				const data: RefreshResponse = await response.json();

				if (!data?.accessToken) {
					this.currentError = normalizeApiError("UNAUTHORIZED", "Sessão expirada.");
					this.currentExpiresAt = null;
					return false;
				}

				console.log("Sessão renovada com sucesso.");

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
		if (typeof window === "undefined") return true;

		const expiresAt = this.currentExpiresAt ?? pageExpiresAt;
		if (!expiresAt) return true;

		const nowInMs = Date.now();
		const BUFFER_MS = 60 * 1000;

		if (expiresAt - nowInMs <= BUFFER_MS) {
			console.log("Token próximo do vencimento.");
			return await AuthService.silentRefresh(fetchImpl);
		}

		return true;
	}
}
