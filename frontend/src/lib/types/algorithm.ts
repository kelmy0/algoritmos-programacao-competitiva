export interface Algorithm {
	publicId: string;
	slug: string;
	name: string;
	category: string;
	content: string;
	difficulty: string;
	status: string;
	authorId: string;
	createdAt: Date;
	updatedAt: Date;
}

export interface ListAlgorithmsResponse {
	page: number;
	limit: number;
	hasMore: boolean;
	algorithms: Algorithm[];
}

export type AlgorithmDetailResponse = {
	algorithm: Algorithm & {
		contentHtml: string;
	};
};
