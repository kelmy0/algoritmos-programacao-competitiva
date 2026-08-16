<script lang="ts">
	import { customFetch } from "$lib/api/client";
	import type { ApiError } from "$lib/types/api";
	import { focusTrap } from "$lib/utils/a11y";
	import { normalizeApiError } from "$lib/utils/errors";
	import { ADMIN_ALGORITHMS_ERRORS } from "$lib/errors/admin/algorithms";
	import type { PageData } from "./$types";
	import Button from "$lib/components/ui/Button.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";

	let { data }: { data: PageData } = $props();

	const algorithm = $derived(data.algorithm);
	let isLoading = $state(false);
	let apiError = $state<ApiError | null>(null);
	let isInfoModalOpen = $state(false);
	let localStatus = $state<string>();
	let status = $derived(localStatus ?? data.algorithm?.status);

	const modalStatusClass = $derived(
		isLoading
			? "bg-blue-950/80 border-blue-900/60 text-blue-400"
			: apiError
				? "bg-red-950/80 border-red-900/60 text-red-400"
				: "bg-emerald-950/80 border-emerald-900/60 text-emerald-400"
	);

	const modalTitle = $derived(
		isLoading ? "Restaurando..." : apiError ? apiError.code : "Algoritmo restaurado!"
	);

	const modalDescription = $derived(
		isLoading
			? "Aguarde enquanto restauramos o algoritmo."
			: apiError
				? apiError.message
				: "Algoritmo restaurado com sucesso!"
	);

	function closeInfoModal() {
		isInfoModalOpen = false;
	}

	async function handleRestore() {
		if (isLoading) return;

		isInfoModalOpen = true;

		if (!algorithm?.slug || !algorithm.publicId) {
			apiError = normalizeApiError(
				"ALGORITHM_INVALID_PUBLIC_ID",
				"Id público do algoritmo não é valido!",
				ADMIN_ALGORITHMS_ERRORS
			);
			return;
		}

		isLoading = true;

		const { error, status } = await customFetch<null>(
			window.fetch,
			`/api/admin/algorithms/restore/${algorithm.slug}-${algorithm.publicId}`,
			{
				method: "PATCH"
			},
			ADMIN_ALGORITHMS_ERRORS
		);

		isLoading = false;

		if (error || status !== 204) {
			apiError = error || normalizeApiError("INTERNAL_SERVER_ERROR");
			return;
		}

		apiError = null;
		localStatus = "pending";
	}
</script>

<svelte:head>
	{#if algorithm}
		<title>{algorithm.name}</title>
		<meta name="robots" content="noindex, nofollow" />
	{:else}
		<title>Algoritmo</title>
	{/if}
</svelte:head>

<div class="mx-auto max-w-4xl space-y-6 font-inter">
	<a
		href="/admin/dashboard"
		class="inline-flex items-center gap-2 text-xs font-medium text-text-secondary hover:text-text-primary transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand rounded-md py-1 px-1.5 -ml-1.5 w-fit"
	>
		<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
		</svg>
		<span>Voltar para a dashboard</span>
	</a>

	{#if algorithm}
		<article
			class="bg-app-surface border-app-border space-y-6 rounded-xl border p-6 shadow-xl md:p-8"
		>
			{#if status === "deleted"}
				<div
					class="flex flex-wrap items-center justify-between gap-4 rounded-lg border border-red-500/20 bg-red-500/10 p-4 text-red-200"
				>
					<div class="flex items-center gap-2 text-sm">
						<svg
							class="h-5 w-5 text-red-400 shrink-0"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
							/>
						</svg>
						<span>Este item está na <strong>Lixeira</strong>.</span>
					</div>

					<Button
						size="sm"
						onclick={() => handleRestore()}
						disabled={isLoading}
						variant="success-soft"
					>
						<svg
							class="w-4 h-4"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
							/>
						</svg>
						<span>Restaurar item</span>
					</Button>
				</div>
			{/if}

			<header class="border-app-border space-y-3 border-b pb-6">
				<div class="flex flex-wrap items-center justify-between gap-2">
					<span
						class="text-text-brand bg-text-brand/10 border-text-brand/20 rounded-md border px-2.5 py-1 font-mono text-xs"
						aria-label="Categoria: {algorithm.category}"
					>
						{algorithm.category}
					</span>
					<span class="font-mono text-xs text-text-muted">
						<span class="sr-only">Identificador: </span>ID: {algorithm.publicId}
					</span>
				</div>

				<h1 class="text-text-primary font-montserrat text-3xl font-bold tracking-tight">
					{algorithm.name}
				</h1>
			</header>

			<section class="space-y-4" aria-labelledby="heading-implementacao">
				<h2
					id="heading-implementacao"
					class="text-sm font-semibold uppercase tracking-wider text-text-secondary"
				>
					Implementação
				</h2>

				<div class="prose prose-invert max-w-none font-mono text-sm text-text-secondary">
					{@html algorithm.contentHtml}
				</div>
			</section>

			{#if status === "deleted"}
				<footer class="border-t border-app-border pt-6 flex justify-end">
					<Button onclick={() => handleRestore()} disabled={isLoading} variant="success-soft">
						<svg
							class="w-4 h-4"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
							/>
						</svg>
						<span>Restaurar item</span>
					</Button>
				</footer>
			{/if}
		</article>
	{/if}
</div>

{#if isInfoModalOpen}
	<Modal
		isOpen={isInfoModalOpen}
		title={modalTitle}
		description={modalDescription}
		variant={apiError ? "danger" : "success"}
		{isLoading}
		onClose={() => closeInfoModal()}
		{focusTrap}
	>
		{#snippet footer()}
			<button
				type="button"
				onclick={() => closeInfoModal()}
				disabled={isLoading}
				class="w-full sm:w-auto px-4 py-2 rounded-lg bg-app-overlay text-text-primary border border-app-border text-sm font-medium hover:bg-app-border focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
			>
				Fechar
			</button>
		{/snippet}
	</Modal>
{/if}
