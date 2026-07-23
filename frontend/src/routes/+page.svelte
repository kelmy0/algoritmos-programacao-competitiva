<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const difficultyStyles: Record<string, string> = {
		beginner: 'bg-emerald-950 text-emerald-300 border-emerald-800',
		intermediate: 'bg-amber-950 text-amber-300 border-amber-800',
		advanced: 'bg-orange-950 text-orange-300 border-orange-800',
		expert: 'bg-red-950 text-red-300 border-red-800'
	};

	const difficultyLabels: Record<string, string> = {
		beginner: 'Iniciante',
		intermediate: 'Intermediário',
		advanced: 'Avançado',
		expert: 'Especialista'
	};
</script>

<svelte:head>
	<title>Algoritmos para programação competitiva</title>

	<meta
		name="description"
		content="Coleção de algoritmos e estruturas de dados otimizados para maratonas de programação e competições."
	/>
	<meta
		name="keywords"
		content="algoritmos, maratona de programação, c++, estruturas de dados, competitiva"
	/>

	<meta property="og:type" content="website" />
	<meta property="og:title" content="Algoritmos para programação competitiva" />
	<meta
		property="og:description"
		content="Explore a coleção de algoritmos e estruturas de dados com implementações prontas para uso em competições."
	/>
	<!--<meta property="og:image" content="/og-image.png" />-->

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="Algoritmos para programação competitiva" />
	<meta
		name="twitter:description"
		content="Coleção de algoritmos e estruturas de dados para programação competitiva."
	/>
	<meta name="twitter:image" content="/og-image.png" />
</svelte:head>

<div class="space-y-6 font-inter">
	<header
		class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-gray-800"
	>
		<div>
			<h1 class="font-montserrat font-bold text-2xl md:text-3xl text-text-primary tracking-tight">
				Algoritmos
			</h1>
			<p class="text-sm text-gray-400 mt-1">
				Explore a coleção de algoritmos e estruturas de dados.
			</p>
		</div>
	</header>

	<!-- Algorithm Grid -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
		{#each data.algorithms as item (item.Id)}
			<article
				class="relative bg-app-surface border border-gray-800 rounded-xl p-5 shadow-lg flex flex-col justify-between hover:border-gray-700 hover:shadow-xl transition-all duration-200 group"
			>
				<div class="space-y-3">
					<div class="flex items-start justify-between gap-3">
						<h2
							class="font-montserrat font-semibold text-base text-text-primary group-hover:text-text-brand transition-colors line-clamp-2"
							title={item.Name}
						>
							<a
								href="/algorithms/{item.Slug}-{item.PublicId}"
								class="after:absolute after:inset-0 focus:outline-none focus:ring-2 focus:ring-text-brand focus:ring-offset-2 focus:ring-offset-app-surface rounded-xl"
							>
								{item.Name}
							</a>
						</h2>

						<span
							class="text-xs font-semibold uppercase tracking-wider px-2.5 py-1 rounded-md border shrink-0 relative z-10 {difficultyStyles[
								item.Difficulty
							] ?? 'bg-gray-800 text-gray-300 border-gray-700'}"
						>
							{difficultyLabels[item.Difficulty] ?? item.Difficulty}
						</span>
					</div>

					<p class="text-xs font-medium text-gray-400 flex items-center gap-2">
						<span class="inline-block w-2 h-2 rounded-full bg-text-brand" aria-hidden="true"></span>
						{item.Category}
					</p>
				</div>

				<div class="pt-4 mt-4 border-t border-gray-800/80 flex items-center justify-between">
					<span
						class="text-xs font-mono text-gray-400 truncate max-w-30 relative z-10 select-all"
						title="Copiar ID"
					>
						{item.PublicId}
					</span>

					<div
						aria-hidden="true"
						class="text-xs font-medium text-text-brand group-hover:underline flex items-center gap-1 transition-all pointer-events-none"
					>
						<span>Ver detalhes</span>
						<svg
							class="w-4 h-4 group-hover:translate-x-0.5 transition-transform"
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
		{:else}
			<div
				role="status"
				class="col-span-full py-12 text-center bg-app-surface border border-gray-800 rounded-xl"
			>
				<p class="text-gray-400 text-sm">Nenhum algoritmo encontrado.</p>
			</div>
		{/each}
	</div>
</div>
