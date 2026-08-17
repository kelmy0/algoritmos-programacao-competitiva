<script lang="ts">
	import Alert from "$lib/components/ui/Alert.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import Input from "$lib/components/ui/Input.svelte";
	import MarkdownEditor from "$lib/components/ui/MarkdownEditor.svelte";
	import Select, { type SelectOption } from "$lib/components/ui/Select.svelte";
	import { AlgorithmEditor } from "$lib/states/editor.svelte";
	import { scrollToAndFocus } from "$lib/utils/errors";
	import { NewAlgorithmController } from "./newAlgorithm.svelte";

	const editor = new AlgorithmEditor(() => null);
	const controller = new NewAlgorithmController();

	let alertDiv = $state<HTMLDivElement | null>(null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		await controller.save(editor);
		await scrollToAndFocus(alertDiv);
	}

	const difficultyOptions: SelectOption[] = [
		{ value: "beginner", label: "Iniciante" },
		{ value: "intermediate", label: "Intermediário" },
		{ value: "advanced", label: "Avançado" },
		{ value: "expert", label: "Especialista" }
	];
</script>

<svelte:head>
	<title>Criar Algoritmo</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6 font-inter p-6">
	<header class="border-b border-app-border pb-4">
		<h1 class="font-montserrat text-2xl font-bold text-text-primary">Criar Novo Algoritmo</h1>
		<p class="text-sm text-text-secondary mt-1">
			Preencha os metadados e escreva o conteúdo em Markdown com o preview ao lado.
		</p>
	</header>

	<form onsubmit={handleSubmit} class="space-y-6" aria-label="Formulário de criação de algoritmo">
		<fieldset
			class="grid grid-cols-1 md:grid-cols-3 gap-4 bg-app-surface p-5 rounded-xl border border-app-border"
		>
			<legend class="sr-only">Metadados do Algoritmo</legend>

			<Input
				id="name"
				label="Nome do Algoritmo *"
				placeholder="Ex: Busca Binária"
				required
				autocomplete="off"
				disabled={controller.isLoading}
				bind:value={editor.name}
				bind:inputRef={editor.nameInput}
				touched={editor.hasNameError}
				error="O nome precisa ter no mínimo 3 letras válidas."
				onblur={() => editor.onNameBlur()}
			/>

			<Input
				id="category"
				label="Categoria"
				placeholder="Ex: Grafos, Busca, DP"
				required
				autocomplete="off"
				disabled={controller.isLoading}
				bind:value={editor.category}
				bind:inputRef={editor.categoryInput}
				touched={editor.hasCategoryError}
				error="A categoria precisa ter no mínimo 3 letras válidas."
				onblur={() => editor.onCategoryBlur()}
			/>

			<Select
				id="difficulty"
				label="Dificuldade"
				required
				options={difficultyOptions}
				bind:value={editor.difficulty}
				bind:selectRef={editor.difficultyInput}
				disabled={controller.isLoading}
				touched={editor.hasDifficultyError}
				error="A dificuldade precisa ser uma das 4 opções."
				onblur={() => editor.onDifficultyBlur()}
			/>
		</fieldset>

		<div class="invisible" bind:this={alertDiv}></div>
		{#if controller.apiError}
			<Alert
				type="error"
				title="Erro ao salvar algoritmo"
				message={controller.apiError.message ||
					"Ocorreu um erro ao tentar salvar. Tente novamente."}
				isLoading={controller.isLoading}
			/>
		{/if}

		{#if controller.isSuccess}
			<Alert type="success" title="Algoritmo salvo com sucesso!">
				Suas alterações foram enviadas e já estão em espera para aprovação.
				{#if controller.link}
					<a href={controller.link} class="underline hover:text-emerald-200">
						Visualizar o algoritmo enviado
					</a>.
				{/if}
			</Alert>
		{/if}

		<MarkdownEditor
			bind:content={editor.content}
			bind:contentInput={editor.contentInput}
			hasError={editor.hasContentError}
			disabled={controller.isLoading}
			previewPromise={editor.previewPromise}
			insertSnippet={(before, after, placeholder) =>
				editor.insertSnippet(before, after, placeholder)}
			onContentBlur={() => editor.onContentBlur()}
		/>

		<div
			class="flex flex-col md:flex-row {editor.hasContentError
				? 'justify-between'
				: 'justify-end'} items-stretch md:items-center gap-4 pt-4 border-t border-app-border/60"
		>
			{#if editor.hasContentError}
				<p id="content-error" role="alert" class="text-xs text-amber-500 self-center">
					O conteúdo precisa de no mínimo 10 letras.
				</p>
			{/if}
			<Button type="submit" isLoading={controller.isLoading} disabled={controller.isLoading}>
				{controller.isLoading ? "Salvando..." : "Salvar algoritmo"}
			</Button>
		</div>
	</form>
</div>
