package main

import (
	"os"

	"github.com/appscode/go/log"
	"github.com/appscode/go/runtime"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func main() {
	filename := runtime.GOPath() + "/src/github.com/kubedb/apimachinery/apis/kubedb/v1alpha1/crds.yaml"

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	crds := []*crd_api.CustomResourceDefinition{
		api.Postgres{}.CustomResourceDefinition(),
		api.Elasticsearch{}.CustomResourceDefinition(),
		api.MySQL{}.CustomResourceDefinition(),
		api.MongoDB{}.CustomResourceDefinition(),
		api.Redis{}.CustomResourceDefinition(),
		api.Memcached{}.CustomResourceDefinition(),
		api.Snapshot{}.CustomResourceDefinition(),
		api.DormantDatabase{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		crdutils.MarshallCrd(f, crd, "yaml")
	}
}