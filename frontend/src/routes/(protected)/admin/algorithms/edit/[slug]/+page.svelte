<script lang="ts">
	import Alert from "$lib/components/ui/Alert.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import Input from "$lib/components/ui/Input.svelte";
	import MarkdownEditor from "$lib/components/ui/MarkdownEditor.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";
	import type { SelectOption } from "$lib/components/ui/Select.svelte";
	import Select from "$lib/components/ui/Select.svelte";
	import { focusTrap } from "$lib/utils/a11y";
	import { AlgorithmEditor } from "$lib/utils/editor.svelte";
	import type { PageData } from "./$types";
	import { EditAlgorithmController } from "./editAlgorithm.svelte";

	let { data }: { data: PageData } = $props();

	const editor = new AlgorithmEditor();
	const controller = new EditAlgorithmController();

	let isInitialized = $state(false);
	let localStatus = $state<string>();
	let status = $derived(localStatus ?? data.algorithm?.status);

	$effect(() => {
		if (data.algorithm && !isInitialized) {
			editor.load(data.algorithm);
			controller.publicId = data.algorithm.publicId;
			controller.slug = data.algorithm.slug;
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

	type ActionType = "save" | "delete" | "restore";

	const errorLabels: Record<ActionType, string> = {
		save: "salvar",
		delete: "deletar",
		restore: "restaurar"
	};

	const successLabels: Record<ActionType, string> = {
		save: "salvo",
		delete: "deletado",
		restore: "restaurado"
	};

	const successText: Record<ActionType, string> = {
		save: "Suas alterações foram enviadas e já estão em espera para aprovação.",
		delete: "Algoritmo movido para a lixeira! Você tem até 7 dias para restaurá-lo.",
		restore: "Algoritmo restaurado com sucesso! Ele já está em espera para aprovação."
	};

	const difficultyOptions: SelectOption[] = [
		{ value: "beginner", label: "Iniciante" },
		{ value: "intermediate", label: "Intermediário" },
		{ value: "advanced", label: "Avançado" },
		{ value: "expert", label: "Especialista" }
	];
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

			<Input
				id="name"
				label="Nome do Algoritmo *"
				placeholder="Ex: Busca Binária"
				required
				autocomplete="off"
				disabled={controller.isLoading || controller.isDeleting}
				bind:value={editor.name}
				bind:inputRef={editor.nameInput}
				touched={editor.hasNameError || controller.hasNameError}
				error="O nome precisa ter no mínimo 3 letras válidas."
				oninput={() => editor.onNameInput()}
				onblur={() => editor.onNameBlur()}
			/>

			<Input
				id="category"
				label="Categoria"
				placeholder="Ex: Grafos, Busca, DP"
				required
				autocomplete="off"
				disabled={controller.isLoading || controller.isDeleting}
				bind:value={editor.category}
				bind:inputRef={editor.categoryInput}
				touched={editor.hasCategoryError || controller.hasCategoryError}
				error="A categoria precisa ter no mínimo 3 letras válidas."
				oninput={() => editor.onCategoryInput()}
				onblur={() => editor.onCategoryBlur()}
			/>

			<Select
				id="difficulty"
				label="Dificuldade"
				required
				options={difficultyOptions}
				bind:value={editor.difficulty}
				bind:selectRef={editor.difficultyInput}
				disabled={controller.isLoading || controller.isDeleting}
				touched={editor.hasDifficultyError}
				error="A dificuldade precisa ser uma das 4 opções."
				onchange={() => editor.onDifficultyInput()}
				onblur={() => editor.onDifficultyBlur()}
			/>
		</fieldset>

		<div class="invisible" bind:this={controller.alertDiv}></div>
		{#if controller.apiError}
			<Alert
				type="error"
				title="Erro ao {errorLabels[controller.lastAction] ?? controller.lastAction} algoritmo"
				message={controller.apiError.message ||
					"Ocorreu um erro ao processar sua solicitação. Tente novamente."}
				isLoading={controller.isLoading || controller.isDeleting}
			/>
		{/if}

		{#if controller.isSuccess}
			{@const label = successLabels[controller.lastAction] ?? controller.lastAction}
			{@const text = successText[controller.lastAction] ?? controller.lastAction}
			{@const isDelete = controller.lastAction === "delete"}

			<Alert type={isDelete ? "warning" : "success"} title="Algoritmo {label} com sucesso!">
				{text}
				{#if controller.link}
					<a
						href={controller.link}
						class="underline hover:opacity-80 transition-opacity"
						onclick={() => (controller.isSuccess = false)}
					>
						Visualizar o algoritmo {label}
					</a>.
				{/if}
			</Alert>
		{/if}

		<MarkdownEditor
			bind:content={editor.content}
			bind:contentInput={editor.contentInput}
			hasError={editor.hasContentError || controller.hasContentError}
			disabled={controller.isLoading || controller.isDeleting}
			previewPromise={editor.previewPromise}
			previewUrl={controller.slug
				? `/admin/algorithms/my-algorithms/${controller.slug}-${controller.publicId}`
				: undefined}
			insertSnippet={(before, after, placeholder) =>
				editor.insertSnippet(before, after, placeholder)}
			onContentInput={() => editor.onContentInput()}
			onContentBlur={() => editor.onContentBlur()}
		/>

		<div
			class="flex flex-col-reverse sm:flex-row justify-between items-stretch sm:items-center gap-4 pt-4 border-t border-gray-800/60"
		>
			{#if status !== "deleted"}
				<Button
					variant="danger-soft"
					onclick={() => controller.openDeleteModal()}
					disabled={controller.isLoading || controller.isDeleting}
					aria-haspopup="dialog"
					aria-expanded={controller.isDeleteModalOpen}
				>
					<svg
						class="w-4 h-4 shrink-0"
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
				</Button>
			{:else}
				<Button
					variant="success-soft"
					onclick={() => handleRestore()}
					disabled={controller.isLoading || controller.isDeleting}
				>
					<svg
						class="w-4 h-4 shrink-0"
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
				</Button>
			{/if}

			<div class="flex flex-col md:flex-row items-stretch md:items-center gap-4 justify-end">
				{#if editor.hasContentError || controller.hasContentError}
					<p id="content-error" role="alert" class="text-xs text-amber-500 self-center">
						O conteúdo precisa de no mínimo 10 letras.
					</p>
				{/if}
				<Button isLoading={controller.isLoading} disabled={controller.isLoading}>
					{controller.isLoading ? "Salvando..." : "Salvar algoritmo"}
				</Button>
			</div>
		</div>
	</form>
</div>

{#if controller.isDeleteModalOpen}
	<Modal
		isOpen={controller.isDeleteModalOpen}
		onClose={() => controller.closeDeleteModal()}
		title="Excluir Algoritmo?"
		description="Tem certeza que deseja mover para lixeira este algoritmo? Você terá 7 dias para restaurá-lo."
		variant="danger"
		isLoading={controller.isDeleting || controller.isLoading}
		{focusTrap}
	>
		{#snippet footer()}
			<Button
				variant="outline"
				onclick={() => controller.closeDeleteModal()}
				disabled={controller.isDeleting || controller.isLoading}
			>
				Cancelar
			</Button>

			<Button
				variant="danger"
				onclick={handleDelete}
				isLoading={controller.isDeleting}
				disabled={controller.isLoading}
			>
				{controller.isDeleting ? "Excluindo..." : "Sim, Excluir"}
			</Button>
		{/snippet}
	</Modal>
{/if}
