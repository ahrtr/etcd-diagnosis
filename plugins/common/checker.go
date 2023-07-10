package common

import "github.com/ahrtr/etcd-diagnosis/agent"

type Checker struct {
	agent.GlobalConfig
	Name string
}
