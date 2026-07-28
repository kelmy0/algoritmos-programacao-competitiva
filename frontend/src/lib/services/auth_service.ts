import type { ApiError } from "$lib/types/api";
import { normalizeApiError } from "$lib/utils/errors";
import { AUTH_ERRORS } from "../../routes/(public)/auth/login/login.svelte";
import { customFetch } from "$lib/api/client";

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
	private static activeExpiresAt: number | null = null;

	static getLastError(): ApiError | null {
		return this.currentError;
	}

	static clearAutoRefreshTimer(): void {
		if (typeof window === "undefined") return;

		this.activeExpiresAt = null;
		if (window.__refreshTimeoutId) {
			clearTimeout(window.__refreshTimeoutId);
			window.__refreshTimeoutId = undefined;
		}
	}

	static async silentRefresh(fetchFn: typeof fetch = fetch): Promise<boolean> {
		if (typeof window === "undefined") return false;
		if (window.__activeRefreshPromise) {
			return window.__activeRefreshPromise;
		}

		AuthService.clearAutoRefreshTimer();

		window.__activeRefreshPromise = (async () => {
			this.currentError = null;

			try {
				const { data, error } = await customFetch<RefreshResponse>(
					window.fetch,
					"/api/auth/refresh",
					{ method: "POST" },
					AUTH_ERRORS
				);

				if (error || !data?.accessToken) {
					this.currentError = error || normalizeApiError("UNAUTHORIZED", "Sessão expirada.");
					return false;
				}

				console.log("Sessão renovada com sucesso");

				const newExpiresAt = data.expiresAt ?? Date.now() + 15 * 60 * 1000;
				AuthService.startAutoRefreshTimer(newExpiresAt);

				return true;
			} catch (error) {
				this.currentError = normalizeApiError(error, "Não foi possível renovar a sessão.");
				return false;
			} finally {
				window.__activeRefreshPromise = undefined;
			}
		})();

		return window.__activeRefreshPromise;
	}

	static startAutoRefreshTimer(expiresAt?: number | null): void {
		if (typeof window === "undefined") return;

		if (!expiresAt) {
			AuthService.clearAutoRefreshTimer();
			return;
		}

		if (this.activeExpiresAt === expiresAt && window.__refreshTimeoutId) {
			return;
		}

		AuthService.clearAutoRefreshTimer();
		this.activeExpiresAt = expiresAt;

		const nowInMs = Date.now();
		const BUFFER_MS = 60 * 1000;
		const timeUntilRefresh = expiresAt - nowInMs - BUFFER_MS;

		if (timeUntilRefresh <= 0) {
			console.log("Token próximo da expiração ou expirado. Disparando refresh imediato");
			AuthService.silentRefresh();
		} else {
			console.log(
				`Iniciada a contagem do refresh, faltam exatamente: ${Math.round(timeUntilRefresh / 1000)}s`
			);
			window.__refreshTimeoutId = setTimeout(() => {
				AuthService.silentRefresh();
			}, timeUntilRefresh);
		}
	}
}
