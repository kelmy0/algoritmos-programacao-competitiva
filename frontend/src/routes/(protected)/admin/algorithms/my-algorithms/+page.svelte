<script lang="ts">
	import Card from "$lib/components/ui/Card.svelte";
	import { DIFFICULTY_MAP, STATUS_MAP } from "$lib/constants/algorithms";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Meus Algoritmos</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="space-y-6 font-inter">
	<header
		class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-gray-800"
	>
		<div>
			<h1 class="font-montserrat font-bold text-2xl md:text-3xl text-text-primary tracking-tight">
				Meus Algoritmos
			</h1>
			<p class="text-sm text-gray-400 mt-1">Veja seus algoritmos.</p>
		</div>
	</header>

	<!-- Algorithm Grid -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
		{#each data.algorithms as item (item.publicId)}
			{@const diff = DIFFICULTY_MAP[item.difficulty] ?? {
				label: item.difficulty,
				style: "bg-gray-800/60 text-gray-300 border-gray-700"
			}}
			{@const st = STATUS_MAP[item.status] ?? {
				label: item.status,
				style: "bg-gray-800/60 text-gray-300 border-gray-700"
			}}

			<Card title={item.name} href="/admin/algorithms/my-algorithms/{item.slug}-{item.publicId}">
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

				{#snippet headerRightBottom()}
					<span
						class="text-xs font-semibold uppercase tracking-wider px-2.5 py-1 rounded-md border {st.style}"
					>
						{st.label}
					</span>
				{/snippet}

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
</div>
