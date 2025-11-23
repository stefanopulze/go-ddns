package cloudflare

import (
	"fmt"
	"go-ddns/internal/config"
	"go-ddns/internal/provider"
	"log/slog"
	"net/http"
)

func NewCloudflareClient(cfg config.CloudflareConfig) (provider.Client, error) {
	return &client{
		apiToken: cfg.ApiToken,
		zoneId:   cfg.ZoneId,
	}, nil
}

type client struct {
	apiToken   string
	zoneId     string
	domainName string
	record     *RecordSummary
}

func (c *client) UpdateIp(domain string, ip string) error {
	slog.Info("Updating cloudflare dns record",
		slog.String("domain", domain),
		slog.String("ip", ip))

	if c.record == nil {
		record, err := c.getRecordId(domain)
		if err != nil {
			return err
		}

		slog.Info("Got dns record summary")
		c.record = record
	}

	return c.updateDnsRecord(c.record, ip)
}

func (c *client) prepareClient() {

}

func (c *client) addRequestHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
}
