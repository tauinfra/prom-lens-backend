CREATE TABLE IF NOT EXISTS valyria_authn_menu (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  parent_id BIGINT DEFAULT NULL COMMENT '父菜单ID',
  name VARCHAR(64) NOT NULL COMMENT '路由名称',
  path VARCHAR(255) NOT NULL COMMENT '路由路径',
  component VARCHAR(255) DEFAULT NULL COMMENT '组件路径',
  title VARCHAR(64) NOT NULL COMMENT '菜单标题',
  icon VARCHAR(64) DEFAULT NULL COMMENT '图标',
  rank INT DEFAULT 0 COMMENT '排序',
  show_parent TINYINT(1) DEFAULT 0 COMMENT '显示父级',
  show_link TINYINT(1) DEFAULT 1 COMMENT '显示菜单',
  creator VARCHAR(64) NOT NULL COMMENT '创建人',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'
);

INSERT INTO valyria_authn_menu
  (parent_id, name, path, component, title, icon, rank, show_parent, show_link, creator, created_at, updated_at)
VALUES
  (NULL, 'Kubernetes', '/kubernetes', NULL, 'Kubernetes', 'IconKubernetes', 1, 0, 1, 'system', NOW(3), NOW(3));

SET @k8s_menu_id = LAST_INSERT_ID();

INSERT INTO valyria_authn_menu
  (parent_id, name, path, component, title, icon, rank, show_parent, show_link, creator, created_at, updated_at)
VALUES
  (@k8s_menu_id, 'Cluster', '/kubernetes/clusters', '@/views/kubernetes/cluster/index.vue', '集群管理', 'IconCluster', 1, 1, 1, 'system', NOW(3), NOW(3)),
  (@k8s_menu_id, 'Node', '/kubernetes/nodes', '@/views/kubernetes/node/index.vue', '节点管理', 'IconServerOutline', 2, 1, 1, 'system', NOW(3), NOW(3));

