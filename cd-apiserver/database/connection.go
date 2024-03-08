package database

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdHost = "localhost"
	etcdPort = 2379
)

func SettingConfigForEtcd() clientv3.Config {
	// default
	config := clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)},
		DialTimeout: 5 * time.Second,
	}
	return config
}

type Client struct {
	client *clientv3.Client
}

func NewClientForEtcd(config clientv3.Config) (*Client, error) {
	client, err := clientv3.New(config)
	if err != nil {

		return nil, err
	}
	return &Client{client: client}, nil
}
