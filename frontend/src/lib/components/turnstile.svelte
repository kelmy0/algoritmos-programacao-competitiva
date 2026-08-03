<script lang="ts">
	import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";

	let { onsuccess, onexpire }: { onsuccess?: (token: string) => void; onexpire?: () => void } =
		$props();

	let container: HTMLDivElement;
	let widgetId: string | null = null;

	export function reset() {
		if (widgetId && window.turnstile) {
			window.turnstile.reset(widgetId);
		}
	}

	$effect(() => {
		if (!container) return;

		const renderWidget = () => {
			if (window.turnstile && container && !widgetId) {
				widgetId = window.turnstile.render(container, {
					sitekey: PUBLIC_TURNSTILE_SITE_KEY,
					theme: "dark",
					callback: (token: string) => {
						onsuccess?.(token);
					},
					"expired-callback": () => {
						onexpire?.();
					}
				});
			}
		};

		if (window.turnstile) {
			renderWidget();
		} else {
			const interval = setInterval(() => {
				if (window.turnstile) {
					clearInterval(interval);
					renderWidget();
				}
			}, 100);
		}

		return () => {
			if (widgetId && window.turnstile) {
				window.turnstile.remove(widgetId);
				widgetId = null;
			}
		};
	});
</script>

<div
	bind:this={container}
	class="my-4 flex h-16.5 w-75 shrink-0 overflow-hidden bg-gray-800"
	aria-label="Verificação de segurança Turnstile"
></div>
