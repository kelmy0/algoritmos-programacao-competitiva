import type { ApiError } from '$lib/types/api';

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
	USER_NOT_FOUND: 'Usuário não encontrado.',
	SESSION_EXPIRED: 'Sua sessão expirou ou ocorreu uma falha de segurança. Faça login novamente.',
	MISSING_USER_ID: 'Identificação do usuário ausente na requisição.',
	MISSING_COOKIE: 'Cookie de autenticação obrigatório está ausente.',
	MISSING_HEADER: 'Cabeçalho de segurança obrigatório não foi enviado.',
	INVALID_HEADER_FORMAT: 'O formato do cabeçalho de autenticação está incorreto.',
	INVALID_ACCESS_TOKEN: 'Seu token de acesso é inválido ou expirou.',
	RESTRICTED_AREA: 'Esta é uma área restrita do sistema.',
	NO_PERMISSION: 'Você não tem permissão para realizar esta ação.',
	AUTH_UNEXPECTED_ERROR: 'Falha inesperada no login. Tente novamente!',

	// Two Factor
	'2FA_NOT_INITIATED': 'A autenticação em dois fatores não está ativa.',
	'2FA_INVALID_CODE': 'Código inválido ou expirado! Caso o erro persista, faça login novamente.',
	MISSING_PRE_TOKEN: 'Está faltando o token da autenticação em dois fatores',

	// OAuth Integrations (Google and GitHub)
	MISSING_OAUTH_CODE: 'Código de autenticação do provedor social ausente.',
	MISSING_TOKEN_ID: 'Token ID de validação social ausente.',
	INVALID_GOOGLE_TOKEN: 'O token de autenticação do Google é inválido ou expirou.',
	MISSING_GOOGLE_EMAIL: 'Não foi possível obter seu e-mail através da conta do Google.',
	UNVERIFIED_GOOGLE_EMAIL:
		'Sua conta do Google precisa ter o e-mail verificado para ser utilizada.',
	GITHUB_EMAIL_UNVERIFIED:
		'Sua conta do GitHub precisa ter o e-mail verificado para ser utilizada.',
	EMAIL_MISMATCH_SOCIAL_LINK: 'Os e-mails precisam ser iguais para vincular na mesma conta.',
	LINK_SOCIAL_ACCOUNT_FAILED: 'Erro ao vincular contas. Tente novamente!',

	// Internal User Operation Failures
	AUTH_QUERY_USER_FAILED: 'Erro interno ao consultar dados cadastrais.'
};

export function normalizeApiError(
	error: unknown,
	fallbackMessage = 'Ocorreu um erro inesperado.',
	localErrors?: Record<string, string>
): ApiError {
	if (error instanceof TypeError && error.message.includes('fetch')) {
		return {
			code: 'NETWORK_ERROR',
			message: GLOBAL_ERRORS.NETWORK_ERROR
		};
	}

	if (typeof error === 'string') {
		const translatedMessage = localErrors?.[error] ?? GLOBAL_ERRORS[error] ?? fallbackMessage;

		return {
			code: error,
			message: translatedMessage
		};
	}

	if (typeof error === 'object' && error !== null && 'code' in error) {
		const apiErr = error as { code: string; message?: string };
		const code = apiErr.code;

		const translatedMessage = localErrors?.[code] ?? GLOBAL_ERRORS[code] ?? fallbackMessage;

		return {
			code,
			message: translatedMessage
		};
	}

	return {
		code: 'UNEXPECTED_ERROR',
		message: fallbackMessage
	};
}
