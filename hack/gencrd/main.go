package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appscode/go/log"
	gort "github.com/appscode/go/runtime"
	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/client-go/openapi"
	"kubedb.dev/apimachinery/apis"
	cataloginstall "kubedb.dev/apimachinery/apis/catalog/install"
	catalogv1alpha1 "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	kubedbinstall "kubedb.dev/apimachinery/apis/kubedb/install"
	kubedbv1alpha1 "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
)

func generateCRDDefinitions() {
	apis.EnableStatusSubresource = true

	filename := gort.GOPath() + "/src/kubedb.dev/apimachinery/apis/kubedb/v1alpha1/crds.yaml"
	os.Remove(filename)

	err := os.MkdirAll(filepath.Join(gort.GOPath(), "/src/kubedb.dev/apimachinery/api/crds"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	crds := []*crd_api.CustomResourceDefinition{
		kubedbv1alpha1.DormantDatabase{}.CustomResourceDefinition(),
		kubedbv1alpha1.Elasticsearch{}.CustomResourceDefinition(),
		kubedbv1alpha1.Etcd{}.CustomResourceDefinition(),
		kubedbv1alpha1.Memcached{}.CustomResourceDefinition(),
		kubedbv1alpha1.MongoDB{}.CustomResourceDefinition(),
		kubedbv1alpha1.MySQL{}.CustomResourceDefinition(),
		kubedbv1alpha1.Postgres{}.CustomResourceDefinition(),
		kubedbv1alpha1.Redis{}.CustomResourceDefinition(),
		kubedbv1alpha1.Snapshot{}.CustomResourceDefinition(),

		catalogv1alpha1.ElasticsearchVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.EtcdVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.MemcachedVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.MongoDBVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.MySQLVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.PostgresVersion{}.CustomResourceDefinition(),
		catalogv1alpha1.RedisVersion{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		filename := filepath.Join(gort.GOPath(), "/src/kubedb.dev/apimachinery/api/crds", crd.Spec.Names.Singular+".yaml")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		crdutils.MarshallCrd(f, crd, "yaml")
		f.Close()
	}
}

func generateSwaggerJson() {
	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	kubedbinstall.Install(Scheme)
	cataloginstall.Install(Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "KubeDB",
			Version: "v0",
			Contact: &spec.ContactInfo{
				Name:  "AppsCode Inc.",
				URL:   "https://appscode.com",
				Email: "hello@appscode.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			kubedbv1alpha1.GetOpenAPIDefinitions,
			catalogv1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []openapi.TypeInfo{
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralDormantDatabase, kubedbv1alpha1.ResourceKindDormantDatabase, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralElasticsearch, kubedbv1alpha1.ResourceKindElasticsearch, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralEtcd, kubedbv1alpha1.ResourceKindEtcd, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralMariaDB, kubedbv1alpha1.ResourceKindMariaDB, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralMemcached, kubedbv1alpha1.ResourceKindMemcached, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralMongoDB, kubedbv1alpha1.ResourceKindMongoDB, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralMySQL, kubedbv1alpha1.ResourceKindMySQL, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralPerconaXtraDB, kubedbv1alpha1.ResourceKindPerconaXtraDB, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralPgBouncer, kubedbv1alpha1.ResourceKindPgBouncer, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralPostgres, kubedbv1alpha1.ResourceKindPostgres, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralProxySQL, kubedbv1alpha1.ResourceKindProxySQL, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralRedis, kubedbv1alpha1.ResourceKindRedis, true},
			{kubedbv1alpha1.SchemeGroupVersion, kubedbv1alpha1.ResourcePluralSnapshot, kubedbv1alpha1.ResourceKindSnapshot, true},

			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralElasticsearchVersion, catalogv1alpha1.ResourceKindElasticsearchVersion, true},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralEtcdVersion, catalogv1alpha1.ResourceKindEtcdVersion, false},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralMemcachedVersion, catalogv1alpha1.ResourceKindMemcachedVersion, true},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralMongoDBVersion, catalogv1alpha1.ResourceKindMongoDBVersion, true},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralMySQLVersion, catalogv1alpha1.ResourceKindMySQLVersion, true},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralPerconaXtraDBVersion, catalogv1alpha1.ResourceKindPerconaXtraDBVersion, false},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralPostgresVersion, catalogv1alpha1.ResourceKindPostgresVersion, false},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralProxySQLVersion, catalogv1alpha1.ResourceKindProxySQLVersion, false},
			{catalogv1alpha1.SchemeGroupVersion, catalogv1alpha1.ResourcePluralRedisVersion, catalogv1alpha1.ResourceKindRedisVersion, false},
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/kubedb.dev/apimachinery/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	// generateCRDDefinitions()
	generateSwaggerJson()
}
