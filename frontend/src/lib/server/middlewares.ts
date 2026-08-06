import { normalizeApiError } from "$lib/utils/errors";
import { error, json, type RequestEvent, type RequestHandler } from "@sveltejs/kit";
import { redis } from "./redis";

interface TokenBucketOptions {
	capacity: number;
	fillRate: number;
}

export type Middleware = (event: RequestEvent) => Promise<Response | void> | Response | void;

export function useMiddlewares(...middlewares: Middleware[]) {
	return (handler: RequestHandler): RequestHandler => {
		return async (event) => {
			for (const middleware of middlewares) {
				const result = await middleware(event);

				if (result instanceof Response) {
					return result;
				}
			}

			return handler(event);
		};
	};
}

export const requireAuth: Middleware = (event) => {
	if (!event.locals.user) {
		error(401, normalizeApiError("UNAUTHORIZED", "Não autenticado."));
	}
};

export function requirePermission(permission?: string): Middleware {
	return (event) => {
		const user = event.locals.user;

		if (!user?.isEmployee) {
			error(404, normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada."));
		}

		if (permission && !user.permissions?.includes(permission)) {
			error(404, normalizeApiError("PAGE_NOT_FOUND", "Página não encontrada."));
		}
	};
}

export const createAlgorithms = requirePermission("create:algorithms");
export const moderateAlgorithms = requirePermission("moderate:algorithms");

const TOKEN_BUCKET_LUA = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local fill_rate_per_ms = tonumber(ARGV[2]) / 1000.0
local now_ms = tonumber(ARGV[3])

local fill_time_seconds = math.ceil(capacity / (fill_rate_per_ms * 1000.0))
local ttl_seconds = math.max(10, fill_time_seconds * 1.2)

local data = redis.call("HMGET", key, "tokens", "last_updated")
local tokens = tonumber(data[1])
local last_updated = tonumber(data[2])

if not tokens then
    tokens = capacity
    last_updated = now_ms
else
    local delta_ms = math.max(0, now_ms - last_updated)
    tokens = math.min(capacity, tokens + (delta_ms * fill_rate_per_ms))
    last_updated = now_ms
end

if tokens >= 1 then
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "last_updated", last_updated)
    redis.call("EXPIRE", key, ttl_seconds)
    return {1, math.floor(tokens)}
else
    redis.call("HMSET", key, "tokens", tokens, "last_updated", last_updated)
    redis.call("EXPIRE", key, ttl_seconds)
    return {0, math.floor(tokens)}
end
`;

export function rateLimit(options: TokenBucketOptions = { capacity: 60, fillRate: 1 }): Middleware {
	return async (event) => {
		try {
			const userId = event.locals.user?.id;
			let clientIp = "127.0.0.1";

			try {
				clientIp = event.getClientAddress();
			} catch {
				clientIp =
					event.request.headers.get("x-forwarded-for")?.split(",")[0].trim() || "127.0.0.1";
			}

			const identifier = userId ? `usr_${userId}` : `ip_${clientIp}`;
			const redisKey = `ratelimit:tb:${event.url.pathname}:${identifier}`;
			const nowMs = Date.now();

			const result = (await redis.eval(
				TOKEN_BUCKET_LUA,
				1,
				redisKey,
				options.capacity.toString(),
				options.fillRate.toString(),
				nowMs.toString()
			)) as [number, number];

			const [allowed, remainingTokens] = result;

			event.setHeaders({
				"X-RateLimit-Limit": options.capacity.toString(),
				"X-RateLimit-Remaining": Math.max(0, remainingTokens).toString()
			});

			if (allowed === 0) {
				const normalizedError = normalizeApiError("TOO_MANY_REQUESTS");
				const retryAfterSeconds = Math.max(1, Math.ceil(1 / options.fillRate));

				return json(normalizedError, {
					status: 429,
					headers: {
						"Retry-After": retryAfterSeconds.toString(),
						"X-RateLimit-Remaining": "0"
					}
				});
			}
		} catch (err) {
			console.error("[RateLimit Error - Redis offline?]:", err);
		}
	};
}

export const standardApiLimiter = rateLimit({ capacity: 5, fillRate: 5 });
export const authFlowLimiter = rateLimit({ capacity: 5, fillRate: 0.1 });
export const strictAbuseLimiter = rateLimit({ capacity: 2, fillRate: 0.0055 });

export function limitBodySize(maxBytes: number): Middleware {
	return async (event) => {
		const contentLength = event.request.headers.get("content-length");

		if (contentLength) {
			const length = parseInt(contentLength, 10);
			if (!isNaN(length) && length > maxBytes) {
				const normalizedError = normalizeApiError("PAYLOAD_TOO_LARGE");

				return json(normalizedError, { status: 413 });
			}
		}

		if (event.request.body) {
			try {
				const clonedRequest = event.request.clone();
				const buffer = await clonedRequest.arrayBuffer();

				if (buffer.byteLength > maxBytes) {
					const normalizedError = normalizeApiError("PAYLOAD_TOO_LARGE");

					return json(normalizedError, { status: 413 });
				}
			} catch (err) {
				console.error("[LimitBodySize Error]:", err);
			}
		}
	};
}

export const hundredKbBodySize = limitBodySize(1024 * 128);
export const tenMbBodySize = limitBodySize(10 * 1024 * 1024);

export function limitQueryParamsSize(maxChars: number = 2048): Middleware {
	return async (event: RequestEvent) => {
		const fullUrl = event.url.pathname + event.url.search;

		if (fullUrl.length > maxChars) {
			const normalizedError = normalizeApiError("URL_TOO_LARGE");
			return json(normalizedError, { status: 400 });
		}
	};
}

export const fiveHundredQuerySize = limitQueryParamsSize(512);
export const thousandQuerySize = limitQueryParamsSize(1024);
export const twoThousandUrlSize = limitQueryParamsSize(2048);
