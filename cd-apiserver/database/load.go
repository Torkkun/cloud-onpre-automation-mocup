package database

import (
	"context"
	"fmt"
)

func (etcd *Client) GetPlaybookTemplate(name string) ([]byte, error) {
	prefix := etcdPlaybookTemplatePrefix + name
	resp, err := etcd.client.Get(context.Background(), prefix)
	if err != nil {
		return nil, fmt.Errorf("get playbook/template function failed: %v", err)
	}
	kv := resp.Kvs[0]
	if kv == nil {
		return nil, fmt.Errorf("nil request, not value")
	}
	return kv.Value, nil
}

func (etcd *Client) GetHostIP(hostname string) (string, error) {
	prefix := etcdIpPrefix + hostname
	resp, err := etcd.client.Get(context.Background(), prefix)
	if err != nil {
		return "", fmt.Errorf("get playbook/template function failed: %v", err)
	}
	kv := resp.Kvs[0]
	if kv == nil {
		return "", fmt.Errorf("nil request, not value")
	}
	return string(kv.Value), nil
}
