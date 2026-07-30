INSERT INTO roles(name, is_employee) VALUES
('user', FALSE), ('admin', TRUE), ('moderator', TRUE), ('creator', TRUE), ('manager', TRUE),
('viewer', TRUE);

INSERT INTO permissions(slug) VALUES
('create:algorithms'), 
('moderate:algorithms'),
('manage:categories'),
('manage:users'),
('manage:settings'),
('view:reports'),
('view:audit_logs');

INSERT INTO role_permissions(role_id, permission_id) VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7),
(3, 2),
(4, 1),
(5, 3), (5, 4), (5, 5),
(6, 6), (6, 7);

INSERT INTO users(name, username, email, password_hash, role_id) VALUES
('user teste', 'user_teste', 'user@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 1),
('Adm teste', 'adm_teste', 'admin@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 2),
('Moderator teste', 'moderator_teste', 'moderator@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 3),
('Creator teste', 'creator_teste', 'creator@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 4),
('Manager teste', 'manager_teste', 'manager@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 5),
('Viewer teste', 'viewer_teste', 'viewer@gmail.com', '$argon2id$v=19$m=65536,t=3,p=4$IIZofmSeiiATyVyGC3cmgg$A4tCBHsy869mMWKkL8Cmj7z+Hfzjwsxaly2x7AmXrCA', 6);
