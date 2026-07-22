import { jwtDecode } from 'jwt-decode';
import type { ApiError } from '$lib/types/api';
import { getErrorMessage } from '$lib/utils/errors';
import { AUTH_ERRORS } from '../../routes/auth/login/login.svelte';

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
				const response = await fetch('/auth/refresh', { method: 'POST' });

				if (!response.ok) {
					const errorData: ApiError = await response.json();
					this.currentError = {
						code: errorData.code || 'AUTH_UNEXPECTED_ERROR',
						message: getErrorMessage(errorData.code, errorData.message, AUTH_ERRORS)
					};
					return false;
				}

				const data = await response.json();
				if (data.accessToken) {
					console.log('Renovado');
					AuthService.startAutoRefreshTimer(data.accessToken);
				}

				return true;
			} catch (error) {
				this.currentError = {
					code: 'AUTH_UNEXPECTED_ERROR',
					message: getErrorMessage('AUTH_UNEXPECTED_ERROR', 'Falha na rede ou no servidor.')
				};
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
				console.log(`Iniciado a contagem, faltam: ${timeUntilRefresh / 1000}s`);
				this.refreshTimeoutId = setTimeout(() => {
					AuthService.silentRefresh();
				}, timeUntilRefresh);
			}
		} catch (err) {
			console.error('Error scheduling the next refresh:', err);
		}
	}
}
