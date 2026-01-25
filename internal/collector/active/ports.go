package active

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spf13/viper"
)

type PortCollector struct{}

func init() {
	collector.Register(&PortCollector{})
}

func (c *PortCollector) Name() string {
	return "ports"
}

func (c *PortCollector) Description() string {
	return "Active TCP port scanner for common services"
}

func (c *PortCollector) IsActive() bool {
	return true
}

func (c *PortCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	mode := viper.GetString("collectors.ports.mode")
	var ports []int

	switch mode {
	case "top-100":
		ports = getTop100Ports()
	case "custom":
		ports = viper.GetIntSlice("collectors.ports.custom_ports")
	default:
		// Default list (common ports)
		ports = []int{20, 21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5432, 5900, 8080, 8443}
	}

	results := make(map[int]string)

	for _, port := range ports {
		// Apply rate limit
		if err := ethics.Wait("ports"); err != nil {
			// If context is cancelled, we should probably stop
			return nil, err
		}

		address := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			results[port] = "open"
			conn.Close()
		}
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, err
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("ports_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "ports",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
			"count":  len(results),
			"mode":   mode,
		},
	}

	return []core.Evidence{evidence}, nil
}

func getTop100Ports() []int {
	return []int{
		20, 21, 22, 23, 25, 53, 80, 81, 88, 110, 111, 113, 119, 135, 137, 138, 139, 143, 161, 179,
		389, 443, 445, 465, 513, 514, 515, 548, 554, 587, 631, 636, 873, 990, 993, 995, 1025, 1026, 1027, 1028,
		1029, 1110, 1433, 1521, 1720, 1723, 1755, 1900, 2000, 2001, 2049, 2121, 2717, 3000, 3128, 3306, 3389, 3690, 3999, 4444,
		4899, 5000, 5009, 5051, 5060, 5101, 5190, 5357, 5432, 5631, 5666, 5800, 5900, 6000, 6001, 6646, 6667, 7000, 7070, 8000,
		8008, 8009, 8080, 8081, 8443, 8888, 9000, 9090, 9100, 9102, 9999, 10000, 27017, 32768, 49152, 49153, 49154, 49155, 50000,
	}
}