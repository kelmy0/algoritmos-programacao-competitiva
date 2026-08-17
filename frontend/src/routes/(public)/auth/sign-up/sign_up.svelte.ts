import { goto } from "$app/navigation";
import { customFetch } from "$lib/api/client";
import { BaseAuthController } from "$lib/controllers/base_auth_controller.svelte";
import { SIGN_UP_ERRORS } from "$lib/errors/auth/sign-up";
import type { SignUpServerResponse } from "$lib/types/auth/sign-up";
import { normalizeApiError } from "$lib/utils/errors";
import { isValidEmail, sanitizeHumanName, sanitizeUsername } from "$lib/utils/sanitize";

export class SignUpController extends BaseAuthController {
	#name = $state("");
	#username = $state("");
	#email = $state("");
	#password = $state("");
	#confirmPassword = $state("");

	get name() {
		return this.#name;
	}

	set name(value: string) {
		this.#name = value;
		this.clearApiError("REGISTRATION_INVALID_NAME");
	}

	get username() {
		return this.#username;
	}

	set username(value: string) {
		this.#username = sanitizeUsername(value);
		this.clearApiError("REGISTRATION_INVALID_USERNAME");
	}

	get email() {
		return this.#email;
	}

	set email(value: string) {
		this.#email = value;
		this.clearApiError("REGISTRATION_INVALID_EMAIL");
	}

	get password() {
		return this.#password;
	}

	set password(value: string) {
		this.#password = value;
		this.clearApiError(["USER_PASSWORDS_DONT_MATCH", "USER_PASSWORD_NOT_VALID"]);
	}

	get confirmPassword() {
		return this.#confirmPassword;
	}

	set confirmPassword(value: string) {
		this.#confirmPassword = value;
		this.clearApiError("USER_PASSWORDS_DONT_MATCH");
	}

	get cleanName() {
		return sanitizeHumanName(this.#name);
	}
	get isNameValid() {
		return this.cleanName.length >= 6;
	}

	get cleanUsername() {
		return sanitizeUsername(this.#username);
	}
	get isUsernameValid() {
		return this.cleanUsername.length >= 6;
	}

	get isEmailValid() {
		return isValidEmail(this.#email);
	}

	get hasMinLength() {
		return this.#password.length >= 8;
	}

	get hasUppercase() {
		return /[A-Z]/.test(this.#password);
	}

	get hasLowercase() {
		return /[a-z]/.test(this.#password);
	}

	get hasNumber() {
		return /\d/.test(this.#password);
	}

	get hasSpecialChar() {
		return /[@$!%*?&]/.test(this.#password);
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
		return this.#password === this.#confirmPassword;
	}

	onNameBlur() {
		this.#name = this.cleanName;
	}

	async signUp(): Promise<boolean> {
		if (this._isLoading || !this.validateTurnstile()) return false;

		if (!this.isNameValid) {
			this._apiError = {
				code: "REGISTRATION_INVALID_NAME",
				message: "O nome deve conter pelo menos 6 letras."
			};
			return false;
		}

		if (!this.isUsernameValid) {
			this._apiError = {
				code: "REGISTRATION_INVALID_USERNAME",
				message: "Username deve ter pelo menos 6 caracteres válidos."
			};
			return false;
		}

		if (!this.isEmailValid) {
			this._apiError = {
				code: "REGISTRATION_INVALID_EMAIL",
				message: "Digite um endereço de e-mail válido."
			};
			return false;
		}

		if (!this.isPasswordValid) {
			this._apiError = {
				code: "USER_PASSWORD_NOT_VALID",
				message: SIGN_UP_ERRORS.USER_PASSWORD_NOT_VALID
			};
			return false;
		}

		if (!this.isPasswordsMatching) {
			this._apiError = {
				code: "USER_PASSWORDS_DONT_MATCH",
				message: SIGN_UP_ERRORS.USER_PASSWORDS_DONT_MATCH
			};
			return false;
		}

		this._isLoading = true;

		const { data, error } = await customFetch<SignUpServerResponse>(
			window.fetch,
			"/api/auth/sign-up",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-CF-Turnstile-Response": this._turnstileToken
				},
				body: JSON.stringify({
					name: this.#name,
					username: this.#username,
					email: this.#email,
					password: this.#password,
					confirmPassword: this.#confirmPassword
				})
			},
			SIGN_UP_ERRORS
		);

		this._isLoading = false;

		if (error) {
			this._apiError = error;
			this.resetTurnstile();
			return false;
		}

		if (!data) {
			this._apiError = normalizeApiError("INTERNAL_SERVER_ERROR");
			return false;
		}

		if (data.autoLogin) {
			this._apiError = null;
			await goto("/", { invalidateAll: true });
		} else if (data.success) {
			this._apiError = null;
			await goto("/login");
		} else {
			this._apiError = {
				code: "REGISTRATION_UNEXPECTED_ERROR",
				message: SIGN_UP_ERRORS.REGISTRATION_UNEXPECTED_ERROR
			};
			return false;
		}

		return true;
	}
}
