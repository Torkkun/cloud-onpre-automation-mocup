package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDynamicInventory(t *testing.T) {
	LoadDynamicInventoryAllHost()
}
func TestSingleLoad(t *testing.T) {
	config := settingConfigForEtcd()
	etcd, err := newClientForEtcd(config)
	if err != nil {
		log.Printf("failed to create etcd client: %v\n", err)
	}
	defer etcd.client.Close()

	hostvars = map[string]map[string]string{}
	group = []string{}

	if err := etcd.getHostandIP("dtest"); err != nil {
		log.Printf("failed get host and ip function: %v", err)
		return
	}

	if err := etcd.getHostVars("dtest"); err != nil {
		log.Printf("failed gt hostvars: %v", err)
		return
	}
	fmt.Println(hostvars, group)
}

func TestStore(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to create etcd client: %v\n", err)
	}
	defer client.Close()

	key := fmt.Sprintf("/test/ansible/dynamic/testvars/%s", "dtest")

	body := map[string]string{
		"ansible_user": "test",
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		fmt.Println("json marshal err: ", err)
		return
	}
	client.Put(context.Background(), key, string(bytes))
}

func TestLoad(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to create etcd client: %v\n", err)
	}
	defer client.Close()

	key := fmt.Sprintf("/test/ansible/dynamic/testvars/%s", "testvmid2")

	resp, err := client.Get(context.Background(), key)
	if err != nil {
		fmt.Println("failed etcd get: ", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Key))
		fmt.Println(string(kv.Value))
		var hostvar map[string]string
		if err = json.Unmarshal(kv.Value, &hostvar); err != nil {
			fmt.Printf("fatal Unmarshl Value:%v", err)
			return
		}
		fmt.Println(hostvar)
	}

}

// --host {host_name or group_name}
func ForTestLoadDynamicInventoryHostOnly(hostname string) (map[string]map[string]string, []string, error) {
	config := settingConfigForEtcd()
	etcd, err := newClientForEtcd(config)
	if err != nil {
		log.Printf("failed to create etcd client: %v\n", err)
	}
	defer etcd.client.Close()
	hostvars = map[string]map[string]string{}
	group = []string{}

	if err := etcd.TestgetHostandIP(hostname); err != nil {
		return nil, nil, err
	}

	if err := etcd.TestgetHostVars(hostname); err != nil {
		return nil, nil, err
	}
	return hostvars, group, nil
}

// --list
func ForTestLoadDynamicInventoryAllHost() (map[string]map[string]string, []string, error) {
	config := settingConfigForEtcd()
	etcd, err := newClientForEtcd(config)
	if err != nil {
		log.Printf("failed to create etcd client: %v\n", err)
	}
	defer etcd.client.Close()

	hostvars = map[string]map[string]string{}
	group = []string{}

	if err := etcd.TestgetAllHostandIP(); err != nil {
		return nil, nil, err
	}

	if err := etcd.TestgetAllHostVars(); err != nil {
		return nil, nil, err
	}
	return hostvars, group, nil
}

func (etcd *Client) TestgetAllHostandIP() error {
	prefix := "/test/ansible/dynamic/ip/"
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

func (etcd *Client) TestgetHostandIP(hostname string) error {
	prefix := "/test/ansible/dynamic/ip/" + hostname
	ipresp, err := etcd.client.Get(context.Background(), prefix)
	if err != nil {
		return fmt.Errorf("fatal get %s's key values: %v", prefix, err)
	}
	kv := ipresp.Kvs[0]
	if kv == nil {
		return nil
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

func (etcd *Client) TestgetAllHostVars() error {
	prefix := "/test/ansible/dynamic/testvars/"
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

func (etcd *Client) TestgetHostVars(hostname string) error {
	prefix := "/test/ansible/dynamic/testvars/" + hostname
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
