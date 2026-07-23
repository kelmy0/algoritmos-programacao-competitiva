<script lang="ts">
	import { AlgorithmEditor } from './editor.svelte';

	const editor = new AlgorithmEditor();

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const payload = editor.getPayload();
		console.log(payload);
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
					bind:value={editor.name}
					placeholder="Ex: Busca Binária"
					required
					aria-required="true"
					class="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand"
				/>
			</div>

			<div class="space-y-2">
				<label for="category" class="block text-sm font-medium text-gray-200">
					Categoria <span class="text-red-400" aria-hidden="true">*</span>
				</label>
				<input
					id="category"
					type="text"
					bind:value={editor.category}
					placeholder="Ex: Grafos, Busca, DP"
					required
					aria-required="true"
					class="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand"
				/>
			</div>

			<div class="space-y-2">
				<label for="difficulty" class="block text-sm font-medium text-gray-200">Dificuldade</label>
				<select
					id="difficulty"
					bind:value={editor.difficulty}
					class="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand"
				>
					<option value="beginner">Iniciante</option>
					<option value="intermediate">Intermediário</option>
					<option value="advanced">Avançado</option>
					<option value="expert">Especialista</option>
				</select>
			</div>
		</fieldset>

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
						onclick={() => editor.insertSnippet('## ', '', 'Título')}
						aria-label="Inserir Título Nível 2"
						class="btn-toolbar">H2</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('### ', '', 'Subtítulo')}
						aria-label="Inserir Subtítulo Nível 3"
						class="btn-toolbar">H3</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('**', '**', 'negrito')}
						aria-label="Texto em Negrito"
						class="btn-toolbar"><b aria-hidden="true">B</b></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('*', '*', 'itálico')}
						aria-label="Texto em Itálico"
						class="btn-toolbar"><i aria-hidden="true">I</i></button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('\n```cpp\n', '\n```\n', '// seu código C++ aqui')}
						aria-label="Inserir bloco de código C++"
						class="btn-toolbar font-mono text-text-brand">C++ Code</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('> ', '', 'Nota importante')}
						aria-label="Inserir citação"
						class="btn-toolbar">Quote</button
					>
					<button
						type="button"
						onclick={() => editor.insertSnippet('1. ', '', 'Item')}
						aria-label="Inserir lista numerada"
						class="btn-toolbar">Lista</button
					>
				</div>

				<label for="content-editor" class="sr-only">Conteúdo em Markdown</label>
				<textarea
					id="content-editor"
					bind:value={editor.content}
					placeholder="Escreva o conteúdo em Markdown aqui..."
					required
					aria-required="true"
					class="w-full flex-1 p-4 bg-transparent text-gray-200 font-mono text-sm resize-none focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand leading-relaxed"
				></textarea>
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

		<div class="flex justify-end">
			<button
				type="submit"
				class="px-6 py-2.5 rounded-lg bg-text-brand text-gray-950 font-semibold hover:bg-blue-400 focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-text-brand focus-visible:ring-offset-gray-900 transition-colors"
			>
				Salvar Algoritmo
			</button>
		</div>
	</form>
</div>
