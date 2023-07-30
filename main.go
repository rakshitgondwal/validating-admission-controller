package main

import (
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

const (
	valCon = "validatin-controller"
)

func NewDefautlOptions() *Options {
	o := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	o.SecureServingOptions.BindPort = 8443
	o.SecureServingOptions.ServerCert.PairName = valCon
}

func main() {
	options := NewDefautlOptions()
	fs := pflag.NewFlagSet(valCon, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valCon)

	options.AddFlagSet(fs)
}
