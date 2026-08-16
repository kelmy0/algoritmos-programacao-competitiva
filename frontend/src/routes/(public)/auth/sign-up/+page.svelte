<script lang="ts">
	import { SignUpController } from "./sign_up.svelte";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import Turnstile from "$lib/components/turnstile.svelte";
	import Input from "$lib/components/ui/Input.svelte";
	import ValidationCard from "$lib/components/ui/ValidationCard.svelte";
	import Alert from "$lib/components/ui/Alert.svelte";
	import Button from "$lib/components/ui/Button.svelte";

	const controller = new SignUpController();

	onMount(() => {
		const error = page.url.searchParams.get("error");

		if (error == "social_auth_failed") {
		}
	});
</script>

<svelte:head>
	<title>Criar conta</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta
		name="description"
		content="Criar conta. Crie uma conta para salvar seus algoritmos favoritos!"
	/>
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4 scroll-smooth">
	<div class="w-full max-w-md p-8 bg-app-surface border border-app-border rounded-xl shadow-xl">
		<!-- Header do Card -->
		<div class="mb-8 text-center">
			<h2 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
				Bem-vindo!
			</h2>
			<p class="text-sm text-text-muted font-inter">
				Crie uma conta para salvar seus algoritmos favoritos.
			</p>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => controller.signUp(e)} class="space-y-5 font-inter">
			<!-- Name -->
			<Input
				id="name"
				name="name"
				type="text"
				label="Nome completo"
				placeholder="Pedro da Silva"
				autocomplete="name"
				required
				bind:value={controller.name}
				bind:inputRef={controller.nameInput}
				touched={controller.touched.name}
				error={!controller.isNameValid ? "O nome precisa ter no mínimo 6 letras." : undefined}
				disabled={controller.isLoading}
				oninput={() => controller.onNameInput()}
				onblur={() => controller.onNameBlur()}
			/>

			<!-- Username -->
			<Input
				id="username"
				name="username"
				type="text"
				label="Nome de usuário"
				placeholder="usuario123"
				autocomplete="username"
				required
				bind:value={controller.username}
				bind:inputRef={controller.usernameInput}
				touched={controller.touched.username}
				error={!controller.isUsernameValid
					? "Usuário inválido (mínimo 3 caracteres, apenas letras, números e _)."
					: undefined}
				disabled={controller.isLoading}
				oninput={() => controller.onUsernameInput()}
				onblur={() => controller.onUsernameBlur()}
			/>

			<!-- Email -->
			<Input
				id="email"
				name="email"
				type="email"
				label="E-mail"
				placeholder="seu@email.com"
				autocomplete="email"
				required
				bind:value={controller.email}
				bind:inputRef={controller.emailInput}
				touched={controller.touched.email}
				error={!controller.isEmailValid ? "Digite um e-mail válido." : undefined}
				disabled={controller.isLoading}
				oninput={() => controller.onEmailInput()}
				onblur={() => controller.onEmailBlur()}
			/>

			<!-- Password -->
			<Input
				id="password"
				name="password"
				type={controller.showPassword ? "text" : "password"}
				label="Senha"
				placeholder="••••••••"
				autocomplete="new-password"
				required
				bind:value={controller.password}
				bind:inputRef={controller.passwordInput}
				touched={controller.touched.password}
				error={!controller.isPasswordValid
					? "A senha não atende aos requisitos mínimos."
					: undefined}
				disabled={controller.isLoading}
				oninput={() => controller.onPasswordInput()}
				onblur={() => controller.onPasswordBlur()}
			>
				{#snippet suffixIcon()}
					<button
						type="button"
						onclick={() => controller.togglePassword()}
						class="p-1 rounded text-text-muted hover:text-text-primary transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand"
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

			{#if controller.password.length > 0}
				<ValidationCard
					title="A senha precisa conter:"
					requirements={[
						{ label: "Pelo menos 8 caracteres", met: controller.hasMinLength },
						{
							label: "Letras maiúsculas e minúsculas",
							met: controller.hasUppercase && controller.hasLowercase
						},
						{ label: "Pelo menos um número (0-9)", met: controller.hasNumber },
						{ label: "Pelo menos um símbolo (@$!%*?&)", met: controller.hasSpecialChar }
					]}
				/>
			{/if}

			<!-- Confirm Password -->
			<Input
				id="confirmPassword"
				name="confirmPassword"
				type={controller.showConfirmPassword ? "text" : "password"}
				label="Confirmar senha"
				placeholder="••••••••"
				autocomplete="new-password"
				required
				bind:value={controller.confirmPassword}
				bind:inputRef={controller.confirmPasswordInput}
				touched={controller.touched.confirmPassword}
				error={controller.apiError?.code === "USER_PASSWORDS_DONT_MATCH"
					? "As senhas não coincidem."
					: controller.touched.confirmPassword &&
						  controller.confirmPassword.length > 0 &&
						  !controller.isPasswordsMatching
						? "As senhas não coincidem."
						: undefined}
				disabled={controller.isLoading}
				oninput={() => controller.onPasswordInput()}
				onblur={() => controller.onConfirmPasswordBlur()}
			>
				{#snippet suffixIcon()}
					<button
						type="button"
						onclick={() => controller.toggleConfirmPassword()}
						class="p-1 rounded text-text-muted hover:text-text-primary transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand"
						aria-label={controller.showConfirmPassword
							? "Ocultar confirmação de senha"
							: "Mostrar confirmação de senha"}
					>
						{#if controller.showConfirmPassword}
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

			<!-- Submit button -->
			<Button
				type="submit"
				class="w-full"
				isLoading={controller.isLoading}
				disabled={controller.isLoading}
			>
				{controller.isLoading ? "Criando..." : "Criar"}
			</Button>

			<!-- Login -->
			<p class="text-center text-sm text-text-muted pt-2">
				Já tem uma conta?
				<a
					href="/auth/login"
					class="text-text-brand font-medium hover:underline transition-all ml-1"
				>
					Faça login
				</a>
			</p>
		</form>
		<!-- Divider -->
		<div class="relative my-6 flex items-center justify-center">
			<div class="absolute inset-0 flex items-center">
				<div class="w-full border-t border-app-border"></div>
			</div>
			<div class="relative bg-app-bg px-3 text-xs text-text-muted uppercase tracking-wider">
				ou continue com
			</div>
		</div>

		<!-- Social Login Buttons -->
		<div class="grid grid-cols-2 gap-3">
			<!-- Google -->
			<a
				href="/api/auth/social/google"
				class="flex items-center justify-center gap-2 py-2.5 px-3 bg-app-bg/50 border border-app-border rounded-lg text-xs font-medium text-text-secondary hover:bg-app-bg/40 hover:border-app-border-hover hover:text-text-primary transition-all focus:outline-none focus:ring-1 focus:ring-text-brand"
				aria-label="Entrar com o Google"
			>
				<svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" aria-hidden="true">
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
				href="/api/auth/social/github"
				class="flex items-center justify-center gap-2 py-2.5 px-3 bg-app-bg/50 border border-app-border rounded-lg text-xs font-medium text-text-secondary hover:bg-app-bg/40 hover:border-app-border-hover hover:text-text-primary transition-all focus:outline-none focus:ring-1 focus:ring-text-brand"
				aria-label="Entrar com o GitHub"
			>
				<svg class="w-4 h-4 shrink-0 fill-current" viewBox="0 0 24 24" aria-hidden="true">
					<path
						d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
					/>
				</svg>
				<span>GitHub</span>
			</a>
		</div>
	</div>
</div>
