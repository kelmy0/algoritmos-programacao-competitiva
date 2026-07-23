<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/state';
	import { AuthService } from '$lib/services/auth_service';
	import { fade } from 'svelte/transition';

	let { children } = $props();

	let isSidebarOpen = $state(false);
	let isProfileMenuOpen = $state(false);

	$effect(() => {
		if (page.data.accessToken) {
			AuthService.startAutoRefreshTimer(page.data.accessToken);
		}
	});

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeMobileMenu();
			closeProfileMenu();
		}
	}

	function closeMobileMenu() {
		isSidebarOpen = false;
	}

	function closeProfileMenu() {
		isProfileMenuOpen = false;
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<svelte:head>
	<title>Algoritmos para Maratona de Programação</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-app-bg text-text-primary font-inter flex flex-col">
	<!--Topbar-->
	<header
		class="bg-app-surface h-16 border-b border-gray-800 fixed top-0 left-0 right-0 flex items-center justify-between px-6 z-50"
	>
		<!--Title-->
		<a
			class="font-bold text-xl text-text-primary font-montserrat tracking-tight hover:text-text-brand transition-colors flex items-center gap-2"
			href="/"
		>
			<span class="text-text-brand">&lt;/&gt;</span>
			<span>Algoritmos para Maratona</span>
		</a>

		<!--Search field-->
		<div class="flex-1 max-w-md mx-8 hidden md:block">
			<form role="search" onsubmit={(e) => e.preventDefault()}>
				<div class="relative">
					<span
						class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-gray-500"
						aria-hidden="true"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0"
							/>
						</svg>
					</span>
					<input
						type="search"
						id="global-search"
						aria-label="Pesquisar algoritmos e tópicos"
						placeholder="Pesquisar..."
						class="w-full pl-10 pr-4 py-2 bg-app-bg/50 border border-gray-800 rounded-lg text-text-primary placeholder-gray-500 text-sm focus:bg-app-bg focus:border-text-brand focus:ring-1 focus:ring-text-brand focus:outline-none transition-all"
					/>
				</div>
			</form>
		</div>

		<!--Login button-->
		<div class="items-center hidden md:flex relative">
			{#if page.data.user}
				{#if isProfileMenuOpen}
					<div
						aria-hidden="true"
						tabindex="-1"
						onclick={closeProfileMenu}
						class="fixed inset-0 z-40 bg-transparent cursor-default"
					></div>
				{/if}

				<!-- Profile Button -->
				<button
					type="button"
					onclick={() => (isProfileMenuOpen = !isProfileMenuOpen)}
					class="flex items-center gap-2.5 p-1.5 rounded-lg hover:bg-app-bg border border-transparent hover:border-gray-800 transition-all focus:outline-none focus:ring-1 focus:ring-text-brand z-50 cursor-pointer"
					aria-expanded={isProfileMenuOpen}
					aria-haspopup="true"
					aria-label="Menu do usuário para {page.data.user.username ?? 'Usuário'}"
					aria-controls="user-profile-menu"
				>
					<!-- Avatar / Icon -->
					<div
						aria-hidden="true"
						class="w-8 h-8 rounded-lg bg-text-brand/10 border border-text-brand/30 text-text-brand flex items-center justify-center font-bold text-sm font-montserrat"
					>
						{page.data.user.username ? page.data.user.username[0].toUpperCase() : 'U'}
					</div>

					<!-- Username -->
					<span class="text-sm font-medium text-text-primary max-w-30 truncate">
						{page.data.user.username ?? 'Usuário'}
					</span>

					<!-- Dropdown indicator -->
					<svg
						class="w-4 h-4 text-gray-400 transition-transform duration-200 {isProfileMenuOpen
							? 'rotate-180'
							: ''}"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M19 9l-7 7-7-7"
						/>
					</svg>
				</button>

				<!-- Dropdown Menu -->
				{#if isProfileMenuOpen}
					<div
						transition:fade={{ duration: 120 }}
						id="user-profile-menu"
						role="menu"
						aria-label="Opções do perfil"
						class="absolute right-0 top-full mt-2 w-56 bg-app-surface border border-gray-800 rounded-xl shadow-2xl py-2 z-50 flex flex-col text-sm"
					>
						<!-- Menu header -->
						<div class="px-4 py-2 border-b border-gray-800/80 mb-1">
							<p class="text-xs text-gray-400">Conectado como</p>
							<p class="font-semibold text-text-primary truncate">
								{page.data.user.username ?? 'Usuário'}
							</p>
						</div>

						<!-- Links / Options -->
						<a
							href="/"
							role="menuitem"
							onclick={closeProfileMenu}
							class="flex items-center gap-2.5 px-4 py-2 text-gray-300 hover:text-text-primary hover:bg-app-bg transition-colors"
						>
							<svg
								aria-hidden="true"
								class="w-4 h-4 text-gray-400"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
								/>
							</svg>
							Meu Perfil
						</a>

						<a
							href="/"
							role="menuitem"
							onclick={closeProfileMenu}
							class="flex items-center gap-2.5 px-4 py-2 text-gray-300 hover:text-text-primary hover:bg-app-bg transition-colors"
						>
							<svg
								aria-hidden="true"
								class="w-4 h-4 text-gray-400"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
								/>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
								/>
							</svg>
							Configurações
						</a>

						<div class="my-1 border-t border-gray-800"></div>

						<!-- Logout -->
						<form method="POST" action="/auth/logout">
							<button
								type="submit"
								role="menuitem"
								class="w-full flex items-center gap-2.5 px-4 py-2 text-red-400 hover:text-red-300
						hover:bg-red-500/10 transition-colors text-left font-medium"
							>
								<svg
									aria-hidden="true"
									class="w-4 h-4"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
									/>
								</svg>
								Sair
							</button>
						</form>
					</div>
				{/if}
			{:else}
				<a
					href="/auth/login"
					class="px-4 py-1.5 border border-text-brand text-text-brand hover:bg-text-brand hover:text-app-bg font-medium text-sm rounded-lg transition-all focus:outline-none focus:ring-2 focus:ring-text-brand"
				>
					Entrar
				</a>
			{/if}
		</div>

		<!--Mobile Menu Button-->
		<div class="flex items-center gap-3 md:hidden">
			<button
				type="button"
				onclick={() => (isSidebarOpen = !isSidebarOpen)}
				class="p-2 text-text-brand hover:bg-app-bg rounded-lg transition-colors cursor-pointer focus:outline-none focus:ring-2 focus:ring-text-brand"
				aria-expanded={isSidebarOpen}
				aria-controls="mobile-sidebar"
				aria-label={isSidebarOpen ? 'Fechar menu de navegação' : 'Abrir menu de navegação'}
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
						d="M4 6h16M4 12h16M4 18H16"
					/>
				</svg>
			</button>
		</div>
	</header>

	<div class="flex flex-1 pt-16">
		<!--Mobile Overlay Backdrop-->
		{#if isSidebarOpen}
			<div
				onclick={closeMobileMenu}
				aria-hidden="true"
				tabindex="-1"
				class="fixed inset-0 bg-black/60 z-40 md:hidden backdrop-blur-sm border-0 cursor-default"
			></div>
		{/if}
		<!-- ========================================== -->
		<!--  		SIDEBAR MOBILE     					-->
		<!-- ========================================== -->
		<aside
			id="mobile-sidebar"
			aria-label="Navegação mobile"
			aria-hidden={!isSidebarOpen}
			class="fixed top-16 bottom-0 right-0 z-50 w-64 bg-app-surface border-l border-gray-800 transform transition-transform duration-200 ease-in-out md:hidden
    		{isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}"
		>
			<nav class="p-4 flex flex-col h-full justify-between">
				<div class="space-y-1">
					<a
						href="/"
						aria-current={page.url.pathname === '/' ? 'page' : undefined}
						class="flex items-center gap-3 px-4 py-2.5 rounded-r-lg font-medium transition-colors
                {page.url.pathname === '/'
							? 'bg-app-bg/50 text-text-brand border-l-2 border-text-brand'
							: 'text-gray-400 border-transparent hover:bg-app-bg/30 hover:text-text-primary'}"
						onclick={() => (isSidebarOpen = false)}
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
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
						<span>Início</span>
					</a>

					<!-- Other mobile links will go here -->
				</div>

				<div class="pt-4 border-t border-gray-800 mt-auto">
					{#if page.data.user}
						<div class="flex flex-col gap-3">
							<div class="flex items-center gap-3 px-2">
								<div
									class="w-8 h-8 rounded-full bg-text-brand/20 text-text-brand flex items-center justify-center font-bold text-sm"
								>
									{page.data.user.username?.[0]?.toUpperCase() ?? 'U'}
								</div>
								<div class="flex flex-col min-w-0">
									<span class="text-sm font-medium text-text-primary truncate">
										{page.data.user.username ?? 'Usuário'}
									</span>
									<span class="text-xs text-gray-500 truncate">
										{page.data.user.email ?? ''}
									</span>
								</div>
							</div>

							<div class="flex flex-col gap-1 pt-2">
								<a
									href="/"
									onclick={() => (isSidebarOpen = false)}
									class="flex items-center gap-2.5 px-3 py-2 text-sm text-gray-400 hover:text-text-primary hover:bg-app-bg/30 rounded-lg transition-colors"
								>
									<svg
										class="w-4 h-4"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
										/>
									</svg>
									Meu Perfil
								</a>

								<a
									href="/"
									onclick={() => (isSidebarOpen = false)}
									class="flex items-center gap-2.5 px-3 py-2 text-sm text-gray-400 hover:text-text-primary hover:bg-app-bg/30 rounded-lg transition-colors"
								>
									<svg
										aria-hidden="true"
										class="w-4 h-4 text-gray-400"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
										/>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
										/>
									</svg>
									Configurações
								</a>

								<!-- Logout -->
								<form method="POST" action="/auth/logout">
									<button
										type="submit"
										class="flex items-center gap-2.5 px-3 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-colors w-full text-left cursor-pointer"
									>
										<svg
											class="w-4 h-4"
											fill="none"
											stroke="currentColor"
											viewBox="0 0 24 24"
											aria-hidden="true"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
											/>
										</svg>
										Sair
									</button>
								</form>
							</div>
						</div>
					{:else}
						<a
							href="/auth/login"
							onclick={() => (isSidebarOpen = false)}
							class="w-full py-2.5 border border-text-brand text-text-brand hover:bg-text-brand hover:text-app-bg font-semibold text-sm rounded-lg transition-all flex items-center justify-center gap-2"
						>
							Entrar
						</a>
					{/if}
				</div>
			</nav>
		</aside>

		<!-- ========================================== -->
		<!-- 		SIDEBAR DESKTOP     				-->
		<!-- ========================================== -->

		<aside
			aria-label="Navegação principal"
			class="hidden md:block w-64 bg-app-surface border-r border-gray-800 sticky top-16 h-[calc(100vh-4rem)]"
		>
			<nav class="p-4 flex flex-col h-full justify-between">
				<div class="space-y-6">
					<div>
						<!-- MAIN MENU -->
						<p class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
							Menu
						</p>
						<div class="space-y-1">
							<a
								href="/"
								aria-current={page.url.pathname === '/' ? 'page' : undefined}
								class="flex items-center gap-3 px-4 py-2.5 rounded-r-lg font-medium transition-colors
                        {page.url.pathname === '/'
									? 'bg-app-bg/50 text-text-brand border-l-2 border-text-brand'
									: 'text-gray-400 border-transparent hover:bg-app-bg/30 hover:text-text-primary'}"
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
										d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
									/>
								</svg>
								<span>Início</span>
							</a>
						</div>
					</div>

					<!-- OTHERS -->
					<div>
						<p class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
							Categorias
						</p>

						<div class="space-y-1">
							<a
								href="/?category=tecnologia"
								aria-current={page.url.searchParams.get('category') === 'tecnologia'
									? 'page'
									: undefined}
								class="flex items-center justify-between px-4 py-2 rounded-lg text-sm font-medium transition-all
            {page.url.searchParams.get('category') === 'tecnologia'
									? 'bg-text-brand/10 text-text-brand'
									: 'text-gray-400 hover:bg-app-bg/40 hover:text-text-primary'}"
							>
								<div class="flex items-center gap-2.5 truncate">
									<svg
										class="w-4 h-4 opacity-70"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
										/>
									</svg>
									<span class="truncate">Grafos</span>
								</div>

								<!-- Count -->
								<span class="text-xs px-2 py-0.5 rounded-full bg-app-bg/80 text-gray-400 font-mono">
									5
								</span>
							</a>

							<a
								href="/?category=design"
								aria-current={page.url.searchParams.get('category') === 'design'
									? 'page'
									: undefined}
								class="flex items-center justify-between px-4 py-2 rounded-lg text-sm font-medium transition-all
            {page.url.searchParams.get('category') === 'design'
									? 'bg-text-brand/10 text-text-brand'
									: 'text-gray-400 hover:bg-app-bg/40 hover:text-text-primary'}"
							>
								<div class="flex items-center gap-2.5 truncate">
									<svg
										class="w-4 h-4 opacity-70"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
										/>
									</svg>
									<span class="truncate">Programação Dinâmica</span>
								</div>
								<span class="text-xs px-2 py-0.5 rounded-full bg-app-bg/80 text-gray-400 font-mono">
									2
								</span>
							</a>
						</div>
					</div>
				</div>

				<div class="pt-4 border-t border-gray-800 text-xs text-gray-500 text-center">
					v1.0.0 • Algoritmos de Programação
				</div>
			</nav>
		</aside>

		<!--SPA container-->
		<main class="flex-1 p-6 min-w-0 overflow-y-auto">
			<div class="max-w-7xl mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
</div>
