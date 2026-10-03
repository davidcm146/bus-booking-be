-- +goose Up
-- Roles are reference data the application needs at runtime (signup resolves
-- the passenger role), so they live in a migration rather than in seeds.
INSERT INTO roles (id, name) VALUES
  (1, 'passenger'),
  (2, 'operator'),
  (3, 'admin')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('roles', 'id'), (SELECT max(id) FROM roles));

-- +goose Down
DELETE FROM roles WHERE id IN (1, 2, 3);
