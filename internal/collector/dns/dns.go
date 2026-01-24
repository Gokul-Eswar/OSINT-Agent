package dns

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

type DNSCollector struct{}

func init() {
	collector.Register(&DNSCollector{})
}

func (d *DNSCollector) Name() string {
	return "dns"
}

func (d *DNSCollector) Description() string {
	return "Passive DNS lookup for A, MX, and NS records"
}

func (d *DNSCollector) IsActive() bool {
	return false
}

func (d *DNSCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	results := make(map[string][]string)

	// A Records
	ips, _ := net.LookupHost(target)
	results["A"] = ips

	// MX Records
	mxs, _ := net.LookupMX(target)
	for _, mx := range mxs {
		results["MX"] = append(results["MX"], mx.Host)
	}

	// NS Records
	nss, _ := net.LookupNS(target)
	for _, ns := range nss {
		results["NS"] = append(results["NS"], ns.Host)
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

	fileName := fmt.Sprintf("dns_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "dns",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
			"types":  []string{"A", "MX", "NS"},
		},
		RawData: results,
	}

	return []core.Evidence{evidence}, nil
}
