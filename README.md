# etcd-diagnosis

## Overview
etcd-diagnosis is a comprehensive tool for etcd diagnosis. It diagnoses running etcd clusters and generates a
report with just one command. It reuses most of the `etcdctl` global flags, so users follow the same experience
as `etcdctl` to use `etcd-diagnosis`. See the complete flags below,

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
      --password string                password for authentication (if this option is used, --user option shouldn't include password)
      --user string                    username[:password] for authentication (prompt if password is not supplied)
      --version                        print the version and exit
```

## Examples
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

## Design
It's simple: one generic diagnosis engine + extensible plugins. Each plugin performs a diagnosis, and implements the
[Plugin](https://github.com/ahrtr/etcd-diagnosis/blob/67a33648b652430735af7b4b79037dc59171400c/engine/intf/plugin.go#L3)
interface. Currently, there are 4 plugins, see table below,

| Name              | Description                                                                                          |
|-------------------|------------------------------------------------------------------------------------------------------|
| membershipChecker | It checks whethere each endpoint has the same member list                                            |
| epStatusChecker   | It checks each endpoint's status, and verify whether their status is consistent                      |
| readCheck         | It checks each endpoint can serve serialiable read requests and the duration to serve a read request |
| metricsChecker    | It collects some prometheus metrics from each endpoint                                               |
| any else?         | You are welcome to contribute new plugins!                                                           |

## Contributing
Any contribution (e.g. new plugins) is welcome!
