-- K8s global admin roles via scope-binding-permissions
-- Usage: set subject_type / subject_id before running
-- subject_type allowed: 'user' | 'group'

-- ====== 1) Global scope (cluster="*", namespace="*") ======
USE valyria;
INSERT IGNORE INTO valyria_kubernetes_scope (cluster, namespace, creator, created_at, updated_at)
VALUES ('*', '*', 'system', NOW(3), NOW(3));

SET @scope_id := (
  SELECT id FROM valyria_kubernetes_scope
  WHERE cluster = '*' AND namespace = '*'
  LIMIT 1
);

-- ====== 2) Permissions ======
-- Global read-only (list/get)
INSERT IGNORE INTO valyria_kubernetes_permission
  (resource, action, code, creator, created_at, updated_at)
VALUES
  ('*', 'list', 'kubernetes:*:list', 'system', NOW(3), NOW(3)),
  ('*', 'get',  'kubernetes:*:get',  'system', NOW(3), NOW(3));

-- Global operator (create/update/delete/exec/logs)
INSERT IGNORE INTO valyria_kubernetes_permission
  (resource, action, code, creator, created_at, updated_at)
VALUES
  ('*', 'create', 'kubernetes:*:create', 'system', NOW(3), NOW(3)),
  ('*', 'update', 'kubernetes:*:update', 'system', NOW(3), NOW(3)),
  ('*', 'delete', 'kubernetes:*:delete', 'system', NOW(3), NOW(3)),
  ('*', 'exec',   'kubernetes:*:exec',   'system', NOW(3), NOW(3)),
  ('*', 'logs',   'kubernetes:*:logs',   'system', NOW(3), NOW(3));

-- Global admin (all)
INSERT IGNORE INTO valyria_kubernetes_permission
  (resource, action, code, creator, created_at, updated_at)
VALUES
  ('*', '*', 'kubernetes:*:*', 'system', NOW(3), NOW(3));

-- ====== 3) Bindings & permissions ======
-- Skip binding in this seed. Bind later manually.
