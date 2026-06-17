INSERT INTO algorithms (id, name, category, difficulty) VALUES
('1', 'Busca Binária', 'Busca', 'beginner'),
('2', 'Algoritmo de Dijkstra', 'Grafos', 'intermediate'),
('3', 'Segment Tree', 'Estrutura de Dados', 'advanced'),
('4', 'Heavy-Light Decomposition', 'Grafos', 'expert')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    category = EXCLUDED.category,
    difficulty = EXCLUDED.difficulty,
    updated_at = CURRENT_TIMESTAMP;