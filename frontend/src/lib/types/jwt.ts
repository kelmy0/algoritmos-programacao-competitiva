import type { JWTPayload } from "jose";

export interface RefreshServerResponse {
	accessToken: boolean;
	expiresAt?: number;
}

export interface RefreshResponse {
	accessToken: string;
}

export interface BaseCustomJwtPayload extends JWTPayload {
	jti: string;
	sub: string;
	exp: number;
	iss: string;
}

export interface AccessJwtPayload extends BaseCustomJwtPayload {
	username: string;
	email: string;
	permissions: string[];
	isEmployee: boolean;
}

export interface RefreshJwtPayload extends BaseCustomJwtPayload {
	familyId: string;
}

export interface PreAuthToken extends BaseCustomJwtPayload {}

export interface JwtValided<T extends BaseCustomJwtPayload = AccessJwtPayload> {
	claims: T | null;
	valid: boolean;
}
