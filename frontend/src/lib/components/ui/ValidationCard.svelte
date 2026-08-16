<script lang="ts">
	export interface Requirement {
		label: string;
		met: boolean;
	}

	interface Props {
		title?: string;
		requirements: Requirement[];
	}

	let { title = "Requisitos mínimos:", requirements }: Props = $props();
</script>

<div
	aria-live="polite"
	class="p-3 bg-app-bg/30 border border-gray-800/80 rounded-lg space-y-1.5 text-xs mt-2 transition-all"
>
	{#if title}
		<p class="font-medium text-gray-400 mb-1">{title}</p>
	{/if}

	{#each requirements as req}
		<div
			class="flex items-center gap-2 transition-colors {req.met
				? 'text-emerald-400'
				: 'text-gray-400'}"
		>
			<svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				{#if req.met}
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="3"
						d="M5 13l4 4L19 7"
					/>
				{:else}
					<circle cx="12" cy="12" r="3" fill="currentColor" />
				{/if}
			</svg>
			<span>{req.label}</span>
		</div>
	{/each}
</div>
