import { PUBLIC_API_URL } from '$env/static/public';
import type { PageLoad } from './$types';

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
	const response = await fetch(`${PUBLIC_API_URL}/api/algorithms`);

	if (!response.ok) {
		throw new Error(`Erro ao buscar algoritmos: ${response.statusText}`);
	}

	const result: ApiResponse = await response.json();

	return {
		algorithms: result.data,
		pagination: {
			page: result.page,
			limit: result.limit
		}
	};
}
