// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { Turnstile } from "cloudflare-turnstile";

declare global {
	interface Window {
		turnstile?: Turnstile.TurnstileObject;
	}

	namespace App {
		interface Error {
			message: string;
			code?: string;
		}
		interface Locals {
			user: {
				id: string;
				name: string;
				username: string;
				email: string;
				permissions: string[];
				isEmployee: boolean;
				is2FAEnabled: boolean;
			} | null;
			accessToken: string | null;
		}
		interface PageData {
			user?: {
				id: string;
				name: string;
				username: string;
				email: string;
				permissions: string[];
				isEmployee: boolean;
				is2FAEnabled: boolean;
			} | null;
			accessToken?: string | null;
		}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
