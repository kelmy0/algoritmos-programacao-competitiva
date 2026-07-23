import { renderMarkdown } from '$lib/services/markdown';

export class AlgorithmEditor {
	name = $state('');
	category = $state('');
	difficulty = $state('beginner');
	content = $state('');

	previewHtml = $state('');
	isReady = $state(false);

	previewPromise = $derived.by(() => {
		if (!this.content.trim()) return Promise.resolve('');
		return renderMarkdown(this.content);
	});

	constructor() {
		this.init();
	}

	private async init() {
		await renderMarkdown('');
		this.isReady = true;
		this.updatePreview();
	}

	async updatePreview() {
		if (!this.content) {
			this.previewHtml = '';
			return;
		}
		this.previewHtml = await renderMarkdown(this.content);
	}

	insertSnippet(startTag: string, endTag: string = '', defaultText: string = '') {
		const textarea = document.getElementById('content-editor') as HTMLTextAreaElement;
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

	getPayload() {
		return {
			name: this.name,
			category: this.category,
			difficulty: this.difficulty,
			content: this.content
		};
	}
}
