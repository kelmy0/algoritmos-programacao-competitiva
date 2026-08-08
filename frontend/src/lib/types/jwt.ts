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
	name: string;
	username: string;
	email: string;
	permissions: string[];
	isEmployee: boolean;
	is2FAEnabled: boolean;
}

export interface RefreshJwtPayload extends BaseCustomJwtPayload {
	familyId: string;
	dvh: string;
}

export interface PreAuthToken extends BaseCustomJwtPayload {
	dvh: string;
	is2FAEnabled: boolean;
}

export interface JwtValided<T extends BaseCustomJwtPayload = AccessJwtPayload> {
	claims: T | null;
	valid: boolean;
}
