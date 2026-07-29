export function focusTrap(node: HTMLElement) {
	const previouslyFocusedElement = document.activeElement as HTMLElement | null;
	const getFocusableElements = () => {
		return Array.from(
			node.querySelectorAll<HTMLElement>(
				'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
			)
		).filter((el) => !el.hasAttribute("disabled"));
	};

	const focusables = getFocusableElements();
	if (focusables.length > 0) {
		focusables[0].focus();
	} else {
		node.focus();
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key !== "Tab") return;

		const currentFocusables = getFocusableElements();
		if (currentFocusables.length === 0) return;

		const firstElement = currentFocusables[0];
		const lastElement = currentFocusables[currentFocusables.length - 1];

		if (e.shiftKey && document.activeElement === firstElement) {
			e.preventDefault();
			lastElement.focus();
		} else if (!e.shiftKey && document.activeElement === lastElement) {
			e.preventDefault();
			firstElement.focus();
		}
	}

	node.addEventListener("keydown", handleKeyDown);

	return {
		destroy() {
			node.removeEventListener("keydown", handleKeyDown);
			if (previouslyFocusedElement && typeof previouslyFocusedElement.focus === "function") {
				setTimeout(() => {
					previouslyFocusedElement.focus();
				}, 0);
			}
		}
	};
}

export function clickOutside(node: HTMLElement, callback?: () => void) {
	function handleClick(event: MouseEvent) {
		if (node && !node.contains(event.target as Node)) {
			callback?.();
		}
	}

	document.addEventListener("click", handleClick, true);

	return {
		destroy() {
			document.removeEventListener("click", handleClick, true);
		}
	};
}
