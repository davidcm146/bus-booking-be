INSERT INTO roles (id, name) VALUES
  (1, 'passenger'),
  (2, 'operator'),
  (3, 'admin')
ON CONFLICT (id) DO NOTHING;
