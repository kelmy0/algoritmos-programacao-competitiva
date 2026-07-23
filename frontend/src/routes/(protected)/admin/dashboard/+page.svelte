<script>
	import { page } from '$app/state';

	const adminActions = [
		{
			title: 'Novo Algoritmo',
			description: 'Cadastre um novo algoritmo ou estrutura de dados no sistema.',
			href: '/admin/algorithms/new',
			badge: 'Criar',
			badgeColor: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
			iconPath: 'M12 4v16m8-8H4',
			permission: 'create:algorithms'
		},
		{
			title: 'Gerenciar & Editar',
			description: 'Edite informações, corrija conteúdos ou atualize os algoritmos existentes.',
			href: '/admin/algorithms/edit',
			badge: 'Editar',
			badgeColor: 'bg-amber-500/10 text-amber-400 border-amber-500/30',
			iconPath:
				'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
			permission: 'update:algorithms'
		},
		{
			title: 'Lixeira & Removidos',
			description: 'Visualize ou exclua permanentemente itens desativados da plataforma.',
			href: '/admin/algorithms/trash',
			badge: 'Deletar',
			badgeColor: 'bg-rose-500/10 text-rose-400 border-rose-500/30',
			iconPath:
				'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16',
			permission: 'delete:algorithms'
		}
	];
</script>

<svelte:head>
	<title>Painel Administrativo</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="space-y-6 font-inter">
	<header
		class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-gray-800"
	>
		<div>
			<h1 class="font-montserrat font-bold text-2xl md:text-3xl text-text-primary tracking-tight">
				Painel Administrativo
			</h1>
			<p class="text-sm text-gray-400 mt-1">
				Gerencie os conteúdos, algoritmos e configurações gerais da plataforma.
			</p>
		</div>
	</header>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
		{#each adminActions as card}
			{#if page.data.user?.permissions.includes(card.permission)}
				<article
					class="relative bg-app-surface border border-gray-800 rounded-xl p-5 shadow-lg flex flex-col justify-between hover:border-gray-700 hover:shadow-xl transition-all duration-200 group"
				>
					<div class="space-y-4">
						<div class="flex items-start justify-between gap-3">
							<div
								class="p-2.5 rounded-lg bg-app-bg/60 border border-gray-800 text-text-brand group-hover:border-text-brand/40 transition-colors"
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
										d={card.iconPath}
									/>
								</svg>
							</div>

							<span
								class="text-xs font-semibold uppercase tracking-wider px-2.5 py-1 rounded-md border shrink-0 relative z-10 {card.badgeColor}"
							>
								{card.badge}
							</span>
						</div>

						<div>
							<h2
								class="font-montserrat font-semibold text-lg text-text-primary group-hover:text-text-brand transition-colors"
							>
								<a
									href={card.href}
									class="after:absolute after:inset-0 focus:outline-none focus:ring-2 focus:ring-text-brand focus:ring-offset-2 focus:ring-offset-app-surface rounded-xl"
								>
									{card.title}
								</a>
							</h2>
							<p class="text-xs text-gray-400 mt-2 leading-relaxed">
								{card.description}
							</p>
						</div>
					</div>

					<div class="pt-4 mt-4 border-t border-gray-800/80 flex items-center justify-end">
						<div
							aria-hidden="true"
							class="text-xs font-medium text-text-brand group-hover:underline flex items-center gap-1 transition-all pointer-events-none"
						>
							<span>Acessar rota</span>
							<svg
								class="w-4 h-4 group-hover:translate-x-0.5 transition-transform"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
							</svg>
						</div>
					</div>
				</article>
			{/if}
		{/each}
	</div>
</div>
