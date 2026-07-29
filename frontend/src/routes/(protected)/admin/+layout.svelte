<script lang="ts">
	import type { LayoutData } from "./$types";
	import { AdminController } from "./admin_controller.svelte";

	let { data, children }: { data: LayoutData; children: any } = $props();

	const controller = new AdminController();
</script>

{#if !data.hasAdminSecret}
	<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4">
		<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
			<!-- Header do Card -->
			<div class="mb-8 text-center">
				<h2 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
					Chave de Segurança Administrador
				</h2>
				<p class="text-sm text-gray-400 font-inter">Digite a senha das rotas administradoras.</p>
			</div>

			<!-- Form -->
			<form onsubmit={(e) => controller.sendPassword(e)} class="space-y-5 font-inter">
				<!-- Password -->
				<div class="space-y-2">
					<label for="password" class="block text-sm font-medium text-gray-300">Senha</label>
					<div class="relative flex items-center">
						<input
							type={controller.showPassword ? "text" : "password"}
							id="password"
							name="password"
							autocomplete="current-password"
							minlength="8"
							bind:value={controller.password}
							oninput={() => controller.onInput()}
							onblur={() => (controller.touched.password = true)}
							aria-required="true"
							aria-invalid={controller.touched.password && !controller.isPasswordValid}
							aria-describedby={controller.touched.password && !controller.isPasswordValid
								? "password-error"
								: undefined}
							placeholder="••••••••"
							required
							disabled={controller.isLoading}
							class="w-full px-4 pr-10 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
                {(controller.touched.password && !controller.isPasswordValid) ||
							controller.apiError?.code === 'AUTH_INVALID_EMAIL_PASSWORD'
								? 'border-red-500 focus:border-red-500 focus:ring-red-500'
								: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
						/>
						<button
							type="button"
							onclick={() => controller.togglePassword()}
							class="absolute right-3 p-1 rounded text-zinc-400 hover:text-white transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand"
							aria-label={controller.showPassword ? "Ocultar senha" : "Mostrar senha"}
						>
							{#if controller.showPassword}
								<svg
									xmlns="http://www.w3.org/2000/svg"
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
									xmlns="http://www.w3.org/2000/svg"
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
					</div>
					{#if controller.touched.password && !controller.isPasswordValid}
						<p id="password-error" role="alert" class="text-xs text-red-400">
							A senha deve conter no mínimo 8 caracteres.
						</p>
					{/if}
				</div>

				<!-- Dynamic API Error Box -->
				{#if controller.apiError}
					<div
						role="alert"
						aria-live="assertive"
						class="p-3 bg-red-950/30 border border-red-900/50 rounded-lg text-red-400 text-sm flex items-start gap-2
            {controller.isLoading ? 'opacity-40 pointer-events-none' : 'opacity-100'}"
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

				<!-- Submit button -->
				<button
					type="submit"
					disabled={controller.isLoading}
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
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						<span>Entrando...</span>
					{:else}
						<span>Entrar</span>
					{/if}
				</button>
			</form>
		</div>
	</div>
{:else}
	{@render children()}
{/if}
