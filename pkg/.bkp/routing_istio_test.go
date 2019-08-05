package pkg

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedclientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"io/ioutil"
	"log"
	"os"

	"github.com/stretchr/testify/assert"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var namespace string
var mockedDestinationRuleName string
var mockedVirtualServiceName string

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp() {
	istioClient = versionedclientFake.NewSimpleClientset()
	namespace = "test-namespace"
	mockedDestinationRuleName = "api-unit-test-destinationrule"
	mockedVirtualServiceName = "api-unit-test-virtualservice"

	var err error

	_, err = CreateMockedDestinationRule(
		mockedDestinationRuleName,
		namespace,
		map[string]string{
			"app":     "api-xpto",
			"version": "2.1.3",
		})
	if err != nil {
		log.Fatal("Could not create mocked DestinationRule")
	}

	_, err = CreateMockedVirtualService(mockedVirtualServiceName, namespace)
	if err != nil {
		log.Fatal("Could not create mocked VirtualService")
	}
}

func CreateMockedDestinationRule(resourceName string, resourceNamespace string, subsetLabels map[string]string) (mockedDestinationRule *v1alpha32.DestinationRule, error error) {
	mockedDr := &v1alpha32.DestinationRule{}
	mockedDr.Name = resourceName
	mockedDr.Namespace = resourceNamespace
	mockedDr.Labels = map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

	mockedDr.Spec.Subsets = append(mockedDr.Spec.Subsets, &v1alpha3.Subset{
		Name:   "subset-test",
		Labels: subsetLabels,
	})

	_, err := istioClient.NetworkingV1alpha3().DestinationRules(resourceNamespace).Create(mockedDr)

	return mockedDr, err
}

func CreateMockedVirtualService(resourceName string, resourceNamespace string) (mockedVirtualService *v1alpha32.VirtualService, error error) {
	mockedVs := &v1alpha32.VirtualService{}
	mockedVs.Name = resourceName
	mockedVs.Namespace = resourceNamespace
	mockedVs.Labels = map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

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

	return mockedVs, err
}

func TestSanitizeVersionString(t *testing.T) {
	versions := []string{"2.0.0", "2/0/0", "2-0_0"}

	for _, version := range versions {
		sanitizedVersion, err := SanitizeVersionString(version)
		assert.NoError(t, err)
		assert.Equal(t, "200", sanitizedVersion)
	}

}

func TestGetAllDestinationRules(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedDrs, err := GetAllDestinationRules("random-cid", namespace, listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.DestinationRuleList{}, *mockedDrs)
	assert.EqualValues(t, "api-unit-test-destinationrule", mockedDrs.Items[0].Name)
}

func TestGetDestinationRule(t *testing.T) {
	getOptions := metav1.GetOptions{}
	dr, err := GetDestinationRule("random-cid", mockedDestinationRuleName, namespace, getOptions)
	if err != nil {
	}

	assert.NoError(t, err)
	assert.NotNil(t, mockedDestinationRuleName, dr.Name)
	assert.IsType(t, v1alpha32.DestinationRule{}, *dr)
	assert.EqualValues(t, mockedDestinationRuleName, dr.Name)
}

func TestGetAllVirtualServices(t *testing.T) {
	listOptions := metav1.ListOptions{}
	mockedVss, err := GetAllVirtualServices("random-cid", namespace, listOptions)

	assert.NoError(t, err)
	assert.IsType(t, v1alpha32.VirtualServiceList{}, *mockedVss)
	assert.EqualValues(t, "api-unit-test-virtualservice", mockedVss.Items[0].Name)
}

func TestGetVirtualService(t *testing.T) {
	getOptions := metav1.GetOptions{}
	vs, err := GetVirtualService("random-cid", mockedVirtualServiceName, namespace, getOptions)
	if err != nil {
	}

	assert.NoError(t, err)
	assert.NotNil(t, mockedVirtualServiceName, vs.Name)
	assert.IsType(t, v1alpha32.VirtualService{}, *vs)
	assert.EqualValues(t, mockedVirtualServiceName, vs.Name)
}

func TestGetResourcesToUpdate(t *testing.T) {
	v := IstioValues{
		"api-xpto",
		"2.1.3",
		123,
		namespace,
	}

	userLabelSelector := map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

	// test case: happy path
	mockedDrs, mockedVrs, err := GetResourcesToUpdate("random-cid", v, userLabelSelector)
	assert.NoError(t, err)
	assert.NotNil(t, mockedDrs)
	assert.NotNil(t, mockedVrs)

	// test case: missing labelSelector
	mockedDrsMissingLabels, mockedVssMissingLabels, err := GetResourcesToUpdate("random-cid", v, map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockedDrsMissingLabels.Items))
	assert.Equal(t, 1, len(mockedVssMissingLabels.Items))

	// test case: missing IstioValues params
	mockedDrsMissingValues, mockedVssMissingValues, err := GetResourcesToUpdate("random-cid", IstioValues{}, userLabelSelector)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockedDrsMissingValues.Items))
	assert.Equal(t, 1, len(mockedVssMissingValues.Items))
}

func TestUpdateDestinationRule(t *testing.T) {
	newMockedDestinationRule := "custom-destinationrule"
	updatedHost := "updated-host"

	dr, _ := CreateMockedDestinationRule(newMockedDestinationRule, namespace, map[string]string{"app": "new-app"})
	dr.Spec.Host = updatedHost

	err := UpdateDestinationRule("random-cid", namespace, dr)
	assert.NoError(t, err)
	assert.EqualValues(t, newMockedDestinationRule, dr.Name)
	assert.EqualValues(t, updatedHost, dr.Spec.Host)

}

func TestUpdateVirtualService(t *testing.T) {
	newMockedVirtualService := "custom-virtualservice"
	updatedHosts := []string{"updated-host-1", "updated-host-2"}

	vs, _ := CreateMockedVirtualService(newMockedVirtualService, namespace)
	vs.Spec.Hosts = updatedHosts

	err := UpdateVirtualService("random-cid", namespace, vs)
	assert.NoError(t, err)
	assert.EqualValues(t, newMockedVirtualService, vs.Name)
	assert.IsType(t, []string{}, vs.Spec.Hosts)
	assert.EqualValues(t, updatedHosts, vs.Spec.Hosts)

}
