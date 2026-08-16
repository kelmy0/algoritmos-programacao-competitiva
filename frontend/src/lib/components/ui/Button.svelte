<script lang="ts">
	import type { HTMLButtonAttributes } from "svelte/elements";
	import type { Snippet } from "svelte";

	type Variant =
		| "primary"
		| "secondary"
		| "success"
		| "success-soft"
		| "warning"
		| "danger"
		| "danger-soft"
		| "dark"
		| "outline"
		| "ghost";
	type Size = "sm" | "md" | "lg";

	interface Props extends HTMLButtonAttributes {
		variant?: Variant;
		size?: Size;
		isLoading?: boolean;
		children?: Snippet;
	}

	let {
		type = "button",
		variant = "primary",
		size = "md",
		isLoading = false,
		disabled = false,
		class: className = "",
		children,
		...restProps
	}: Props = $props();

	const variantClasses: Record<Variant, string> = {
		primary:
			"bg-text-brand text-app-bg font-semibold hover:opacity-90 focus-visible:ring-text-brand",
		secondary:
			"bg-app-overlay text-text-primary border border-app-border hover:bg-app-border hover:border-app-border-hover focus-visible:ring-text-brand",
		success:
			"bg-emerald-600 text-white font-medium hover:bg-emerald-500 focus-visible:ring-emerald-500",
		"success-soft":
			"border border-emerald-900/50 bg-emerald-950/30 text-emerald-400 font-semibold hover:bg-emerald-900/50 hover:text-emerald-300 focus-visible:ring-emerald-500",
		warning:
			"bg-amber-500 text-app-bg font-semibold hover:bg-amber-400 focus-visible:ring-amber-400",
		danger: "bg-red-600 text-white font-medium hover:bg-red-500 focus-visible:ring-red-500",
		"danger-soft":
			"border border-red-900/50 bg-red-950/30 text-red-400 font-semibold hover:bg-red-900/50 hover:text-red-300 focus-visible:ring-red-500",
		dark: "bg-app-bg text-text-primary border border-app-border hover:bg-app-surface focus-visible:ring-text-brand",
		outline:
			"bg-transparent border border-app-border text-text-secondary hover:text-text-primary hover:border-app-border-hover focus-visible:ring-text-brand",
		ghost:
			"bg-transparent text-text-secondary hover:bg-app-overlay hover:text-text-primary focus-visible:ring-text-brand"
	};

	const sizeClasses: Record<Size, string> = {
		sm: "py-1.5 px-3 text-xs rounded-md gap-1.5",
		md: "py-2.5 px-5 text-sm rounded-lg gap-2",
		lg: "py-3.5 px-6 text-base rounded-xl gap-2.5"
	};
</script>

<button
	{type}
	disabled={disabled || isLoading}
	aria-busy={isLoading}
	class="inline-flex items-center justify-center font-inter cursor-pointer active:scale-[0.98] transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900 disabled:opacity-50 disabled:pointer-events-none disabled:cursor-not-allowed {sizeClasses[
		size
	]} {variantClasses[variant]} {className}"
	{...restProps}
>
	{#if isLoading}
		<svg
			class="animate-spin h-4 w-4 currentColor shrink-0"
			fill="none"
			viewBox="0 0 24 24"
			aria-hidden="true"
		>
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
			></circle>
			<path
				class="opacity-75"
				fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>
	{/if}

	{#if children}
		{@render children()}
	{/if}
</button>
