import { PUBLIC_API_URL } from '$env/static/public';
import { normalizeApiError } from '$lib/utils/errors';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { ApiError } from '$lib/types/api';

interface Algorithm {
	Id: string;
	PublicId: string;
	Slug: string;
	Name: string;
	Category: string;
	Difficulty: string;
}

interface ApiResponse {
	page: number;
	limit: number;
	data: Algorithm[];
}

export async function load({ fetch }: Parameters<PageLoad>[0]) {
	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/algorithms`);

		if (!response.ok) {
			const errorData: ApiError = await response.json().catch(() => ({}));
			const apiError = normalizeApiError(
				errorData,
				'Não foi possível carregar a lista de algoritmos.'
			);
			error(response.status, {
				message: apiError.message,
				code: apiError.code
			});
		}

		const result: ApiResponse = await response.json();

		return {
			algorithms: result.data,
			pagination: {
				page: result.page,
				limit: result.limit
			}
		};
	} catch (err) {
		if (typeof err === 'object' && err !== null && 'status' in err) {
			throw err;
		}

		const apiError = normalizeApiError(err, 'Erro ao conectar ao servidor.');

		error(500, {
			message: apiError.message,
			code: apiError.code
		});
	}
}
