import { goto, invalidateAll } from '$app/navigation';
import type { ApiError } from '$lib/types/api';
import { normalizeApiError } from '$lib/utils/errors';
import { isValidEmail } from '../sign-up/sign_up.svelte';

interface LoginResponse {
	access_token?: string;
	requires_2fa: boolean;
	pre_auth_token?: string;
}

export const AUTH_ERRORS: Record<string, string> = {
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

	constructor(initialError: ApiError | null = null) {
		this.apiError = initialError;
	}

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

	async login(event: SubmitEvent) {
		event.preventDefault();
		this.touched.email = true;
		this.touched.password = true;

		if (!this.isEmailValid || !this.isPasswordValid) return;

		this.isLoading = true;
		this.apiError = null;

		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					email: this.email,
					password: this.password
				})
			});

			if (!response.ok) {
				const errorData: ApiError = await response.json().catch(() => ({}));
				this.apiError = normalizeApiError(errorData, 'Falha no login.', AUTH_ERRORS);
				return;
			}

			const data: LoginResponse = await response.json();

			if (data.requires_2fa) {
				goto(`/auth/verify-2fa?token=${data.pre_auth_token}`);
				return;
			}

			if (data.access_token) {
				await invalidateAll();
				await goto('/');
			}
		} catch (err) {
			this.apiError = normalizeApiError(err, 'Falha ao se conectar com o servidor.', AUTH_ERRORS);
		} finally {
			this.isLoading = false;
		}
	}
}
