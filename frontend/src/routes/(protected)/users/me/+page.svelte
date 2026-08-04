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
						value="dev@exemplo.com"
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
							class="text-[10px] uppercase font-bold px-2 py-0.5 bg-amber-950/60 border border-amber-800/80 text-amber-400 rounded-md"
						>
							Desativado
						</span>
					</div>
					<p class="text-xs text-gray-400 max-w-md">
						Adicione uma camada extra de segurança usando um aplicativo autenticador (Google
						Authenticator, Authy, etc).
					</p>
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
							d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
						/>
					</svg>
					<span>Ativar 2FA</span>
				</a>
			</div>
		</div>
	</section>
</div>
