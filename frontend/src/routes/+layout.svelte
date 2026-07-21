<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/state';

	let { children } = $props();
	let isSidebarOpen = $state(false);
</script>

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
			<span class="text-text-brand">&lt;/&gt;</span> Algoritmos para Maratona
		</a>

		<!--Search field-->
		<div class="flex-1 max-w-md mx-8 hidden md:block">
			<div class="relative">
				<span
					class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-gray-500"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
				</span>
				<input
					type="search"
					placeholder="Pesquisar..."
					class="w-full pl-10 pr-4 py-2 bg-app-bg/50 border border-gray-800 rounded-lg text-text-primary
					placeholder-gray-500 text-sm focus:bg-app-bg focus:border-text-brand focus:ring-1
					focus:ring-text-brand focus:outline-none transition-all"
				/>
			</div>
		</div>

		<!--Login button-->
		<div class="items-center hidden md:block">
			<a
				href="/login"
				class="px-4 py-1.5 border border-text-brand
			text-text-brand hover:bg-text-brand hover:text-app-bg font-medium text-sm rounded-lg transition-all"
				aria-label="Sign up"
			>
				Entrar
			</a>
		</div>

		<!--Mobile Menu Button-->
		<div class="flex items-center gap-3 md:hidden">
			<button
				type="button"
				onclick={() => (isSidebarOpen = !isSidebarOpen)}
				class="p-2 text-text-brand hover:bg-app-surface rounded-lg transition-colors cursor-pointer"
				aria-label="Toggle Menu"
			>
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
			<button
				onclick={() => (isSidebarOpen = false)}
				class="fixed inset-0 bg-black/60 z-40 md:hidden cursor-default w-full h-full
				border-0 p-0 appearance-none backdrop-blur-sm"
				aria-label="Close menu"
			></button>
		{/if}

		<aside
			class="fixed top-16 bottom-0 right-0 md:z-40 z-50 w-64 bg-app-surface border-l border-gray-800 md:border-l-0
    md:border-r transform transition-transform duration-200 ease-in-out md:sticky md:left-0 md:right-auto md:translate-x-0
    {isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}"
		>
			<nav class="p-4 flex flex-col h-full justify-between">
				<!-- Top links -->
				<div class="space-y-1">
					<a
						href="/"
						class="flex items-center gap-3 px-4 py-2.5
                bg-app-bg/50 text-text-brand border-l-2 border-text-brand rounded-r-lg font-medium transition-colors
					{page.url.pathname === '/'
							? 'bg-app-bg/50 text-text-brand border-text-brand'
							: 'text-gray-400 border-transparent hover:bg-app-bg/30 hover:text-text-primary'}
				"
						onclick={() => (isSidebarOpen = false)}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
						<span>Início</span>
					</a>

					<!-- Other future navigation links would go here. -->
				</div>

				<div class="pt-4 border-t border-gray-800 md:hidden mt-auto">
					<a
						href="/login"
						onclick={() => (isSidebarOpen = false)}
						class="w-full py-2.5 border border-text-brand text-text-brand
                hover:bg-text-brand hover:text-app-bg font-semibold text-sm rounded-lg
                transition-all flex items-center justify-center gap-2"
						aria-label="Sign in"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
							/>
						</svg>
						Login
					</a>
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
