<script lang="ts">
	import { SignUpController } from './sign_up.svelte';

	const controller = new SignUpController();
</script>

<div class="flex items-center justify-center min-h-[calc(100vh-10rem)] px-4 scroll-smooth">
	<div class="w-full max-w-md p-8 bg-app-surface border border-gray-800 rounded-xl shadow-xl">
		<!-- Header do Card -->
		<div class="mb-8 text-center">
			<h2 class="font-montserrat font-bold text-2xl text-text-primary tracking-tight mb-2">
				Bem-vindo!
			</h2>
			<p class="text-sm text-gray-400 font-inter">
				Crie uma conta para salvar seus algoritmos favoritos.
			</p>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => controller.signUp(e)} class="space-y-5 font-inter">
			<!-- Name -->
			<div class="space-y-2">
				<label for="name" class="block text-sm font-medium text-gray-300">Nome</label>
				<input
					type="text"
					id="name"
					name="name"
					autocomplete="name"
					bind:value={controller.name}
					bind:this={controller.nameInput}
					oninput={() => controller.onNameInput()}
					onblur={() => controller.onNameBlur()}
					aria-required="true"
					aria-invalid={controller.touched.name && !controller.isNameValid}
					aria-describedby={controller.touched.name &&
					controller.name.length > 0 &&
					!controller.isNameValid
						? 'name-error'
						: undefined}
					placeholder="Pedro da Silva"
					required
					disabled={controller.isLoading}
					class="w-full px-4 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
        {controller.apiError?.code === 'REGISTRATION_INVALID_NAME' ||
					(controller.touched.name && !controller.isNameValid)
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
				/>
				{#if controller.touched.name && controller.name.length > 0 && !controller.isNameValid}
					<p id="name-error" role="alert" class="text-xs text-amber-500">
						O nome precisa ter no mínimo 6 letras válidas.
					</p>
				{/if}
			</div>

			<!-- Username -->
			<div class="space-y-2">
				<label for="username" class="block text-sm font-medium text-gray-300">Nome de usuário</label
				>
				<input
					bind:value={controller.username}
					bind:this={controller.usernameInput}
					type="text"
					id="username"
					name="username"
					autocomplete="username"
					aria-required="true"
					aria-invalid={(controller.touched.username &&
						controller.username.length > 0 &&
						!controller.isUsernameValid) ||
						controller.apiError?.code === 'USERNAME_ALREADY_USED' ||
						controller.apiError?.code === 'REGISTRATION_INVALID_USERNAME'}
					aria-describedby={(controller.touched.username &&
						controller.username.length > 0 &&
						!controller.isUsernameValid) ||
					controller.apiError?.code === 'USERNAME_ALREADY_USED' ||
					controller.apiError?.code === 'REGISTRATION_INVALID_USERNAME'
						? 'username-error'
						: undefined}
					oninput={() => controller.onUsernameInput()}
					onblur={() => controller.onUsernameBlur()}
					placeholder="pedro_silva"
					required
					disabled={controller.isLoading}
					class="w-full px-4 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
        {controller.apiError?.code === 'REGISTRATION_INVALID_USERNAME' ||
					(controller.touched.username && !controller.isUsernameValid) ||
					controller.apiError?.code === 'USERNAME_ALREADY_USED'
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
				/>
				{#if (controller.touched.username && controller.username.length > 0 && !controller.isUsernameValid) || controller.apiError?.code === 'USERNAME_ALREADY_USED'}
					<p
						id="username-error"
						role="alert"
						class="text-xs {controller.apiError?.code === 'USERNAME_ALREADY_USED'
							? 'text-red-400'
							: 'text-amber-500'}"
					>
						{controller.apiError?.code === 'USERNAME_ALREADY_USED'
							? 'Nome de usuário já usado, tente outro nome de usuário.'
							: 'Apenas letras, números, _ e - (mínimo 6 caracteres).'}
					</p>
				{/if}
			</div>

			<!-- Email -->
			<div class="space-y-2">
				<label for="email" class="block text-sm font-medium text-gray-300">E-mail</label>
				<input
					type="email"
					id="email"
					name="email"
					autocomplete="email"
					bind:value={controller.email}
					bind:this={controller.emailInput}
					oninput={() => controller.onEmailInput()}
					onblur={() => controller.onEmailBlur()}
					aria-required="true"
					aria-invalid={(controller.touched.email && !controller.isEmailValid) ||
						controller.apiError?.code === 'EMAIL_ALREADY_USED' ||
						controller.apiError?.code === 'REGISTRATION_INVALID_EMAIL'}
					aria-describedby={(controller.touched.email && !controller.isEmailValid) ||
					controller.apiError?.code === 'EMAIL_ALREADY_USED'
						? 'email-error'
						: undefined}
					placeholder="seu@email.com"
					required
					disabled={controller.isLoading}
					class="w-full px-4 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
        {controller.apiError?.code === 'REGISTRATION_INVALID_EMAIL' ||
					(controller.touched.email && !controller.isEmailValid) ||
					controller.apiError?.code === 'EMAIL_ALREADY_USED'
						? 'border-red-500 focus:border-red-500 focus:ring-red-500'
						: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
				/>
				{#if (controller.touched.email && !controller.isEmailValid) || controller.apiError?.code === 'EMAIL_ALREADY_USED'}
					<p id="email-error" role="alert" class="text-xs text-red-400">
						{controller.apiError?.code === 'EMAIL_ALREADY_USED'
							? 'Email já cadastrado em uma conta, tente fazer login.'
							: 'Digite um endereço de e-mail válido.'}
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
						autocomplete="new-password"
						bind:value={controller.password}
						bind:this={controller.passwordInput}
						aria-required="true"
						aria-invalid={controller.apiError?.code === 'USER_PASSWORD_NOT_VALID'}
						oninput={() => controller.onPasswordInput()}
						onblur={() => controller.onPasswordBlur()}
						placeholder="••••••••"
						required
						disabled={controller.isLoading}
						class="w-full px-4 pr-10 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
                {controller.apiError?.code === 'USER_PASSWORD_NOT_VALID'
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

				{#if controller.password.length > 0}
					<div
						aria-live="polite"
						class="p-3 bg-app-bg/30 border border-gray-800/80 rounded-lg space-y-1.5 text-xs mt-2 transition-all"
					>
						<p class="font-medium text-gray-400 mb-1">A senha precisa conter:</p>

						<div
							class="flex items-center gap-2 transition-colors {controller.hasMinLength
								? 'text-emerald-400'
								: 'text-gray-500'}"
						>
							<svg
								class="w-3.5 h-3.5 shrink-0"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								{#if controller.hasMinLength}
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="3"
										d="M5 13l4 4L19 7"
									/>
								{:else}
									<circle cx="12" cy="12" r="3" fill="currentColor" />
								{/if}
							</svg>
							<span>Pelo menos 8 caracteres</span>
						</div>

						<div
							class="flex items-center gap-2 transition-colors {controller.hasUppercase &&
							controller.hasLowercase
								? 'text-emerald-400'
								: 'text-gray-500'}"
						>
							<svg
								class="w-3.5 h-3.5 shrink-0"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								{#if controller.hasUppercase && controller.hasLowercase}
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="3"
										d="M5 13l4 4L19 7"
									/>
								{:else}
									<circle cx="12" cy="12" r="3" fill="currentColor" />
								{/if}
							</svg>
							<span>Letras maiúsculas e minúsculas</span>
						</div>

						<div
							class="flex items-center gap-2 transition-colors {controller.hasNumber
								? 'text-emerald-400'
								: 'text-gray-500'}"
						>
							<svg
								class="w-3.5 h-3.5 shrink-0"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								{#if controller.hasNumber}
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="3"
										d="M5 13l4 4L19 7"
									/>
								{:else}
									<circle cx="12" cy="12" r="3" fill="currentColor" />
								{/if}
							</svg>
							<span>Pelo menos um número (0-9)</span>
						</div>

						<div
							class="flex items-center gap-2 transition-colors {controller.hasSpecialChar
								? 'text-emerald-400'
								: 'text-gray-500'}"
						>
							<svg
								class="w-3.5 h-3.5 shrink-0"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								{#if controller.hasSpecialChar}
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="3"
										d="M5 13l4 4L19 7"
									/>
								{:else}
									<circle cx="12" cy="12" r="3" fill="currentColor" />
								{/if}
							</svg>
							<span>Pelo menos um símbolo (@$!%*?&)</span>
						</div>
					</div>
				{/if}
			</div>

			<!-- Confirm Password -->
			<div class="space-y-2">
				<label for="confirmPassword" class="block text-sm font-medium text-gray-300">
					Confirmar senha
				</label>
				<div class="relative flex items-center">
					<input
						type={controller.showConfirmPassword ? 'text' : 'password'}
						id="confirmPassword"
						name="confirmPassword"
						autocomplete="new-password"
						bind:value={controller.confirmPassword}
						bind:this={controller.confirmPasswordInput}
						aria-required="true"
						aria-invalid={(controller.touched.confirmPassword && !controller.isPasswordsMatching) ||
							controller.apiError?.code === 'USER_PASSWORDS_DONT_MATCH'}
						aria-describedby={(controller.touched.confirmPassword &&
							controller.confirmPassword.length > 0 &&
							!controller.isPasswordsMatching) ||
						controller.apiError?.code === 'USER_PASSWORDS_DONT_MATCH'
							? 'confirm-password-error'
							: undefined}
						oninput={() => controller.onPasswordInput()}
						onblur={() => controller.onConfirmPasswordBlur()}
						placeholder="••••••••"
						required
						disabled={controller.isLoading}
						class="w-full px-4 pr-10 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:ring-1 focus:outline-none transition-all disabled:opacity-50
                {controller.apiError?.code === 'USER_PASSWORDS_DONT_MATCH' ||
						(controller.touched.confirmPassword && !controller.isPasswordsMatching)
							? 'border-red-500 focus:border-red-500 focus:ring-red-500'
							: 'border-gray-800 focus:border-text-brand focus:ring-text-brand'}"
					/>
					<button
						type="button"
						onclick={() => controller.toggleConfirmPassword()}
						class="absolute right-3 p-1 rounded text-zinc-400 hover:text-white transition-colors focus:outline-none focus:ring-1 focus:ring-text-brand"
						aria-label={controller.showConfirmPassword
							? 'Ocultar confirmação de senha'
							: 'Mostrar confirmação de senha'}
					>
						{#if controller.showConfirmPassword}
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
				{#if controller.touched.confirmPassword && controller.confirmPassword.length > 0 && !controller.isPasswordsMatching}
					<p id="confirm-password-error" role="alert" class="text-xs text-red-400">
						As senhas não coincidem.
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
						<!-- Api errors -->
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
					<!-- Loading spinner -->
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
					<span>Criando...</span>
				{:else}
					<span>Criar</span>
				{/if}
			</button>

			<!-- Login -->
			<p class="text-center text-sm text-gray-400 pt-2">
				Já tem uma conta?
				<a href="/login" class="text-text-brand font-medium hover:underline transition-all ml-1">
					Faça login
				</a>
			</p>
		</form>
	</div>
</div>
