import type { ApiError } from '$lib/types/api';
import { getErrorMessage } from '$lib/utils/errors';
import { AUTH_ERRORS } from '../../routes/auth/login/login.svelte';

export class AuthService {
	private static currentError: ApiError | null = null;

	static getLastError(): ApiError | null {
		return this.currentError;
	}

	static async silentRefresh(): Promise<boolean> {
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

			return true;
		} catch (error) {
			this.currentError = {
				code: 'AUTH_UNEXPECTED_ERROR',
				message: getErrorMessage('AUTH_UNEXPECTED_ERROR', 'Falha na rede ou no servidor.')
			};

			return false;
		}
	}

	/**
	 * @param intervalMinutes Intervalo em minutos (padrão: 14 min)
	 */
	static startAutoRefreshTimer(intervalMinutes: number = 14): void {
		if (typeof window === 'undefined') return;

		const intervalMs = intervalMinutes * 60 * 1000;

		setInterval(async () => {
			await AuthService.silentRefresh();
		}, intervalMs);
	}
}
