<script lang="ts">
	import Turnstile from "$lib/components/turnstile.svelte";
	import Alert from "$lib/components/ui/Alert.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import CodeInput from "$lib/components/ui/CodeInput.svelte";
	import { TwoFactorController } from "./two_factor_verify.svelte";

	const controller = new TwoFactorController();
</script>

<svelte:head>
	<title>Verificação 2FA</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4">
	<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
		<!-- Header -->
		<div class="mb-8 text-center">
			<h1 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
				Autenticação em dois fatores
			</h1>
			<p id="code-hint" class="text-sm text-gray-400 font-inter">
				Digite o código de 6 dígitos gerado pelo seu aplicativo autenticador.
			</p>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => controller.sendCode(e)} class="space-y-5 font-inter">
			<CodeInput
				bind:value={controller.code}
				touched={controller.touched.code}
				error={!controller.isCodeValid || controller.apiError?.code === "2FA_INVALID_CODE"
					? "O código deve conter exatamente 6 números."
					: undefined}
				disabled={controller.isLoading}
				oninput={(e) => controller.onInput(e)}
				onblur={() => (controller.touched.code = true)}
			/>

			<div class="flex justify-center">
				<Turnstile
					bind:this={controller.turnstileComponent}
					onsuccess={(token) => controller.onTurnstileSuccess(token)}
					onexpire={() => controller.onTurnstileExpire()}
				/>
			</div>

			<!-- Dynamic API Error Box -->
			{#if controller.apiError}
				<Alert
					title="Erro de Autenticação"
					message={controller.apiError.message}
					isLoading={controller.isLoading}
				/>
			{/if}

			<!-- Submit -->
			<Button
				type="submit"
				class="w-full"
				isLoading={controller.isLoading}
				disabled={controller.isLoading}
			>
				{controller.isLoading ? "Verificando..." : "Verificar"}
			</Button>
		</form>
	</div>
</div>
