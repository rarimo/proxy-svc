package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"time"
)

type Proxy struct {
	ApiKey string `fig:"api_key,required"`
	Url    string `fig:"url,required"`
	ListId string `fig:"list_id,required"`
	Client *http.Client
}

func (c *config) Proxy() *Proxy {
	return c.proxy.Do(func() interface{} {
		var config Proxy

		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(c.getter, "proxy")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure proxy"))
		}

		config.Client = &http.Client{
			Timeout: time.Minute * 2,
		}

		return &config
	}).(*Proxy)
}
