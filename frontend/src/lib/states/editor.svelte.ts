import { renderMarkdown } from "$lib/services/markdown";
import type { ApiError } from "$lib/types/api";
import { normalizeApiError, scrollToAndFocus } from "$lib/utils/errors";
import { sanitizeTitle } from "$lib/utils/sanitize";
import { DIFFICULTIES, type Algorithm, type Difficulty } from "$lib/types/algorithm";
import type { AlgorithmPayload } from "$lib/schemas/algorithm";

export const EDITOR_ERRORS: Record<string, string> = {
	INVALID_NAME: "Nome inválido.",
	INVALID_CATEGORY: "Categoria inválida.",
	INVALID_DIFFICULTY: "O nome da dificuldade não é valido.",
	INVALID_CONTENT: "Conteúdo inválido ou muito pequeno."
};

export class AlgorithmEditor {
	#name = $state("");
	#category = $state("");
	#difficulty = $state("beginner");
	#content = $state("");
	#editorError = $state<ApiError | null>(null);

	nameInput = $state<HTMLInputElement | null>(null);
	categoryInput = $state<HTMLInputElement | null>(null);
	difficultyInput = $state<HTMLSelectElement | null>(null);
	contentInput = $state<HTMLTextAreaElement | null>(null);

	touched = $state({
		name: false,
		category: false,
		difficulty: false,
		content: false
	});

	previewPromise = $derived.by(() => {
		const text = this.#content.trim();
		if (!text) return Promise.resolve("");
		return renderMarkdown(text);
	});

	constructor(getAlgorithm: () => Algorithm | null) {
		$effect.pre(() => {
			const algo = getAlgorithm();
			if (algo) this.load(algo);
		});

		renderMarkdown("").then(() => {});
	}

	get name() {
		return this.#name;
	}

	set name(v: string) {
		this.#name = v;
		this.clearEditorError(["INVALID_NAME", "ALGORITHM_INVALID_NAME"]);
	}

	get category() {
		return this.#category;
	}

	set category(v: string) {
		this.#category = v;
		this.clearEditorError(["INVALID_CATEGORY", "ALGORITHM_INVALID_CATEGORY"]);
	}

	get difficulty() {
		return this.#difficulty;
	}

	set difficulty(v: string) {
		this.#difficulty = v;
		this.clearEditorError(["INVALID_DIFFICULTY"]);
	}

	get content() {
		return this.#content;
	}

	set content(v: string) {
		this.#content = v;
		this.clearEditorError(["INVALID_CONTENT", "ALGORITHM_INVALID_CONTENT"]);
	}

	get editorError() {
		return this.#editorError;
	}

	get isNameValid() {
		return this.#name.trim().length >= 3;
	}

	get isCategoryValid() {
		return this.#category.trim().length >= 3;
	}

	get isDifficultyValid() {
		return (DIFFICULTIES as readonly string[]).includes(this.#difficulty);
	}

	get isContentValid() {
		return this.#content.trim().length >= 10;
	}

	get hasNameError() {
		return (
			this.editorError?.code === "INVALID_NAME" ||
			this.editorError?.code === "ALGORITHM_INVALID_NAME" ||
			(this.touched.name && !this.isNameValid)
		);
	}

	get hasCategoryError() {
		return (
			this.editorError?.code === "INVALID_CATEGORY" ||
			this.editorError?.code === "ALGORITHM_INVALID_CATEGORY" ||
			(this.touched.category && !this.isCategoryValid)
		);
	}

	get hasDifficultyError() {
		return (
			this.editorError?.code === "INVALID_DIFFICULTY" ||
			(this.touched.difficulty && !this.isDifficultyValid)
		);
	}

	get hasContentError() {
		return (
			this.editorError?.code === "INVALID_CONTENT" ||
			this.editorError?.code === "ALGORITHM_INVALID_CONTENT" ||
			(this.touched.content && !this.isContentValid)
		);
	}

	onNameBlur() {
		this.touched.name = true;
		this.#name = sanitizeTitle(this.#name);
	}

	onCategoryBlur() {
		this.touched.category = true;
		this.#category = sanitizeTitle(this.#category);
	}

	onDifficultyBlur() {
		this.touched.difficulty = true;
	}
	onContentBlur() {
		this.touched.content = true;
	}

	clearEditorError(codes: string[]) {
		if (this.editorError && codes.includes(this.editorError.code)) {
			this.#editorError = null;
		}
	}

	insertSnippet(startTag: string, endTag = "", defaultText = "") {
		const textarea = this.contentInput;
		if (!textarea) return;

		const start = textarea.selectionStart;
		const end = textarea.selectionEnd;
		const selectedText = this.#content.substring(start, end) || defaultText;

		const replacement = `${startTag}${selectedText}${endTag}`;
		this.#content = this.#content.substring(0, start) + replacement + this.#content.substring(end);

		setTimeout(() => {
			textarea.focus();
			textarea.setSelectionRange(
				start + startTag.length,
				start + startTag.length + selectedText.length
			);
		}, 0);
	}

	load(data: Algorithm) {
		this.#name = data.name || "";
		this.#category = data.category || "";
		this.#difficulty = data.difficulty || "beginner";
		this.#content = data.content || "";
	}

	setApiError(error: ApiError) {
		this.#editorError = error;

		if (error.code === "ALGORITHM_INVALID_NAME") scrollToAndFocus(this.nameInput);
		else if (error.code === "ALGORITHM_INVALID_CATEGORY") scrollToAndFocus(this.categoryInput);
		else if (error.code === "ALGORITHM_INVALID_CONTENT") scrollToAndFocus(this.contentInput);
	}

	getPayload(): AlgorithmPayload | null {
		this.touched = { name: true, category: true, content: true, difficulty: true };

		if (!this.isNameValid) {
			this.#editorError = normalizeApiError("INVALID_NAME", "Nome inválido.", EDITOR_ERRORS);
			scrollToAndFocus(this.nameInput);
			return null;
		}

		if (!this.isCategoryValid) {
			this.#editorError = normalizeApiError(
				"INVALID_CATEGORY",
				"Categoria inválida.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.categoryInput);
			return null;
		}

		if (!this.isDifficultyValid) {
			this.#editorError = normalizeApiError(
				"INVALID_DIFFICULTY",
				"Dificuldade inválida.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.difficultyInput);
			return null;
		}

		if (!this.isContentValid) {
			this.#editorError = normalizeApiError(
				"INVALID_CONTENT",
				"Conteúdo muito pequeno.",
				EDITOR_ERRORS
			);
			scrollToAndFocus(this.contentInput);
			return null;
		}

		return {
			name: this.#name,
			category: this.#category,
			difficulty: this.#difficulty as Difficulty,
			content: this.#content
		};
	}
}
