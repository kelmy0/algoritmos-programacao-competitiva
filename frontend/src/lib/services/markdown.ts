import { Marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
import { createHighlighter, type Highlighter } from 'shiki';

let highlighterPromise: Promise<Highlighter> | null = null;
let markedInstance: Marked | null = null;

export async function getHighlighter(): Promise<Highlighter> {
	if (!highlighterPromise) {
		highlighterPromise = createHighlighter({
			themes: ['github-dark'],
			langs: ['cpp', 'python']
		});
	}
	return highlighterPromise;
}

export async function getMarked(): Promise<Marked> {
	if (markedInstance) return markedInstance;

	const highlighter = await getHighlighter();

	markedInstance = new Marked(
		markedHighlight({
			async: true,
			highlight(code, lang) {
				const loadedLangs = highlighter.getLoadedLanguages();
				const language = loadedLangs.includes(lang) ? lang : 'text';

				return highlighter.codeToHtml(code, {
					lang: language,
					theme: 'github-dark'
				});
			}
		})
	);

	return markedInstance;
}

export async function renderMarkdown(content: string): Promise<string> {
	if (!content.trim()) return '';
	const marked = await getMarked();
	return await marked.parse(content);
}
