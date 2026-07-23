<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const algorithm = $derived(data.algorithm);
</script>

<div class="mx-auto max-w-4xl space-y-6 font-inter">
	<a
		href="/"
		class="inline-flex items-center gap-2 text-xs font-medium text-gray-300 transition-colors hover:text-white focus:outline-none focus:ring-2 focus:ring-brand focus:ring-offset-2 focus:ring-offset-gray-900"
	>
		<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
		</svg>
		<span>Voltar para o início</span>
	</a>

	{#if algorithm}
		<article
			class="bg-app-surface border-gray-800 space-y-6 rounded-xl border p-6 shadow-xl md:p-8"
		>
			<header class="border-gray-800 space-y-3 border-b pb-6">
				<div class="flex flex-wrap items-center justify-between gap-2">
					<span
						class="text-text-brand bg-text-brand/10 border-text-brand/20 rounded-md border px-2.5 py-1 font-mono text-xs"
						aria-label="Categoria: {algorithm.Category}"
					>
						{algorithm.Category}
					</span>
					<span class="font-mono text-xs text-gray-400">
						<span class="sr-only">Identificador: </span>ID: {algorithm.PublicId}
					</span>
				</div>

				<h1 class="text-text-primary font-montserrat text-3xl font-bold tracking-tight">
					{algorithm.Name}
				</h1>
			</header>

			<section class="space-y-4" aria-labelledby="heading-implementacao">
				<h2
					id="heading-implementacao"
					class="text-sm font-semibold uppercase tracking-wider text-gray-200"
				>
					Implementação
				</h2>

				<div class="prose prose-invert max-w-none font-mono text-sm text-gray-200">
					{@html algorithm.contentHtml}
				</div>
			</section>
		</article>
	{/if}
</div>
