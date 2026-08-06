import { API_DOMAIN } from "$env/static/private";
import type { AccessJwtPayload, BaseCustomJwtPayload, JwtValided } from "$lib/types/jwt";
import { importSPKI, jwtVerify } from "jose";

const publicKeyCache = new Map<string, CryptoKey>();

async function getPublicKey(publicKeyB64: string): Promise<CryptoKey> {
	if (publicKeyCache.has(publicKeyB64)) {
		return publicKeyCache.get(publicKeyB64)!;
	}

	const pem = Buffer.from(publicKeyB64, "base64").toString("utf-8");
	const cryptoKey = await importSPKI(pem, "EdDSA");

	publicKeyCache.set(publicKeyB64, cryptoKey);
	return cryptoKey;
}

export async function validateJWT<T extends BaseCustomJwtPayload = AccessJwtPayload>(
	publicKeyB64: string,
	token: string
): Promise<JwtValided<T>> {
	if (!token) {
		return { claims: null, valid: false };
	}

	try {
		const publicKey = await getPublicKey(publicKeyB64);

		const { payload } = await jwtVerify(token, publicKey, {
			issuer: API_DOMAIN
		});

		return {
			claims: payload as unknown as T,
			valid: true
		};
	} catch (error) {
		return { claims: null, valid: false };
	}
}
