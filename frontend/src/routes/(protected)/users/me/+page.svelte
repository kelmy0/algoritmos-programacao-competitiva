<script lang="ts">
	import { page } from "$app/state";
	import { focusTrap } from "$lib/utils/a11y";
	import { slide } from "svelte/transition";
	import { MeController } from "./me_controller.svelte";

	const controller = new MeController();
	controller.is2FAEnabled = page.data.user?.is2FAEnabled || false;

	const twoFactorLabels: Record<string, string> = {
		generateCode:
			"Você precisará de um aplicativo autenticador. Insira sua senha e clique em gerar chave.",
		saveCode: "Salve a chave em seu aplicativo autenticador e insira o código gerado nele."
	};
</script>

<svelte:head>
	<title>Minha conta</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta name="description" content="Informações e configurações da conta." />
</svelte:head>

<div class="max-w-4xl space-y-8 font-inter">
	<header
		class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-gray-800"
	>
		<div>
			<h1 class="font-montserrat font-bold text-2xl md:text-3xl text-text-primary tracking-tight">
				Minha conta
			</h1>
			<p class="text-sm text-gray-400 mt-1">
				Ajuste suas informações pessoais e configurações de segurança.
			</p>
		</div>
	</header>

	<!-- Account -->
	<section class="bg-app-surface border border-gray-800 rounded-xl p-6 shadow-xl space-y-6">
		<div class="border-b border-gray-800/80 pb-4">
			<h2 class="font-montserrat font-semibold text-lg text-text-primary">Perfil</h2>
			<p class="text-xs text-gray-400 mt-0.5">Informações visíveis do seu perfil.</p>
		</div>

		<form onsubmit={(e) => e.preventDefault()} class="space-y-5">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-5">
				<!-- Name -->
				<div class="space-y-2">
					<label for="name" class="block text-sm font-medium text-gray-300">Nome completo</label>
					<input
						type="text"
						id="name"
						name="name"
						placeholder="Seu nome"
						value={page.data.user?.name}
						class="w-full px-4 py-2.5 bg-app-bg/50 border border-gray-800 rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:border-text-brand focus:ring-1 focus:ring-text-brand focus:outline-none transition-all"
					/>
				</div>

				<!-- Username -->
				<div class="space-y-2">
					<label for="username" class="block text-sm font-medium text-gray-300"
						>Nome de usuário</label
					>
					<div class="relative flex items-center">
						<span
							class="absolute left-3.5 text-gray-500 text-sm font-medium leading-none pointer-events-none select-none"
							>@</span
						>
						<input
							type="text"
							id="username"
							name="username"
							placeholder="dev_user"
							value={page.data.user?.username}
							class="w-full pl-8 pr-4 py-2.5 bg-app-bg/50 border border-gray-800 rounded-lg text-text-primary placeholder-gray-600 text-sm focus:bg-app-bg focus:border-text-brand focus:ring-1 focus:ring-text-brand focus:outline-none transition-all"
						/>
					</div>
				</div>
			</div>

			<!-- Email -->
			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<label for="email" class="block text-sm font-medium text-gray-300">E-mail</label>
					<span class="text-xs text-gray-500">O e-mail não pode ser alterado</span>
				</div>
				<div class="relative flex items-center">
					<input
						type="email"
						id="email"
						name="email"
						value={page.data.user?.email}
						disabled
						class="w-full px-4 py-2.5 bg-app-bg/30 border border-gray-800/60 rounded-lg text-gray-400 text-sm cursor-not-allowed select-none opacity-80"
					/>
					<svg
						class="w-4 h-4 text-gray-500 absolute right-3.5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
						/>
					</svg>
				</div>
			</div>

			<div class="flex justify-end pt-2">
				<button
					type="submit"
					class="px-5 py-2.5 bg-text-brand text-app-bg font-semibold text-sm rounded-lg hover:opacity-90 active:scale-[0.98] transition-all cursor-pointer"
				>
					Salvar perfil
				</button>
			</div>
		</form>
	</section>

	<!-- Social accounts -->
	<section class="bg-app-surface border border-gray-800 rounded-xl p-6 shadow-xl space-y-6">
		<div class="border-b border-gray-800/80 pb-4">
			<h2 class="font-montserrat font-semibold text-lg text-text-primary">Contas vinculadas</h2>
			<p class="text-xs text-gray-400 mt-0.5">Conecte suas redes para facilitar o login.</p>
		</div>

		<div class="space-y-3">
			<!-- Google -->
			<div
				class="flex items-center justify-between p-4 bg-app-bg/40 border border-gray-800 rounded-lg transition-all hover:border-gray-700"
			>
				<div class="flex items-center gap-3">
					<div class="p-2 bg-app-bg border border-gray-800 rounded-md shrink-0">
						<svg class="w-5 h-5" viewBox="0 0 24 24">
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
					</div>
					<div>
						<p class="text-sm font-medium text-text-primary">Google</p>
						<p class="text-xs text-gray-400">Conectado como dev@exemplo.com</p>
					</div>
				</div>

				<!-- Connected state -->
				<button
					type="button"
					class="px-3 py-1.5 border border-red-900/60 text-red-400 hover:bg-red-950/30 hover:border-red-800 rounded-md text-xs font-medium transition-all cursor-pointer"
				>
					Desvincular
				</button>
			</div>

			<!-- GitHub -->
			<div
				class="flex items-center justify-between p-4 bg-app-bg/40 border border-gray-800 rounded-lg transition-all hover:border-gray-700"
			>
				<div class="flex items-center gap-3">
					<div class="p-2 bg-app-bg border border-gray-800 rounded-md shrink-0">
						<svg class="w-5 h-5 fill-current text-text-primary" viewBox="0 0 24 24">
							<path
								d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm font-medium text-text-primary">GitHub</p>
						<p class="text-xs text-gray-500">Não conectado</p>
					</div>
				</div>

				<!-- Disconnected state -->
				<button
					type="button"
					class="px-3 py-1.5 bg-app-bg border border-gray-700 text-gray-300 hover:text-white hover:border-gray-600 rounded-md text-xs font-medium transition-all cursor-pointer"
				>
					Vincular
				</button>
			</div>
		</div>
	</section>

	<!-- Security -->
	<section class="bg-app-surface border border-gray-800 rounded-xl p-6 shadow-xl space-y-6">
		<div class="border-b border-gray-800/80 pb-4">
			<h2 class="font-montserrat font-semibold text-lg text-text-primary">Segurança</h2>
			<p class="text-xs text-gray-400 mt-0.5">Gerencie sua senha e autenticação de dois fatores.</p>
		</div>

		<div class="space-y-6">
			<!-- Change password -->
			<div
				class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 bg-app-bg/40 border border-gray-800 rounded-lg"
			>
				<div>
					<p class="text-sm font-medium text-text-primary">Senha de acesso</p>
					<p class="text-xs text-gray-400 mt-0.5">Última alteração há mais de 30 dias</p>
				</div>

				<a
					href="/"
					class="px-4 py-2 bg-app-bg border border-gray-700 text-gray-300 hover:text-white hover:border-gray-600 rounded-lg text-xs font-medium transition-all shrink-0 cursor-pointer flex items-center justify-center gap-2"
				>
					<svg
						class="w-4 h-4 text-gray-400"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 0121 9z"
						/>
					</svg>
					<span>Alterar senha</span>
				</a>
			</div>

			<!-- 2FA -->
			<div
				class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 bg-app-bg/40 border border-gray-800 rounded-lg"
			>
				<div class="space-y-1">
					<div class="flex items-center gap-2">
						<p class="text-sm font-medium text-text-primary">Autenticação de Dois Fatores (2FA)</p>
						<span
							class="text-xs uppercase font-bold px-2 py-0.5 rounded-md border invisible ml:visible
							{page.data.user?.is2FAEnabled
								? 'bg-emerald-950 border-emerald-800 text-emerald-300'
								: 'bg-amber-950 border-amber-800 text-amber-300'}
							 "
						>
							{page.data.user?.is2FAEnabled ? "Ativado" : "Desativado"}
						</span>
					</div>
					<p class="text-xs text-gray-400 max-w-md">
						Adicione uma camada extra de segurança usando um aplicativo autenticador (Google
						Authenticator, Authy, etc).
					</p>
				</div>
				{#if page.data.user?.is2FAEnabled}
					<button
						type="button"
						onclick={() => controller.open2FAModal()}
						disabled={controller.isLoading}
						aria-haspopup="dialog"
						aria-expanded={controller.is2FAModalOpen}
						class="px-4 py-2 rounded-lg border font-semibold text-sm focus:outline-none focus-visible:ring-2
						focus-visible:ring-offset-2 transition-all flex items-center justify-center gap-2 disabled:opacity-50
						cursor-pointer disabled:cursor-not-allowed shrink-0 w-full sm:w-auto
						border-red-900/50 bg-red-950/30 text-red-400 hover:bg-red-900/50 hover:text-red-300
						focus-visible:ring-red-500 focus-visible:ring-offset-gray-900"
					>
						<svg
							class="w-4 h-4"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
							/>
						</svg>
						<span>Desativar 2FA</span>
					</button>
				{:else}
					<button
						type="button"
						onclick={() => controller.open2FAModal()}
						disabled={controller.isLoading}
						aria-haspopup="dialog"
						aria-expanded={controller.is2FAModalOpen}
						class="px-4 py-2 rounded-lg border font-semibold text-sm focus:outline-none focus-visible:ring-2
						focus-visible:ring-offset-2 transition-all flex items-center justify-center gap-2 disabled:opacity-50
						cursor-pointer disabled:cursor-not-allowed shrink-0 w-full sm:w-auto
						border-emerald-900/50 bg-emerald-950/30 text-emerald-400 hover:bg-emerald-900/50 hover:text-emerald-300
						focus-visible:ring-emerald-500 focus-visible:ring-offset-gray-900"
					>
						<svg
							class="w-4 h-4"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
							/>
						</svg>
						<span>Ativar 2FA</span>
					</button>
				{/if}
			</div>
		</div>
	</section>
</div>

{#if controller.is2FAModalOpen}
	<div
		use:focusTrap
		class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm overflow-y-auto flex min-h-full items-center justify-center p-4"
		role="dialog"
		aria-modal="true"
		aria-busy={controller.isLoading}
		aria-labelledby="two-factor-modal-title"
		aria-describedby="two-factor-modal-description"
		onkeydown={(e) => e.key === "Escape" && controller.close2FAModal()}
		tabindex="-1"
	>
		<div
			class="bg-app-surface border border-gray-800 rounded-xl p-6 max-w-md w-full flex flex-col gap-5 shadow-2xl animate-in fade-in zoom-in-95 duration-150 relative my-auto"
		>
			<!-- Modal Header -->
			<div class="flex items-start gap-3">
				<div
					class="p-2.5 self-start rounded-lg shrink-0 border bg-emerald-950/80 border-emerald-900/60 text-emerald-400"
				>
					<svg
						class="w-6 h-6"
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
				</div>

				<div class="flex-1 pr-6">
					<h2 id="two-factor-modal-title" class="text-lg font-bold text-gray-100 font-montserrat">
						Ativar autenticação em dois fatores
					</h2>
					<p
						id="two-factor-modal-description"
						class="text-sm text-gray-300 mt-1 leading-relaxed"
						aria-live="polite"
					>
						{#if !controller.twoFactorSecret}
							{twoFactorLabels.generateCode}
						{:else}
							{twoFactorLabels.saveCode}
						{/if}
					</p>
				</div>

				<button
					type="button"
					onclick={() => controller.close2FAModal()}
					disabled={controller.isLoading}
					aria-label="Fechar modal"
					class="hover:cursor-pointer absolute top-4 right-4 text-gray-400 hover:text-white p-1 rounded-lg focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 disabled:opacity-50 transition-colors"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			<!-- FORMS -->
			{#if !controller.twoFactorSecret}
				<form onsubmit={(e) => controller.generate2FA(e)} class="space-y-5 font-inter">
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
								controller.apiError?.code === 'AUTH_INCORRECT_PASSWORD'
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
										<path
											d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"
										/>
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

					<div
						class="flex flex-col-reverse sm:flex-row justify-end gap-3 pt-4 border-t border-gray-800"
					>
						<button
							type="button"
							onclick={() => controller.close2FAModal()}
							disabled={controller.isLoading}
							class="w-full sm:w-auto px-4 py-2 rounded-lg bg-gray-900 text-gray-300 hover:text-white border border-gray-700 text-sm font-medium hover:bg-gray-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
						>
							Cancelar
						</button>
						<button
							type="submit"
							disabled={controller.isLoading}
							aria-busy={controller.isLoading}
							class="
					w-full sm:w-auto px-4 py-2 rounded-lg text-sm font-semibold focus:outline-none focus-visible:ring-2
					transition-colors disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed border
					border-emerald-900/60 bg-emerald-950/70 text-emerald-400 hover:bg-emerald-900/50 hover:text-emerald-300
					focus-visible:ring-emerald-500 focus-visible:ring-offset-gray-900 flex items-center justify-center gap-2"
						>
							{#if controller.isLoading}
								<svg
									class="animate-spin h-4 w-4 text-emerald-400"
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
								<span>Gerando...</span>
							{:else}
								<span>Gerar chave</span>
							{/if}
						</button>
					</div>
				</form>
			{:else}
				<form onsubmit={(e) => controller.save2FA(e)} class="space-y-5 font-inter">
					<div
						transition:slide={{ duration: 250 }}
						class="flex flex-col items-center gap-5 p-4 bg-black/20 border border-gray-800 rounded-lg"
					>
						<div class="bg-white p-2.5 rounded-xl shadow-inner shrink-0">
							<img
								src={controller.qrCodeUrl}
								alt="QR Code para Autenticação 2FA"
								class="w-48 h-48 rounded"
							/>
						</div>

						<div class="w-full space-y-2">
							<label for="twoFactorSecret" class="text-xs font-medium text-gray-300">
								Chave secreta (se não conseguir escanear):
							</label>
							<div class="flex items-center gap-2">
								<div class="flex-1 p-2.5 bg-app-bg border border-gray-700 rounded-lg shadow-inner">
									<code
										id="twoFactorSecret"
										aria-label="Chave secreta de configuração"
										class="text-sm font-mono font-medium text-emerald-300 break-all leading-relaxed tracking-wider"
									>
										{controller.twoFactorSecret || "Gerando chave..."}
									</code>
								</div>

								<button
									type="button"
									title="Copiar chave secreta"
									aria-label="Copiar chave secreta para a área de transferência"
									class="p-2.5 rounded-lg border border-gray-700 bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-emerald-500"
								>
									<svg
										class="w-5 h-5"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>
					<div class="space-y-2">
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
							oninput={(e) => controller.on2FAInput(e)}
							onblur={() => (controller.touched.code = true)}
							aria-required="true"
							aria-invalid={controller.touched.code && !controller.isCodeValid}
							aria-describedby="code-hint {controller.touched.code && !controller.isCodeValid
								? 'code-error'
								: ''}"
							placeholder="000000"
							required
							disabled={controller.isLoading}
							class="w-full px-4 pr-10 py-2.5 bg-app-bg/50 border rounded-lg text-text-primary placeholder-gray-600
                           text-center font-mono text-lg tracking-[0.5em] focus:bg-app-bg focus:ring-2 focus:outline-none
                           transition-all disabled:opacity-50
                    {(controller.touched.code && !controller.isCodeValid) ||
							controller.apiError?.code === '2FA_INVALID_CODE'
								? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
								: 'border-gray-800 focus:border-text-brand focus:ring-text-brand/20'}"
						/>
					</div>
					<div
						class="flex flex-col-reverse sm:flex-row justify-end gap-3 pt-4 border-t border-gray-800"
					>
						<button
							type="button"
							onclick={() => controller.close2FAModal()}
							disabled={controller.isLoading}
							class="w-full sm:w-auto px-4 py-2 rounded-lg bg-gray-900 text-gray-300 hover:text-white border border-gray-700 text-sm font-medium hover:bg-gray-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
						>
							Cancelar
						</button>
						<button
							type="submit"
							disabled={controller.isLoading}
							aria-busy={controller.isLoading}
							class="
					w-full sm:w-auto px-4 py-2 rounded-lg text-sm font-semibold focus:outline-none focus-visible:ring-2
					transition-colors disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed border
					border-emerald-900/60 bg-emerald-950/70 text-emerald-400 hover:bg-emerald-900/50 hover:text-emerald-300
					focus-visible:ring-emerald-500 focus-visible:ring-offset-gray-900 flex items-center justify-center gap-2"
						>
							{#if controller.isLoading}
								<svg
									class="animate-spin h-4 w-4 text-emerald-400"
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
								<span>Salvando...</span>
							{:else}
								<span>Salvar chave</span>
							{/if}
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
{/if}
