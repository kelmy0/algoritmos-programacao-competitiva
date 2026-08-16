<script lang="ts">
	import { onMount } from "svelte";
	import type { LayoutData } from "./$types";
	import { AdminController } from "./admin_controller.svelte";
	import { createActivityKeeper } from "$lib/utils/idle-ping";
	import { invalidateAll } from "$app/navigation";
	import { focusTrap } from "$lib/utils/a11y";
	import Input from "$lib/components/ui/Input.svelte";
	import Alert from "$lib/components/ui/Alert.svelte";
	import Button from "$lib/components/ui/Button.svelte";

	let { data, children }: { data: LayoutData; children: any } = $props();

	const controller = new AdminController();

	let hasBeenAuthenticated = $state(false);

	$effect(() => {
		if (data.hasAdminSecret) {
			hasBeenAuthenticated = true;
		}
	});

	let showAuthOverlay = $derived(!data.hasAdminSecret);

	onMount(() => {
		const keeper = createActivityKeeper({
			intervalMinutes: 5,
			onUnauthorized: () => {
				invalidateAll();
			}
		});
		keeper.start();

		return () => keeper.stop();
	});
</script>

{#if !hasBeenAuthenticated && showAuthOverlay}
	<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4">
		<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
			{@render AuthForm("auth-title-login")}
		</div>
	</div>
{:else}
	<div class="relative min-h-screen">
		<div aria-hidden={showAuthOverlay}>
			{@render children()}
		</div>

		{#if showAuthOverlay}
			<div
				use:focusTrap
				role="dialog"
				aria-modal="true"
				aria-labelledby="auth-title-overlay"
				tabindex="-1"
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md animate-in fade-in duration-200"
			>
				<div
					class="w-full max-w-md p-8 bg-app-surface/95 border border-gray-800 rounded-xl shadow-2xl backdrop-blur-xl"
				>
					<div class="mb-4 text-center">
						<span
							class="inline-block px-3 py-1 mb-2 text-xs font-semibold text-amber-400 bg-amber-950/40 border border-amber-800/50 rounded-full"
						>
							Sessão Expirada por Inatividade
						</span>
					</div>
					{@render AuthForm("auth-title-overlay")}
				</div>
			</div>
		{/if}
	</div>
{/if}

{#snippet AuthForm(headingId: string)}
	<div class="mb-8 text-center">
		<h2
			id={headingId}
			class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2"
		>
			Chave de Segurança Administrador
		</h2>
		<p class="text-sm text-gray-400 font-inter">
			Digite sua senha para desbloquear e continuar de onde parou.
		</p>
	</div>

	<form onsubmit={(e) => controller.sendPassword(e)} class="space-y-5 font-inter">
		<!-- Password -->
		<Input
			id="password"
			name="password"
			type={controller.showPassword ? "text" : "password"}
			label="Senha"
			placeholder="••••••••"
			autocomplete="current-password"
			minlength={8}
			required
			disabled={controller.isLoading}
			bind:value={controller.password}
			touched={controller.touched.password}
			error={!controller.isPasswordValid
				? "A senha deve conter no mínimo 8 caracteres."
				: undefined}
			oninput={() => controller.onInput()}
			onblur={() => (controller.touched.password = true)}
		>
			{#snippet suffixIcon()}
				<button
					type="button"
					onclick={() => controller.togglePassword()}
					class="p-1 rounded text-zinc-400 hover:text-white transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand"
					aria-label={controller.showPassword ? "Ocultar senha" : "Mostrar senha"}
				>
					{#if controller.showPassword}
						<svg
							class="h-5 w-5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
						>
							<path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
							<circle cx="12" cy="12" r="3" />
						</svg>
					{:else}
						<svg
							class="h-5 w-5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
						>
							<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
							<path
								d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
							/>
							<path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
							<line x1="2" x2="22" y1="2" y2="22" />
						</svg>
					{/if}
				</button>
			{/snippet}
		</Input>

		<!-- Dynamic API Error Box -->
		{#if controller.apiError}
			<Alert
				title="Erro de Autenticação"
				message={controller.apiError.message}
				isLoading={controller.isLoading}
			/>
		{/if}

		<!-- Submit button -->
		<Button
			type="submit"
			class="w-full"
			isLoading={controller.isLoading}
			disabled={controller.isLoading}
		>
			{controller.isLoading ? "Entrando..." : "Entrar"}
		</Button>
	</form>
{/snippet}
