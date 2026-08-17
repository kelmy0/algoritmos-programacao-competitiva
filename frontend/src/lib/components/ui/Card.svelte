<script lang="ts">
	import type { Snippet } from "svelte";
	import { fade } from "svelte/transition";

	interface Props {
		title: string;
		href: string;
		description?: string;
		footerLeft?: Snippet;
		actionLabel?: string;
		icon?: Snippet;
		headerRightTop?: Snippet;
		headerRightBottom?: Snippet;
		children?: Snippet;
	}

	let {
		title,
		href,
		description,
		footerLeft,
		actionLabel = "Ver detalhes",
		icon,
		headerRightTop,
		headerRightBottom,
		children
	}: Props = $props();
</script>

<article
	class="relative bg-app-surface border border-app-border rounded-xl p-5 shadow-lg flex flex-col justify-between hover:border-app-border-hover hover:shadow-xl transition-all duration-200 group"
>
	<div class="space-y-3">
		{#if icon}
			<div class="flex items-center justify-between gap-3">
				<div
					class="p-2.5 rounded-lg bg-app-bg/60 border border-app-border text-text-brand group-hover:border-text-brand/40 transition-colors flex items-center justify-center shrink-0"
				>
					{@render icon()}
				</div>

				{#if headerRightTop}
					<div class="shrink-0 relative z-10 flex items-center">
						{@render headerRightTop()}
					</div>
				{/if}
			</div>

			<div>
				<h2
					class="font-montserrat font-semibold text-lg text-text-primary group-hover:text-text-brand transition-colors leading-tight"
					{title}
				>
					<a
						{href}
						class="after:absolute after:inset-0 focus:outline-none focus:ring-2 focus:ring-text-brand focus:ring-offset-2 focus:ring-offset-app-surface rounded-xl"
					>
						{title}
					</a>
				</h2>

				{#if description}
					<p class="text-xs text-text-secondary mt-2 leading-relaxed">
						{description}
					</p>
				{/if}

				{#if children}
					<div class="mt-3">
						{@render children()}
					</div>
				{/if}
			</div>
		{:else}
			<div class="flex items-center justify-between gap-3 min-h-8">
				<h2
					class="font-montserrat font-semibold text-base text-text-primary group-hover:text-text-brand transition-colors line-clamp-2 leading-snug"
					{title}
				>
					<a
						{href}
						class="after:absolute after:inset-0 focus:outline-none focus:ring-2 focus:ring-text-brand focus:ring-offset-2 focus:ring-offset-app-surface rounded-xl"
					>
						{title}
					</a>
				</h2>

				{#if headerRightTop}
					<div class="shrink-0 relative z-10 flex items-center">
						{@render headerRightTop()}
					</div>
				{/if}
			</div>

			{#if children || headerRightBottom || description}
				<div class="flex items-center justify-between gap-3 min-h-7">
					<div class="flex-1 flex items-center">
						{#if description}
							<p class="text-xs text-text-secondary leading-relaxed line-clamp-2">
								{description}
							</p>
						{/if}

						{#if children}
							{@render children()}
						{/if}
					</div>

					{#if headerRightBottom}
						<div class="shrink-0 relative z-10 flex items-center">
							{@render headerRightBottom()}
						</div>
					{/if}
				</div>
			{/if}
		{/if}
	</div>

	<div class="pt-4 mt-4 border-t border-app-border flex items-center justify-between">
		<div class="relative z-10 flex items-center">
			{#if footerLeft}
				{@render footerLeft()}
			{/if}
		</div>

		<div
			aria-hidden="true"
			class="text-xs font-medium text-text-brand group-hover:underline flex items-center gap-1 transition-all pointer-events-none ml-auto"
		>
			<span>{actionLabel}</span>
			<svg
				class="w-4 h-4 group-hover:translate-x-0.5 transition-transform shrink-0"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
			>
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
			</svg>
		</div>
	</div>
</article>
