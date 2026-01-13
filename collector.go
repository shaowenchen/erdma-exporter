package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// ErdmaCollector collects ERDMA metrics
type ErdmaCollector struct {
	// Version info
	versionDesc *prometheus.Desc

	// Device info
	deviceGUIDDesc *prometheus.Desc

	// Statistics metrics
	listenCreateCntDesc         *prometheus.Desc
	listenIpv6CntDesc           *prometheus.Desc
	listenSuccessCntDesc        *prometheus.Desc
	listenFailedCntDesc         *prometheus.Desc
	listenDestroyCntDesc        *prometheus.Desc
	acceptTotalCntDesc          *prometheus.Desc
	acceptSuccessCntDesc        *prometheus.Desc
	acceptFailedCntDesc         *prometheus.Desc
	rejectCntDesc               *prometheus.Desc
	rejectFailedCntDesc         *prometheus.Desc
	connectTotalCntDesc         *prometheus.Desc
	connectSuccessCntDesc       *prometheus.Desc
	connectFailedCntDesc        *prometheus.Desc
	connectTimeoutCntDesc       *prometheus.Desc
	connectResetCntDesc         *prometheus.Desc
	cmdqSubmittedCntDesc        *prometheus.Desc
	cmdqCompCntDesc             *prometheus.Desc
	cmdqEqNotifyCntDesc         *prometheus.Desc
	cmdqEqEventCntDesc          *prometheus.Desc
	cmdqCqArmedCntDesc          *prometheus.Desc
	erdmaAeqEventCntDesc        *prometheus.Desc
	erdmaAeqNotifyCntDesc       *prometheus.Desc
	verbsAllocMrCntDesc         *prometheus.Desc
	verbsAllocMrFailedCntDesc   *prometheus.Desc
	verbsAllocPdCntDesc         *prometheus.Desc
	verbsAllocPdFailedCntDesc   *prometheus.Desc
	verbsAllocUctxCntDesc       *prometheus.Desc
	verbsAllocUctxFailedCntDesc *prometheus.Desc
	verbsCreateCqCntDesc        *prometheus.Desc
	verbsCreateCqFailedCntDesc  *prometheus.Desc
	verbsCreateQpCntDesc        *prometheus.Desc
	verbsCreateQpFailedCntDesc  *prometheus.Desc
	verbsDeallocPdCntDesc       *prometheus.Desc
	verbsDeallocUctxCntDesc     *prometheus.Desc
	verbsDeregMrCntDesc         *prometheus.Desc
	verbsDeregMrFailedCntDesc   *prometheus.Desc
	verbsDestroyCqCntDesc       *prometheus.Desc
	verbsDestroyCqFailedCntDesc *prometheus.Desc
	verbsDestroyQpCntDesc       *prometheus.Desc
	verbsDestroyQpFailedCntDesc *prometheus.Desc
	verbsGetDmaMrCntDesc        *prometheus.Desc
	verbsGetDmaMrFailedCntDesc  *prometheus.Desc
	verbsRegUsrMrCntDesc        *prometheus.Desc
	verbsRegUsrMrFailedCntDesc  *prometheus.Desc
	hwTxReqsCntDesc             *prometheus.Desc
	hwTxPacketsCntDesc          *prometheus.Desc
	hwTxBytesCntDesc            *prometheus.Desc
	hwDisableDropCntDesc        *prometheus.Desc
	hwBpsLimitDropCntDesc       *prometheus.Desc
	hwPpsLimitDropCntDesc       *prometheus.Desc
	hwRxPacketsCntDesc          *prometheus.Desc
	hwRxBytesCntDesc            *prometheus.Desc
	hwRxDisableDropCntDesc      *prometheus.Desc
	hwRxBpsLimitDropCntDesc     *prometheus.Desc
	hwRxPpsLimitDropCntDesc     *prometheus.Desc
}

// NewErdmaCollector creates a new ERDMA collector
func NewErdmaCollector() (*ErdmaCollector, error) {
	return &ErdmaCollector{
		versionDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "driver", "version"),
			"ERDMA kernel driver version",
			[]string{"version", "node"},
			nil,
		),
		deviceGUIDDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "device", "info"),
			"ERDMA device information",
			[]string{"device", "node_guid", "node"},
			nil,
		),
		listenCreateCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "listen", "create_total"),
			"Total number of listen create operations",
			[]string{"device", "node"},
			nil,
		),
		listenIpv6CntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "listen", "ipv6_total"),
			"Total number of IPv6 listen operations",
			[]string{"device", "node"},
			nil,
		),
		listenSuccessCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "listen", "success_total"),
			"Total number of successful listen operations",
			[]string{"device", "node"},
			nil,
		),
		listenFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "listen", "failed_total"),
			"Total number of failed listen operations",
			[]string{"device", "node"},
			nil,
		),
		listenDestroyCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "listen", "destroy_total"),
			"Total number of listen destroy operations",
			[]string{"device", "node"},
			nil,
		),
		acceptTotalCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "accept", "total"),
			"Total number of accept operations",
			[]string{"device", "node"},
			nil,
		),
		acceptSuccessCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "accept", "success_total"),
			"Total number of successful accept operations",
			[]string{"device", "node"},
			nil,
		),
		acceptFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "accept", "failed_total"),
			"Total number of failed accept operations",
			[]string{"device", "node"},
			nil,
		),
		rejectCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "reject", "total"),
			"Total number of reject operations",
			[]string{"device", "node"},
			nil,
		),
		rejectFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "reject", "failed_total"),
			"Total number of failed reject operations",
			[]string{"device", "node"},
			nil,
		),
		connectTotalCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "connect", "total"),
			"Total number of connect operations",
			[]string{"device", "node"},
			nil,
		),
		connectSuccessCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "connect", "success_total"),
			"Total number of successful connect operations",
			[]string{"device", "node"},
			nil,
		),
		connectFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "connect", "failed_total"),
			"Total number of failed connect operations",
			[]string{"device", "node"},
			nil,
		),
		connectTimeoutCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "connect", "timeout_total"),
			"Total number of connect timeout operations",
			[]string{"device", "node"},
			nil,
		),
		connectResetCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "connect", "reset_total"),
			"Total number of connect reset operations",
			[]string{"device", "node"},
			nil,
		),
		cmdqSubmittedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cmdq", "submitted_total"),
			"Total number of submitted command queue operations",
			[]string{"device", "node"},
			nil,
		),
		cmdqCompCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cmdq", "completed_total"),
			"Total number of completed command queue operations",
			[]string{"device", "node"},
			nil,
		),
		cmdqEqNotifyCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cmdq", "eq_notify_total"),
			"Total number of command queue event queue notifications",
			[]string{"device", "node"},
			nil,
		),
		cmdqEqEventCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cmdq", "eq_event_total"),
			"Total number of command queue event queue events",
			[]string{"device", "node"},
			nil,
		),
		cmdqCqArmedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cmdq", "cq_armed_total"),
			"Total number of command queue completion queue armed operations",
			[]string{"device", "node"},
			nil,
		),
		erdmaAeqEventCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "aeq", "event_total"),
			"Total number of async event queue events",
			[]string{"device", "node"},
			nil,
		),
		erdmaAeqNotifyCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "aeq", "notify_total"),
			"Total number of async event queue notifications",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocMrCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_mr_total"),
			"Total number of verbs memory region allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocMrFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_mr_failed_total"),
			"Total number of failed verbs memory region allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocPdCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_pd_total"),
			"Total number of verbs protection domain allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocPdFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_pd_failed_total"),
			"Total number of failed verbs protection domain allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocUctxCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_uctx_total"),
			"Total number of verbs user context allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsAllocUctxFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "alloc_uctx_failed_total"),
			"Total number of failed verbs user context allocations",
			[]string{"device", "node"},
			nil,
		),
		verbsCreateCqCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "create_cq_total"),
			"Total number of verbs completion queue creations",
			[]string{"device", "node"},
			nil,
		),
		verbsCreateCqFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "create_cq_failed_total"),
			"Total number of failed verbs completion queue creations",
			[]string{"device", "node"},
			nil,
		),
		verbsCreateQpCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "create_qp_total"),
			"Total number of verbs queue pair creations",
			[]string{"device", "node"},
			nil,
		),
		verbsCreateQpFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "create_qp_failed_total"),
			"Total number of failed verbs queue pair creations",
			[]string{"device", "node"},
			nil,
		),
		verbsDeallocPdCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "dealloc_pd_total"),
			"Total number of verbs protection domain deallocations",
			[]string{"device", "node"},
			nil,
		),
		verbsDeallocUctxCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "dealloc_uctx_total"),
			"Total number of verbs user context deallocations",
			[]string{"device", "node"},
			nil,
		),
		verbsDeregMrCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "dereg_mr_total"),
			"Total number of verbs memory region deregistrations",
			[]string{"device", "node"},
			nil,
		),
		verbsDeregMrFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "dereg_mr_failed_total"),
			"Total number of failed verbs memory region deregistrations",
			[]string{"device", "node"},
			nil,
		),
		verbsDestroyCqCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "destroy_cq_total"),
			"Total number of verbs completion queue destructions",
			[]string{"device", "node"},
			nil,
		),
		verbsDestroyCqFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "destroy_cq_failed_total"),
			"Total number of failed verbs completion queue destructions",
			[]string{"device", "node"},
			nil,
		),
		verbsDestroyQpCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "destroy_qp_total"),
			"Total number of verbs queue pair destructions",
			[]string{"device", "node"},
			nil,
		),
		verbsDestroyQpFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "destroy_qp_failed_total"),
			"Total number of failed verbs queue pair destructions",
			[]string{"device", "node"},
			nil,
		),
		verbsGetDmaMrCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "get_dma_mr_total"),
			"Total number of verbs DMA memory region get operations",
			[]string{"device", "node"},
			nil,
		),
		verbsGetDmaMrFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "get_dma_mr_failed_total"),
			"Total number of failed verbs DMA memory region get operations",
			[]string{"device", "node"},
			nil,
		),
		verbsRegUsrMrCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "reg_usr_mr_total"),
			"Total number of verbs user memory region registrations",
			[]string{"device", "node"},
			nil,
		),
		verbsRegUsrMrFailedCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "verbs", "reg_usr_mr_failed_total"),
			"Total number of failed verbs user memory region registrations",
			[]string{"device", "node"},
			nil,
		),
		hwTxReqsCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "tx_requests_total"),
			"Total number of hardware transmit requests",
			[]string{"device", "node"},
			nil,
		),
		hwTxPacketsCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "tx_packets_total"),
			"Total number of hardware transmit packets",
			[]string{"device", "node"},
			nil,
		),
		hwTxBytesCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "tx_bytes_total"),
			"Total number of hardware transmit bytes",
			[]string{"device", "node"},
			nil,
		),
		hwDisableDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "disable_drop_total"),
			"Total number of hardware disable drop operations",
			[]string{"device", "node"},
			nil,
		),
		hwBpsLimitDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "bps_limit_drop_total"),
			"Total number of hardware BPS limit drops",
			[]string{"device", "node"},
			nil,
		),
		hwPpsLimitDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "pps_limit_drop_total"),
			"Total number of hardware PPS limit drops",
			[]string{"device", "node"},
			nil,
		),
		hwRxPacketsCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "rx_packets_total"),
			"Total number of hardware receive packets",
			[]string{"device", "node"},
			nil,
		),
		hwRxBytesCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "rx_bytes_total"),
			"Total number of hardware receive bytes",
			[]string{"device", "node"},
			nil,
		),
		hwRxDisableDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "rx_disable_drop_total"),
			"Total number of hardware receive disable drops",
			[]string{"device", "node"},
			nil,
		),
		hwRxBpsLimitDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "rx_bps_limit_drop_total"),
			"Total number of hardware receive BPS limit drops",
			[]string{"device", "node"},
			nil,
		),
		hwRxPpsLimitDropCntDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "hw", "rx_pps_limit_drop_total"),
			"Total number of hardware receive PPS limit drops",
			[]string{"device", "node"},
			nil,
		),
	}, nil
}

// Describe implements prometheus.Collector
func (c *ErdmaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.versionDesc
	ch <- c.deviceGUIDDesc
	ch <- c.listenCreateCntDesc
	ch <- c.listenIpv6CntDesc
	ch <- c.listenSuccessCntDesc
	ch <- c.listenFailedCntDesc
	ch <- c.listenDestroyCntDesc
	ch <- c.acceptTotalCntDesc
	ch <- c.acceptSuccessCntDesc
	ch <- c.acceptFailedCntDesc
	ch <- c.rejectCntDesc
	ch <- c.rejectFailedCntDesc
	ch <- c.connectTotalCntDesc
	ch <- c.connectSuccessCntDesc
	ch <- c.connectFailedCntDesc
	ch <- c.connectTimeoutCntDesc
	ch <- c.connectResetCntDesc
	ch <- c.cmdqSubmittedCntDesc
	ch <- c.cmdqCompCntDesc
	ch <- c.cmdqEqNotifyCntDesc
	ch <- c.cmdqEqEventCntDesc
	ch <- c.cmdqCqArmedCntDesc
	ch <- c.erdmaAeqEventCntDesc
	ch <- c.erdmaAeqNotifyCntDesc
	ch <- c.verbsAllocMrCntDesc
	ch <- c.verbsAllocMrFailedCntDesc
	ch <- c.verbsAllocPdCntDesc
	ch <- c.verbsAllocPdFailedCntDesc
	ch <- c.verbsAllocUctxCntDesc
	ch <- c.verbsAllocUctxFailedCntDesc
	ch <- c.verbsCreateCqCntDesc
	ch <- c.verbsCreateCqFailedCntDesc
	ch <- c.verbsCreateQpCntDesc
	ch <- c.verbsCreateQpFailedCntDesc
	ch <- c.verbsDeallocPdCntDesc
	ch <- c.verbsDeallocUctxCntDesc
	ch <- c.verbsDeregMrCntDesc
	ch <- c.verbsDeregMrFailedCntDesc
	ch <- c.verbsDestroyCqCntDesc
	ch <- c.verbsDestroyCqFailedCntDesc
	ch <- c.verbsDestroyQpCntDesc
	ch <- c.verbsDestroyQpFailedCntDesc
	ch <- c.verbsGetDmaMrCntDesc
	ch <- c.verbsGetDmaMrFailedCntDesc
	ch <- c.verbsRegUsrMrCntDesc
	ch <- c.verbsRegUsrMrFailedCntDesc
	ch <- c.hwTxReqsCntDesc
	ch <- c.hwTxPacketsCntDesc
	ch <- c.hwTxBytesCntDesc
	ch <- c.hwDisableDropCntDesc
	ch <- c.hwBpsLimitDropCntDesc
	ch <- c.hwPpsLimitDropCntDesc
	ch <- c.hwRxPacketsCntDesc
	ch <- c.hwRxBytesCntDesc
	ch <- c.hwRxDisableDropCntDesc
	ch <- c.hwRxBpsLimitDropCntDesc
	ch <- c.hwRxPpsLimitDropCntDesc
}

// Collect implements prometheus.Collector
func (c *ErdmaCollector) Collect(ch chan<- prometheus.Metric) {
	// Get node name
	nodeName := getNodeName()

	// Get version
	version, err := getVersion()
	if err == nil {
		ch <- prometheus.MustNewConstMetric(
			c.versionDesc,
			prometheus.GaugeValue,
			1.0,
			version,
			nodeName,
		)
	}

	// Get devices
	devices, err := getDevices()
	if err != nil {
		return
	}

	// Collect metrics for each device
	for _, device := range devices {
		// Device info
		ch <- prometheus.MustNewConstMetric(
			c.deviceGUIDDesc,
			prometheus.GaugeValue,
			1.0,
			device.Name,
			device.GUID,
			nodeName,
		)

		// Get statistics for this device
		stats, err := getDeviceStats(device.Name)
		if err != nil {
			continue
		}

		// Emit all statistics
		c.emitStats(ch, device.Name, nodeName, stats)
	}
}

func (c *ErdmaCollector) emitStats(ch chan<- prometheus.Metric, device string, nodeName string, stats map[string]uint64) {
	emitMetric := func(desc *prometheus.Desc, key string) {
		if val, ok := stats[key]; ok {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(val), device, nodeName)
		}
	}

	emitMetric(c.listenCreateCntDesc, "listen_create_cnt")
	emitMetric(c.listenIpv6CntDesc, "listen_ipv6_cnt")
	emitMetric(c.listenSuccessCntDesc, "listen_success_cnt")
	emitMetric(c.listenFailedCntDesc, "listen_failed_cnt")
	emitMetric(c.listenDestroyCntDesc, "listen_destroy_cnt")
	emitMetric(c.acceptTotalCntDesc, "accept_total_cnt")
	emitMetric(c.acceptSuccessCntDesc, "accept_success_cnt")
	emitMetric(c.acceptFailedCntDesc, "accept_failed_cnt")
	emitMetric(c.rejectCntDesc, "reject_cnt")
	emitMetric(c.rejectFailedCntDesc, "reject_failed_cnt")
	emitMetric(c.connectTotalCntDesc, "connect_total_cnt")
	emitMetric(c.connectSuccessCntDesc, "connect_success_cnt")
	emitMetric(c.connectFailedCntDesc, "connect_failed_cnt")
	emitMetric(c.connectTimeoutCntDesc, "connect_timeout_cnt")
	emitMetric(c.connectResetCntDesc, "connect_reset_cnt")
	emitMetric(c.cmdqSubmittedCntDesc, "cmdq_submitted_cnt")
	emitMetric(c.cmdqCompCntDesc, "cmdq_comp_cnt")
	emitMetric(c.cmdqEqNotifyCntDesc, "cmdq_eq_notify_cnt")
	emitMetric(c.cmdqEqEventCntDesc, "cmdq_eq_event_cnt")
	emitMetric(c.cmdqCqArmedCntDesc, "cmdq_cq_armed_cnt")
	emitMetric(c.erdmaAeqEventCntDesc, "erdma_aeq_event_cnt")
	emitMetric(c.erdmaAeqNotifyCntDesc, "erdma_aeq_notify_cnt")
	emitMetric(c.verbsAllocMrCntDesc, "verbs_alloc_mr_cnt")
	emitMetric(c.verbsAllocMrFailedCntDesc, "verbs_alloc_mr_failed_cnt")
	emitMetric(c.verbsAllocPdCntDesc, "verbs_alloc_pd_cnt")
	emitMetric(c.verbsAllocPdFailedCntDesc, "verbs_alloc_pd_failed_cnt")
	emitMetric(c.verbsAllocUctxCntDesc, "verbs_alloc_uctx_cnt")
	emitMetric(c.verbsAllocUctxFailedCntDesc, "verbs_alloc_uctx_failed_cnt")
	emitMetric(c.verbsCreateCqCntDesc, "verbs_create_cq_cnt")
	emitMetric(c.verbsCreateCqFailedCntDesc, "verbs_create_cq_failed_cnt")
	emitMetric(c.verbsCreateQpCntDesc, "verbs_create_qp_cnt")
	emitMetric(c.verbsCreateQpFailedCntDesc, "verbs_create_qp_failed_cnt")
	emitMetric(c.verbsDeallocPdCntDesc, "verbs_dealloc_pd_cnt")
	emitMetric(c.verbsDeallocUctxCntDesc, "verbs_dealloc_uctx_cnt")
	emitMetric(c.verbsDeregMrCntDesc, "verbs_dereg_mr_cnt")
	emitMetric(c.verbsDeregMrFailedCntDesc, "verbs_dereg_mr_failed_cnt")
	emitMetric(c.verbsDestroyCqCntDesc, "verbs_destroy_cq_cnt")
	emitMetric(c.verbsDestroyCqFailedCntDesc, "verbs_destroy_cq_failed_cnt")
	emitMetric(c.verbsDestroyQpCntDesc, "verbs_destroy_qp_cnt")
	emitMetric(c.verbsDestroyQpFailedCntDesc, "verbs_destroy_qp_failed_cnt")
	emitMetric(c.verbsGetDmaMrCntDesc, "verbs_get_dma_mr_cnt")
	emitMetric(c.verbsGetDmaMrFailedCntDesc, "verbs_get_dma_mr_failed_cnt")
	emitMetric(c.verbsRegUsrMrCntDesc, "verbs_reg_usr_mr_cnt")
	emitMetric(c.verbsRegUsrMrFailedCntDesc, "verbs_reg_usr_mr_failed_cnt")
	emitMetric(c.hwTxReqsCntDesc, "hw_tx_reqs_cnt")
	emitMetric(c.hwTxPacketsCntDesc, "hw_tx_packets_cnt")
	emitMetric(c.hwTxBytesCntDesc, "hw_tx_bytes_cnt")
	emitMetric(c.hwDisableDropCntDesc, "hw_disable_drop_cnt")
	emitMetric(c.hwBpsLimitDropCntDesc, "hw_bps_limit_drop_cnt")
	emitMetric(c.hwPpsLimitDropCntDesc, "hw_pps_limit_drop_cnt")
	emitMetric(c.hwRxPacketsCntDesc, "hw_rx_packets_cnt")
	emitMetric(c.hwRxBytesCntDesc, "hw_rx_bytes_cnt")
	emitMetric(c.hwRxDisableDropCntDesc, "hw_rx_disable_drop_cnt")
	emitMetric(c.hwRxBpsLimitDropCntDesc, "hw_rx_bps_limit_drop_cnt")
	emitMetric(c.hwRxPpsLimitDropCntDesc, "hw_rx_pps_limit_drop_cnt")
}

// Device represents an ERDMA device
type Device struct {
	Name string
	GUID string
}

// findCommand finds a command in PATH (container has erdma-tools installed)
func findCommand(cmdName string) string {
	// Try standard PATH (container has erdma-tools installed)
	if path, err := exec.LookPath(cmdName); err == nil {
		return path
	}
	// Fallback to command name (will use PATH)
	return cmdName
}

// getVersion gets the ERDMA driver version
func getVersion() (string, error) {
	eadmPath := findCommand("eadm")
	log.Printf("Debug: Using eadm path: %s", eadmPath)
	
	cmd := exec.Command(eadmPath, "ver")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		log.Printf("Debug: eadm ver command failed: %v, stderr: %s", err, stderr.String())
		return "", fmt.Errorf("failed to execute eadm ver: %w, stderr: %s", err, stderr.String())
	}
	
	output := stdout.Bytes()
	if stderr.Len() > 0 {
		log.Printf("Debug: eadm ver stderr: %s", stderr.String())
	}

	// Parse output: "Query kernel driver version: 0.2.38"
	re := regexp.MustCompile(`Query kernel driver version:\s+(\S+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("failed to parse version from output: %s", string(output))
}

// getDevices gets the list of ERDMA devices
func getDevices() ([]Device, error) {
	ibvDevicesPath := findCommand("ibv_devices")
	log.Printf("Debug: Using ibv_devices path: %s", ibvDevicesPath)
	
	// Check if command exists and is executable
	if _, err := os.Stat(ibvDevicesPath); err != nil {
		log.Printf("Debug: Command file check failed: %v", err)
	}
	
	cmd := exec.Command(ibvDevicesPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		log.Printf("Debug: ibv_devices command failed with error: %v", err)
		log.Printf("Debug: stderr output: %s", stderr.String())
		return nil, fmt.Errorf("failed to execute ibv_devices command: %w, stderr: %s", err, stderr.String())
	}

	output := stdout.Bytes()
	// Debug: log raw output
	log.Printf("Debug: ibv_devices raw output (stdout):\n%s", string(output))
	if stderr.Len() > 0 {
		log.Printf("Debug: ibv_devices stderr output:\n%s", stderr.String())
	}

	var devices []Device
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Skip header lines
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip header lines (first 2 lines)
		if lineNum <= 2 {
			log.Printf("Debug: Skipping header line %d: %s", lineNum, line)
			continue
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse device line: "erdma_0             02163efffe5030b3"
		fields := strings.Fields(line)
		log.Printf("Debug: Parsing line %d: %s (fields: %d)", lineNum, line, len(fields))
		if len(fields) >= 2 {
			devices = append(devices, Device{
				Name: fields[0],
				GUID: fields[1],
			})
			log.Printf("Debug: Found device: %s (GUID: %s)", fields[0], fields[1])
		} else {
			log.Printf("Debug: Line %d does not match device format (expected 2+ fields, got %d)", lineNum, len(fields))
		}
	}

	if err := scanner.Err(); err != nil {
		return devices, fmt.Errorf("error scanning ibv_devices output: %w", err)
	}

	log.Printf("Debug: Total devices found: %d", len(devices))
	return devices, nil
}

// getDeviceStats gets statistics for a specific device
func getDeviceStats(device string) (map[string]uint64, error) {
	eadmPath := findCommand("eadm")
	log.Printf("Debug: Getting stats for device %s using eadm path: %s", device, eadmPath)
	
	cmd := exec.Command(eadmPath, "stat", "-d", device)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		log.Printf("Debug: eadm stat command failed for device %s: %v, stderr: %s", device, err, stderr.String())
		return nil, fmt.Errorf("failed to execute eadm stat: %w, stderr: %s", err, stderr.String())
	}
	
	output := stdout.Bytes()
	if stderr.Len() > 0 {
		log.Printf("Debug: eadm stat stderr for device %s: %s", device, stderr.String())
	}

	stats := make(map[string]uint64)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse line: "listen_create_cnt : 0"
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])

		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			continue
		}

		stats[key] = val
	}

	return stats, scanner.Err()
}

// getNodeName gets the node name from environment variable or hostname
func getNodeName() string {
	// Try to get from environment variable first (Kubernetes)
	if nodeName := os.Getenv("NODE_NAME"); nodeName != "" {
		return nodeName
	}

	// Fallback to hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
