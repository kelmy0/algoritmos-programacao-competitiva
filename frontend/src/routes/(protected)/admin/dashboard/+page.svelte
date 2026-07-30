<script>
	import { page } from "$app/state";

	const adminActions = [
		{
			title: "Novo Algoritmo",
			description: "Cadastre um novo algoritmo ou estrutura de dados no sistema.",
			href: "/admin/algorithms/new",
			badge: "Criar",
			badgeColor: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
			iconPath: "M12 4v16m8-8H4",
			permission: "create:algorithms"
		},
		{
			title: "Gerenciar & Editar",
			description: "Edite informações, corrija conteúdos ou atualize os algoritmos existentes.",
			href: "/admin/algorithms/edit",
			badge: "Editar",
			badgeColor: "bg-amber-500/10 text-amber-400 border-amber-500/30",
			iconPath:
				"M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
			permission: "create:algorithms"
		},
		{
			title: "Meus Algoritmos",
			description:
				"Visualize seus algoritmos cadastrados exatamente como aparecem para os leitores na plataforma.",
			href: "/admin/algorithms/my-algorithms",
			badge: "Leitura",
			badgeColor: "bg-teal-500/10 text-teal-400 border-teal-500/30",
			iconPath:
				"M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z",
			permission: "create:algorithms"
		},
		{
			title: "Fila de Moderação",
			description:
				"Análise submissões, aprove algoritmos de usuários ou rejeite conteúdos inadequados.",
			href: "/admin/algorithms/moderation",
			badge: "Moderar",
			badgeColor: "bg-blue-500/10 text-blue-400 border-blue-500/30",
			iconPath:
				"M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z",
			permission: "moderate:algorithms"
		},
		{
			title: "Categorias & Tags",
			description:
				"Crie, edite e organize linguagens de programação, tópicos e tags dos algoritmos.",
			href: "/admin/categories",
			badge: "Organizar",
			badgeColor: "bg-indigo-500/10 text-indigo-400 border-indigo-500/30",
			iconPath:
				"M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z",
			permission: "manage:categories"
		},
		{
			title: "Relatórios & Denúncias",
			description: "Analise conteúdos denunciados pela comunidade por bugs, erros ou violações.",
			href: "/admin/reports",
			badge: "Alertas",
			badgeColor: "bg-orange-500/10 text-orange-400 border-orange-500/30",
			iconPath:
				"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z",
			permission: "view:reports"
		},
		{
			title: "Gestão de Usuários",
			description: "Gerencie contas, altere cargos (Employee/User) e ajuste permissões de acesso.",
			href: "/admin/users",
			badge: "Acessos",
			badgeColor: "bg-purple-500/10 text-purple-400 border-purple-500/30",
			iconPath:
				"M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z",
			permission: "manage:users"
		},
		{
			title: "Logs de Auditoria",
			description: "Acompanhe o histórico de ações realizadas por funcionários e moderadores.",
			href: "/admin/audit-logs",
			badge: "Logs",
			badgeColor: "bg-sky-500/10 text-sky-400 border-sky-500/30",
			iconPath:
				"M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z",
			permission: "view:audit_logs"
		},
		{
			title: "Lixeira & Removidos",
			description: "Visualize seus itens excluídos ou desativados da plataforma.",
			href: "/admin/algorithms/trash",
			badge: "Deletar",
			badgeColor: "bg-rose-500/10 text-rose-400 border-rose-500/30",
			iconPath:
				"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
			permission: "create:algorithms"
		},
		{
			title: "Configurações Globais",
			description: "Ajuste parâmetros do sistema, integrações com APIs e modo de manutenção.",
			href: "/admin/settings",
			badge: "Sistema",
			badgeColor: "bg-zinc-500/10 text-zinc-400 border-zinc-500/30",
			iconPath:
				"M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
			permission: "manage:settings"
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
