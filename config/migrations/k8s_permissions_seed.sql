-- K8s permissions seed (resource + action)
-- Uses INSERT IGNORE to avoid duplicates

INSERT IGNORE INTO valyria_kubernetes_permission
  (resource, action, code, creator, created_at, updated_at)
VALUES
  -- cluster
  ('cluster', 'list', 'kubernetes:cluster:list', 'system', NOW(3), NOW(3)),
  ('cluster', 'get', 'kubernetes:cluster:get', 'system', NOW(3), NOW(3)),
  ('cluster', 'create', 'kubernetes:cluster:create', 'system', NOW(3), NOW(3)),
  ('cluster', 'update', 'kubernetes:cluster:update', 'system', NOW(3), NOW(3)),
  ('cluster', 'delete', 'kubernetes:cluster:delete', 'system', NOW(3), NOW(3)),

  -- node
  ('node', 'list', 'kubernetes:node:list', 'system', NOW(3), NOW(3)),
  ('node', 'get', 'kubernetes:node:get', 'system', NOW(3), NOW(3)),
  ('node', 'update', 'kubernetes:node:update', 'system', NOW(3), NOW(3)),

  -- namespace
  ('namespace', 'list', 'kubernetes:namespace:list', 'system', NOW(3), NOW(3)),
  ('namespace', 'create', 'kubernetes:namespace:create', 'system', NOW(3), NOW(3)),
  ('namespace', 'update', 'kubernetes:namespace:update', 'system', NOW(3), NOW(3)),
  ('namespace', 'delete', 'kubernetes:namespace:delete', 'system', NOW(3), NOW(3)),

  -- deployment
  ('deployment', 'list', 'kubernetes:deployment:list', 'system', NOW(3), NOW(3)),
  ('deployment', 'get', 'kubernetes:deployment:get', 'system', NOW(3), NOW(3)),
  ('deployment', 'create', 'kubernetes:deployment:create', 'system', NOW(3), NOW(3)),
  ('deployment', 'update', 'kubernetes:deployment:update', 'system', NOW(3), NOW(3)),
  ('deployment', 'delete', 'kubernetes:deployment:delete', 'system', NOW(3), NOW(3)),

  -- pod
  ('pod', 'list', 'kubernetes:pod:list', 'system', NOW(3), NOW(3)),
  ('pod', 'get', 'kubernetes:pod:get', 'system', NOW(3), NOW(3)),
  ('pod', 'exec', 'kubernetes:pod:exec', 'system', NOW(3), NOW(3)),
  ('pod', 'logs', 'kubernetes:pod:logs', 'system', NOW(3), NOW(3)),

  -- configmap
  ('configmap', 'list', 'kubernetes:configmap:list', 'system', NOW(3), NOW(3)),
  ('configmap', 'get', 'kubernetes:configmap:get', 'system', NOW(3), NOW(3)),
  ('configmap', 'create', 'kubernetes:configmap:create', 'system', NOW(3), NOW(3)),
  ('configmap', 'update', 'kubernetes:configmap:update', 'system', NOW(3), NOW(3)),
  ('configmap', 'delete', 'kubernetes:configmap:delete', 'system', NOW(3), NOW(3)),

  -- statefulset
  ('statefulset', 'list', 'kubernetes:statefulset:list', 'system', NOW(3), NOW(3)),
  ('statefulset', 'get', 'kubernetes:statefulset:get', 'system', NOW(3), NOW(3)),
  ('statefulset', 'create', 'kubernetes:statefulset:create', 'system', NOW(3), NOW(3)),
  ('statefulset', 'update', 'kubernetes:statefulset:update', 'system', NOW(3), NOW(3)),
  ('statefulset', 'delete', 'kubernetes:statefulset:delete', 'system', NOW(3), NOW(3)),

  -- daemonset
  ('daemonset', 'list', 'kubernetes:daemonset:list', 'system', NOW(3), NOW(3)),
  ('daemonset', 'get', 'kubernetes:daemonset:get', 'system', NOW(3), NOW(3)),
  ('daemonset', 'create', 'kubernetes:daemonset:create', 'system', NOW(3), NOW(3)),
  ('daemonset', 'update', 'kubernetes:daemonset:update', 'system', NOW(3), NOW(3)),
  ('daemonset', 'delete', 'kubernetes:daemonset:delete', 'system', NOW(3), NOW(3)),

  -- replicaset
  ('replicaset', 'list', 'kubernetes:replicaset:list', 'system', NOW(3), NOW(3)),
  ('replicaset', 'get', 'kubernetes:replicaset:get', 'system', NOW(3), NOW(3)),
  ('replicaset', 'update', 'kubernetes:replicaset:update', 'system', NOW(3), NOW(3)),
  ('replicaset', 'delete', 'kubernetes:replicaset:delete', 'system', NOW(3), NOW(3)),

  -- service
  ('service', 'list', 'kubernetes:service:list', 'system', NOW(3), NOW(3)),
  ('service', 'get', 'kubernetes:service:get', 'system', NOW(3), NOW(3)),
  ('service', 'create', 'kubernetes:service:create', 'system', NOW(3), NOW(3)),
  ('service', 'update', 'kubernetes:service:update', 'system', NOW(3), NOW(3)),
  ('service', 'delete', 'kubernetes:service:delete', 'system', NOW(3), NOW(3)),

  -- secret
  ('secret', 'list', 'kubernetes:secret:list', 'system', NOW(3), NOW(3)),
  ('secret', 'get', 'kubernetes:secret:get', 'system', NOW(3), NOW(3)),
  ('secret', 'create', 'kubernetes:secret:create', 'system', NOW(3), NOW(3)),
  ('secret', 'update', 'kubernetes:secret:update', 'system', NOW(3), NOW(3)),
  ('secret', 'delete', 'kubernetes:secret:delete', 'system', NOW(3), NOW(3)),

  -- role
  ('role', 'list', 'kubernetes:role:list', 'system', NOW(3), NOW(3)),
  ('role', 'get', 'kubernetes:role:get', 'system', NOW(3), NOW(3)),
  ('role', 'create', 'kubernetes:role:create', 'system', NOW(3), NOW(3)),
  ('role', 'update', 'kubernetes:role:update', 'system', NOW(3), NOW(3)),
  ('role', 'delete', 'kubernetes:role:delete', 'system', NOW(3), NOW(3)),

  -- rolebinding
  ('rolebinding', 'list', 'kubernetes:rolebinding:list', 'system', NOW(3), NOW(3)),
  ('rolebinding', 'get', 'kubernetes:rolebinding:get', 'system', NOW(3), NOW(3)),
  ('rolebinding', 'create', 'kubernetes:rolebinding:create', 'system', NOW(3), NOW(3)),
  ('rolebinding', 'update', 'kubernetes:rolebinding:update', 'system', NOW(3), NOW(3)),
  ('rolebinding', 'delete', 'kubernetes:rolebinding:delete', 'system', NOW(3), NOW(3)),

  -- clusterrole
  ('clusterrole', 'list', 'kubernetes:clusterrole:list', 'system', NOW(3), NOW(3)),
  ('clusterrole', 'get', 'kubernetes:clusterrole:get', 'system', NOW(3), NOW(3)),
  ('clusterrole', 'create', 'kubernetes:clusterrole:create', 'system', NOW(3), NOW(3)),
  ('clusterrole', 'update', 'kubernetes:clusterrole:update', 'system', NOW(3), NOW(3)),
  ('clusterrole', 'delete', 'kubernetes:clusterrole:delete', 'system', NOW(3), NOW(3)),

  -- clusterrolebinding
  ('clusterrolebinding', 'list', 'kubernetes:clusterrolebinding:list', 'system', NOW(3), NOW(3)),
  ('clusterrolebinding', 'get', 'kubernetes:clusterrolebinding:get', 'system', NOW(3), NOW(3)),
  ('clusterrolebinding', 'create', 'kubernetes:clusterrolebinding:create', 'system', NOW(3), NOW(3)),
  ('clusterrolebinding', 'update', 'kubernetes:clusterrolebinding:update', 'system', NOW(3), NOW(3)),
  ('clusterrolebinding', 'delete', 'kubernetes:clusterrolebinding:delete', 'system', NOW(3), NOW(3)),

  -- serviceaccount
  ('serviceaccount', 'list', 'kubernetes:serviceaccount:list', 'system', NOW(3), NOW(3)),
  ('serviceaccount', 'get', 'kubernetes:serviceaccount:get', 'system', NOW(3), NOW(3)),
  ('serviceaccount', 'create', 'kubernetes:serviceaccount:create', 'system', NOW(3), NOW(3)),
  ('serviceaccount', 'update', 'kubernetes:serviceaccount:update', 'system', NOW(3), NOW(3)),
  ('serviceaccount', 'delete', 'kubernetes:serviceaccount:delete', 'system', NOW(3), NOW(3)),

  -- ingress
  ('ingress', 'list', 'kubernetes:ingress:list', 'system', NOW(3), NOW(3)),
  ('ingress', 'get', 'kubernetes:ingress:get', 'system', NOW(3), NOW(3)),
  ('ingress', 'create', 'kubernetes:ingress:create', 'system', NOW(3), NOW(3)),
  ('ingress', 'update', 'kubernetes:ingress:update', 'system', NOW(3), NOW(3)),
  ('ingress', 'delete', 'kubernetes:ingress:delete', 'system', NOW(3), NOW(3)),

  -- ingressclass
  ('ingressclass', 'list', 'kubernetes:ingressclass:list', 'system', NOW(3), NOW(3)),
  ('ingressclass', 'get', 'kubernetes:ingressclass:get', 'system', NOW(3), NOW(3)),
  ('ingressclass', 'create', 'kubernetes:ingressclass:create', 'system', NOW(3), NOW(3)),
  ('ingressclass', 'update', 'kubernetes:ingressclass:update', 'system', NOW(3), NOW(3)),
  ('ingressclass', 'delete', 'kubernetes:ingressclass:delete', 'system', NOW(3), NOW(3)),

  -- storageclass
  ('storageclass', 'list', 'kubernetes:storageclass:list', 'system', NOW(3), NOW(3)),
  ('storageclass', 'get', 'kubernetes:storageclass:get', 'system', NOW(3), NOW(3)),
  ('storageclass', 'create', 'kubernetes:storageclass:create', 'system', NOW(3), NOW(3)),
  ('storageclass', 'update', 'kubernetes:storageclass:update', 'system', NOW(3), NOW(3)),
  ('storageclass', 'delete', 'kubernetes:storageclass:delete', 'system', NOW(3), NOW(3)),

  -- persistentvolume
  ('persistentvolume', 'list', 'kubernetes:persistentvolume:list', 'system', NOW(3), NOW(3)),
  ('persistentvolume', 'get', 'kubernetes:persistentvolume:get', 'system', NOW(3), NOW(3)),
  ('persistentvolume', 'create', 'kubernetes:persistentvolume:create', 'system', NOW(3), NOW(3)),
  ('persistentvolume', 'update', 'kubernetes:persistentvolume:update', 'system', NOW(3), NOW(3)),
  ('persistentvolume', 'delete', 'kubernetes:persistentvolume:delete', 'system', NOW(3), NOW(3)),

  -- persistentvolumeclaim
  ('persistentvolumeclaim', 'list', 'kubernetes:persistentvolumeclaim:list', 'system', NOW(3), NOW(3)),
  ('persistentvolumeclaim', 'get', 'kubernetes:persistentvolumeclaim:get', 'system', NOW(3), NOW(3)),
  ('persistentvolumeclaim', 'create', 'kubernetes:persistentvolumeclaim:create', 'system', NOW(3), NOW(3)),
  ('persistentvolumeclaim', 'update', 'kubernetes:persistentvolumeclaim:update', 'system', NOW(3), NOW(3)),
  ('persistentvolumeclaim', 'delete', 'kubernetes:persistentvolumeclaim:delete', 'system', NOW(3), NOW(3)),

  -- event
  ('event', 'list', 'kubernetes:event:list', 'system', NOW(3), NOW(3)),

  -- task
  ('task', 'list', 'kubernetes:task:list', 'system', NOW(3), NOW(3)),
  ('task', 'get', 'kubernetes:task:get', 'system', NOW(3), NOW(3)),
  ('task', 'create', 'kubernetes:task:create', 'system', NOW(3), NOW(3)),
  ('task', 'update', 'kubernetes:task:update', 'system', NOW(3), NOW(3)),
  ('task', 'delete', 'kubernetes:task:delete', 'system', NOW(3), NOW(3)),

  -- taskrun
  ('taskrun', 'list', 'kubernetes:taskrun:list', 'system', NOW(3), NOW(3)),
  ('taskrun', 'get', 'kubernetes:taskrun:get', 'system', NOW(3), NOW(3)),
  ('taskrun', 'create', 'kubernetes:taskrun:create', 'system', NOW(3), NOW(3)),
  ('taskrun', 'update', 'kubernetes:taskrun:update', 'system', NOW(3), NOW(3)),
  ('taskrun', 'delete', 'kubernetes:taskrun:delete', 'system', NOW(3), NOW(3)),

  -- pipeline
  ('pipeline', 'list', 'kubernetes:pipeline:list', 'system', NOW(3), NOW(3)),
  ('pipeline', 'get', 'kubernetes:pipeline:get', 'system', NOW(3), NOW(3)),
  ('pipeline', 'create', 'kubernetes:pipeline:create', 'system', NOW(3), NOW(3)),
  ('pipeline', 'update', 'kubernetes:pipeline:update', 'system', NOW(3), NOW(3)),
  ('pipeline', 'delete', 'kubernetes:pipeline:delete', 'system', NOW(3), NOW(3)),

  -- pipelinerun
  ('pipelinerun', 'list', 'kubernetes:pipelinerun:list', 'system', NOW(3), NOW(3)),
  ('pipelinerun', 'get', 'kubernetes:pipelinerun:get', 'system', NOW(3), NOW(3)),
  ('pipelinerun', 'create', 'kubernetes:pipelinerun:create', 'system', NOW(3), NOW(3)),
  ('pipelinerun', 'update', 'kubernetes:pipelinerun:update', 'system', NOW(3), NOW(3)),
  ('pipelinerun', 'delete', 'kubernetes:pipelinerun:delete', 'system', NOW(3), NOW(3));
