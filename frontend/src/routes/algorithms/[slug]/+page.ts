import { PUBLIC_API_URL } from '$env/static/public';
import { normalizeApiError } from '$lib/utils/errors';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { ALGORITHMS_ERRORS } from '../algorithms';
import { Marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
import { createHighlighter } from 'shiki';

const highlighter = await createHighlighter({
	themes: ['github-dark'],
	langs: ['cpp']
});

const marked = new Marked(
	markedHighlight({
		async: true,
		highlight(code, lang) {
			const language = highlighter.getLoadedLanguages().includes(lang) ? lang : 'text';
			return highlighter.codeToHtml(code, {
				lang: language,
				theme: 'github-dark'
			});
		}
	})
);

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
		const contentHtml = await marked.parse(algorithm.Content || '');

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
