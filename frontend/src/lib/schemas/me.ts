import z from "zod";

const sanitizedPassword = (minLen: number, errorCode: string) =>
	z.union([z.literal(""), z.string().min(minLen, errorCode)]);

const sanitized2FACode = (len: number, errorCode: string) =>
	z
		.string()
		.transform((v) => v.replaceAll(/\D/g, ""))
		.pipe(z.string().length(len, errorCode));

export const generate2FASchema = z.object({
	password: sanitizedPassword(8, "INVALID_PASSWORD")
});

export const enable2FASchema = z.object({
	code: sanitized2FACode(6, "INVALID_2FA_CODE")
});

export const disable2FASchema = generate2FASchema;
