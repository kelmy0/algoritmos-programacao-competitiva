<script lang="ts">
	import { TwoFactorController } from './two_factor_verify.svelte';
	import { onMount } from 'svelte';

	const controller = new TwoFactorController();

	onMount(() => controller.getToken());
</script>

<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4">
	<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
		<!-- Header -->
		<div class="mb-8 text-center">
			<h1 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
				Autenticação em dois fatores
			</h1>
			<p class="text-sm text-gray-400 font-inter">
				Digite o código de 6 dígitos gerado pelo seu aplicativo autenticador.
			</p>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => controller.sendCode(e)} class="space-y-5 font-inter">
			<div class="space-y-2">
				<label for="code" class="block text-sm font-medium text-gray-300">
					Código de verificação
				</label>

				<input
					type="text"
					inputmode="numeric"
					pattern="[0-9]*"
					id="code"
					name="code"
					minlength="6"
					maxlength="6"
					autocomplete="one-time-code"
					bind:value={controller.code}
					oninput={(e) => controller.onInput(e)}
					onblur={() => (controller.touched.code = true)}
					aria-required="true"
					aria-invalid={controller.touched.code && !controller.isCodeValid}
					aria-describedby={controller.touched.code && !controller.isCodeValid
						? 'code-error'
						: undefined}
					placeholder="000000"
					required
					disabled={controller.isLoading}
					class="w-full px-4 py-3 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600
                           text-center font-mono text-xl tracking-[0.5em] focus:bg-app-bg focus:ring-2 focus:outline-none
                           transition-all disabled:opacity-50
                    {(controller.touched.code && !controller.isCodeValid) ||
					controller.apiError?.code === '2FA_INVALID_CODE'
						? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
						: 'border-gray-800 focus:border-text-brand focus:ring-text-brand/20'}"
				/>

				{#if controller.touched.code && !controller.isCodeValid}
					<p id="code-error" role="alert" class="text-xs text-red-400">
						O código deve conter exatamente 6 números.
					</p>
				{/if}
			</div>

			<!-- Box de Erro da API -->
			{#if controller.apiError}
				<div
					role="alert"
					aria-live="assertive"
					class="p-3 bg-red-950/30 border border-red-900/50 rounded-lg text-red-400 text-sm flex items-start gap-2"
				>
					<svg
						class="w-5 h-5 shrink-0 mt-0.5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
						/>
					</svg>
					<div>
						<span class="font-semibold block mb-0.5">Erro de Autenticação</span>
						<p class="text-xs text-red-300/90">{controller.apiError.message}</p>
					</div>
				</div>
			{/if}

			<!-- Botão de Submit -->
			<button
				type="submit"
				disabled={controller.isLoading || !controller.isCodeValid}
				aria-busy={controller.isLoading}
				class="w-full py-2.5 bg-text-brand text-app-bg font-semibold text-sm rounded-lg cursor-pointer
                       hover:opacity-90 active:scale-[0.98] transition-all disabled:opacity-50 disabled:pointer-events-none
                       flex items-center justify-center gap-2"
			>
				{#if controller.isLoading}
					<svg
						class="animate-spin h-4 w-4 text-app-bg"
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
					<span>Verificando...</span>
				{:else}
					<span>Confirmar e Entrar</span>
				{/if}
			</button>
		</form>
	</div>
</div>
