package client

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"

	"sns/app/post/internal/conf"
)

func NewDiscoveryConsul(cfg *conf.Registry) (*consul.Registry, error) {
	c := api.DefaultConfig()
	c.Address = cfg.Consul.Address
	c.Token = cfg.Consul.Token
	if cfg.Consul.Scheme != "" {
		c.Scheme = cfg.Consul.Scheme
	}
	if c.Scheme == "https" {
		c.TLSConfig.InsecureSkipVerify = true
	}

	cli, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}

	r := consul.New(cli)

	return r, nil
}
