INSERT INTO users (username, password, email, status, role, first_name, last_name, bio, profile_pic)
VALUES
('admin', '$2a$12$examplehashedpassword', 'admin@example.com', 'active', 'admin', 'Admin', 'User', 'System Administrator', 'https://example.com/admin.png'),
('user1', '$2a$12$examplehashedpassword', 'user1@example.com', 'active', 'user', 'John', 'Doe', 'A regular user.', 'https://example.com/user1.png');
