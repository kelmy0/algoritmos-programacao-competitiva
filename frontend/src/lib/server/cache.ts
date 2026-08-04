type CacheEntry<T> = {
	data: T;
	headers: Record<string, string>;
	expiresAt: number;
};

class MemoryCache {
	private cache = new Map<string, CacheEntry<any>>();
	private readonly maxItems: number;

	constructor(maxItems = 100) {
		this.maxItems = maxItems;
	}

	get<T>(key: string): CacheEntry<T> | null {
		const entry = this.cache.get(key);
		if (!entry) return null;

		if (Date.now() > entry.expiresAt) {
			this.cache.delete(key);
			return null;
		}

		return entry;
	}

	set<T>(key: string, data: T, headers: Record<string, string>, ttlSeconds: number): void {
		if (this.cache.size >= this.maxItems && !this.cache.has(key)) {
			const oldestKey = this.cache.keys().next().value;
			if (oldestKey) this.cache.delete(oldestKey);
		}

		this.cache.set(key, {
			data,
			headers,
			expiresAt: Date.now() + ttlSeconds * 1000
		});
	}

	delete(key: string): void {
		this.cache.delete(key);
	}
}

export const svelteServerCache = new MemoryCache(100);
