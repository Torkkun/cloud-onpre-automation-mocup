package database

import (
	"context"
	"fmt"
)

func (etcd *Client) PutHostIP(hostname, address string) error {
	prefix := etcdIpPrefix + hostname
	if _, err := etcd.client.Put(context.Background(), prefix, address); err != nil {
		return fmt.Errorf("putHostIp function failed: %s", err)
	}
	return nil
}
