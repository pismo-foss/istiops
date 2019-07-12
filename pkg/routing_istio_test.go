package pkg

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedclientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"log"
	"os"

	"github.com/stretchr/testify/assert"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var namespace string

func TestMain(m *testing.M) {
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp() {
	istioClient = versionedclientFake.NewSimpleClientset()
	namespace := "test-namespace"
	var err error

	err = CreateMockedDestinationRule("api-unit-test-destinationrule", namespace)
	if err != nil {
		log.Fatal("Could not create mocked DestinationRule")
	}

	err = CreateMockedVirtualService("api-unit-test-virtualservice", namespace)
	if err != nil {
		log.Fatal("Could not create mocked VirtualService")
	}
}

func CreateMockedDestinationRule(resourceName string, resourceNamespace string) (error error) {
	mockedDr := &v1alpha32.DestinationRule{}
	mockedDr.Name = resourceName
	mockedDr.Namespace = resourceNamespace

	mockedDr.Spec.Subsets = append(mockedDr.Spec.Subsets, &v1alpha3.Subset{
		Name: "subset-test",
		Labels: map[string]string{
			"app":     "api-xpto",
			"version": "2.1.3",
		},
	})

	_, err := istioClient.NetworkingV1alpha3().DestinationRules(resourceNamespace).Create(mockedDr)

	return err
}

func CreateMockedVirtualService(resourceName string, resourceNamespace string) (error error) {
	mockedVs := &v1alpha32.VirtualService{}
	mockedVs.Name = resourceName
	mockedVs.Namespace = resourceNamespace

	mockedVs.Spec.Hosts = []string{"api-unit-test.domain.io"}
	mockedVs.Spec.Gateways = []string{"unit-test-gateway"}

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, &v1alpha3.HTTPRoute{})

	defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}
	defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: "api-xpto", Subset: "subset-test", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: 8080}}}}
	defaultRoute := &v1alpha3.HTTPRoute{}
	defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
	defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, defaultRoute)

	_, err := istioClient.NetworkingV1alpha3().VirtualServices(resourceNamespace).Create(mockedVs)

	return err
}

func TestGetAllDestinationRules(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedDrs, err := GetAllDestinationRules("random-cid", namespace, listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.DestinationRuleList{}, *mockedDrs)
	assert.EqualValues(t, "api-unit-test-destinationrule", mockedDrs.Items[0].Name)
}

func TestGetAllVirtualServices(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedVss, err := GetAllVirtualServices("random-cid", namespace, listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.VirtualServiceList{}, *mockedVss)
	assert.EqualValues(t, "api-unit-test-virtualservice", mockedVss.Items[0].Name)
}

func TestGetResourcesToUpdate(t *testing.T) {
	v := IstioValues{
		"api-xpto",
		"2.0.0",
		123,
		namespace,
	}

	userLabelSelector := map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

	// test case: happy path
	mockedResourcesToUpdate, err := GetResourcesToUpdate("random-cid", v, userLabelSelector)
	assert.NoError(t, err)
	assert.NotNil(t, mockedResourcesToUpdate)

	// test case: missing labelSelector
	mockedMissingLabels, err := GetResourcesToUpdate("random-cid", v, map[string]string{})
	assert.NoError(t, err)
	assert.Nil(t, mockedMissingLabels)

	// test case: missing IstioValues params
	mockedMissingValues, err := GetResourcesToUpdate("random-cid", IstioValues{}, userLabelSelector)
	assert.NoError(t, err)
	assert.NotNil(t, mockedMissingValues)
}
