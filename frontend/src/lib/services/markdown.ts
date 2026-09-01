import type { Marked } from "marked";
import type { HighlighterCore } from "shiki";
import type DOMPurifyType from "isomorphic-dompurify";

let highlighterPromise: Promise<HighlighterCore> | null = null;
let markedInstance: Marked | null = null;
let purifyInstance: typeof DOMPurifyType | null = null;

export async function getHighlighter(): Promise<HighlighterCore> {
	if (!highlighterPromise) {
		const [
			{ createHighlighterCore },
			{ createJavaScriptRegexEngine },
			langCpp,
			langPython,
			themeGithubDark
		] = await Promise.all([
			import("shiki/core"),
			import("shiki/engine/javascript"),
			import("@shikijs/langs/cpp"),
			import("@shikijs/langs/python"),
			import("@shikijs/themes/github-dark")
		]);

		highlighterPromise = createHighlighterCore({
			themes: [themeGithubDark],
			langs: [langCpp, langPython],
			engine: createJavaScriptRegexEngine()
		});
	}
	return highlighterPromise;
}

export async function getMarked(): Promise<Marked> {
	if (markedInstance) return markedInstance;

	const [{ Marked }, { markedHighlight }, highlighter] = await Promise.all([
		import("marked"),
		import("marked-highlight"),
		getHighlighter()
	]);

	markedInstance = new Marked(
		markedHighlight({
			async: true,
			highlight(code, lang) {
				const loadedLangs = highlighter.getLoadedLanguages();
				const language = loadedLangs.includes(lang) ? lang : "text";

				return highlighter.codeToHtml(code, {
					lang: language,
					theme: "github-dark"
				});
			}
		})
	);

	return markedInstance;
}

export async function renderMarkdown(content: string): Promise<string> {
	if (!content.trim()) return "";

	if (!purifyInstance) {
		const DOMPurifyModule = await import("isomorphic-dompurify");
		purifyInstance = DOMPurifyModule.default;
	}

	const marked = await getMarked();
	const rawHtml = await marked.parse(content);
	const sanitizedHtml = purifyInstance.sanitize(rawHtml, { ADD_ATTR: ["target"] });
	return sanitizedHtml;
}
