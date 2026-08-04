export interface SignUpServerResponse {
	success: boolean;
	autoLogin: boolean;
}

export interface SignUpResponse {
	accessToken?: string;
	success: boolean;
	autoLogin: boolean;
}
