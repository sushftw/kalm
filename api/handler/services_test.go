package handler

import (
	"github.com/kalmhq/kalm/api/resources"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

type ServicesHandlerTestSuite struct {
	WithControllerTestSuite
	namespace string
}

func (suite *ServicesHandlerTestSuite) SetupSuite() {
	suite.WithControllerTestSuite.SetupSuite()
	suite.namespace = "ns-test-services"
	suite.ensureNamespaceExist(suite.namespace)
}

func (suite *ServicesHandlerTestSuite) TestServicesHandler() {
	// create a service
	service := v1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-services",
			Namespace: suite.namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Port:     80,
					Protocol: v1.ProtocolTCP,
				},
			},
		},
	}

	err := suite.Create(&service)
	suite.Nil(err)

	// list all services
	var services []*resources.Service
	rec := suite.NewRequest(http.MethodGet, "/v1alpha1/services", "")
	rec.BodyAsJSON(&services)
	suite.NotNil(rec)
	suite.EqualValues(2, len(services)) // default service named kubernetes
	suite.EqualValues("test-services", services[1].Name)
}

func TestServicesHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ServicesHandlerTestSuite))
}
