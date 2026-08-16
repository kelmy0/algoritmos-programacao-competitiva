<script lang="ts">
	import type { HTMLInputAttributes } from "svelte/elements";
	import type { Snippet } from "svelte";

	interface Props extends HTMLInputAttributes {
		label?: string;
		error?: string;
		touched?: boolean;
		helperText?: string;
		value?: string | number;
		prefixIcon?: Snippet;
		suffixIcon?: Snippet;
		inputRef?: HTMLInputElement | null;
	}

	let {
		id,
		label,
		error,
		touched = false,
		helperText,
		value = $bindable(""),
		prefixIcon,
		suffixIcon,
		class: className = "",
		required,
		disabled,
		inputRef = $bindable(null),
		...restProps
	}: Props = $props();

	const hasError = $derived(Boolean(touched && error));
	const errorId = $derived(id ? `${id}-error` : undefined);
</script>

<div class="space-y-2">
	{#if label}
		<div class="flex items-center justify-between">
			<label for={id} class="block text-sm font-medium text-text-secondary">
				{label}
			</label>
			{#if helperText}
				<span class="text-xs text-text-muted">{helperText}</span>
			{/if}
		</div>
	{/if}

	<div class="relative flex items-center">
		{#if prefixIcon}
			<div
				class="absolute left-3.5 flex items-center pointer-events-none select-none text-text-muted"
			>
				{@render prefixIcon()}
			</div>
		{/if}

		<input
			{id}
			{disabled}
			{required}
			bind:value
			bind:this={inputRef}
			aria-required={required ? "true" : undefined}
			aria-invalid={hasError}
			aria-describedby={hasError ? errorId : undefined}
			class="w-full py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-text-muted text-sm
			focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed
      {prefixIcon ? 'pl-8' : 'px-4'} 
      {suffixIcon ? 'pr-10' : 'px-4'} 
      {hasError
				? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
				: 'border-app-border focus:border-text-brand focus:ring-text-brand/20'} 
      {className}"
			{...restProps}
		/>

		{#if suffixIcon}
			<div class="absolute right-3 flex items-center">
				{@render suffixIcon()}
			</div>
		{/if}
	</div>

	{#if hasError}
		<p id={errorId} role="alert" class="text-xs text-red-400">
			{error}
		</p>
	{/if}
</div>
