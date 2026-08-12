export interface LoginServerResponse {
	accessToken: boolean;
	requires2FA: boolean;
	preAuthToken: boolean;
}

export interface LoginResponse {
	accessToken?: string;
	requires2FA: boolean;
	preAuthToken?: string;
}
