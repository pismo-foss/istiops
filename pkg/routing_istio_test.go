package pkg

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	versionedclientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"

	"testing"
)

func CreateMockedDestinationRule(t *testing.T, istioClient versioned.Interface) (mocked *v1alpha32.DestinationRule) {
	mockedDr := &v1alpha32.DestinationRule{}
	mockedDr.Name = "api-unit-test-destinationrule"
	mockedDr.Namespace = "default"

	mockedDr.Spec.Subsets = append(mockedDr.Spec.Subsets, &v1alpha3.Subset{
		Name: "subset-test",
		Labels: map[string]string{
			"app": "api-xpto",
		},
	})

	mockedDr, err := istioClient.NetworkingV1alpha3().DestinationRules("default").Create(mockedDr)
	if mockedDr == nil {
		t.Error("Could not create mocked destinationRule")
		t.Error(err)
	}
	return mockedDr
}

func CreateMockedVirtualService(t *testing.T, istioClient versioned.Interface) (mocked *v1alpha32.VirtualService) {
	mockedVs := &v1alpha32.VirtualService{}
	mockedVs.Name = "api-unit-test-virtualservice"
	mockedVs.Namespace = "default"

	mockedVs.Spec.Hosts = []string{"api-unit-test.domain.io"}
	mockedVs.Spec.Gateways = []string{"unit-test-gateway"}

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, &v1alpha3.HTTPRoute{})

	defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}
	defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: "api-xpto", Subset: "subset-name", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: 8080}}}}
	defaultRoute := &v1alpha3.HTTPRoute{}
	defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
	defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, defaultRoute)

	mockedVs, err := istioClient.NetworkingV1alpha3().VirtualServices("default").Create(mockedVs)
	if mockedVs == nil {
		t.Error("Could not create mocked virtualservice")
		t.Error(err)
	}

	return mockedVs
}

func TestGetAllVirtualServices(t *testing.T) {
	istioClient := versionedclientFake.NewSimpleClientset()

	mockedDr := CreateMockedDestinationRule(t, istioClient)
	assert.NotNil(t, mockedDr)

	mockedVs := CreateMockedVirtualService(t, istioClient)
	assert.NotNil(t, mockedVs)

	listOptions := metav1.ListOptions{}
	_, err := GetAllVirtualServices("random-cid", "default", listOptions)
	assert.NoError(t, err)
}

//
//func TestGetAllDestinationRules(t *testing.T) {
//	_, err := GetAllDestinationRules("random-cid", "default")
//	assert.NoError(t, err)
//}
//
//func TestGetResourcesToUpdate(t *testing.T) {
//	v := IstioValues{
//		"api-xpto",
//		"2.0.0",
//		123,
//		"default",
//	}
//	userLabelSelector := map[string]string{
//		"app": "api-xpto",
//	}
//	_, err := GetResourcesToUpdate("random-cid", v, userLabelSelector)
//	assert.NoError(t, err)
//
//}
