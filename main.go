package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/ahrtr/etcd-diagnosis/agent"
	"github.com/ahrtr/etcd-diagnosis/engine"
	"github.com/ahrtr/etcd-diagnosis/engine/intf"
	"github.com/ahrtr/etcd-diagnosis/offline"
	"github.com/ahrtr/etcd-diagnosis/plugins/epstatus"
	"github.com/ahrtr/etcd-diagnosis/plugins/membership"
	"github.com/ahrtr/etcd-diagnosis/plugins/metrics"
	"github.com/ahrtr/etcd-diagnosis/plugins/read"
)

var (
	globalCfg = agent.GlobalConfig{}
)

func newDiagnosisCommand() *cobra.Command {
	diagnosisCmd := &cobra.Command{
		Use:   "etcd-diagnosis",
		Short: "An one-stop etcd diagnosis tool",
		Run:   diagnosisCommandFunc,
	}

	diagnosisCmd.Flags().StringSliceVar(&globalCfg.Endpoints, "endpoints", []string{"127.0.0.1:2379"}, "comma separated etcd endpoints")
	diagnosisCmd.Flags().BoolVar(&globalCfg.UseClusterEndpoints, "cluster", false, "use all endpoints from the cluster member list")

	diagnosisCmd.Flags().DurationVar(&globalCfg.DialTimeout, "dial-timeout", 2*time.Second, "dial timeout for client connections")
	diagnosisCmd.Flags().DurationVar(&globalCfg.CommandTimeout, "command-timeout", 5*time.Second, "command timeout (excluding dial timeout)")
	diagnosisCmd.Flags().DurationVar(&globalCfg.KeepAliveTime, "keepalive-time", 2*time.Second, "keepalive time for client connections")
	diagnosisCmd.Flags().DurationVar(&globalCfg.KeepAliveTimeout, "keepalive-timeout", 5*time.Second, "keepalive timeout for client connections")

	diagnosisCmd.Flags().BoolVar(&globalCfg.Insecure, "insecure-transport", true, "disable transport security for client connections")

	diagnosisCmd.Flags().BoolVar(&globalCfg.InsecureSkipVerify, "insecure-skip-tls-verify", false, "skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
	diagnosisCmd.Flags().StringVar(&globalCfg.CertFile, "cert", "", "identify secure client using this TLS certificate file")
	diagnosisCmd.Flags().StringVar(&globalCfg.KeyFile, "key", "", "identify secure client using this TLS key file")
	diagnosisCmd.Flags().StringVar(&globalCfg.CaFile, "cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")

	diagnosisCmd.Flags().StringVar(&globalCfg.Username, "user", "", "username[:password] for authentication (prompt if password is not supplied)")
	diagnosisCmd.Flags().StringVar(&globalCfg.Password, "password", "", "password for authentication (if this option is used, --user option shouldn't include password)")

	diagnosisCmd.Flags().StringVarP(&globalCfg.DNSDomain, "discovery-srv", "d", "", "domain name to query for SRV records describing cluster endpoints")
	diagnosisCmd.Flags().StringVarP(&globalCfg.DNSService, "discovery-srv-name", "", "", "service name to query when using DNS discovery")
	diagnosisCmd.Flags().BoolVar(&globalCfg.InsecureDiscovery, "insecure-discovery", true, "accept insecure SRV records describing cluster endpoints")

	diagnosisCmd.Flags().IntVar(&globalCfg.DbQuotaBytes, "etcd-storage-quota-bytes", 2*1024*1024*1024, "etcd storage quota in bytes (the value passed to etcd instance by flag --quota-backend-bytes)")

	diagnosisCmd.Flags().BoolVar(&globalCfg.PrintVersion, "version", false, "print the version and exit")

	diagnosisCmd.Flags().BoolVar(&globalCfg.Offline, "offline", false, "offline analysis")
	diagnosisCmd.Flags().StringVar(&globalCfg.DataDir, "data-dir", "", "path to data directory")

	return diagnosisCmd
}

func main() {
	diagnosisCmd := newDiagnosisCommand()
	if err := diagnosisCmd.Execute(); err != nil {
		if diagnosisCmd.SilenceErrors {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
		os.Exit(1)
	}
}

func printVersion(printVersion bool) {
	if printVersion {
		fmt.Printf("etcd-diagnosis Version: %s\n", Version)
		fmt.Printf("Git SHA: %s\n", GitSHA)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
}

func diagnosisCommandFunc(_ *cobra.Command, _ []string) {
	printVersion(globalCfg.PrintVersion)

	if globalCfg.Offline {
		offline.AnalyzeOffline(globalCfg.DataDir)
		return
	}

	onlineAnalysis()

	log.Println("etcd diagnosis done!")
}

func onlineAnalysis() {
	log.Println("etcd diagnosis starting...")

	plugins := []intf.Plugin{
		membership.NewPlugin(globalCfg),
		epstatus.NewPlugin(globalCfg),
		read.NewPlugin(globalCfg, false),
		read.NewPlugin(globalCfg, true),
		metrics.NewPlugin(globalCfg),
	}

	engine.Diagnose(globalCfg, plugins)
}
