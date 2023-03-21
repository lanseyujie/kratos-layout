package conf

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/contrib/config/kubernetes/v2"
	"github.com/go-kratos/kratos/v2/config"
	"k8s.io/client-go/util/homedir"
)

type SourceKubernetes config.Source

func NewSourceKubernetes() SourceKubernetes {
	// load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined
	opts := []kubernetes.Option{
		kubernetes.Namespace("humpback"),
		kubernetes.LabelSelector("app=" + Name),
	}

	kubeConfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if f, err := os.Stat(kubeConfig); err == nil && !f.IsDir() {
		opts = append(opts, kubernetes.KubeConfig(kubeConfig))
	}

	// kubectl create configmap <config-name> \
	//   --from-file=path/to/config/files \
	//   --namespace=<namespace> \
	//   --dry-run=client -o yaml \
	//   | sed '/^\s*creationTimestamp:/d' \
	//   | sed '/^\s*resourceVersion:/d' \
	//   | sed '/^\s*uid:/d' \
	//   | sed 's/\(namespace: <namespace>\)/\1\n  labels:\n    app: <config-name>/' \
	//   | kubectl apply -f -
	source := kubernetes.NewSource(opts...)

	return source
}
