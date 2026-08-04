<script lang="ts">
	import { customFetch } from "$lib/api/client";
	import type { ApiError } from "$lib/types/api";
	import { focusTrap } from "$lib/utils/a11y";
	import { normalizeApiError } from "$lib/utils/errors";
	import { ADMIN_ALGORITHMS_ERRORS } from "../../new/newAlgorithm.svelte";
	import type { PageData } from "./$types";

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
		class="inline-flex items-center gap-2 text-xs font-medium text-gray-300 hover:text-text-primary transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand rounded-md py-1 px-1.5 -ml-1.5 w-fit"
	>
		<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
		</svg>
		<span>Voltar para a dashboard</span>
	</a>

	{#if algorithm}
		<article
			class="bg-app-surface border-gray-800 space-y-6 rounded-xl border p-6 shadow-xl md:p-8"
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

					<button
						type="button"
						onclick={() => handleRestore()}
						disabled={isLoading}
						class="px-4 py-2 rounded-md border border-emerald-500/30 bg-emerald-600/20 text-emerald-400 font-mono font-semibold text-xs hover:bg-emerald-600/30 hover:text-emerald-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
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
						Restaurar item
					</button>
				</div>
			{/if}

			<header class="border-gray-800 space-y-3 border-b pb-6">
				<div class="flex flex-wrap items-center justify-between gap-2">
					<span
						class="text-text-brand bg-text-brand/10 border-text-brand/20 rounded-md border px-2.5 py-1 font-mono text-xs"
						aria-label="Categoria: {algorithm.category}"
					>
						{algorithm.category}
					</span>
					<span class="font-mono text-xs text-gray-400">
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
					class="text-sm font-semibold uppercase tracking-wider text-gray-200"
				>
					Implementação
				</h2>

				<div class="prose prose-invert max-w-none font-mono text-sm text-gray-200">
					{@html algorithm.contentHtml}
				</div>
			</section>

			{#if status === "deleted"}
				<footer class="border-t border-gray-800 pt-6 flex justify-end">
					<button
						type="button"
						onclick={() => handleRestore()}
						disabled={isLoading}
						class="px-5 py-2.5 rounded-md border border-emerald-500/30 bg-emerald-600/20 text-emerald-400 font-mono font-semibold text-xs hover:bg-emerald-600/30 hover:text-emerald-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
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
						Restaurar item
					</button>
				</footer>
			{/if}
		</article>
	{/if}
</div>

{#if isInfoModalOpen}
	<div
		use:focusTrap
		class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4"
		role="dialog"
		aria-modal="true"
		aria-busy={isLoading}
		aria-labelledby="info-modal-title"
		aria-describedby="info-modal-description"
		onkeydown={(e) => e.key === "Escape" && closeInfoModal()}
		tabindex="-1"
	>
		<div
			class="bg-app-surface border border-gray-800 rounded-xl p-6 max-w-md w-full flex flex-col gap-5 shadow-2xl animate-in fade-in zoom-in-95 duration-150 relative"
		>
			<div class="flex items-start gap-3">
				<div class="p-2.5 self-start rounded-lg shrink-0 border {modalStatusClass}">
					{#if isLoading}
						<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<g>
								<circle cx="12" cy="3" r="1">
									<animate
										id="spinner_7Z73"
										begin="0;spinner_tKsu.end-0.5s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="16.50" cy="4.21" r="1">
									<animate
										id="spinner_Wd87"
										begin="spinner_7Z73.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="7.50" cy="4.21" r="1">
									<animate
										id="spinner_tKsu"
										begin="spinner_9Qlc.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="19.79" cy="7.50" r="1">
									<animate
										id="spinner_lMMO"
										begin="spinner_Wd87.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="4.21" cy="7.50" r="1">
									<animate
										id="spinner_9Qlc"
										begin="spinner_Khxv.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="21.00" cy="12.00" r="1">
									<animate
										id="spinner_5L9t"
										begin="spinner_lMMO.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="3.00" cy="12.00" r="1">
									<animate
										id="spinner_Khxv"
										begin="spinner_ld6P.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="19.79" cy="16.50" r="1">
									<animate
										id="spinner_BfTD"
										begin="spinner_5L9t.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="4.21" cy="16.50" r="1">
									<animate
										id="spinner_ld6P"
										begin="spinner_XyBs.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="16.50" cy="19.79" r="1">
									<animate
										id="spinner_7gAK"
										begin="spinner_BfTD.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="7.50" cy="19.79" r="1">
									<animate
										id="spinner_XyBs"
										begin="spinner_HiSl.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<circle cx="12" cy="21" r="1">
									<animate
										id="spinner_HiSl"
										begin="spinner_7gAK.begin+0.1s"
										attributeName="r"
										calcMode="spline"
										dur="0.6s"
										values="1;2;1"
										keySplines=".27,.42,.37,.99;.53,0,.61,.73"
									/>
								</circle>
								<animateTransform
									attributeName="transform"
									type="rotate"
									dur="6s"
									values="360 12 12;0 12 12"
									repeatCount="indefinite"
								/>
							</g>
						</svg>
					{:else if apiError}
						<svg
							class="w-6 h-6"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
							/>
						</svg>
					{:else}
						<svg
							class="w-6 h-6"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
							></path>
						</svg>
					{/if}
				</div>

				<div class="flex-1 pr-6">
					<h2 id="info-modal-title" class="text-lg font-bold text-gray-100 font-montserrat">
						{modalTitle}
					</h2>
					<p
						id="info-modal-description"
						aria-live="polite"
						class="text-sm text-gray-300 mt-1 leading-relaxed"
					>
						{#if isLoading}
							Aguarde enquanto restauramos o algoritmo.
						{:else if apiError}
							{apiError.message}
						{:else}
							Algoritmo restaurado com sucesso!
						{/if}
					</p>
				</div>

				<button
					type="button"
					onclick={closeInfoModal}
					aria-label="Fechar modal"
					class="hover:cursor-pointer absolute top-4 right-4 text-gray-400 hover:text-white p-1 rounded-lg focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			<div class="flex justify-end gap-3 pt-4 border-t border-gray-800">
				<button
					type="button"
					onclick={closeInfoModal}
					disabled={isLoading}
					class="w-full sm:w-auto px-4 py-2 rounded-lg bg-gray-900 text-gray-300 hover:text-white border border-gray-700 text-sm font-medium hover:bg-gray-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
				>
					Fechar
				</button>
			</div>
		</div>
	</div>
{/if}
