<script lang="ts">
	import type { Snippet } from "svelte";

	type ModalVariant = "danger" | "warning" | "info" | "success";

	interface Props {
		isOpen: boolean;
		title: string;
		description?: string;
		variant?: ModalVariant;
		isLoading?: boolean;
		onClose: () => void;
		focusTrap?: (node: HTMLElement) => { destroy?: () => void };
		icon?: Snippet;
		children?: Snippet;
		footer?: Snippet;
	}

	let {
		isOpen,
		title,
		description,
		variant = "info",
		isLoading = false,
		onClose,
		focusTrap,
		icon,
		children,
		footer
	}: Props = $props();

	const titleId = $derived(`modal-title-${Math.random().toString(36).substring(2, 9)}`);
	const descId = $derived(`modal-desc-${Math.random().toString(36).substring(2, 9)}`);

	const variantStyles = $derived(
		{
			danger: {
				iconBg: "bg-red-950/80 border-red-900/60 text-red-400",
				iconPath:
					"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
			},
			warning: {
				iconBg: "bg-amber-950/80 border-amber-900/60 text-amber-400",
				iconPath:
					"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
			},
			info: {
				iconBg: "bg-blue-950/80 border-blue-900/60 text-blue-400",
				iconPath: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			},
			success: {
				iconBg: "bg-emerald-950/80 border-emerald-900/60 text-emerald-400",
				iconPath: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
			}
		}[variant]
	);

	function useTrap(node: HTMLElement) {
		if (focusTrap) {
			return focusTrap(node);
		}
	}
</script>

{#if isOpen}
	<div
		use:useTrap
		class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm overflow-y-auto flex min-h-full items-center justify-center p-4"
		role="dialog"
		aria-modal="true"
		aria-busy={isLoading}
		aria-labelledby={titleId}
		aria-describedby={description ? descId : undefined}
		onkeydown={(e) => e.key === "Escape" && !isLoading && onClose()}
		tabindex="-1"
	>
		<div
			class="bg-app-surface border border-gray-800 rounded-xl p-6 max-w-md w-full flex flex-col gap-5 shadow-2xl animate-in fade-in zoom-in-95 duration-150 relative my-auto"
		>
			<div class="flex items-start gap-3">
				<div class="p-2.5 self-start rounded-lg shrink-0 border {variantStyles.iconBg}">
					{#if icon}
						{@render icon()}
					{:else}
						<svg
							class="w-6 h-6"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d={variantStyles.iconPath}
							/>
						</svg>
					{/if}
				</div>

				<div class="flex-1 pr-6">
					<h2 id={titleId} class="text-lg font-bold text-gray-100 font-montserrat">
						{title}
					</h2>
					{#if description}
						<p id={descId} class="text-sm text-gray-300 mt-1 leading-relaxed">
							{description}
						</p>
					{/if}
				</div>

				<button
					type="button"
					onclick={onClose}
					disabled={isLoading}
					aria-label="Fechar modal"
					class="hover:cursor-pointer absolute top-4 right-4 text-gray-400 hover:text-white p-1 rounded-lg focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 disabled:opacity-50 transition-colors"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			{#if children}
				<div class="flex flex-col gap-4">
					{@render children()}
				</div>
			{/if}

			{#if footer}
				<div
					class="flex flex-col-reverse sm:flex-row justify-end gap-3 pt-4 border-t border-gray-800"
				>
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
