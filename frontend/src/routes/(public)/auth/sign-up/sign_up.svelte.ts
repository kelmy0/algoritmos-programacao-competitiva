import { goto } from "$app/navigation";
import { customFetch } from "$lib/api/client";
import { SIGN_UP_ERRORS } from "$lib/errors/auth/sign-up";
import type { ApiError } from "$lib/types/api";
import type { SignUpServerResponse } from "$lib/types/auth/sign-up";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";
import { isValidEmail, sanitizeHumanName, sanitizeUsername } from "$lib/utils/sanitize";
import { tick } from "svelte";

export class SignUpController {
	name = $state("");
	username = $state("");
	email = $state("");
	password = $state("");
	confirmPassword = $state("");
	turnstileToken = $state("");
	isLoading = $state(false);
	apiError = $state<ApiError | null>(null);
	showPassword = $state(false);
	showConfirmPassword = $state(false);

	nameInput = $state<HTMLInputElement | null>(null);
	usernameInput = $state<HTMLInputElement | null>(null);
	emailInput = $state<HTMLInputElement | null>(null);
	passwordInput = $state<HTMLInputElement | null>(null);
	confirmPasswordInput = $state<HTMLInputElement | null>(null);

	turnstileComponent: { reset: () => void } | null = null;

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
		this.clearApiError(["REGISTRATION_INVALID_NAME"]);
	}

	onNameBlur() {
		this.touched.name = true;
		this.name = this.cleanName;
	}

	onUsernameInput() {
		this.username = sanitizeUsername(this.username);
		this.clearApiError(["REGISTRATION_INVALID_USERNAME"]);
	}

	onUsernameBlur() {
		this.touched.username = true;
	}

	onEmailInput() {
		this.clearApiError(["REGISTRATION_INVALID_EMAIL"]);
	}

	onEmailBlur() {
		this.touched.email = true;
	}

	onPasswordInput() {
		this.clearApiError(["USER_PASSWORDS_DONT_MATCH", "USER_PASSWORD_NOT_VALID"]);
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

	onTurnstileSuccess(token: string) {
		this.turnstileToken = token;
	}

	onTurnstileExpire() {
		this.turnstileToken = "";
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

		if (this.isLoading) return;

		if (!this.isNameValid) {
			this.apiError = {
				code: "REGISTRATION_INVALID_NAME",
				message: "O nome deve conter pelo menos 6 letras."
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isUsernameValid) {
			this.apiError = {
				code: "REGISTRATION_INVALID_USERNAME",
				message: "Username deve ter pelo menos 6 caracteres válidos."
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isEmailValid) {
			this.apiError = {
				code: "REGISTRATION_INVALID_EMAIL",
				message: "Digite um endereço de e-mail válido."
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isPasswordValid) {
			this.apiError = {
				code: "USER_PASSWORD_NOT_VALID",
				message: SIGN_UP_ERRORS.USER_PASSWORD_NOT_VALID
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.isPasswordsMatching) {
			this.apiError = {
				code: "USER_PASSWORDS_DONT_MATCH",
				message: SIGN_UP_ERRORS.USER_PASSWORDS_DONT_MATCH
			};
			await this.focusFirstInvalidField();
			return;
		}

		if (!this.turnstileToken) {
			this.apiError = normalizeApiError("CAPTCHA_REQUIRED");
			return;
		}

		this.isLoading = true;

		const { data, error } = await customFetch<SignUpServerResponse>(
			window.fetch,
			"/api/auth/sign-up",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this.turnstileToken
				},
				body: JSON.stringify({
					name: this.name,
					username: this.username,
					email: this.email,
					password: this.password,
					confirmPassword: this.confirmPassword
				})
			},
			SIGN_UP_ERRORS
		);

		this.isLoading = false;

		if (error) {
			this.apiError = error;
			this.turnstileToken = "";
			this.turnstileComponent?.reset();
			return;
		}

		if (!data) {
			this.apiError = normalizeApiError(
				"INTERNAL_SERVER_ERROR",
				"Falha ao processar resposta do servidor.",
				SIGN_UP_ERRORS
			);
			return;
		}

		if (data.autoLogin) {
			this.apiError = null;
			await goto("/", { invalidateAll: true });
		} else if (data.success && !data.autoLogin) {
			this.apiError = null;
			goto("/login");
		} else {
			this.apiError = {
				code: "REGISTRATION_UNEXPECTED_ERROR",
				message: SIGN_UP_ERRORS.REGISTRATION_UNEXPECTED_ERROR
			};
		}
	}

	private async focusFirstInvalidField() {
		await tick();

		if (!this.isNameValid) {
			scrollToAndFocus(this.nameInput);
		} else if (!this.isUsernameValid) {
			scrollToAndFocus(this.usernameInput);
		} else if (!this.isEmailValid) {
			scrollToAndFocus(this.emailInput);
		} else if (!this.isPasswordValid) {
			scrollToAndFocus(this.passwordInput);
		} else if (!this.isPasswordsMatching) {
			scrollToAndFocus(this.confirmPasswordInput);
		}
	}
}
