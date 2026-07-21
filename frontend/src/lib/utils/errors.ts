export const GLOBAL_ERRORS: Record<string, string> = {
	// Network and Operating System (Captured by the Fetch chat)
	NETWORK_ERROR: 'Não foi possível conectar ao servidor. Verifique sua conexão com a internet.',

	// General Server Errors
	INTERNAL_SERVER_ERROR:
		'Ocorreu um erro interno inesperado no servidor. Tente novamente mais tarde.',
	INVALID_REQUEST_BODY: 'Os dados enviados estão corrompidos ou em formato inválido.',
	PAGE_NOT_FOUND: 'A página ou recurso solicitado não foi encontrado.',
	TOO_MANY_REQUESTS: 'Você fez muitas requisições seguidas. Aguarde um momento e tente novamente.',

	// Generic Authentication and Session
	SESSION_EXPIRED: 'Sua sessão expirou ou ocorreu uma falha de segurança. Faça login novamente.',
	MISSING_USER_ID: 'Identificação do usuário ausente na requisição.',
	MISSING_COOKIE: 'Cookie de autenticação obrigatório está ausente.',
	MISSING_HEADER: 'Cabeçalho de segurança obrigatório não foi enviado.',
	INVALID_HEADER_FORMAT: 'O formato do cabeçalho de autenticação está incorreto.',
	INVALID_ACCESS_TOKEN: 'Seu token de acesso é inválido ou expirou.',
	RESTRICTED_AREA: 'Esta é uma área restrita do sistema.',
	NO_PERMISSION: 'Você não tem permissão para realizar esta ação.',

	// OAuth Integrations (Google and GitHub)
	MISSING_OAUTH_CODE: 'Código de autenticação do provedor social ausente.',
	MISSING_TOKEN_ID: 'Token ID de validação social ausente.',
	INVALID_GOOGLE_TOKEN: 'O token de autenticação do Google é inválido ou expirou.',
	MISSING_GOOGLE_EMAIL: 'Não foi possível obter seu e-mail através da conta do Google.',
	UNVERIFIED_GOOGLE_EMAIL:
		'Sua conta do Google precisa ter o e-mail verificado para ser utilizada.',
	GITHUB_EMAIL_UNVERIFIED: 'Sua conta do GitHub precisa ter o e-mail verificado para ser utilizada.'
};

export function getErrorMessage(
	code: string,
	fallbackMessage: string,
	localErrors?: Record<string, string>
): string {
	if (localErrors && code in localErrors) return localErrors[code];
	if (code in GLOBAL_ERRORS) return GLOBAL_ERRORS[code];

	return fallbackMessage || 'Ocorreu um erro desconhecido.';
}
