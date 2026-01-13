package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "erdma"
)

// printInitialInfo prints initial debug information
func printInitialInfo() {
	log.Println("================================================================================")
	log.Println("ERDMA Exporter Initial Information")
	log.Println("================================================================================")

	// Get node name
	nodeName := getNodeName()
	log.Printf("Node Name: %s", nodeName)

	// Get version
	version, err := getVersion()
	if err != nil {
		log.Printf("Failed to get ERDMA driver version: %v", err)
	} else {
		log.Printf("ERDMA Driver Version: %s", version)
	}

	// Get devices
	log.Printf("Attempting to discover ERDMA devices...")

	// Check required directories and files
	checkPaths := []string{
		"/dev/infiniband",
		"/sys/class/infiniband",
		"/sys/bus/pci",
		"/sys/devices",
		"/usr/bin",
		"/usr/sbin",
	}
	log.Printf("Checking required paths:")
	for _, path := range checkPaths {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				log.Printf("  ✓ %s (directory exists)", path)
				// List contents for key directories
				if path == "/sys/class/infiniband" {
					if entries, err := os.ReadDir(path); err == nil {
						if len(entries) == 0 {
							log.Printf("    Warning: %s is empty (no devices found)", path)
						} else {
							log.Printf("    Found %d entries in %s:", len(entries), path)
							for _, entry := range entries {
								log.Printf("      - %s", entry.Name())
							}
						}
					}
				}
			} else {
				log.Printf("  ✓ %s (file exists)", path)
			}
		} else {
			log.Printf("  ✗ %s (not found: %v)", path, err)
		}
	}

	ibvDevicesPath := findCommand("ibv_devices")
	log.Printf("Using ibv_devices command: %s", ibvDevicesPath)

	// Check if command exists and is executable
	if info, err := os.Stat(ibvDevicesPath); err == nil {
		log.Printf("Command file info: mode=%v, size=%d", info.Mode(), info.Size())
		if info.Mode().Perm()&0111 == 0 {
			log.Printf("Warning: Command file is not executable")
		}
	} else {
		log.Printf("Warning: ibv_devices command file check failed: %v", err)
		log.Printf("Checking container paths...")
		containerPaths := []string{
			"/usr/bin/ibv_devices",
			"/usr/sbin/ibv_devices",
		}
		found := false
		for _, path := range containerPaths {
			if info, err := os.Stat(path); err == nil {
				log.Printf("Found ibv_devices at: %s (mode=%v)", path, info.Mode())
				found = true
				break
			}
		}
		if !found {
			log.Printf("Error: ibv_devices command not found in container (erdma-tools should be installed)")
		}
	}

	devices, err := getDevices()
	if err != nil {
		log.Printf("Failed to get ERDMA devices: %v", err)
		log.Println("================================================================================")
		return
	}

	log.Printf("Found %d ERDMA device(s):", len(devices))
	for i, device := range devices {
		log.Printf("  Device %d: %s (GUID: %s)", i+1, device.Name, device.GUID)

		// Get initial statistics for this device
		stats, err := getDeviceStats(device.Name)
		if err != nil {
			log.Printf("    Failed to get statistics: %v", err)
			continue
		}

		log.Printf("    Statistics:")
		log.Printf("      Listen: create=%d, success=%d, failed=%d, destroy=%d",
			getStatValue(stats, "listen_create_cnt"),
			getStatValue(stats, "listen_success_cnt"),
			getStatValue(stats, "listen_failed_cnt"),
			getStatValue(stats, "listen_destroy_cnt"))
		log.Printf("      Accept: total=%d, success=%d, failed=%d",
			getStatValue(stats, "accept_total_cnt"),
			getStatValue(stats, "accept_success_cnt"),
			getStatValue(stats, "accept_failed_cnt"))
		log.Printf("      Connect: total=%d, success=%d, failed=%d, timeout=%d, reset=%d",
			getStatValue(stats, "connect_total_cnt"),
			getStatValue(stats, "connect_success_cnt"),
			getStatValue(stats, "connect_failed_cnt"),
			getStatValue(stats, "connect_timeout_cnt"),
			getStatValue(stats, "connect_reset_cnt"))
		log.Printf("      Hardware TX: requests=%d, packets=%d, bytes=%d",
			getStatValue(stats, "hw_tx_reqs_cnt"),
			getStatValue(stats, "hw_tx_packets_cnt"),
			getStatValue(stats, "hw_tx_bytes_cnt"))
		log.Printf("      Hardware RX: packets=%d, bytes=%d",
			getStatValue(stats, "hw_rx_packets_cnt"),
			getStatValue(stats, "hw_rx_bytes_cnt"))
	}

	log.Println("================================================================================")
}

// getStatValue safely gets a statistic value from the map
func getStatValue(stats map[string]uint64, key string) uint64 {
	if val, ok := stats[key]; ok {
		return val
	}
	return 0
}

var (
	listenAddress = flag.String("web.listen-address", ":9101", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	// Create a new ERDMA collector
	collector, err := NewErdmaCollector()
	if err != nil {
		log.Fatalf("Failed to create ERDMA collector: %v", err)
	}

	// Register the collector
	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	// Add the standard process and Go metrics
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	reg.MustRegister(prometheus.NewGoCollector())

	// Setup HTTP server
	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>ERDMA Exporter</title></head>
			<body>
			<h1>ERDMA Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	// Print initial debug information
	printInitialInfo()

	log.Printf("Starting ERDMA exporter on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
