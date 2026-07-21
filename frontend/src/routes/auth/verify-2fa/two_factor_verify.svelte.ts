import { goto } from '$app/navigation';
import type { ApiError } from '$lib/types/api';
import { setCookie } from '$lib/utils/cookie';
import { page } from '$app/state';
import { getErrorMessage } from '$lib/utils/errors';
import { PUBLIC_API_URL, PUBLIC_ENV } from '$env/static/public';

interface TwoFactorRequest {
	pre_auth_token: string;
	code: string;
}

interface TwoFactorResponse {
	access_token: string;
	requires_2fa: string;
}

const TWO_FACTOR_ERRORS: Record<string, string> = {
	INVALID_SESSION_DATA: 'Está faltando o id do usuário no token. Faça login novamente!'
};

export class TwoFactorController {
	token = '';
	code = $state('');
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);

	touched = $state({
		code: false
	});

	get isCodeValid() {
		return this.code.length === 6;
	}

	onInput(event: Event) {
		const input = event.target as HTMLInputElement;

		this.code = input.value.replace(/\D/g, '');

		if (this.apiError) {
			this.apiError = null;
		}

		if (this.code.length === 6 && !this.isLoading) {
			const form = input.closest('form');
			if (form) {
				form.requestSubmit();
			}
		}
	}

	getToken() {
		const token = page.url.searchParams.get('token');

		if (!token) {
			goto('/auth/login?error=MISSING_PRE_TOKEN');
			return;
		}

		this.token = token;
	}

	async sendCode(event: SubmitEvent) {
		event.preventDefault();

		this.touched.code = true;

		if (!this.isCodeValid) {
			return;
		}

		this.isLoading = true;
		this.apiError = null;

		const bodyRequest: TwoFactorRequest = {
			pre_auth_token: this.token,
			code: this.code
		};

		try {
			const twoFactorUrl = `${PUBLIC_API_URL}/api/auth/verify-2fa`;
			const response = await fetch(twoFactorUrl, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(bodyRequest)
			});

			if (!response.ok) {
				this.touched.code = false;
				this.code = '';
				const errorData: ApiError = await response.json();
				this.apiError = {
					code: errorData.code,
					message: getErrorMessage(errorData.code, errorData.message, TWO_FACTOR_ERRORS)
				};
				return;
			}

			const data: TwoFactorResponse = await response.json();

			if (data.requires_2fa) {
				goto(`/auth/login?error=AUTH_UNEXPECTED_ERROR`);
				return;
			}

			if (data.access_token) {
				setCookie('access_token', data.access_token, 15, PUBLIC_ENV !== 'development');
				goto('/');
			}
		} catch {
			this.apiError = {
				code: 'NETWORK_ERROR',
				message: getErrorMessage('NETWORK_ERROR', '')
			};
		} finally {
			this.isLoading = false;
		}
	}
}
