import { goto } from '$app/navigation';
import type { ApiError } from '$lib/types/api';
import { setCookie } from '$lib/utils/cookie';
import { page } from '$app/state';
import { getErrorMessage } from '$lib/utils/errors';
import { PUBLIC_API_URL, PUBLIC_ENV } from '$env/static/public';
import { isValidEmail } from '../sign-up/sign_up.svelte';

interface LoginRequest {
	email: string;
	password: string;
}

interface LoginResponse {
	access_token?: string;
	requires_2fa: boolean;
	pre_auth_token?: string;
}

const AUTH_ERRORS: Record<string, string> = {
	AUTH_INVALID_EMAIL_PASSWORD: 'E-mail ou senha incorretos. Verifique seus dados.',
	USER_ALREADY_EXISTS:
		'Este e-mail já está cadastrado. Tente entrar por outro método ou use um email diferente.',
	SOCIAL_ACCOUNT_ALREADY_LINKED:
		'Este email já esta ligado a outra conta. Tente entrar por outro método ou use um email diferente. '
};

export class LoginController {
	email = $state('');
	password = $state('');
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	showPassword = $state(false);

	touched = $state({
		email: false,
		password: false
	});

	get isEmailValid() {
		return isValidEmail(this.email);
	}

	get isPasswordValid() {
		return this.password.length >= 8;
	}

	onInput() {
		if (this.apiError) {
			this.apiError = null;
		}
	}

	togglePassword() {
		this.showPassword = !this.showPassword;
	}

	checkErrors() {
		const error = page.url.searchParams.get('error');

		if (error) {
			this.apiError = {
				code: error,
				message: getErrorMessage(error, 'Erro desconhecido.', AUTH_ERRORS)
			};
		}
	}

	async login(event: SubmitEvent) {
		event.preventDefault();
		this.touched.email = true;
		this.touched.password = true;

		if (!this.isEmailValid || !this.isPasswordValid) {
			return;
		}

		this.isLoading = true;
		this.apiError = null;

		const bodyRequest: LoginRequest = {
			email: this.email,
			password: this.password
		};

		try {
			const loginUrl = `${PUBLIC_API_URL}/api/auth/login`;
			const response = await fetch(loginUrl, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(bodyRequest)
			});

			if (!response.ok) {
				const errorData: ApiError = await response.json();
				this.apiError = {
					code: errorData.code,
					message: getErrorMessage(errorData.code, errorData.message, AUTH_ERRORS)
				};
				return;
			}

			const data: LoginResponse = await response.json();

			if (data.requires_2fa) {
				goto(`/auth/verify-2fa?token=${data.pre_auth_token}`);
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
