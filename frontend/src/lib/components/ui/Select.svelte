<script lang="ts">
	import type { HTMLSelectAttributes } from "svelte/elements";
	import type { Snippet } from "svelte";

	export interface SelectOption {
		value: string | number;
		label: string;
		disabled?: boolean;
	}

	interface Props extends HTMLSelectAttributes {
		label?: string;
		error?: string;
		touched?: boolean;
		helperText?: string;
		value?: string | number;
		options?: SelectOption[];
		selectRef?: HTMLSelectElement | null;
		children?: Snippet;
	}

	let {
		id,
		label,
		error,
		touched = false,
		helperText,
		value = $bindable(""),
		options = [],
		selectRef = $bindable(null),
		class: className = "",
		required,
		disabled,
		children,
		...restProps
	}: Props = $props();

	const hasError = $derived(Boolean(touched && error));
	const errorId = $derived(id ? `${id}-error` : undefined);
</script>

<div class="space-y-2">
	{#if label}
		<div class="flex items-center justify-between">
			<label for={id} class="block text-sm font-medium text-gray-300">
				{label}
				{#if required}
					<span class="text-red-400" aria-hidden="true">*</span>
				{/if}
			</label>
			{#if helperText}
				<span class="text-xs text-gray-500">{helperText}</span>
			{/if}
		</div>
	{/if}

	<div class="relative flex items-center">
		<select
			{id}
			{disabled}
			{required}
			bind:value
			bind:this={selectRef}
			aria-required={required ? "true" : undefined}
			aria-invalid={hasError}
			aria-describedby={hasError ? errorId : undefined}
			class="w-full py-2.5 px-3 bg-app-bg/50 border rounded-lg text-text-primary text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all hover:cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed
            {hasError
				? 'border-red-500 focus:border-red-500 focus:ring-red-500'
				: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'} 
            {className}"
			{...restProps}
		>
			{#if children}
				{@render children()}
			{:else}
				{#each options as opt}
					<option value={opt.value} disabled={opt.disabled}>
						{opt.label}
					</option>
				{/each}
			{/if}
		</select>
	</div>

	{#if hasError}
		<p id={errorId} role="alert" class="text-xs text-amber-500">
			{error}
		</p>
	{/if}
</div>
