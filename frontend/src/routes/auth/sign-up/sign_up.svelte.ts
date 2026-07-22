import { goto, invalidateAll } from '$app/navigation';
import type { ApiError } from '$lib/types/api';
import { getErrorMessage } from '$lib/utils/errors';
import { tick } from 'svelte';

interface SignUpRequest {
	name: string;
	username: string;
	email: string;
	password: string;
	confirm_password: string;
}

interface SignUpResponse {
	access_token?: string;
	success: boolean;
	auto_login: boolean;
}

export class SignUpController {
	name = $state('');
	username = $state('');
	email = $state('');
	password = $state('');
	confirmPassword = $state('');
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	showPassword = $state(false);
	showConfirmPassword = $state(false);

	nameInput = $state<HTMLInputElement | null>(null);
	usernameInput = $state<HTMLInputElement | null>(null);
	emailInput = $state<HTMLInputElement | null>(null);
	passwordInput = $state<HTMLInputElement | null>(null);
	confirmPasswordInput = $state<HTMLInputElement | null>(null);

	touched = $state({
		name: false,
		username: false,
		email: false,
		password: false,
		confirmPassword: false
	});

	get hasMinLength() {
		return this.password.length >= 8;
	}
	get hasUppercase() {
		return /[A-Z]/.test(this.password);
	}
	get hasLowercase() {
		return /[a-z]/.test(this.password);
	}
	get hasNumber() {
		return /\d/.test(this.password);
	}
	get hasSpecialChar() {
		return /[@$!%*?&]/.test(this.password);
	}

	get isPasswordValid() {
		return (
			this.hasMinLength &&
			this.hasUppercase &&
			this.hasLowercase &&
			this.hasNumber &&
			this.hasSpecialChar
		);
	}

	get isPasswordsMatching() {
		return this.password === this.confirmPassword;
	}

	get cleanName() {
		return sanitizeHumanName(this.name);
	}
	get isNameValid() {
		return this.cleanName.length >= 6;
	}

	get cleanUsername() {
		return sanitizeUsername(this.username);
	}
	get isUsernameValid() {
		return this.cleanUsername.length >= 6;
	}

	get isEmailValid() {
		return isValidEmail(this.email);
	}

	clearApiError(codes: string[]) {
		if (this.apiError && codes.includes(this.apiError.code)) {
			this.apiError = null;
		}
	}

	onNameInput() {
		this.clearApiError(['REGISTRATION_INVALID_NAME']);
	}

	onNameBlur() {
		this.touched.name = true;
		this.name = this.cleanName;
	}

	onUsernameInput() {
		this.username = sanitizeUsername(this.username);
		this.clearApiError(['REGISTRATION_INVALID_USERNAME']);
	}

	onUsernameBlur() {
		this.touched.username = true;
	}

	onEmailInput() {
		this.clearApiError(['REGISTRATION_INVALID_EMAIL']);
	}

	onEmailBlur() {
		this.touched.email = true;
	}

	onPasswordInput() {
		this.clearApiError(['USER_PASSWORDS_DONT_MATCH', 'USER_PASSWORD_NOT_VALID']);
	}

	onPasswordBlur() {
		this.touched.password = true;
	}

	onConfirmPasswordBlur() {
		this.touched.confirmPassword = true;
	}

	togglePassword() {
		this.showPassword = !this.showPassword;
	}

	toggleConfirmPassword() {
		this.showConfirmPassword = !this.showConfirmPassword;
	}

	async signUp(event: SubmitEvent) {
		event.preventDefault();

		this.touched = {
			name: true,
			username: true,
			email: true,
			password: true,
			confirmPassword: true
		};

		if (!this.isNameValid) {
			this.apiError = {
				code: 'REGISTRATION_INVALID_NAME',
				message: 'O nome deve conter pelo menos 6 letras.'
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isUsernameValid) {
			this.apiError = {
				code: 'REGISTRATION_INVALID_USERNAME',
				message: 'Username deve ter pelo menos 6 caracteres válidos.'
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isEmailValid) {
			this.apiError = {
				code: 'REGISTRATION_INVALID_EMAIL',
				message: 'Digite um endereço de e-mail válido.'
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isPasswordValid) {
			this.apiError = {
				code: 'USER_PASSWORD_NOT_VALID',
				message: SIGN_UP_ERRORS.USER_PASSWORD_NOT_VALID
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isPasswordsMatching) {
			this.apiError = {
				code: 'USER_PASSWORDS_DONT_MATCH',
				message: SIGN_UP_ERRORS.USER_PASSWORDS_DONT_MATCH
			};
			await this.focusFirstInvalidField();
			return;
		}

		this.isLoading = true;

		try {
			const response = await fetch('/api/auth/sign-up', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: this.name,
					username: this.username,
					email: this.email,
					password: this.password,
					confirm_password: this.confirmPassword
				})
			});

			if (!response.ok) {
				const errorData: ApiError = await response.json();

				this.apiError = {
					code: errorData.code,
					message: getErrorMessage(errorData.code, errorData.message, SIGN_UP_ERRORS)
				};

				this.isLoading = false;
				await tick();

				if (errorData.code === 'USERNAME_ALREADY_USED') {
					this.scrollToAndFocus(this.usernameInput);
				} else if (errorData.code === 'EMAIL_ALREADY_USED') {
					this.scrollToAndFocus(this.emailInput);
				}
				return;
			}

			const data: SignUpResponse = await response.json();

			if (data.access_token) {
				this.apiError = null;
				await invalidateAll();
				await goto('/');
			} else if (data.success && !data.auto_login) {
				this.apiError = null;
				goto('/login');
			} else {
				this.apiError = {
					code: 'REGISTRATION_UNEXPECTED_ERROR',
					message: SIGN_UP_ERRORS.REGISTRATION_UNEXPECTED_ERROR
				};
			}
		} catch (err) {
			this.apiError = {
				code: 'NETWORK_ERROR',
				message: getErrorMessage('NETWORK_ERROR', '')
			};
		} finally {
			this.isLoading = false;
		}
	}

	private async focusFirstInvalidField() {
		await tick();

		if (!this.isNameValid) {
			this.scrollToAndFocus(this.nameInput);
		} else if (!this.isUsernameValid) {
			this.scrollToAndFocus(this.usernameInput);
		} else if (!this.isEmailValid) {
			this.scrollToAndFocus(this.emailInput);
		} else if (!this.isPasswordValid) {
			this.scrollToAndFocus(this.passwordInput);
		} else if (!this.isPasswordsMatching) {
			this.scrollToAndFocus(this.confirmPasswordInput);
		}
	}

	private scrollToAndFocus(element: HTMLInputElement | null) {
		if (!element) return;
		element.scrollIntoView({ behavior: 'smooth', block: 'center' });
		element.focus({ preventScroll: true });
	}
}

const SIGN_UP_ERRORS: Record<string, string> = {
	// Sign-up / Registration
	EMAIL_ALREADY_USED: 'Este endereço de e-mail já está cadastrado.',
	USERNAME_ALREADY_USED: 'Este nome de usuário já está cadastrado.',
	USER_PASSWORDS_DONT_MATCH: 'As senhas digitadas não coincidem.',
	USER_PASSWORD_NOT_VALID: 'A senha é muito fraca.',
	REGISTRATION_INVALID_NAME: 'O campo nome está inválido ou mal preenchido.',
	REGISTRATION_INVALID_USERNAME: 'O campo nome de usuário está inválido ou mal preenchido.',
	REGISTRATION_INVALID_EMAIL: 'O formato do e-mail digitado não é válido.',
	REGISTRATION_UNEXPECTED_ERROR: 'Ocorreu um erro inesperado ao criar sua conta.',

	// General Tokens
	TOKEN_CRYPT_FAILED: 'Erro interno de criptografia de segurança.'
};

export function sanitizeHumanName(name: string): string {
	const clean = name.replace(/[^\p{L}\s.'-]/gu, '');

	const words = clean.replace(/\s+/g, ' ').trim().split(' ');

	if (words.length === 0 || words[0] === '') return '';

	return words.map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()).join(' ');
}

export function sanitizeUsername(username: string): string {
	const clean = username.replace(/[^\p{L}\p{N}_-]/gu, '');
	return clean.replace(/\s+/g, '').toLowerCase();
}

export function isValidEmail(email: string): boolean {
	const clean = email.trim().toLowerCase();
	const atIndex = clean.indexOf('@');
	const lastDotIndex = clean.lastIndexOf('.');

	return atIndex > 0 && lastDotIndex > atIndex + 1 && lastDotIndex < clean.length - 1;
}
