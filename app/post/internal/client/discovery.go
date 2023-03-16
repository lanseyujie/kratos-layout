package client

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"

	"sns/app/post/internal/conf"
)

func NewDiscoveryConsul(cfg *conf.Registry) (*consul.Registry, error) {
	if cfg.Consul == nil || cfg.Consul.Address == "" {
		return nil, nil
	}

	c := api.DefaultConfig()
	c.Address = cfg.Consul.Address

	if cfg.Consul.Scheme != nil && *cfg.Consul.Scheme != "" {
		c.Scheme = *cfg.Consul.Scheme
	}

	if cfg.Consul.Token != nil && *cfg.Consul.Token != "" {
		c.Token = *cfg.Consul.Token
	}

	if c.Scheme == "https" {
		c.TLSConfig.InsecureSkipVerify = true
	}

	cli, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}

	r := consul.New(cli, consul.WithHealthCheck(false))

	return r, nil
}
