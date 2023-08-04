package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
	admv1 "k8s.io/api/admission/v1beta1"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

type Config struct {
	SecureServingInfo *server.SecureServingInfo
}

func (o *Options) Config() *Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}
	c := Config{}
	o.SecureServingOptions.ApplyTo(&c.SecureServingInfo)
	return &c
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
	return o
}

func main() {
	options := NewDefautlOptions()
	fs := pflag.NewFlagSet(valCon, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valCon)

	options.AddFlagSet(fs)

	if err := fs.Parse(os.Args); err != nil {
		panic(err)
	}

	c := options.Config()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(DeploymentValdiation))

	stopCh := server.SetupSignalHandler()

	ch, _, err := c.SecureServingInfo.Serve(mux, 30*time.Second, stopCh)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}
}

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func DeploymentValdiation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeploymentValidation was called")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error %s, reading the bdoy", err.Error())
	}

	gvk := admv1.SchemeGroupVersion.WithKind("AdmissionReview")
	var admissionReview admv1.AdmissionReview

	_, _, err = codecs.UniversalDeserializer().Decode(body, &gvk, &admissionReview)
	if err != nil {
		fmt.Printf("Error %s, converting req body to admission review type", err.Error())
	}

	gvkDeployment := appv1.SchemeGroupVersion.WithKind("Deployment")
	var d appv1.Deployment
	_, _, err = codecs.UniversalDeserializer().Decode(admissionReview.Request.Object.Raw, &gvkDeployment, &d)
	if err != nil {
		fmt.Printf("Error %s, while getting deployement type from admissionreview", err.Error())
	}

	fmt.Printf("deployment resource that we have is %+v\n", d)
}
