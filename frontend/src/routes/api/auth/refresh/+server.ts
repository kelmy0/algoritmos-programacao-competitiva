import { json, type RequestHandler } from "@sveltejs/kit";
import { jwtDecode, type JwtPayload as BaseJwtPayload } from "jwt-decode";
import { normalizeApiError } from "$lib/utils/errors";
import { deleteAuthCookie, setAuthCookie } from "$lib/server/cookies";
import { customFetch } from "$lib/api/client";
import { API_URL } from "$env/static/private";
import { rateLimit, useMiddlewares } from "$lib/server/middlewares";

interface JwtPayload extends BaseJwtPayload {
	username: string;
	email: string;
	permissions?: string[];
	isEmployee: boolean;
}

interface RefreshResponse {
	access_token: string;
}

const refreshToken: RequestHandler = async (event) => {
	const cookieHeader = event.request.headers.get("cookie") || "";
	const clientIp = event.getClientAddress();

	const { data, error, status, headers } = await customFetch<RefreshResponse>(
		event.fetch,
		`${API_URL}/api/auth/refresh`,
		{
			method: "POST",
			headers: { cookie: cookieHeader, "X-Forwarded-For": clientIp }
		}
	);

	if (error) {
		deleteAuthCookie(event.cookies, "access_token");
		return json(error, { status });
	}

	if (!data || !data.access_token) {
		return json(normalizeApiError("INTERNAL_SERVER_ERROR", "Resposta inválida do servidor."), {
			status: 500
		});
	}

	const decoded = jwtDecode<JwtPayload>(data.access_token);
	setAuthCookie(event.cookies, "access_token", data.access_token, 15);

	const response = json({
		accessToken: true,
		expiresAt: decoded.exp ? decoded.exp * 1000 : undefined,
		user: {
			id: decoded.sub,
			username: decoded.username,
			email: decoded.email,
			permissions: decoded.permissions || [],
			is_employee: decoded.isEmployee
		}
	});

	const setCookies = headers.getSetCookie();
	for (const cookieString of setCookies) {
		response.headers.append("set-cookie", cookieString);
	}

	return response;
};

export const POST = useMiddlewares(rateLimit({ capacity: 5, fillRate: 0.1 }))(refreshToken);
