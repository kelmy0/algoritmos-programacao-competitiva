import { renderMarkdown } from "$lib/services/markdown";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "./errors";
import { sanitizeTitle } from "./sanitize";
import { DIFFICULTIES, type Algorithm, type Difficulty } from "$lib/types/algorithm";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";

export const EDITOR_ERRORS: Record<string, string> = {
	INVALID_NAME: "Nome inválido.",
	INVALID_CATEGORY: "Categoria inválida.",
	INVALID_DIFFICULTY: "O nome da dificuldade não é valido.",
	INVALID_CONTENT: "Conteúdo inválido ou muito pequeno."
};

export class AlgorithmEditor {
	name = $state("");
	category = $state("");
	difficulty = $state("beginner");
	content = $state("");

	nameInput = $state<HTMLInputElement | null>(null);
	categoryInput = $state<HTMLInputElement | null>(null);
	difficultyInput = $state<HTMLSelectElement | null>(null);
	contentInput = $state<HTMLTextAreaElement | null>(null);

	previewHtml = $state("");
	isReady = $state(false);

	editorError = $state<ApiError | null>(null);

	touched = $state({
		name: false,
		category: false,
		difficulty: false,
		content: false
	});

	previewPromise = $derived.by(() => {
		if (!this.content.trim()) return Promise.resolve("");
		return renderMarkdown(this.content);
	});

	hasNameError = $derived(
		this.editorError?.code === "INVALID_NAME" || (this.touched.name && !this.isNameValid)
	);

	hasCategoryError = $derived(
		this.editorError?.code === "INVALID_CATEGORY" ||
			(this.touched.category && !this.isCategoryValid)
	);

	hasDifficultyError = $derived(
		this.editorError?.code === "INVALID_DIFFICULTY" ||
			(this.touched.difficulty && !this.isDifficultyValid)
	);

	hasContentError = $derived(
		this.editorError?.code === "INVALID_CONTENT" || (this.touched.content && !this.isContentValid)
	);

	constructor() {
		this.init();
	}

	get isNameValid() {
		return this.name.length >= 3;
	}

	get isCategoryValid() {
		return this.category.length >= 3;
	}

	get isDifficultyValid() {
		return (DIFFICULTIES as readonly string[]).includes(this.difficulty);
	}

	get isContentValid() {
		return this.content.length >= 10;
	}

	load(data: Algorithm) {
		this.name = data.name || "";
		this.category = data.category || "";
		this.difficulty = data.difficulty || "beginner";
		this.content = data.content || "";
	}

	onNameBlur() {
		this.touched.name = true;
		this.name = sanitizeTitle(this.name);
	}

	onCategoryBlur() {
		this.touched.category = true;
		this.category = sanitizeTitle(this.category);
	}

	onDifficultyBlur() {
		this.touched.difficulty = true;
	}

	onContentBlur() {
		this.touched.content = true;
	}

	onNameInput() {
		this.clearEditorError(["INVALID_NAME"]);
	}

	onCategoryInput() {
		this.clearEditorError(["INVALID_CATEGORY"]);
	}

	onDifficultyInput() {
		this.clearEditorError(["INVALID_DIFFICULTY"]);
	}

	onContentInput() {
		this.clearEditorError(["INVALID_CONTENT"]);
	}

	clearEditorError(codes: string[]) {
		if (this.editorError && codes.includes(this.editorError.code)) {
			this.editorError = null;
		}
	}

	private async init() {
		await renderMarkdown("");
		this.isReady = true;
		this.updatePreview();
	}

	async updatePreview() {
		if (!this.content) {
			this.previewHtml = "";
			return;
		}
		this.previewHtml = await renderMarkdown(this.content);
	}

	insertSnippet(startTag: string, endTag: string = "", defaultText: string = "") {
		const textarea = document.getElementById("content-editor") as HTMLTextAreaElement;
		if (!textarea) return;

		const start = textarea.selectionStart;
		const end = textarea.selectionEnd;
		const selectedText = this.content.substring(start, end) || defaultText;

		const replacement = `${startTag}${selectedText}${endTag}`;
		this.content = this.content.substring(0, start) + replacement + this.content.substring(end);

		this.updatePreview();

		setTimeout(() => {
			textarea.focus();
			textarea.setSelectionRange(
				start + startTag.length,
				start + startTag.length + selectedText.length
			);
		}, 0);
	}

	getPayload(): AlgorithmPayload | null {
		this.touched = {
			name: true,
			category: true,
			content: true,
			difficulty: true
		};

		if (!this.isNameValid) {
			this.editorError = normalizeApiError("INVALID_NAME", "Nome inválido.", EDITOR_ERRORS);
			scrollToAndFocus(this.nameInput);
			return null;
		}

		if (!this.isCategoryValid) {
			this.editorError = normalizeApiError(
				"INVALID_CATEGORY",
				"Categoria inválida.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.categoryInput);
			return null;
		}

		if (!this.isDifficultyValid) {
			this.editorError = normalizeApiError(
				"INVALID_DIFFICULTY",
				"O nome da dificuldade não é valido.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.difficultyInput);
			return null;
		}

		if (!this.isContentValid) {
			this.editorError = normalizeApiError(
				"INVALID_CONTENT",
				"Conteúdo inválido ou muito pequeno.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.contentInput);
			return null;
		}

		return {
			name: this.name,
			category: this.category,
			difficulty: this.difficulty as Difficulty,
			content: this.content
		};
	}
}
