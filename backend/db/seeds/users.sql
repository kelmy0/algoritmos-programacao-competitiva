INSERT INTO roles(name, is_employee) VALUES
('user', FALSE), ('admin', TRUE);

INSERT INTO permissions(slug) VALUES
('create:algorithms'), 
('update:algorithms'), 
('delete:algorithms');

INSERT INTO role_permissions(role_id, permission_id) VALUES
(2, 1), 
(2, 2), 
(2, 3);

INSERT INTO users(name, username, email, password_hash, role_id) VALUES
('Adm teste', 'adm_teste', 'teste@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 2),
('user teste', 'user_teste', 'uteste@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 1);