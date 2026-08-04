import type { Handle } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";
import { redis } from "../redis";
import { normalizeApiError } from "$lib/utils/errors";

const TOKEN_BUCKET_LUA = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local fill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = math.ceil(capacity / fill_rate)

local data = redis.call("HMGET", key, "tokens", "last_updated")
local tokens = tonumber(data[1])
local last_updated = tonumber(data[2])

if not tokens then
    tokens = capacity
    last_updated = now
else
    local delta = math.max(0, now - last_updated)
    tokens = math.min(capacity, tokens + (delta * fill_rate))
    last_updated = now
end

if tokens >= 1 then
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "last_updated", last_updated)
    redis.call("EXPIRE", key, ttl)
    return {1, math.floor(tokens)}
else
    redis.call("EXPIRE", key, ttl)
    return {0, 0}
end
`;

export const handleGlobalRateLimit: Handle = async ({ event, resolve }) => {
	if (event.url.pathname.startsWith("/api")) {
		try {
			const userId = event.locals.user?.id;
			const identifier = userId ? `usr_${userId}` : `ip_${event.getClientAddress()}`;

			const redisKey = `global_ratelimit:tb:${identifier}`;
			const now = Math.floor(Date.now() / 1000);

			const CAPACITY = "120";
			const FILL_RATE = "2";

			const result = (await redis.eval(
				TOKEN_BUCKET_LUA,
				1,
				redisKey,
				CAPACITY,
				FILL_RATE,
				now.toString()
			)) as [number, number];

			const [allowed, remainingTokens] = result;

			event.setHeaders({
				"X-Global-RateLimit-Limit": CAPACITY,
				"X-Global-RateLimit-Remaining": remainingTokens.toString()
			});

			if (allowed === 0) {
				const normalizedError = normalizeApiError("TOO_MANY_REQUESTS");

				return json(normalizedError, {
					status: 429,
					headers: {
						"Retry-After": "1"
					}
				});
			}
		} catch (err) {
			console.error("[Global RateLimit Error - Redis Offline?]:", err);
		}
	}

	return resolve(event);
};
