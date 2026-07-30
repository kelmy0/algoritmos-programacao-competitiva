<script lang="ts">
	import { focusTrap } from "$lib/utils/a11y";
	import { AlgorithmEditor } from "$lib/utils/editor.svelte";
	import type { PageData } from "./$types";
	import { EditAlgorithmController } from "./editAlgorithm.svelte";

	let { data }: { data: PageData } = $props();

	const editor = new AlgorithmEditor();
	const controller = new EditAlgorithmController();

	let isInitialized = $state(false);
	let localStatus = $state<string>();
	let status = $derived(localStatus ?? data.algorithm?.Status);

	$effect(() => {
		if (data.algorithm && !isInitialized) {
			editor.load(data.algorithm);
			controller.publicId = data.algorithm.PublicId;
			controller.slug = data.algorithm.Slug;
			isInitialized = true;
		}
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const payload = editor.getPayload();
		if (payload) {
			const success = await controller.handleSubmit(payload);

			if (success) {
				localStatus = "pending";
			}
		}
	}

	async function handleDelete() {
		const success = await controller.handleDelete();
		if (success) {
			localStatus = "deleted";
		}
	}

	async function handleRestore() {
		const success = await controller.handleRestore();
		if (success) {
			localStatus = "pending";
		}
	}

	const errorLabels: Record<string, string> = {
		save: "salvar",
		delete: "deletar",
		restore: "restaurar"
	};

	const successLabels: Record<string, string> = {
		save: "salvo",
		delete: "deletado",
		restore: "restaurado"
	};

	const successText: Record<string, string> = {
		save: "Suas alterações foram enviadas e já estão em espera para aprovação.",
		delete: "Algoritmo movido para lixeira! Você tem até 7 dias para restaurá-lo.",
		restore: "Algoritmo restaurado com sucesso! Ele já está em espera para aprovação."
	};
</script>

<svelte:head>
	<title>Editar Algoritmo</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6 font-inter p-6">
	<header class="border-b border-gray-800 pb-4">
		<h1 class="font-montserrat text-2xl font-bold text-text-primary">Editar um Algoritmo</h1>
		<p class="text-sm text-gray-300 mt-1">
			Preencha os metadados e edite o conteúdo em Markdown com o preview ao lado.
		</p>
	</header>

	<form onsubmit={handleSubmit} class="space-y-6" aria-label="Formulário de criação de algoritmo">
		<fieldset
			class="grid grid-cols-1 md:grid-cols-3 gap-4 bg-app-surface p-5 rounded-xl border border-gray-800"
		>
			<legend class="sr-only">Metadados do Algoritmo</legend>

			<div class="space-y-2">
				<label for="name" class="block text-sm font-medium text-gray-200">
					Nome do Algoritmo <span class="text-red-400" aria-hidden="true">*</span>
				</label>
				<input
					id="name"
					type="text"
					placeholder="Ex: Busca Binária"
					required
					autocomplete="off"
					disabled={controller.isLoading || controller.isDeleting}
					bind:value={editor.name}
					bind:this={editor.nameInput}
					oninput={() => editor.onNameInput()}
					onblur={() => editor.onNameBlur()}
					aria-required="true"
					aria-invalid={editor.hasNameError || controller.hasNameError}
					aria-describedby={editor.hasNameError || controller.hasNameError
						? "name-error"
						: undefined}
					class="w-full bg-gray-900 border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2
					{editor.hasNameError || controller.hasNameError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'} disabled:cursor-not-allowed"
				/>
				{#if editor.hasNameError || controller.hasNameError}
					<p id="name-error" role="alert" class="text-xs text-amber-500">
						O nome precisa ter no mínimo 3 letras válidas.
					</p>
				{/if}
			</div>

			<div class="space-y-2">
				<label for="category" class="block text-sm font-medium text-gray-200">
					Categoria <span class="text-red-400" aria-hidden="true">*</span>
				</label>
				<input
					id="category"
					type="text"
					placeholder="Ex: Grafos, Busca, DP"
					required
					autocomplete="off"
					bind:value={editor.category}
					bind:this={editor.categoryInput}
					disabled={controller.isLoading || controller.isDeleting}
					oninput={() => editor.onCategoryInput()}
					onblur={() => editor.onCategoryBlur()}
					aria-required="true"
					aria-invalid={editor.hasCategoryError || controller.hasCategoryError}
					aria-describedby={editor.hasCategoryError || controller.hasCategoryError
						? "category-error"
						: undefined}
					class="w-full bg-gray-900 border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2
					{editor.hasCategoryError || controller.hasCategoryError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'} disabled:cursor-not-allowed"
				/>
				{#if editor.hasCategoryError || controller.hasCategoryError}
					<p id="category-error" role="alert" class="text-xs text-amber-500">
						A categoria precisa ter no mínimo 3 letras válidas.
					</p>
				{/if}
			</div>

			<div class="space-y-2">
				<label for="difficulty" class="block text-sm font-medium text-gray-200">Dificuldade</label>
				<select
					id="difficulty"
					required
					bind:value={editor.difficulty}
					bind:this={editor.difficultyInput}
					disabled={controller.isLoading || controller.isDeleting}
					oninput={() => editor.onDifficultyInput()}
					onblur={() => editor.onDifficultyBlur()}
					aria-required="true"
					aria-invalid={editor.hasDifficultyError}
					aria-describedby={editor.hasDifficultyError ? "difficulty-error" : undefined}
					class="hover:cursor-pointer w-full bg-gray-900 border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2
					{editor.hasDifficultyError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'} disabled:cursor-not-allowed"
				>
					<option value="beginner">Iniciante</option>
					<option value="intermediate">Intermediário</option>
					<option value="advanced">Avançado</option>
					<option value="expert">Especialista</option>
				</select>
				{#if editor.hasDifficultyError}
					<p id="difficulty-error" role="alert" class="text-xs text-amber-500">
						A dificuldade precisa ser uma das 4 opções.
					</p>
				{/if}
			</div>
		</fieldset>

		<div class="invisible" bind:this={controller.alertDiv}></div>
		{#if controller.apiError}
			<div
				role="alert"
				aria-live="assertive"
				class="p-4 bg-red-950/40 border border-red-900/60 rounded-xl text-red-300 text-sm flex items-start gap-3 shadow-lg transition-all
			{controller.isLoading || controller.isDeleting ? 'opacity-50 pointer-events-none' : 'opacity-100'}"
			>
				<svg
					class="w-5 h-5 shrink-0 mt-0.5 text-red-400"
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
				<div class="space-y-0.5">
					<span class="font-semibold block text-red-200"
						>Erro ao {errorLabels[controller.lastAction] ?? controller.lastAction} algoritmo</span
					>
					<p class="text-xs text-red-300/80 leading-relaxed">
						{controller.apiError.message || "Ocorreu um erro ao tentar salvar. Tente novamente."}
					</p>
				</div>
			</div>
		{/if}

		{#if controller.isSuccess}
			<div
				role="status"
				aria-live="polite"
				class="p-4 border text-sm flex items-start gap-3 shadow-lg transition-all rounded-xl
				{controller.lastAction !== 'delete'
					? 'bg-emerald-950/40 border-emerald-900/60 text-red-300'
					: 'bg-red-950/40 border-red-900/60 text-red-300'}"
			>
				<svg
					class="w-5 h-5 shrink-0 mt-0.5 {controller.lastAction !== 'delete'
						? 'text-emerald-400'
						: 'text-red-400'}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d={controller.lastAction !== "delete"
							? "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
							: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"}
					/>
				</svg>
				<div class="space-y-0.5">
					<span
						class="font-semibold block {controller.lastAction !== 'delete'
							? 'text-emerald-200'
							: 'text-red-200'}"
						>Algoritmo {successLabels[controller.lastAction] ?? controller.lastAction} com sucesso!</span
					>
					<p
						class="text-xs leading-relaxed {controller.lastAction !== 'delete'
							? 'text-emerald-300/80'
							: 'text-red-300/80'}"
					>
						{successText[controller.lastAction] ?? controller.lastAction}
						<a
							href={controller.link}
							class="underline {controller.lastAction !== 'delete'
								? 'hover:text-emerald-200'
								: 'hover:text-red-200'}"
							onclick={() => (controller.isSuccess = false)}
						>
							Visualizar o algoritmo {successLabels[controller.lastAction] ?? controller.lastAction}
						</a>.
					</p>
				</div>
			</div>
		{/if}

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<section
				class="flex flex-col h-150 bg-app-surface border border-gray-800 rounded-xl overflow-hidden"
				aria-label="Editor de código Markdown"
			>
				<div
					role="toolbar"
					aria-label="Ferramentas de formatação Markdown"
					class="min-h-12 flex flex-wrap items-center gap-1.5 px-3 py-2 bg-gray-900 border-b border-gray-800 text-xs"
				>
					<button
						type="button"
						onclick={() => editor.insertSnippet("## ", "", "Título")}
						aria-label="Inserir Título Nível 2"
						class="btn-toolbar hover:cursor-pointer">H2</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("### ", "", "Subtítulo")}
						aria-label="Inserir Subtítulo Nível 3"
						class="btn-toolbar hover:cursor-pointer">H3</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("**", "**", "negrito")}
						aria-label="Texto em Negrito"
						class="btn-toolbar hover:cursor-pointer"><b aria-hidden="true">B</b></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("*", "*", "itálico")}
						aria-label="Texto em Itálico"
						class="btn-toolbar hover:cursor-pointer"><i aria-hidden="true">I</i></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("\n```cpp\n", "\n```\n", "// seu código C++ aqui")}
						aria-label="Inserir bloco de código C++"
						class="btn-toolbar hover:cursor-pointer font-mono text-text-brand">C++ Code</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("> ", "", "Nota importante")}
						aria-label="Inserir citação"
						class="btn-toolbar hover:cursor-pointer">Quote</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("1. ", "", "Item")}
						aria-label="Inserir lista numerada"
						class="btn-toolbar hover:cursor-pointer">Lista</button
					>
				</div>

				<label for="content-editor" class="sr-only">Conteúdo em Markdown</label>
				<textarea
					id="content-editor"
					placeholder="Escreva o conteúdo em Markdown aqui..."
					required
					bind:value={editor.content}
					bind:this={editor.contentInput}
					disabled={controller.isLoading || controller.isDeleting}
					oninput={() => editor.onContentInput()}
					onblur={() => editor.onContentBlur()}
					aria-required="true"
					aria-invalid={editor.hasContentError || controller.hasContentError}
					aria-describedby={editor.hasContentError || controller.hasContentError
						? "content-error"
						: undefined}
					class="w-full flex-1 p-4 bg-transparent text-gray-200 font-mono text-sm
            resize-none focus:outline-none focus-visible:ring-2 leading-relaxed border
            {editor.hasContentError || controller.hasContentError
						? 'border-red-500 focus:ring-red-500'
						: 'border-transparent focus:ring-text-brand'} 
						disabled:cursor-not-allowed"></textarea>
			</section>

			<section
				class="flex flex-col h-150 bg-app-surface border border-gray-800 rounded-xl overflow-hidden"
				aria-label="Preview do conteúdo"
			>
				<div
					class="min-h-12 px-4 py-2 bg-gray-900 border-b border-gray-800 flex items-center justify-between gap-2"
				>
					<span class="text-xs font-mono font-medium text-gray-300">Preview em Tempo Real</span>
					{#if controller.slug}
						<a
							href="/admin/algorithms/my-algorithms/{controller.slug ??
								data.algorithm?.Slug}-{controller.publicId ?? data.algorithm?.PublicId}"
							target="_blank"
							rel="noopener noreferrer"
							class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium text-gray-300 bg-gray-800/60 hover:bg-gray-800 hover:text-white border border-gray-700/60 focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand transition-all"
							title="Abrir visualização em nova aba"
						>
							<svg
								class="w-4 h-4 text-gray-400 group-hover:text-white"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								aria-hidden="true"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
								/>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
								/>
							</svg>
							<span>Visualizar em nova aba</span>
						</a>
					{/if}
				</div>

				<div
					aria-live="polite"
					class="p-6 whitespace-pre-wrap wrap-break-word overflow-y-auto prose prose-invert max-w-none font-mono text-sm text-gray-200 prose-pre:whitespace-pre-wrap prose-pre:wrap-break-words"
				>
					{#if editor.content.trim()}
						{#await editor.previewPromise}
							<p role="status" class="text-gray-400 italic font-sans text-sm">Gerando preview...</p>
						{:then html}
							{@html html}
						{/await}
					{:else}
						<p class="text-gray-400 italic font-sans text-sm">
							O preview aparecerá aqui conforme você digita...
						</p>
					{/if}
				</div>
			</section>
		</div>

		<div
			class="flex flex-col-reverse sm:flex-row justify-between items-stretch sm:items-center gap-4 pt-4 border-t border-gray-800/60"
		>
			{#if status !== "deleted"}
				<button
					type="button"
					onclick={() => controller.openDeleteModal()}
					disabled={controller.isLoading || controller.isDeleting}
					aria-haspopup="dialog"
					aria-expanded={controller.isDeleteModalOpen}
					class="px-5 py-2.5 rounded-lg border border-red-900/50 bg-red-950/30 text-red-400 font-semibold text-sm hover:bg-red-900/50 hover:text-red-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 transition-colors flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
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
							d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
						/>
					</svg>
					Mover para lixeira
				</button>
			{:else}
				<button
					type="button"
					onclick={() => handleRestore()}
					disabled={controller.isLoading || controller.isDeleting}
					class="px-5 py-2.5 rounded-lg border border-emerald-900/50 bg-emerald-950/30 text-emerald-400 font-semibold text-sm hover:bg-emerald-900/50 hover:text-emerald-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 transition-colors flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
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
					Restaurar
				</button>
			{/if}

			<div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 justify-end">
				{#if editor.hasContentError || controller.hasContentError}
					<p id="content-error" role="alert" class="text-xs text-amber-500 self-center">
						O conteúdo precisa de no mínimo 10 letras.
					</p>
				{/if}
				<button
					type="submit"
					disabled={controller.isLoading || controller.isDeleting}
					class="px-6 py-2.5 rounded-lg bg-text-brand text-gray-950 font-semibold text-sm hover:bg-blue-400 focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-text-brand focus-visible:ring-offset-gray-900 transition-colors disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
				>
					{controller.isLoading ? "Salvando..." : "Salvar Algoritmo"}
				</button>
			</div>
		</div>
	</form>
</div>

{#if controller.isDeleteModalOpen}
	<div
		use:focusTrap
		class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4"
		role="dialog"
		aria-modal="true"
		aria-busy={controller.isDeleting}
		aria-labelledby="delete-modal-title"
		aria-describedby="delete-modal-description"
		onkeydown={(e) => e.key === "Escape" && controller.closeDeleteModal()}
		tabindex="-1"
	>
		<div
			class="bg-app-surface border border-gray-800 rounded-xl p-6 max-w-md w-full flex flex-col gap-5 shadow-2xl animate-in fade-in zoom-in-95 duration-150 relative"
		>
			<div class="flex items-start gap-3">
				<div
					class="p-2.5 self-start rounded-lg shrink-0 border bg-red-950/80 border-red-900/60 text-red-400"
				>
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
				</div>

				<div class="flex-1 pr-6">
					<h2 id="delete-modal-title" class="text-lg font-bold text-gray-100 font-montserrat">
						Excluir Algoritmo?
					</h2>
					<p id="delete-modal-description" class="text-sm text-gray-300 mt-1 leading-relaxed">
						Tem certeza que deseja mover para lixeira este algoritmo? Você terá 7 dias para
						restaurá-lo.
					</p>
				</div>

				<button
					type="button"
					onclick={() => controller.closeDeleteModal()}
					disabled={controller.isDeleting || controller.isLoading}
					aria-label="Fechar modal"
					class="hover:cursor-pointer absolute top-4 right-4 text-gray-400 hover:text-white p-1 rounded-lg focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 disabled:opacity-50 transition-colors"
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

			<div
				class="flex flex-col-reverse sm:flex-row justify-end gap-3 pt-4 border-t border-gray-800"
			>
				<button
					type="button"
					onclick={() => controller.closeDeleteModal()}
					disabled={controller.isDeleting || controller.isLoading}
					class="w-full sm:w-auto px-4 py-2 rounded-lg bg-gray-900 text-gray-300 hover:text-white border border-gray-700 text-sm font-medium hover:bg-gray-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
				>
					Cancelar
				</button>
				<button
					type="button"
					onclick={handleDelete}
					disabled={controller.isDeleting || controller.isLoading}
					aria-busy={controller.isDeleting}
					class="w-full sm:w-auto px-4 py-2 rounded-lg bg-red-600 hover:bg-red-500/80 text-white text-sm font-semibold focus:outline-none focus-visible:ring-2 focus-visible:ring-red-400 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 transition-colors disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
				>
					{controller.isDeleting ? "Excluindo..." : "Sim, Excluir"}
				</button>
			</div>
		</div>
	</div>
{/if}
