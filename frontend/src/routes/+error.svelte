<script lang="ts">
	import { page } from '$app/state';
</script>

<div
	class="flex min-h-[60vh] flex-col items-center justify-center p-4 text-center font-inter"
	role="alert"
	aria-live="assertive"
>
	<div
		class="bg-app-surface border-gray-800 space-y-5 rounded-xl border p-8 shadow-2xl max-w-md w-full"
	>
		<div class="flex items-center justify-center gap-2">
			<span
				class="text-text-brand bg-text-brand/10 border-text-brand/20 rounded-md border px-3 py-1 text-xs font-mono font-semibold"
				aria-label="Código de status HTTP: {page.status}"
			>
				Erro {page.status}
			</span>
		</div>

		<h1 class="font-montserrat text-xl font-bold text-white tracking-tight">
			{page.error?.message || 'Ocorreu um erro inesperado.'}
		</h1>

		{#if page.error && 'code' in page.error}
			<p class="text-xs font-mono text-gray-400">
				<span class="sr-only">Código interno do erro: </span>
				Código: {page.error.code}
			</p>
		{/if}

		<div class="pt-2">
			<a
				href="/"
				class="inline-flex items-center justify-center rounded-lg bg-gray-800 px-5 py-2.5 text-sm font-medium text-white hover:bg-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-brand focus:ring-offset-2 focus:ring-offset-gray-900"
			>
				Voltar para o início
			</a>
		</div>
	</div>
</div>
