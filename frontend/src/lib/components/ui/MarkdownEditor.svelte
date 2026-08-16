<script lang="ts">
	interface Props {
		content: string;
		contentInput?: HTMLTextAreaElement | null;
		hasError?: boolean;
		disabled?: boolean;
		previewPromise?: Promise<string>;
		previewUrl?: string;
		insertSnippet: (before: string, after: string, placeholder: string) => void;
		onContentInput?: () => void;
		onContentBlur?: () => void;
	}

	let {
		content = $bindable(""),
		contentInput = $bindable(null),
		hasError = false,
		disabled = false,
		previewPromise,
		previewUrl,
		insertSnippet,
		onContentInput,
		onContentBlur
	}: Props = $props();
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
	<section
		class="flex flex-col h-150 bg-app-surface border border-app-border rounded-xl overflow-hidden"
		aria-label="Editor de código Markdown"
	>
		<div
			role="toolbar"
			aria-label="Ferramentas de formatação Markdown"
			class="min-h-12 flex flex-wrap items-center gap-1.5 px-3 py-2 bg-app-overlay border-b border-app-border text-xs"
		>
			<button
				type="button"
				onclick={() => insertSnippet("## ", "", "Título")}
				aria-label="Inserir Título Nível 2"
				class="btn-toolbar hover:cursor-pointer">H2</button
			>
			<button
				type="button"
				onclick={() => insertSnippet("### ", "", "Subtítulo")}
				aria-label="Inserir Subtítulo Nível 3"
				class="btn-toolbar hover:cursor-pointer">H3</button
			>
			<button
				type="button"
				onclick={() => insertSnippet("**", "**", "negrito")}
				aria-label="Texto em Negrito"
				class="btn-toolbar hover:cursor-pointer"><b aria-hidden="true">B</b></button
			>
			<button
				type="button"
				onclick={() => insertSnippet("*", "*", "itálico")}
				aria-label="Texto em Itálico"
				class="btn-toolbar hover:cursor-pointer"><i aria-hidden="true">I</i></button
			>
			<button
				type="button"
				onclick={() => insertSnippet("\n```cpp\n", "\n```\n", "// seu código C++ aqui")}
				aria-label="Inserir bloco de código C++"
				class="btn-toolbar hover:cursor-pointer font-mono text-text-brand">C++ Code</button
			>
			<button
				type="button"
				onclick={() => insertSnippet("> ", "", "Nota importante")}
				aria-label="Inserir citação"
				class="btn-toolbar hover:cursor-pointer">Quote</button
			>
			<button
				type="button"
				onclick={() => insertSnippet("1. ", "", "Item")}
				aria-label="Inserir lista numerada"
				class="btn-toolbar hover:cursor-pointer">Lista</button
			>
		</div>

		<label for="content-editor" class="sr-only">Conteúdo em Markdown</label>
		<textarea
			id="content-editor"
			placeholder="Escreva o conteúdo em Markdown aqui..."
			required
			bind:value={content}
			bind:this={contentInput}
			{disabled}
			oninput={onContentInput}
			onblur={onContentBlur}
			aria-required="true"
			aria-invalid={hasError}
			aria-describedby={hasError ? "content-error" : undefined}
			class="w-full flex-1 p-4 bg-transparent text-text-primary font-mono text-sm resize-none focus:outline-none focus-visible:ring-2 leading-relaxed border disabled:cursor-not-allowed
            {hasError
				? 'border-red-500 focus:ring-red-500'
				: 'border-transparent focus:ring-text-brand'}"></textarea>
	</section>

	<section
		class="flex flex-col h-150 bg-app-surface border border-app-border rounded-xl overflow-hidden"
		aria-label="Preview do conteúdo"
	>
		<div
			class="min-h-12 px-4 py-2 bg-app-overlay border-b border-app-border flex items-center justify-between gap-2"
		>
			<span class="text-xs font-mono font-medium text-text-secondary">Preview em Tempo Real</span>
			{#if previewUrl}
				<a
					href={previewUrl}
					target="_blank"
					rel="noopener noreferrer"
					class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium
					text-text-secondary bg-app-surface/60 hover:bg-app-surface hover:text-white border border-app-border
					focus:outline-none focus-visible:ring-2 focus-visible:ring-text-brand transition-all"
					title="Abrir visualização em nova aba"
				>
					<svg
						class="w-4 h-4 text-text-secondary group-hover:text-white"
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
			class="p-6 whitespace-pre-wrap wrap-break-word overflow-y-auto prose prose-invert max-w-none font-mono text-sm text-text-primary prose-pre:whitespace-pre-wrap prose-pre:wrap-break-words"
		>
			{#if content.trim()}
				{#await previewPromise}
					<p role="status" class="text-text-secondary italic font-sans text-sm">
						Gerando preview...
					</p>
				{:then html}
					{@html html}
				{/await}
			{:else}
				<p class="text-text-secondary italic font-sans text-sm">
					O preview aparecerá aqui conforme você digita...
				</p>
			{/if}
		</div>
	</section>
</div>
