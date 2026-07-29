import { customFetch } from "$lib/api/client";

interface ActivityKeeperOptions {
	intervalMinutes?: number;
	onUnauthorized?: () => void;
}

export function createActivityKeeper({
	intervalMinutes = 5,
	onUnauthorized
}: ActivityKeeperOptions = {}) {
	let timer: ReturnType<typeof setInterval> | null = null;
	let userHasInteracted = false;

	function onUserActivity() {
		userHasInteracted = true;
	}

	function start() {
		if (typeof window === "undefined") return;

		window.addEventListener("keydown", onUserActivity, { passive: true });
		window.addEventListener("click", onUserActivity, { passive: true });

		timer = setInterval(
			async () => {
				if (userHasInteracted) {
					userHasInteracted = false;
					const { status, error } = await customFetch(fetch, "/api/admin/session");
					if (status === 401 || error) {
						onUnauthorized?.();
					}
				}
			},
			intervalMinutes * 60 * 1000
		);
	}

	function stop() {
		if (typeof window === "undefined") return;

		if (timer) clearInterval(timer);
		window.removeEventListener("keydown", onUserActivity);
		window.removeEventListener("click", onUserActivity);
	}

	return { start, stop };
}
