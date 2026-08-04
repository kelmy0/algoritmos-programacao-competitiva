export interface LoginServerResponse {
	accessToken: boolean;
	requires2FA: boolean;
	preAuthToken?: string;
}

export interface LoginResponse {
	accessToken?: string;
	requires2FA: boolean;
	preAuthToken?: string;
}
