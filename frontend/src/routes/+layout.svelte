<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
	let isSidebarOpen = $state(false);
</script>

<svelte:head>
	<title>Algoritmos para Maratona de Programação</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-gray-100 text-gray-800 font-sans flex flex-col">
	<!--Topbar-->
	<header
		class="bg-indigo-950 border-b border-indigo-300 h-16
	fixed top-0 left-0 right-0 flex items-center justify-between px-4 z-10"
	>
		<!--Just the title-->
		<a class="font-bold text-xl text-gray-100" href="/">Algoritmos para Maratona</a>

		<!--Search field-->
		<div class="flex-1 max-w-md mx-4 hidden md:block">
			<div class="relative">
				<input
					type="search"
					placeholder="Pesquisar..."
					class="w-full pl-10 pr-4 py-2 bg-gray-100 border border-transparent rounded-lg focus:bg-white focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 focus:outline-none transition-all text-sm"
				/>
			</div>
		</div>

		<!--Button to open sidebar-->
		<div class="flex items-center gap-3 md:hidden">
			<button
				type="button"
				onclick={() => (isSidebarOpen = !isSidebarOpen)}
				class="p-2 text-gray-400 hover:bg-blue-200/30 rounded-lg"
				aria-label="Toggle Menu"
			>
				<svg class="w-6 h-6 " fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

	<!--Sidebar-->
	<div class="flex flex-1 pt-16 ">
		{#if isSidebarOpen}
			<button
				onclick={() => (isSidebarOpen = false)}
				class="fixed inset-0 bg-black/40 z-40 md:hidden cursor-default
				w-full h-full border-0 p-0 appearance-none"
				aria-label="Fechar menu"
			></button>
		{/if}

		<aside
			class="fixed top-16 bottom-0 right-0 z-50 w-64 bg-gray-200 md:bg-gray-300/50 border-l
			border-gray-200 transform transition-transform duration-200 ease-in-out
        	md:sticky md:left-0 md:right-auto md:border-l-0 md:border-r md:translate-x-0
        {isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}"
		>
			<nav class="p-4 space-y-1">
				<a
					href="/"
					class="flex items-center gap-3 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg font-medium transition-colors"
				>
					<span>Início</span>
				</a>
			</nav>
		</aside>

		<!--SPA container-->
		<main class="flex-1 p-6 min-w-0 overflow-y-auto bg-gray-200/50">
			<div class="max-w-7xl mx-auto ">
				{@render children()}
			</div>
		</main>
	</div>
</div>
