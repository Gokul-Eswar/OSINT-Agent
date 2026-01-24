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
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 993, 995, 3306, 3389, 5432, 8080, 8443}
	
	results := make(map[int]string)

	for _, port := range commonPorts {
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
		},
	}

	return []core.Evidence{evidence}, nil
}
