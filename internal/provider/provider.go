package provider

type Client interface {
	UpdateIp(domain string, ip string) error
}
