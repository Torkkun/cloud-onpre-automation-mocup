package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdHost = "localhost"
	etcdPort = 2379

	etcdIpPrefix       = "/ip/"
	etcdHostVarsPrefix = "/ansible/dynamic/hostvars/"
)

// 後でカスタム可能にできるように
func settingConfigForEtcd() clientv3.Config {
	// 現在はデフォのみ
	config := clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)},
		DialTimeout: 5 * time.Second,
	}
	return config
}

type Client struct {
	client *clientv3.Client
}

func newClientForEtcd(config clientv3.Config) (*Client, error) {
	client, err := clientv3.New(config)
	if err != nil {

		return nil, err
	}
	return &Client{client: client}, nil
}

var hostvars map[string]map[string]string
var group []string

// --host {host_name or group_name}
func LoadDynamicInventoryHostOnly(hostname string) (map[string]map[string]string, []string, error) {
	config := settingConfigForEtcd()
	etcd, err := newClientForEtcd(config)
	if err != nil {
		log.Printf("failed to create etcd client: %v\n", err)
	}
	defer etcd.client.Close()
	hostvars = map[string]map[string]string{}
	group = []string{}

	if err := etcd.getHostandIP(hostname); err != nil {
		return nil, nil, err
	}

	if err := etcd.getHostVars(hostname); err != nil {
		return nil, nil, err
	}
	return hostvars, group, nil
}

// --list
func LoadDynamicInventoryAllHost() (map[string]map[string]string, []string, error) {
	config := settingConfigForEtcd()
	etcd, err := newClientForEtcd(config)
	if err != nil {
		log.Printf("failed to create etcd client: %v\n", err)
	}
	defer etcd.client.Close()

	hostvars = map[string]map[string]string{}
	group = []string{}

	if err := etcd.getAllHostandIP(); err != nil {
		return nil, nil, err
	}

	if err := etcd.getAllHostVars(); err != nil {
		return nil, nil, err
	}
	return hostvars, group, nil
}

// db usecases
func (etcd *Client) getHostandIP(hostname string) error {
	prefix := etcdIpPrefix + hostname
	ipresp, err := etcd.client.Get(context.Background(), prefix)
	if err != nil {
		return fmt.Errorf("fatal get %s's key values: %v", prefix, err)
	}
	kv := ipresp.Kvs[0]
	if kv == nil {
		//  これは必要なのか？
		return fmt.Errorf("nil request, not value")
	}
	_, ok := hostvars[hostname]

	if !ok {
		hostvars[hostname] = map[string]string{"ansible_ssh_host": string(kv.Value)}
		group = append(group, hostname)
	} else {
		log.Fatalln("failed hostvars init ip")
	}
	return nil
}

func (etcd *Client) getAllHostandIP() error {
	prefix := etcdIpPrefix
	ipresp, err := etcd.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		err := fmt.Errorf("fatal get %s's key values: %v", prefix, err)
		return err
	}

	for _, kv := range ipresp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), prefix)
		parts := strings.Split(key, "/")

		// vmid
		hostname := parts[0]
		_, ok := hostvars[hostname]

		if !ok {
			hostvars[hostname] = map[string]string{"ansible_ssh_host": string(kv.Value)}
			group = append(group, hostname)
		} else {
			log.Fatalln("failed hostvars init ip")
		}
	}
	return nil
}

func (etcd *Client) getHostVars(hostname string) error {
	prefix := etcdHostVarsPrefix + hostname
	hostvarsresp, err := etcd.client.Get(context.Background(), prefix)
	if err != nil {
		return fmt.Errorf("fatal get %s's key values: %v", prefix, err)
	}
	kv := hostvarsresp.Kvs[0]
	if kv == nil {
		return nil
	}
	hostvar, ok := hostvars[hostname]

	if !ok {
		log.Fatalf("failed add hostvars")
	}
	// append hostvars
	var newhostvar map[string]string
	if err = json.Unmarshal(kv.Value, &newhostvar); err != nil {
		return fmt.Errorf("fatal Unmarshal hostvar Value:%v", err)
	}
	hostvars[hostname] = merge(hostvar, newhostvar)
	return nil
}

func (etcd *Client) getAllHostVars() error {
	prefix := etcdHostVarsPrefix
	hostvarsresp, err := etcd.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		err := fmt.Errorf("fatal get %s's key values: %v", prefix, err)
		return err
	}
	for _, kv := range hostvarsresp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), prefix)
		parts := strings.Split(key, "/")
		// vmid
		hostname := parts[0]
		hostvar, ok := hostvars[hostname]
		if !ok {
			log.Fatalf("failed add hostvars")
		}
		// append hostvars
		var newhostvar map[string]string
		if err = json.Unmarshal(kv.Value, &newhostvar); err != nil {
			return fmt.Errorf("fatal Unmarshal hostvar Value:%v", err)
		}
		hostvars[hostname] = merge(hostvar, newhostvar)
	}
	return nil
}
