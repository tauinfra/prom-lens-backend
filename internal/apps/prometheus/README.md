Prometheus Target Curl Example

Target Group
- Create:
  curl -X POST "http://localhost:8080/api/v1/prometheus/target-groups" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "k8s-nodes",
      "description": "k8s nodes targets",
      "labels": {
        "job": "node-exporter",
        "env": "prod"
      }
    }'

- Update:
  curl -X PATCH "http://localhost:8080/api/v1/prometheus/target-groups/1" \
    -H "Content-Type: application/json" \
    -d '{
      "description": "k8s nodes targets (prod)",
      "labels": {
        "job": "node-exporter",
        "env": "prod",
        "cluster": "prom-lens"
      }
    }'

Target
- Create:
  curl -X POST "http://localhost:8080/api/v1/prometheus/target-groups/1/targets" \
    -H "Content-Type: application/json" \
    -d '{
      "ipAddress": "10.0.0.10",
      "port": 9100,
      "labels": {
        "node": "node-a"
      },
      "enabled": true
    }'

- Update:
  curl -X PATCH "http://localhost:8080/api/v1/prometheus/target-groups/1/targets/2" \
    -H "Content-Type: application/json" \
    -d '{
      "labels": {
        "node": "node-a",
        "zone": "cn-hz-1"
      },
      "enabled": false
    }'

ConfigMap 同步格式（Prometheus file_sd）

每个目标组写入 `{groupName}.json`，内容为 file_sd 数组；组级 `labels` 与节点 `labels` 合并，节点覆盖组：

```json
[
  {
    "targets": ["10.10.10.10:9100"],
    "labels": {
      "job": "node-exporter",
      "hostname": "node01"
    }
  }
]
```
