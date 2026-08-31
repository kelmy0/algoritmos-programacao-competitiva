package seeds

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

func SeedAlgorithms(db *pgxpool.Pool) {
	algos := []dto.PostAlgorithmRequest{
		{
			Name:       "Busca Binária",
			Category:   "Busca",
			Difficulty: "beginner",
			Content: `# Busca Binária

A **Busca Binária** é um algoritmo eficiente para encontrar um elemento em uma lista que já está **ordenada**.

## Como funciona?
1. Encontre o elemento do meio.
2. Se o valor for igual ao alvo, a busca termina.
3. Se for diferente, descarte a metade irrelevante.

### Exemplo em C++
` + "```cpp" + `
#include <iostream>
#include <vector>

int main() {
    std::vector<int> nums = {1, 3, 5, 7, 9, 11};
    std::cout << "Executando teste da Busca Binária..." << std::endl;
    return 0;
}
` + "```" + `

> **Complexidade de Tempo:** O(log n)`,
		},
		{
			Name:       "Algoritmo de Dijkstra",
			Category:   "Grafos",
			Difficulty: "intermediate",
			Content: `# Algoritmo de Dijkstra

O **Algoritmo de Dijkstra** resolve o problema do caminho mínimo em um grafo com pesos não negativos.

## Aplicações
* Sistemas de GPS e Mapas.
* Roteamento de pacotes na rede.

### Exemplo em C++
` + "```cpp" + `
#include <iostream>
#include <queue>

int main() {
    std::cout << "Inicializando grafo e executando Dijkstra..." << std::endl;
    return 0;
}
` + "```" + `

> **Nota:** Não suporta arestas com peso negativo.`,
		},
		{
			Name:       "Segment Tree",
			Category:   "Estrutura de Dados",
			Difficulty: "advanced",
			Content: `# Segment Tree

A **Segment Tree** permite realizar consultas e atualizações em intervalos de um array de forma muito rápida.

## Operações
* **Point Update:** O(log n)
* **Range Query:** O(log n)

### Exemplo em C++
` + "```cpp" + `
#include <iostream>
#include <vector>

int main() {
    std::cout << "Construindo Segment Tree..." << std::endl;
    return 0;
}
` + "```" + `

> **Memória:** Requer aproximadamente 4n de espaço.`,
		},
		{
			Name:       "Heavy-Light Decomposition",
			Category:   "Grafos",
			Difficulty: "expert",
			Content: `# Heavy-Light Decomposition (HLD)

Técnica avançada usada para decompor árvores em cadeias de caminhos disjuntos.

## Aplicação
Permite responder a consultas de caminho entre nós em árvores usando uma **Segment Tree**.

### Exemplo em C++
` + "```cpp" + `
#include <iostream>

int main() {
    std::cout << "Decompondo árvore com HLD..." << std::endl;
    return 0;
}
` + "```" + `

> **Complexidade:** O(q log^2 n)`,
		},
	}

	userRepo := repositories.NewUserRepository(db)
	algoRepo := repositories.NewAlgorithmRepository(db)
	algoService := services.NewAlgorithmService(algoRepo, userRepo)

	ctx := context.Background()

	user, err := userRepo.GetUserByEmailForAuth(ctx, "admin@gmail.com")
	if err != nil {
		slog.Error("❌ Error querying admin user for seed", "email", "admin@gmail.com", "error", err)
		return
	}

	slog.Info("🌱 Starting algorithm database seed...")

	for i := range algos {
		res, err := algoService.PostAlgorithm(ctx, algos[i], user.Id)
		if err != nil {
			slog.Error("❌ Failed to seed algorithm", "name", algos[i].Name, "error", err)
			continue
		}

		slog.Info("✅ Algorithm seeded successfully", "slug", res.Slug)
	}

	slog.Info("🎉 Seed process completed!")
}
