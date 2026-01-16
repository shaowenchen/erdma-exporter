# ERDMA Exporter

Prometheus exporter for Aliyun ERDMA.

## 功能

- ERDMA 驱动版本信息
- 设备发现和统计
- 连接管理、命令队列、Verbs API、硬件传输/接收统计

## 构建

```bash
make vendor
make build
```

或使用 Docker：

```bash
make deps
docker build -t shaowenchen/erdma-exporter:latest .
```

## 部署

### Kubernetes

```bash
kubectl apply -f deploy/daemonset.yaml
kubectl apply -f deploy/servicemonitor.yaml
```

## 使用

```bash
./erdma-exporter
```

访问 `http://localhost:9101/metrics` 查看指标。

命令行参数：
- `-web.listen-address`: 监听地址（默认: `:9101`）
- `-web.telemetry-path`: metrics 路径（默认: `/metrics`）

## 指标

所有指标以 `erdma_` 为前缀，所有计数器类型指标使用 `_total` 后缀。

### 驱动和设备信息

- `erdma_driver_version` (Gauge): 驱动版本信息
  - Labels: `version`, `node`
- `erdma_device_info` (Gauge): 设备信息
  - Labels: `device`, `node_guid`, `node`

### 监听相关指标

- `erdma_listen_create_total`: 监听创建操作总数
- `erdma_listen_ipv6_total`: IPv6 监听操作总数
- `erdma_listen_success_total`: 成功的监听操作总数
- `erdma_listen_failed_total`: 失败的监听操作总数
- `erdma_listen_destroy_total`: 监听销毁操作总数

### 接受连接相关指标

- `erdma_accept_total`: 接受操作总数
- `erdma_accept_success_total`: 成功的接受操作总数
- `erdma_accept_failed_total`: 失败的接受操作总数

### 拒绝连接相关指标

- `erdma_reject_total`: 拒绝操作总数
- `erdma_reject_failed_total`: 失败的拒绝操作总数

### 连接相关指标

- `erdma_connect_total`: 连接操作总数
- `erdma_connect_success_total`: 成功的连接操作总数
- `erdma_connect_failed_total`: 失败的连接操作总数
- `erdma_connect_timeout_total`: 连接超时操作总数
- `erdma_connect_reset_total`: 连接重置操作总数

### 命令队列相关指标

- `erdma_cmdq_submitted_total`: 提交的命令队列操作总数
- `erdma_cmdq_completed_total`: 完成的命令队列操作总数
- `erdma_cmdq_eq_notify_total`: 命令队列事件队列通知总数
- `erdma_cmdq_eq_event_total`: 命令队列事件队列事件总数
- `erdma_cmdq_cq_armed_total`: 命令队列完成队列武装操作总数

### 异步事件队列相关指标

- `erdma_aeq_event_total`: 异步事件队列事件总数
- `erdma_aeq_notify_total`: 异步事件队列通知总数

### Verbs API 相关指标

- `erdma_verbs_alloc_mr_total`: Verbs 内存区域分配总数
- `erdma_verbs_alloc_mr_failed_total`: 失败的 Verbs 内存区域分配总数
- `erdma_verbs_alloc_pd_total`: Verbs 保护域分配总数
- `erdma_verbs_alloc_pd_failed_total`: 失败的 Verbs 保护域分配总数
- `erdma_verbs_alloc_uctx_total`: Verbs 用户上下文分配总数
- `erdma_verbs_alloc_uctx_failed_total`: 失败的 Verbs 用户上下文分配总数
- `erdma_verbs_create_cq_total`: Verbs 完成队列创建总数
- `erdma_verbs_create_cq_failed_total`: 失败的 Verbs 完成队列创建总数
- `erdma_verbs_create_qp_total`: Verbs 队列对创建总数
- `erdma_verbs_create_qp_failed_total`: 失败的 Verbs 队列对创建总数
- `erdma_verbs_dealloc_pd_total`: Verbs 保护域释放总数
- `erdma_verbs_dealloc_uctx_total`: Verbs 用户上下文释放总数
- `erdma_verbs_dereg_mr_total`: Verbs 内存区域注销总数
- `erdma_verbs_dereg_mr_failed_total`: 失败的 Verbs 内存区域注销总数
- `erdma_verbs_destroy_cq_total`: Verbs 完成队列销毁总数
- `erdma_verbs_destroy_cq_failed_total`: 失败的 Verbs 完成队列销毁总数
- `erdma_verbs_destroy_qp_total`: Verbs 队列对销毁总数
- `erdma_verbs_destroy_qp_failed_total`: 失败的 Verbs 队列对销毁总数
- `erdma_verbs_get_dma_mr_total`: Verbs DMA 内存区域获取总数
- `erdma_verbs_get_dma_mr_failed_total`: 失败的 Verbs DMA 内存区域获取总数
- `erdma_verbs_reg_usr_mr_total`: Verbs 用户内存区域注册总数
- `erdma_verbs_reg_usr_mr_failed_total`: 失败的 Verbs 用户内存区域注册总数

### 硬件传输相关指标

- `erdma_hw_tx_requests_total`: 硬件发送请求总数
- `erdma_hw_tx_packets_total`: 硬件发送数据包总数
- `erdma_hw_tx_bytes_total`: 硬件发送字节总数
- `erdma_hw_disable_drop_total`: 硬件禁用丢弃总数
- `erdma_hw_bps_limit_drop_total`: 硬件 BPS 限速丢弃总数
- `erdma_hw_pps_limit_drop_total`: 硬件 PPS 限速丢弃总数

### 硬件接收相关指标

- `erdma_hw_rx_packets_total`: 硬件接收数据包总数
- `erdma_hw_rx_bytes_total`: 硬件接收字节总数
- `erdma_hw_rx_disable_drop_total`: 硬件接收禁用丢弃总数
- `erdma_hw_rx_bps_limit_drop_total`: 硬件接收 BPS 限速丢弃总数
- `erdma_hw_rx_pps_limit_drop_total`: 硬件接收 PPS 限速丢弃总数

### 标签说明

所有指标都包含 `node` 标签（节点名称），设备相关指标还包含 `device` 标签（设备名称）。

### 使用示例

查询发送速率：
```promql
rate(erdma_hw_tx_bytes_total{device="erdma_0"}[30s])
```

查询连接失败率：
```promql
rate(erdma_connect_failed_total[30s]) / rate(erdma_connect_total[30s])
```

## Grafana Dashboard

导入 `grafana/dashboard.json` 到 Grafana。

Dashboard 变量：
- `datasource`: 数据源（默认: ops）
- `cluster`: 集群选择
- `node`: 节点选择

## License

MIT

