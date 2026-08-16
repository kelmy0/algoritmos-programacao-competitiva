export const DIFFICULTIES = ["beginner", "intermediate", "advanced", "expert"] as const;
export const STATUS = ["pending", "approved", "rejected", "deleted"] as const;

export type Difficulty = (typeof DIFFICULTIES)[number];
export type Status = (typeof STATUS)[number];

export interface Algorithm {
	publicId: string;
	slug: string;
	name: string;
	category: string;
	content: string;
	difficulty: Difficulty;
	status: Status;
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
