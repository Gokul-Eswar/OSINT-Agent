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

	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
)

type Resolver interface {
	LookupHost(host string) (addrs []string, err error)
	LookupMX(name string) ([]*net.MX, error)
	LookupNS(name string) ([]*net.NS, error)
}

type NetResolver struct{}

func (r *NetResolver) LookupHost(host string) (addrs []string, err error) { return net.LookupHost(host) }
func (r *NetResolver) LookupMX(name string) ([]*net.MX, error)          { return net.LookupMX(name) }
func (r *NetResolver) LookupNS(name string) ([]*net.NS, error)          { return net.LookupNS(name) }

type DNSCollector struct {
	resolver Resolver
}

func init() {
	collector.Register(&DNSCollector{resolver: &NetResolver{}})
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
	log.Info().
		Str("collector", "dns").
		Str("case_id", caseID).
		Str("target", target).
		Msg("collection_started")

	results := make(map[string][]string)

	// A Records
	ips, err := d.resolver.LookupHost(target)
	if err != nil {
		log.Debug().Err(err).Str("target", target).Msg("failed to lookup A records")
	}
	results["A"] = ips

	// MX Records
	mxs, err := d.resolver.LookupMX(target)
	if err != nil {
		log.Debug().Err(err).Str("target", target).Msg("failed to lookup MX records")
	}
	for _, mx := range mxs {
		results["MX"] = append(results["MX"], mx.Host)
	}

	// NS Records
	nss, err := d.resolver.LookupNS(target)
	if err != nil {
		log.Debug().Err(err).Str("target", target).Msg("failed to lookup NS records")
	}
	for _, ns := range nss {
		results["NS"] = append(results["NS"], ns.Host)
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal DNS results")
		return nil, fmt.Errorf("failed to marshal DNS results: %w", err)
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", storageDir).Msg("failed to create storage directory")
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	fileName := fmt.Sprintf("dns_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("failed to write DNS evidence file")
		return nil, fmt.Errorf("failed to write DNS evidence file %s: %w", filePath, err)
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

	log.Info().
		Str("collector", "dns").
		Str("case_id", caseID).
		Str("target", target).
		Int("record_count", len(ips)+len(mxs)+len(nss)).
		Msg("collection_completed")

	return []core.Evidence{evidence}, nil
}
