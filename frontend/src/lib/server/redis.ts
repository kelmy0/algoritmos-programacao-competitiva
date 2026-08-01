import { REDIS_ADDR, REDIS_PASSWORD, REDIS_TLS } from "$env/static/private";
import Redis from "ioredis";

const redisAddr = REDIS_ADDR;
const [host, portStr] = redisAddr.split(":");
const port = parseInt(portStr, 10);
const password = REDIS_PASSWORD || undefined;
const useTLS = REDIS_TLS !== "false";

export const redis = new Redis({
	host,
	port,
	password,
	tls: useTLS ? {} : undefined,
	maxRetriesPerRequest: 3
});

redis.on("error", (err) => {
	console.error("❌ [Redis Error]:", err);
});
