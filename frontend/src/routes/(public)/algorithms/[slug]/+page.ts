import { PUBLIC_API_URL } from '$env/static/public';
import { normalizeApiError } from '$lib/utils/errors';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { ALGORITHMS_ERRORS } from '../algorithms';
import { renderMarkdown } from '$lib/services/markdown';

export const load: PageLoad = async ({ fetch, params }) => {
	const { slug } = params;

	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/algorithms/${slug}`);

		if (!response.ok) {
			const rawError = await response.json().catch(() => null);

			const normalized = normalizeApiError(
				rawError,
				'Não foi possível carregar o algoritmo solicitado.',
				ALGORITHMS_ERRORS
			);

			error(response.status, normalized);
		}

		const { data: algorithm } = await response.json();

		const contentHtml = await renderMarkdown(algorithm.Content || '');

		return {
			algorithm: {
				...algorithm,
				contentHtml
			}
		};
	} catch (err) {
		if (err && typeof err === 'object' && 'status' in err) {
			throw err;
		}

		const normalized = normalizeApiError(
			err,
			'Não foi possível conectar ao servidor. Verifique sua conexão.',
			ALGORITHMS_ERRORS
		);

		error(500, normalized);
	}
};
