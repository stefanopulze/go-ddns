package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type DNSRecord struct {
	Type    string `json:"type"`    // A
	Name    string `json:"name"`    // es. "ddns.example.com"
	Content string `json:"content"` // L'indirizzo IP
	TTL     int    `json:"ttl"`     // 1 per automatico
	Proxied bool   `json:"proxied"` // true se vuoi passare attraverso il CDN Cloudflare
}

type CFResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *client) updateDnsRecord(record *RecordSummary, ip string) error {
	slog.Info("Updating zone dns record", slog.String("ip", ip), slog.String("domain", record.Name))
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", c.zoneId, record.ID)

	payload := DNSRecord{
		Type:    record.Type,
		Name:    record.Name,
		Content: ip,
		TTL:     record.TTL,
		Proxied: record.Proxied,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	c.addRequestHeaders(req)

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer func(resp *http.Response) {
		_ = resp.Body.Close()
	}(resp)

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API Error: %d - Body: %s", resp.StatusCode, string(body))
	}

	var cfResp CFResponse
	if err = json.Unmarshal(body, &cfResp); err != nil {
		return fmt.Errorf("errore parsing JSON risposta: %v", err)
	}

	if !cfResp.Success {
		msg := "unknown cloudflare error"
		if len(cfResp.Errors) > 0 {
			msg = cfResp.Errors[0].Message
		}
		return fmt.Errorf("API Error: %s", msg)
	}

	return nil
}
