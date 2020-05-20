package controllers

import (
	"context"
	"fmt"
	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	v1alpha1 "github.com/kapp-staging/kapp/controller/api/v1alpha1"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/suite"
	istioScheme "istio.io/client-go/pkg/clientset/versioned/scheme"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"math/rand"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BasicSuite struct {
	suite.Suite

	Cfg            *rest.Config
	K8sClient      client.Client
	TestEnv        *envtest.Environment
	MgrStopChannel chan struct{}
}

func (suite *BasicSuite) Eventually(condition func() bool, msgAndArgs ...interface{}) {
	waitFor := time.Second * 20
	tick := time.Millisecond * 500
	suite.Suite.Require().Eventually(condition, waitFor, tick, msgAndArgs...)
}

func (suite *BasicSuite) createComponentPlugin(plugin *v1alpha1.ComponentPlugin) {
	suite.Nil(suite.K8sClient.Create(context.Background(), plugin))

	// after the finalizer is set, the plugin won't auto change
	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(context.Background(), getComponentPluginNamespacedName(plugin), plugin)

		if err != nil {
			return false
		}

		for i := range plugin.Finalizers {
			if plugin.Finalizers[i] == finalizerName {
				return true
			}
		}

		return false
	})
}

func (suite *BasicSuite) createApplicationPlugin(plugin *v1alpha1.ApplicationPlugin) {
	suite.Nil(suite.K8sClient.Create(context.Background(), plugin))

	// after the finalizer is set, the plugin won't auto change
	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(context.Background(), getApplicationPluginNamespacedName(plugin), plugin)

		if err != nil {
			return false
		}

		for i := range plugin.Finalizers {
			if plugin.Finalizers[i] == finalizerName {
				return true
			}
		}

		return false
	})
}

func (suite *BasicSuite) createApplication(application *v1alpha1.Application) {
	suite.Nil(suite.K8sClient.Create(context.Background(), application))

	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(context.Background(), getApplicationNamespacedName(application), application)

		if err != nil {
			return false
		}

		for i := range application.Finalizers {
			if application.Finalizers[i] == finalizerName {
				return true
			}
		}

		return false
	}, "Created application has no finalizer.")
}

func getDockerRegistryNamespacedName(registry *v1alpha1.DockerRegistry) types.NamespacedName {
	return types.NamespacedName{Name: registry.Name, Namespace: registry.Namespace}
}

func (suite *BasicSuite) createDockerRegistry(registry *v1alpha1.DockerRegistry) {
	suite.Nil(suite.K8sClient.Create(context.Background(), registry))

	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(context.Background(), getDockerRegistryNamespacedName(registry), registry)

		if err != nil {
			return false
		}

		for i := range registry.Finalizers {
			if registry.Finalizers[i] == finalizerName {
				return true
			}
		}

		return false
	}, "Created Docker registry has no finalizer.")
}

func (suite *HttpsCertIssuerControllerSuite) createHttpsCertIssuer(issuer v1alpha1.HttpsCertIssuer) {
	suite.Nil(suite.K8sClient.Create(context.Background(), &issuer))
}

func (suite *BasicSuite) reloadObject(key client.ObjectKey, obj runtime.Object) {
	suite.Nil(suite.K8sClient.Get(context.Background(), key, obj))
}

func (suite *BasicSuite) updateObject(obj runtime.Object) {
	suite.Nil(suite.K8sClient.Update(context.Background(), obj))
}

func (suite *BasicSuite) createObject(obj runtime.Object) {
	suite.Nil(suite.K8sClient.Create(context.Background(), obj))
}

func (suite *BasicSuite) reloadComponent(component *v1alpha1.Component) {
	suite.reloadObject(types.NamespacedName{Name: component.Name, Namespace: component.Namespace}, component)
}

func (suite *BasicSuite) updateComponent(component *v1alpha1.Component) {
	suite.updateObject(component)
}

func (suite *BasicSuite) createComponent(component *v1alpha1.Component) {
	suite.createObject(component)
}

func (suite *BasicSuite) SetupSuite() {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(ginkgo.GinkgoWriter)))

	// bootstrapping test environment
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
			filepath.Join("config", "crd", "bases"),
			filepath.Join("..", "resources", "istio"),
			filepath.Join("resources", "istio"),
			filepath.Join("resources"),
			filepath.Join("..", "resources"),
		},
	}

	var err error
	cfg, err := testEnv.Start()
	suite.Nil(err)
	suite.NotNil(cfg)
	suite.Nil(scheme.AddToScheme(scheme.Scheme))
	suite.Nil(istioScheme.AddToScheme(scheme.Scheme))
	suite.Nil(v1alpha1.AddToScheme(scheme.Scheme))
	suite.Nil(v1alpha2.AddToScheme(scheme.Scheme))

	// +kubebuilder:scaffold:scheme

	min := 2000
	max := 8000
	port := rand.Intn(max-min) + min

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: fmt.Sprintf("localhost:%d", port),
	})

	suite.NotNil(mgr)
	suite.Nil(err)

	suite.Nil(NewApplicationReconciler(mgr).SetupWithManager(mgr))
	suite.Nil(NewApplicationPluginReconciler(mgr).SetupWithManager(mgr))
	suite.Nil(NewApplicationPluginBindingReconciler(mgr).SetupWithManager(mgr))

	suite.Nil(NewComponentReconciler(mgr).SetupWithManager(mgr))
	suite.Nil(NewComponentPluginReconciler(mgr).SetupWithManager(mgr))
	suite.Nil(NewComponentPluginBindingReconciler(mgr).SetupWithManager(mgr))

	suite.Nil(NewHttpsCertIssuerReconciler(mgr).SetupWithManager(mgr))
	suite.Nil(NewHttpsCertReconciler(mgr).SetupWithManager(mgr))

	suite.Nil(NewDockerRegistryReconciler(mgr).SetupWithManager(mgr))

	mgrStopChannel := make(chan struct{})

	go func() {
		err = mgr.Start(mgrStopChannel)
		suite.Nil(err)
	}()

	k8sClient := mgr.GetClient()
	suite.NotNil(k8sClient)

	suite.TestEnv = testEnv
	suite.K8sClient = k8sClient
	suite.Cfg = cfg
	suite.MgrStopChannel = mgrStopChannel
}

func (suite *BasicSuite) TearDownSuite() {
	suite.MgrStopChannel <- struct{}{}
	suite.Nil(suite.TestEnv.Stop())
}

func (suite *BasicSuite) createHttpsCert(cert v1alpha1.HttpsCert) {
	suite.Nil(suite.K8sClient.Create(context.Background(), &cert))

	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(
			context.Background(),
			types.NamespacedName{Name: cert.Name},
			&cert,
		)

		return err == nil
	})
}

func (suite *BasicSuite) createHttpsCertIssuer(issuer v1alpha1.HttpsCertIssuer) {
	suite.Nil(suite.K8sClient.Create(context.Background(), &issuer))

	suite.Eventually(func() bool {
		err := suite.K8sClient.Get(
			context.Background(),
			types.NamespacedName{
				Name: issuer.Name,
			},
			&issuer,
		)

		return err == nil
	})
}
