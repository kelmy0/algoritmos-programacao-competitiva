import type { JwtPayload as BaseJwtPayload } from "jwt-decode";

export interface RefreshServerResponse {
	accessToken: boolean;
	expiresAt?: number;
}

export interface RefreshResponse {
	accessToken: string;
}

export interface JwtPayload extends BaseJwtPayload {
	sub: string;
	username: string;
	email: string;
	permissions: string[];
	isEmployee: boolean;
	exp?: number;
}
