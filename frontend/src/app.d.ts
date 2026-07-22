// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: {
				id: string;
				username: string;
				email: string;
				permissions: string[];
				is_employee: boolean;
			} | null;
			accessToken: string | null;
		}
		interface PageData {
			user?: {
				id: string;
				username: string;
				email: string;
				permissions: string[];
				is_employee: boolean;
			} | null;
			accessToken?: string | null;
		}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
