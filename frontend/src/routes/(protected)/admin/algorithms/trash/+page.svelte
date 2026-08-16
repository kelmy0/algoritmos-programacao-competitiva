<script lang="ts">
	import Card from "$lib/components/ui/Card.svelte";
	import { DIFFICULTY_MAP, STATUS_MAP } from "$lib/constants/algorithms";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();

	const difficultyStyles: Record<string, string> = {
		beginner: "bg-emerald-950 text-emerald-300 border-emerald-800",
		intermediate: "bg-amber-950 text-amber-300 border-amber-800",
		advanced: "bg-orange-950 text-orange-300 border-orange-800",
		expert: "bg-red-950 text-red-300 border-red-800"
	};

	const difficultyLabels: Record<string, string> = {
		beginner: "Iniciante",
		intermediate: "Intermediário",
		advanced: "Avançado",
		expert: "Especialista"
	};

	const statusStyles: Record<string, string> = {
		approved: difficultyStyles.beginner,
		pending: difficultyStyles.intermediate,
		rejected: difficultyStyles.advanced,
		deleted: difficultyStyles.expert
	};

	const statusLabels: Record<string, string> = {
		approved: "Aprovado",
		pending: "Pendente",
		rejected: "Rejeitado",
		deleted: "Deletado"
	};
</script>

<svelte:head>
	<title>Algoritmos deletados</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="space-y-6 font-inter">
	<header
		class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-app-border"
	>
		<div>
			<h1 class="font-montserrat font-bold text-2xl md:text-3xl text-text-primary tracking-tight">
				Lixeira
			</h1>
			<p class="text-sm text-text-muted mt-1">Veja seus algoritmos que estão na lixeira.</p>
		</div>
	</header>

	<!-- Algorithm Grid -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
		{#each data.algorithms as item (item.publicId)}
			{@const diff = DIFFICULTY_MAP[item.difficulty] ?? {
				label: item.difficulty,
				style: "bg-app-bg/60 text-text-secondary border-app-border"
			}}
			{@const st = STATUS_MAP[item.status] ?? {
				label: item.status,
				style: "bg-app-bg/60 text-text-secondary border-app-border"
			}}

			<Card title={item.name} href="/admin/algorithms/my-algorithms/{item.slug}-{item.publicId}">
				{#snippet headerRightTop()}
					<span
						class="text-xs font-semibold uppercase tracking-wider px-2.5 py-1 rounded-md border {diff.style}"
					>
						{diff.label}
					</span>
				{/snippet}

				<p class="text-xs font-medium text-text-muted flex items-center gap-2">
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
						class="text-xs font-mono text-text-muted truncate max-w-30 select-all"
						title="Copiar ID"
					>
						{item.publicId}
					</span>
				{/snippet}
			</Card>
		{:else}
			<div
				role="status"
				class="col-span-full py-12 text-center bg-app-surface border border-app-border rounded-xl text-text-muted text-sm"
			>
				Nenhum algoritmo encontrado.
			</div>
		{/each}
	</div>
</div>
