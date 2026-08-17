<script lang="ts">
	import { slide } from "svelte/transition";

	interface Props {
		id?: string;
		name?: string;
		label?: string;
		length?: number;
		value: string;
		error?: string;
		disabled?: boolean;
		touched?: boolean;
		oninput?: (e: Event) => void;
		onblur?: (e: FocusEvent) => void;
	}

	let {
		id = "code",
		name = "code",
		label = "Código de verificação",
		length = 6,
		value = $bindable(""),
		error,
		disabled = false,
		touched = false,
		oninput,
		onblur
	}: Props = $props();

	const hasError = $derived(Boolean(touched && error));
</script>

<div class="space-y-2">
	{#if label}
		<label for={id} class="block text-sm font-medium text-text-secondary">
			{label}
		</label>
	{/if}

	<input
		type="text"
		inputmode="numeric"
		pattern="[0-9]*"
		{id}
		{name}
		minlength={length}
		maxlength={length}
		autocomplete="one-time-code"
		bind:value
		{oninput}
		{onblur}
		aria-required="true"
		aria-invalid={hasError}
		aria-describedby={hasError ? `${id}-error` : undefined}
		placeholder={"0".repeat(length)}
		required
		{disabled}
		class="w-full px-4 py-3 bg-app-bg/50 border rounded-lg text-text-primary placeholder:text-text-muted
               text-center font-mono text-xl tracking-[0.5em] focus:bg-app-bg focus:ring-2 focus:outline-none
               transition-all disabled:opacity-50
        {hasError
			? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
			: 'border-app-border focus:border-text-brand focus:ring-text-brand/20'}"
	/>

	{#if hasError}
		<p
			transition:slide={{ duration: 200 }}
			id={`${id}-error`}
			role="alert"
			class="text-xs text-red-400"
		>
			{error}
		</p>
	{/if}
</div>
