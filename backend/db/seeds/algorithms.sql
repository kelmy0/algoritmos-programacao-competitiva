INSERT INTO algorithms (name, category, difficulty, content) VALUES
('Busca Binária', 'Busca', 'beginner', '```cpp cout << "teste Busca Binária"; ```'),
('Algoritmo de Dijkstra', 'Grafos', 'intermediate', '```cpp cout << "teste Dijkstra"; ```'),
('Segment Tree', 'Estrutura de Dados', 'advanced', '```cpp cout << "teste Segment Tree"; ```'),
('Heavy-Light Decomposition', 'Grafos', 'expert', '```cpp cout << "teste Heavy-Light Decomposition"; ```')
ON CONFLICT (name) DO UPDATE SET
    category = EXCLUDED.category,
    difficulty = EXCLUDED.difficulty,
    content = EXCLUDED.content,
    updated_at = CURRENT_TIMESTAMP;