# etcd-diagnosis

## Overview
etcd-diagnosis is a comprehensive tool for diagnosing etcd clusters. It can analyze running clusters and
generate diagnostic reports with a single command. It reuses most of the `etcdctl` global flags, offering a
familiar experience to `etcdctl` users. In addition to live cluster diagnostics, it also supports offline
data analysis, specifically counting the number of revisions for each key. See the complete flags below,

```
$ ./bin/etcd-diagnosis -h
An one-stop etcd diagnosis tool

Usage:
  etcd-diagnosis [flags]

Flags:
      --cacert string                  verify certificates of TLS-enabled secure servers using this CA bundle
      --cert string                    identify secure client using this TLS certificate file
      --cluster                        use all endpoints from the cluster member list
      --command-timeout duration       command timeout (excluding dial timeout) (default 5s)
      --data-dir string                path to data directory
      --dial-timeout duration          dial timeout for client connections (default 2s)
  -d, --discovery-srv string           domain name to query for SRV records describing cluster endpoints
      --discovery-srv-name string      service name to query when using DNS discovery
      --endpoints strings              comma separated etcd endpoints (default [127.0.0.1:2379])
      --etcd-storage-quota-bytes int   etcd storage quota in bytes (the value passed to etcd instance by flag --quota-backend-bytes) (default 2147483648)
  -h, --help                           help for etcd-diagnosis
      --insecure-discovery             accept insecure SRV records describing cluster endpoints (default true)
      --insecure-skip-tls-verify       skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)
      --insecure-transport             disable transport security for client connections (default true)
      --keepalive-time duration        keepalive time for client connections (default 2s)
      --keepalive-timeout duration     keepalive timeout for client connections (default 5s)
      --key string                     identify secure client using this TLS key file
      --offline                        offline analysis
      --password string                password for authentication (if this option is used, --user option shouldn't include password)
      --user string                    username[:password] for authentication (prompt if password is not supplied)
      --version                        print the version and exit
```

## Examples
### Example 1: online analysis
It's pretty simple & straightforward. See the example below, it automatically diagnoses all the endpoints specified by
flag `--endpoints` and output the diagnosis result to both standard output and the file "etcd_diagnosis_report.json"
(see [example report](https://github.com/ahrtr/etcd-diagnosis/blob/main/examples/etcd_diagnosis_report.json))
under the current directory.

```
$ ./etcd-diagnosis --endpoints=https://10.0.1.10:2379,https://10.0.1.11:2379,https://10.0.1.12:2379 --cacert ./ca.crt --key ./etcd-diagnosis.key --cert ./etcd-diagnosis.crt
```

If the communication isn't protected by TLS (e.g. in dev environment), use a command something like below,
```
$ ./etcd-diagnosis --endpoints=http://10.0.1.10:2379,http://10.0.1.11:2379,http://10.0.1.12:2379
```

### Example 2: offline analysis

The following example performs offline analysis of etcd data located at ~/tmp/etcd/data/.
Each line of the output shows a key and its corresponding revision count in the format:
`key: revision_count`. Note: Make sure etcd is not running before executing this command.

The output is sorted in descending order by revision count.

```
$ ./etcd-diagnosis --offline --data-dir ~/tmp/etcd/data/
2025/05/23 16:58:37 etcd diagnosis performs offline analysis...
All key stats:
/registry/leases/kube-system/kube-controller-manager: 171
/registry/leases/kube-system/kube-scheduler: 171
/registry/leases/vmware-system-csi/external-resizer-csi-vsphere-vmware-com: 69
/registry/leases/vmware-system-csi/external-attacher-leader-csi-vsphere-vmware-com: 69
/registry/leases/vmware-system-csi/csi-vsphere-vmware-com: 69
/registry/leases/vmware-system-csi/vsphere-syncer: 69
/registry/nsx.vmware.com/nsxlocks/election-lock-pks-7c0903ef-3027-4634-94de-7e675e027f69: 58
/registry/masterleases/30.1.0.2: 36
/registry/leases/kube-node-lease/a341653e-a558-4e2d-ac86-cd7d828341bc: 35
/registry/leases/kube-node-lease/37ff8ff5-a3ab-454d-a1e7-a30ce20536f5: 35
/registry/leases/kube-node-lease/535d6dce-f28e-4b53-aff5-5eec8fde2786: 34
/registry/minions/a341653e-a558-4e2d-ac86-cd7d828341bc: 2
/registry/minions/37ff8ff5-a3ab-454d-a1e7-a30ce20536f5: 2
compact_rev_key: 2
/registry/minions/535d6dce-f28e-4b53-aff5-5eec8fde2786: 2
/registry/clusterrolebindings/metrics-server:system:auth-delegator: 1
/registry/pods/vmware-system-csi/vsphere-csi-node-jjrtp: 1
```

## Design
It's simple: one generic diagnosis engine + extensible plugins. Each plugin performs a diagnosis, and implements the
[Plugin](https://github.com/ahrtr/etcd-diagnosis/blob/67a33648b652430735af7b4b79037dc59171400c/engine/intf/plugin.go#L3)
interface. Currently, there are 5 plugins, see table below,

| Name                    | Description                                                                                           |
|-------------------------|-------------------------------------------------------------------------------------------------------|
| membershipChecker       | It checks whethere each endpoint has the same member list                                             |
| epStatusChecker         | It checks each endpoint's status, and verify whether their status is consistent                       |
| serializableReadChecker | It checks each endpoint can serve serialiable read requests and the duration to serve a read request  |
| linearizableReadChecker | It checks each endpoint can serve linearizable read requests and the duration to serve a read request |
| metricsChecker          | It collects some prometheus metrics from each endpoint                                                |
| any else?               | You are welcome to contribute new plugins!                                                            |

## Contributing
Any contribution (e.g. new plugins) is welcome!
