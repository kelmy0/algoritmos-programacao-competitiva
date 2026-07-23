import { jwtDecode } from 'jwt-decode';
import type { ApiError } from '$lib/types/api';
import { normalizeApiError } from '$lib/utils/errors';
import { AUTH_ERRORS } from '../../routes/(public)/auth/login/login.svelte';

interface JwtPayload {
	exp?: number;
}

declare global {
	interface Window {
		__activeRefreshPromise?: Promise<boolean>;
	}
}

export class AuthService {
	private static currentError: ApiError | null = null;
	private static refreshTimeoutId: ReturnType<typeof setTimeout> | null = null;

	static getLastError(): ApiError | null {
		return this.currentError;
	}

	static async silentRefresh(): Promise<boolean> {
		if (typeof window === 'undefined') return false;

		if (window.__activeRefreshPromise) {
			return window.__activeRefreshPromise;
		}

		window.__activeRefreshPromise = (async () => {
			this.currentError = null;

			try {
				const response = await fetch('/api/auth/refresh', { method: 'POST' });

				if (!response.ok) {
					const errorData: ApiError = await response.json().catch(() => null);
					this.currentError = normalizeApiError(
						errorData,
						'Sua sessão expirou. Faça login novamente.',
						AUTH_ERRORS
					);
					return false;
				}

				const data = await response.json();
				if (data.accessToken) {
					console.log('Sessão renovada com sucesso');
					AuthService.startAutoRefreshTimer(data.accessToken);
				}

				return true;
			} catch (error) {
				this.currentError = normalizeApiError(
					error,
					'Não foi possível conectar ao servidor.',
					AUTH_ERRORS
				);
				return false;
			} finally {
				window.__activeRefreshPromise = undefined;
			}
		})();

		return window.__activeRefreshPromise;
	}

	static startAutoRefreshTimer(accessToken: string): void {
		if (typeof window === 'undefined') return;

		if (this.refreshTimeoutId) {
			clearTimeout(this.refreshTimeoutId);
		}

		try {
			const decoded = jwtDecode<JwtPayload>(accessToken);
			if (!decoded.exp) return;

			const nowInMs = Date.now();
			const expInMs = decoded.exp * 1000;

			const BUFFER_MS = 60 * 1000;
			const timeUntilRefresh = expInMs - nowInMs - BUFFER_MS;

			if (timeUntilRefresh <= 0) {
				AuthService.silentRefresh();
			} else {
				console.log(
					`Iniciado a contagem do refresh, faltam: ${Math.round(timeUntilRefresh / 1000)}s`
				);
				this.refreshTimeoutId = setTimeout(() => {
					AuthService.silentRefresh();
				}, timeUntilRefresh);
			}
		} catch (err) {
			console.error('Erro ao agendar a renovação do token:', err);
		}
	}
}
