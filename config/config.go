// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Jolokiabeat JolokiabeatConfig
}

type JolokiabeatConfig struct {
	Period string
	Hosts  []string
	Proxy  ProxyConfig
}

type ProxyConfig struct {
	URL      string
	Password string
	User     string
}
