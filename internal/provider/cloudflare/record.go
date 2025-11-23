package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ListResponse struct {
	Success bool            `json:"success"`
	Result  []RecordSummary `json:"result"`
}

type RecordSummary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func (c *client) getRecordId(domainName string) (*RecordSummary, error) {
	slog.Info("Getting dns record summary", slog.String("domain", domainName))
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", c.zoneId, domainName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addRequestHeaders(req)

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(resp *http.Response) {
		_ = resp.Body.Close()
	}(resp)

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Error: %d - %s", resp.StatusCode, string(body))
	}

	var listResp ListResponse
	if err = json.Unmarshal(body, &listResp); err != nil {
		return nil, err
	}

	if len(listResp.Result) == 0 {
		return nil, fmt.Errorf("no dns record found for domain %s", domainName)
	}

	return &listResp.Result[0], nil
}
