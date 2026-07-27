<script lang="ts">
	import { AlgorithmEditor } from "$lib/utils/editor.svelte";
	import { NewAlgorithmController } from "./newAlgorithm.svelte";

	const editor = new AlgorithmEditor();
	const controller = new NewAlgorithmController();

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const payload = editor.getPayload();
		if (payload) {
			await controller.submit(payload);
		}
	}
</script>

<svelte:head>
	<title>Criar Algoritmo | Admin</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6 font-inter p-6">
	<header class="border-b border-gray-800 pb-4">
		<h1 class="font-montserrat text-2xl font-bold text-text-primary">Criar Novo Algoritmo</h1>
		<p class="text-sm text-gray-300 mt-1">
			Preencha os metadados e escreva o conteúdo em Markdown com o preview ao lado.
		</p>
	</header>

	<form onsubmit={handleSubmit} class="space-y-6" aria-label="Formulário de criação de algoritmo">
		<fieldset
			class="grid grid-cols-1 md:grid-cols-4 gap-4 bg-app-surface p-5 rounded-xl border border-gray-800"
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
					disabled={controller.isLoading}
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
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'}"
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
					disabled={controller.isLoading}
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
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'}"
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
					disabled={controller.isLoading}
					oninput={() => editor.onDifficultyInput()}
					onblur={() => editor.onDifficultyBlur()}
					aria-required="true"
					aria-invalid={editor.hasDifficultyError}
					aria-describedby={editor.hasDifficultyError ? "difficulty-error" : undefined}
					class="w-full bg-gray-900 border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2
					{editor.hasDifficultyError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'}"
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

			<div class="space-y-2">
				<label for="password" class="block text-sm font-medium text-gray-200">
					Senha <span class="text-red-400" aria-hidden="true">*</span>
				</label>
				<input
					id="password"
					type="password"
					placeholder="Senha das rotas admin"
					required
					autocomplete="current-password"
					bind:value={controller.password}
					bind:this={controller.passwordInput}
					disabled={controller.isLoading}
					oninput={() => controller.onPasswordInput()}
					onblur={() => controller.onPasswordBlur()}
					aria-required="true"
					aria-invalid={controller.hasPasswordError}
					aria-describedby={controller.hasPasswordError ? "password-error" : undefined}
					class="w-full bg-gray-900 border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2
					{controller.hasPasswordError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'}"
				/>
				{#if controller.hasPasswordError}
					<p id="password-error" role="alert" class="text-xs text-amber-500">
						A senha precisa de no minimo 8 letras.
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
			{controller.isLoading ? 'opacity-50 pointer-events-none' : 'opacity-100'}"
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
					<span class="font-semibold block text-red-200">Erro ao salvar algoritmo</span>
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
				class="p-4 bg-emerald-950/40 border border-emerald-900/60 rounded-xl text-emerald-300 text-sm flex items-start gap-3 shadow-lg transition-all"
			>
				<svg
					class="w-5 h-5 shrink-0 mt-0.5 text-emerald-400"
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
					/>
				</svg>
				<div class="space-y-0.5">
					<span class="font-semibold block text-emerald-200">Algoritmo salvo com sucesso!</span>
					<p class="text-xs text-emerald-300/80 leading-relaxed">
						Suas alterações foram enviadas e já estão em espera para aprovação.
						<a href={controller.link} class="underline hover:text-emerald-200">
							Visualizar o algoritmo enviado
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
					class="flex flex-wrap items-center gap-1.5 p-3 bg-gray-900 border-b border-gray-800 text-xs"
				>
					<button
						type="button"
						onclick={() => editor.insertSnippet("## ", "", "Título")}
						aria-label="Inserir Título Nível 2"
						class="btn-toolbar">H2</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("### ", "", "Subtítulo")}
						aria-label="Inserir Subtítulo Nível 3"
						class="btn-toolbar">H3</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("**", "**", "negrito")}
						aria-label="Texto em Negrito"
						class="btn-toolbar"><b aria-hidden="true">B</b></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("*", "*", "itálico")}
						aria-label="Texto em Itálico"
						class="btn-toolbar"><i aria-hidden="true">I</i></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("\n```cpp\n", "\n```\n", "// seu código C++ aqui")}
						aria-label="Inserir bloco de código C++"
						class="btn-toolbar font-mono text-text-brand">C++ Code</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("> ", "", "Nota importante")}
						aria-label="Inserir citação"
						class="btn-toolbar">Quote</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet("1. ", "", "Item")}
						aria-label="Inserir lista numerada"
						class="btn-toolbar">Lista</button
					>
				</div>

				<label for="content-editor" class="sr-only">Conteúdo em Markdown</label>
				<textarea
					id="content-editor"
					placeholder="Escreva o conteúdo em Markdown aqui..."
					required
					bind:value={editor.content}
					bind:this={editor.contentInput}
					disabled={controller.isLoading}
					oninput={() => editor.onContentInput()}
					onblur={() => editor.onContentBlur()}
					aria-required="true"
					aria-invalid={editor.hasContentError || controller.hasContentError}
					aria-describedby={editor.hasContentError || controller.hasContentError
						? "content-error"
						: undefined}
					class="w-full flex-1 p-4 bg-transparent text-gray-200 font-mono text-sm
					resize-none border focus:outline-none focus-visible:ring-2 leading-relaxed
					{editor.hasContentError || controller.hasContentError
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-700 focus:border-text-brand focus:ring-text-brand'}"></textarea>
			</section>

			<section
				class="flex flex-col h-150 bg-app-surface border border-gray-800 rounded-xl overflow-hidden"
				aria-label="Preview do conteúdo"
			>
				<div class="p-3 bg-gray-900 border-b border-gray-800 text-xs font-mono text-gray-400">
					Preview em Tempo Real
				</div>
				<div
					aria-live="polite"
					class="p-6 overflow-y-auto prose prose-invert max-w-none font-mono text-sm text-gray-200"
				>
					{#if editor.content.trim()}
						{#await editor.previewPromise}
							<p role="status" class="text-gray-400 italic font-sans text-xs">Gerando preview...</p>
						{:then html}
							{@html html}
						{/await}
					{:else}
						<p class="text-gray-400 italic font-sans text-xs">
							O preview aparecerá aqui conforme você digita...
						</p>
					{/if}
				</div>
			</section>
		</div>

		<div class="flex {editor.hasContentError ? 'justify-between' : 'justify-end'}">
			{#if editor.hasContentError || controller.hasContentError}
				<p id="content-error" role="alert" class="text-xs text-amber-500">
					O conteúdo precisa de no mínimo 10 letras.
				</p>
			{/if}
			<button
				type="submit"
				class="px-6 py-2.5 rounded-lg bg-text-brand text-gray-950 font-semibold hover:bg-blue-400 focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-text-brand focus-visible:ring-offset-gray-900 transition-colors"
			>
				Salvar Algoritmo
			</button>
		</div>
	</form>
</div>
