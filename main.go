package main

import (
	"fmt"
	"os"

	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	_ "k8s.io/kubernetes/pkg/version/prometheus" // for version metric registration
	"k8s.io/kubernetes/pkg/version/verflag"

	_ "github.com/cheapRoc/triton-cloud-controller-manager/triton"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

var (
	// Version is the version for this build, set at build time via LDFLAGS
	Version string
	// GitHash is the short-form commit hash of this build, set at build time
	GitHash string
)

func init() {
	healthz.DefaultHealthz()
}

func main() {
	// main context object for the controller manager
	s := options.NewCloudControllerManagerServer()
	s.AddFlags(pflag.CommandLine)

	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	verflag.PrintAndExitIfRequested()

	// create an instance of the named cloud provider
	// cloudConfigFile is the path to the cloud provider configuration file
	cloud, err := cloudprovider.InitCloudProvider("triton", s.CloudConfigFile)
	if err != nil {
		glog.Fatalf("Cloud provider could not be initialized: %v", err)
	}

	// runs the ExternalCMServer and should never exit
	if err := app.Run(s, cloud); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
