package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type RespHost struct {
	Meta      Meta  `json:"_meta"`
	Ungrouped Group `json:"ungrouped"`
}

type Meta struct {
	HostVars map[string]map[string]string `json:"hostvars"`
}

type Group struct {
	Hosts []string `json:"hosts"`
}

func main() {
	// --list もしくは --host {host}の情報を取得する
	listf, hostf := flagoption()
	// list もしくは hostによってレスポンスを変えている
	// listは全て取得
	// hostはグループもしくはホスト名を指定しマッチする情報を取得。現在はホストのみしか対応していない。
	var resp RespHost
	if listf {
		meta, ungroup, err := LoadDynamicInventoryAllHost()
		if err != nil {
			log.Fatalln(err)
		}
		resp = RespHost{
			Meta: Meta{
				HostVars: meta,
			},
			Ungrouped: Group{
				Hosts: ungroup,
			},
		}
	} else if hostf != "" {
		switch hostf {
		// static test, this is none use etcd.
		case "stest":
			resp = ExampleInventory()
		// dynamic test, this is use etcd.
		case "dtest":
			meta, ungroup, err := LoadDynamicInventoryHostOnly("dtest")
			if err != nil {
				log.Fatalln(err)
			}
			resp = RespHost{
				Meta: Meta{
					HostVars: meta,
				},
				Ungrouped: Group{
					Hosts: ungroup,
				},
			}
		default:
			meta, ungroup, err := LoadDynamicInventoryHostOnly(hostf)
			if err != nil {
				log.Fatalln(err)
			}
			resp = RespHost{
				Meta: Meta{
					HostVars: meta,
				},
				Ungrouped: Group{
					Hosts: ungroup,
				},
			}
		}
		// ホスト引数がない場合、Ansible側でエラーが表示されるはず。
	} else {
		resp.Meta = EmptyInventory()
	}
	inventoryJson, err := json.Marshal(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal inventory to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(inventoryJson))
}
