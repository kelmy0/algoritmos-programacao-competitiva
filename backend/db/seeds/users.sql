INSERT INTO roles(name, is_employee) VALUES
('user', FALSE), ('admin', TRUE);

INSERT INTO users(name, username, email, password_hash, role_id) VALUES
('Adm teste', 'adm_teste', 'teste@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$RKCxwjgfKfv+/1CDNn39Qg$dWA+Ijh6T0f/kvT18NMoiOSYszW31chyL5e71cj81x0', 2),
('user teste', 'user_teste', 'uteste@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$RKCxwjgfKfv+/1CDNn39Qg$dWA+Ijh6T0f/kvT18NMoiOSYszW31chyL5e71cj81x0', 1);