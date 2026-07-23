import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { normalizeApiError } from '$lib/utils/errors';
import { AUTH_ERRORS } from './login.svelte';

export async function load({ locals, url }: Parameters<PageServerLoad>[0]) {
	if (locals.user) {
		redirect(303, '/');
	}

	const errorCode = url.searchParams.get('error');

	const initialError = errorCode
		? normalizeApiError(errorCode, 'Erro ao realizar autenticação.', AUTH_ERRORS)
		: null;

	return {
		initialError
	};
}
