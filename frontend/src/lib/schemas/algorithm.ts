import { DIFFICULTIES } from "$lib/types/algorithm";
import { sanitizeTitle } from "$lib/utils/sanitize";
import z from "zod";

const sanitizedString = (minLen: number, errorCode: string) =>
	z
		.string()
		.transform((val) => sanitizeTitle(val))
		.pipe(z.string().min(minLen, errorCode));

export const algorithmSchema = z.object({
	name: sanitizedString(3, "INVALID_NAME"),
	category: sanitizedString(3, "INVALID_CATEGORY"),
	difficulty: z.enum(DIFFICULTIES, "INVALID_DIFFICULTY"),
	content: z.string().trim().min(10, "INVALID_CONTENT")
});

export type AlgorithmPayload = z.infer<typeof algorithmSchema>;
