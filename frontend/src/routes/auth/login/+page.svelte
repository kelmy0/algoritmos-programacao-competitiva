<script lang="ts">
	import { LoginController } from './login.svelte';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { redirect } from '@sveltejs/kit';

	if (page.data.user) {
		redirect(301, '/');
	}

	const controller = new LoginController();

	onMount(() => controller.checkErrors());
</script>

<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4">
	<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
		<!-- Header do Card -->
		<div class="mb-8 text-center">
			<h2 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
				Bem-vindo de volta!
			</h2>
			<p class="text-sm text-gray-400 font-inter">
				Faça login para gerenciar seus algoritmos salvos.
			</p>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => controller.login(e)} class="space-y-5 font-inter">
			<!-- Email -->
			<div class="space-y-2">
				<label for="email" class="block text-sm font-medium text-gray-300">E-mail</label>
				<input
					type="email"
					id="email"
					name="email"
					autocomplete="email"
					bind:value={controller.email}
					oninput={() => controller.onInput()}
					onblur={() => (controller.touched.email = true)}
					aria-required="true"
					aria-invalid={controller.touched.email && !controller.isEmailValid}
					aria-describedby={controller.touched.email && !controller.isEmailValid
						? 'email-error'
						: undefined}
					placeholder="seu@email.com"
					required
					disabled={controller.isLoading}
					class="w-full px-4 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
            {(controller.touched.email && !controller.isEmailValid) ||
					controller.apiError?.code === 'AUTH_INVALID_EMAIL_PASSWORD'
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
				/>
				{#if controller.touched.email && !controller.isEmailValid}
					<p id="email-error" role="alert" class="text-xs text-red-400">
						Digite um endereço de e-mail válido.
					</p>
				{/if}
			</div>

			<!-- Password -->
			<div class="space-y-2">
				<label for="password" class="block text-sm font-medium text-gray-300">Senha</label>
				<div class="relative flex items-center">
					<input
						type={controller.showPassword ? 'text' : 'password'}
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
							? 'password-error'
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
						aria-label={controller.showPassword ? 'Ocultar senha' : 'Mostrar senha'}
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
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
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

			<!-- Sign up -->
			<p class="text-center text-sm text-gray-400 pt-2">
				Não tem uma conta?
				<a
					href="/auth/sign-up"
					class="text-text-brand font-medium hover:underline transition-all ml-1"
				>
					Crie uma conta
				</a>
			</p>
		</form>

		<!-- Divider -->
		<div class="relative my-6 flex items-center justify-center">
			<div class="absolute inset-0 flex items-center">
				<div class="w-full border-t border-gray-800"></div>
			</div>
			<div class="relative bg-app-bg px-3 text-xs text-gray-500 uppercase tracking-wider">
				ou continue com
			</div>
		</div>

		<!-- Social Login Buttons -->
		<div class="grid grid-cols-2 gap-3">
			<!-- Google -->
			<a
				href="{PUBLIC_API_URL}/api/auth/google"
				class="flex items-center justify-center gap-2 py-2.5 px-3 bg-app-bg/50 border border-gray-800 rounded-lg text-xs font-medium text-gray-300 hover:bg-gray-800/40 hover:border-gray-700 hover:text-white transition-all focus:outline-none focus:ring-1 focus:ring-text-brand"
				aria-label="Entrar com o Google"
			>
				<svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24">
					<path
						fill="#EA4335"
						d="M12 5c1.6 0 3 .6 4.1 1.6l3.1-3.1C17.3 1.7 14.8 1 12 1 7.5 1 3.7 3.6 1.9 7.3l3.7 2.9C6.5 7.2 9 5 12 5z"
					/>
					<path
						fill="#4285F4"
						d="M23.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.5h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.8z"
					/>
					<path
						fill="#FBBC05"
						d="M5.6 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.6.4-2.3L1.9 7.3C.7 9.7 0 10.8 0 12s.7 2.3 1.9 4.7l3.7-2.9z"
					/>
					<path
						fill="#34A853"
						d="M12 23c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3 0-5.5-2.2-6.4-5.2L1.9 16C3.7 19.7 7.5 23 12 23z"
					/>
				</svg>
				<span>Google</span>
			</a>

			<!-- GitHub -->
			<a
				href="{PUBLIC_API_URL}/api/auth/github"
				class="flex items-center justify-center gap-2 py-2.5 px-3 bg-app-bg/50 border border-gray-800 rounded-lg text-xs font-medium text-gray-300 hover:bg-gray-800/40 hover:border-gray-700 hover:text-white transition-all focus:outline-none focus:ring-1 focus:ring-text-brand"
				aria-label="Entrar com o GitHub"
			>
				<svg class="w-4 h-4 shrink-0 fill-current" viewBox="0 0 24 24">
					<path
						d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
					/>
				</svg>
				<span>GitHub</span>
			</a>
		</div>
	</div>
</div>
