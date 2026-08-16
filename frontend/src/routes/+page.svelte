<script lang="ts">
	import { browser } from "$app/environment";
	import type { Algorithm } from "$lib/types/algorithm";
	import { untrack } from "svelte";
	import type { PageData } from "./$types";
	import { DIFFICULTY_MAP } from "$lib/constants/algorithms";
	import Card from "$lib/components/ui/Card.svelte";

	let { data }: { data: PageData } = $props();

	let algorithms = $state<Algorithm[]>(untrack(() => data.algorithms ?? []));
	let page = $state(untrack(() => data.pagination?.page ?? 1));
	let hasMore = $state(untrack(() => data.pagination?.hasMore ?? false));
	let isLoading = $state(false);
	let sentinel = $state<HTMLElement | null>(null);

	$effect(() => {
		algorithms = data.algorithms ?? [];
		page = data.pagination?.page ?? 1;
		hasMore = data.pagination?.hasMore ?? false;
	});

	async function loadMore() {
		if (!browser || isLoading || !hasMore) return;

		isLoading = true;
		const nextPage = page + 1;
		const limit = data.pagination?.limit ?? 12;

		try {
			const res = await fetch(`/api/algorithms?page=${nextPage}&limit=${limit}`);
			if (!res.ok) throw new Error("Erro na API");

			const result = await res.json();

			algorithms = [...algorithms, ...(result.algorithms ?? [])];
			page = result.page ?? nextPage;
			hasMore = result.hasMore ?? false;
		} catch {
			hasMore = false;
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		if (!browser || !sentinel) return;

		const observer = new IntersectionObserver(
			(entries) => {
				const [entry] = entries;
				if (entry.isIntersecting && hasMore && !isLoading) {
					loadMore();
				}
			},
			{
				rootMargin: "200px"
			}
		);

		observer.observe(sentinel);

		return () => {
			observer.disconnect();
		};
	});
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
		{#each algorithms as item (item.publicId)}
			{@const diff = DIFFICULTY_MAP[item.difficulty] ?? {
				label: item.difficulty,
				style: "bg-gray-800/60 text-gray-300 border-gray-700"
			}}

			<Card title={item.name} href="/algorithms/{item.slug}-{item.publicId}">
				{#snippet headerRightTop()}
					<span
						class="text-xs font-semibold uppercase tracking-wider px-2.5 py-1 rounded-md border {diff.style}"
					>
						{diff.label}
					</span>
				{/snippet}

				<p class="text-xs font-medium text-gray-400 flex items-center gap-2">
					<span class="inline-block w-2 h-2 rounded-full bg-text-brand" aria-hidden="true"></span>
					{item.category}
				</p>

				{#snippet footerLeft()}
					<span
						class="text-xs font-mono text-gray-400 truncate max-w-30 select-all"
						title="Copiar ID"
					>
						{item.publicId}
					</span>
				{/snippet}
			</Card>
		{:else}
			<div
				role="status"
				class="col-span-full py-12 text-center bg-app-surface border border-gray-800 rounded-xl text-gray-400 text-sm"
			>
				Nenhum algoritmo encontrado.
			</div>
		{/each}
	</div>
	{#if hasMore}
		<div
			bind:this={sentinel}
			class="py-8 flex justify-center items-center gap-2 text-gray-400 text-sm"
		>
			{#if isLoading}
				<svg class="animate-spin h-5 w-5 text-text-brand" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
					></circle>
					<path
						class="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
					></path>
				</svg>
				<span>Carregando mais algoritmos...</span>
			{/if}
		</div>
	{/if}
</div>
