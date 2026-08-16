export const DIFFICULTY_MAP = {
	beginner: {
		label: "Iniciante",
		style: "bg-emerald-950 text-emerald-300 border-emerald-800"
	},
	intermediate: {
		label: "Intermediário",
		style: "bg-amber-950 text-amber-300 border-amber-800"
	},
	advanced: { label: "Avançado", style: "bg-orange-950 text-orange-300 border-orange-800" },
	expert: { label: "Especialista", style: "bg-red-950 text-red-300 border-red-800" }
} as const;

export const STATUS_MAP = {
	approved: {
		label: "Aprovado",
		style: "bg-emerald-950/40 text-emerald-300 border-emerald-800/60"
	},
	pending: { label: "Pendente", style: "bg-amber-950/40 text-amber-300 border-amber-800/60" },
	rejected: { label: "Rejeitado", style: "bg-orange-950/40 text-orange-300 border-orange-800/60" },
	deleted: { label: "Deletado", style: "bg-red-950/40 text-red-300 border-red-800/60" }
} as const;
