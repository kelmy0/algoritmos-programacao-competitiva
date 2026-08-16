<script lang="ts">
	import type { Snippet } from "svelte";

	type AlertVariant = "error" | "success" | "warning" | "info";

	interface Props {
		type?: AlertVariant;
		title: string;
		message?: string;
		isLoading?: boolean;
		children?: Snippet;
	}

	let { type = "error", title, message, isLoading = false, children }: Props = $props();

	const config = $derived(
		{
			error: {
				role: "alert",
				ariaLive: "assertive" as const,
				containerClass: "bg-red-950/30 border-red-900/50 text-red-300",
				titleClass: "text-red-200",
				messageClass: "text-red-300/80",
				iconClass: "text-red-400",
				iconPath:
					"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
			},
			success: {
				role: "status",
				ariaLive: "polite" as const,
				containerClass: "bg-emerald-950/30 border-emerald-900/50 text-emerald-300",
				titleClass: "text-emerald-200",
				messageClass: "text-emerald-300/80",
				iconClass: "text-emerald-400",
				iconPath: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
			},
			warning: {
				role: "alert",
				ariaLive: "polite" as const,
				containerClass: "bg-amber-950/30 border-amber-900/50 text-amber-300",
				titleClass: "text-amber-200",
				messageClass: "text-amber-300/80",
				iconClass: "text-amber-400",
				iconPath:
					"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
			},
			info: {
				role: "status",
				ariaLive: "polite" as const,
				containerClass: "bg-sky-950/30 border-sky-900/50 text-sky-300",
				titleClass: "text-sky-200",
				messageClass: "text-sky-300/80",
				iconClass: "text-sky-400",
				iconPath: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			}
		}[type]
	);
</script>

<div
	role={config.role}
	aria-live={config.ariaLive}
	class="p-4 border rounded-xl text-sm flex items-start gap-3 shadow-lg transition-all {config.containerClass} {isLoading
		? 'opacity-50 pointer-events-none'
		: 'opacity-100'}"
>
	<svg
		class="w-5 h-5 shrink-0 mt-0.5 {config.iconClass}"
		fill="none"
		stroke="currentColor"
		viewBox="0 0 24 24"
		aria-hidden="true"
	>
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={config.iconPath} />
	</svg>
	<div class="space-y-0.5">
		<span class="font-semibold block {config.titleClass}">{title}</span>
		{#if children}
			<div class="text-xs {config.messageClass} leading-relaxed">
				{@render children()}
			</div>
		{:else if message}
			<p class="text-xs {config.messageClass} leading-relaxed">
				{message}
			</p>
		{/if}
	</div>
</div>
