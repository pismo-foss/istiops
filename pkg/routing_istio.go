package pkg

import (
	"crypto/sha256"
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var subsetRouteExists bool

// IstioOperationsInterface set IstiOps interface for handling routing
type IstioOperationsInterface interface {
	SetLabelsVirtualService(cid string, name string, labels map[string]string) error
	SetLabelsDestinationRule(cid string, name string, labels map[string]string) error
	Headers(cid string, labels map[string]string, headers map[string]string) error
	Percentage(cid string, labels map[string]string, percentage int32) error
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(cid string, namespace string, listOptions metav1.ListOptions) (virtualServiceList *v1alpha32.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), cid)
	vss, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return vss, nil
}

// GetVirtualService returns a single virtualService object given a name & namespace
func GetVirtualService(cid string, name string, namespace string, getOptions metav1.GetOptions) (virtualService *v1alpha32.VirtualService, error error) {
	utils.Info(fmt.Sprintf("Getting virtualService '%s' to update...", name), cid)
	vs, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(name, getOptions)
	if err != nil {
		return nil, err
	}
	return vs, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(cid string, namespace string, listOptions metav1.ListOptions) (destinationRuleList *v1alpha32.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Getting all destinationrules..."), cid)
	drs, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return drs, nil
}

// GetDestinationRules returns a single destinationRule object given a name & namespace
func GetDestinationRule(cid string, name string, namespace string, getOptions metav1.GetOptions) (destinationRule *v1alpha32.DestinationRule, error error) {
	utils.Info(fmt.Sprintf("Getting destinationRule '%s' to update...", name), cid)
	dr, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return dr, nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(cid string, namespace string, virtualService *v1alpha32.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rule for virtualService '%s'...", virtualService.Name), cid)
	_, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(cid string, namespace string, destinationRule *v1alpha32.DestinationRule) error {
	utils.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), cid)
	_, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}

// GenerateShaFromMap returns a slice of hashes (sha256) for every key:value in given map[string]string
func GenerateShaFromMap(mapToHash map[string]string) ([]string, error) {
	var mapHashes []string

	for k, v := range mapToHash {
		keyValue := fmt.Sprintf("%s=%s", k, v)
		sha256 := sha256.Sum256([]byte(keyValue))
		mapHashes = append(mapHashes, fmt.Sprintf("%x", sha256))
	}

	return mapHashes, nil
}

func StringfyLabelSelector(cid string, labelSelector map[string]string) (string, error) {
	var labelsPair []string

	for key, value := range labelSelector {
		labelsPair = append(labelsPair, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(labelsPair[:], ","), nil
}

// GetResourcesToUpdate returns a slice of all DestinationRules and/or VirtualServices (based on given labelSelectors to a posterior update
func GetResourcesToUpdate(cid string, v IstioValues, labelSelector map[string]string) (matchedVirtualServices *v1alpha32.VirtualServiceList, matchedDestinationRules *v1alpha32.DestinationRuleList, error error) {
	stringfiedLabelSelector, _ := StringfyLabelSelector(cid, labelSelector)

	listOptions := metav1.ListOptions{
		LabelSelector: stringfiedLabelSelector,
	}

	matchedDrs, err := GetAllDestinationRules(cid, v.Namespace, listOptions)
	if err != nil {
		return nil, nil, err
	}

	matchedVss, err := GetAllVirtualServices(cid, v.Namespace, listOptions)
	if err != nil {
		return nil, nil, err
	}

	if len(matchedDrs.Items) == 0 || len(matchedVss.Items) == 0 {
		utils.Fatal(fmt.Sprintf("Couldn't find any istio resources based on given labelSelector '%s' to update. ", stringfiedLabelSelector), cid)
	}

	return matchedVss, matchedDrs, nil
}

// CreateNewVirtualServiceHttpRoute returns an existent VirtualService with a new basic HTTP route appended to it
func CreateNewVirtualServiceHttpRoute(cid string, labels map[string]string, hostname string, subset string, portNumber uint32) (httpRoute *v1alpha3.HTTPRoute, error error) {
	utils.Info(fmt.Sprintf("Creating new http route for subset '%s'...", subset), cid)

	newMatch := &v1alpha3.HTTPMatchRequest{
		Headers: map[string]*v1alpha3.StringMatch{},
	}

	// append user labels to exact match
	for labelKey, labelValue := range labels {
		newMatch.Headers[labelKey] = &v1alpha3.StringMatch{
			MatchType: &v1alpha3.StringMatch_Exact{
				Exact: labelValue,
			},
		}
	}

	defaultDestination := &v1alpha3.HTTPRouteDestination{
		Destination: &v1alpha3.Destination{
			Host:   hostname,
			Subset: subset,
			Port: &v1alpha3.PortSelector{
				Port: &v1alpha3.PortSelector_Number{
					Number: portNumber,
				},
			},
		},
	}

	newRoute := &v1alpha3.HTTPRoute{}
	newRoute.Match = append(newRoute.Match, newMatch)
	newRoute.Route = append(newRoute.Route, defaultDestination)

	return newRoute, nil
}

// Percentage set percentage as routing-match strategy for istio resources
func (v IstioValues) Percentage(cid string, labels map[string]string, percentage int32) error {

	return nil
}

// Headers set headers as routing-match strategy for istio resources
func (v IstioValues) Headers(cid string, labels map[string]string, headers map[string]string) error {
	var subsetRouteExists bool
	replacer := strings.NewReplacer(".", "", "-", "", "/", "")
	simplifiedVersion := replacer.Replace(v.Version)
	simplifiedVersion = strings.ToLower(simplifiedVersion)
	subsetRuleName := fmt.Sprintf("%s-%d", simplifiedVersion, v.Build)

	vss, drs, err := GetResourcesToUpdate(cid, v, labels)
	if err != nil {
		return err
	}

	for _, ds := range drs.Items {
		for subsetKey, subset := range ds.Spec.Subsets {
			// If an existent subsetName already exists, just update it with the given user labels otherwise create
			if subset.Name == subsetRuleName {
				utils.Info(fmt.Sprintf("Setting user labels to subset '%s", subset.Name), cid)
				ds.Spec.Subsets[subsetKey].Labels = labels
			}
		}

		err := UpdateDestinationRule(cid, v.Namespace, &ds)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could not update destinationRule '%s' due to error '%s'", ds.Name, err), cid)
		}
	}

	//Search for virtualservice's rule which matches subset name to append headers routing to it

	for _, vs := range vss.Items {
		subsetRouteExists = false
		for _, httpRules := range vs.Spec.Http {
			for _, matchValue := range httpRules.Route {
				// in case of a non-existent destination-subset, mark to be create it
				if matchValue.Destination.Subset == subsetRuleName {
					utils.Info(fmt.Sprintf("Subset '%s' already created for vs '%s", subsetRuleName, vs.Name ), cid)
					subsetRouteExists = true
					fmt.Println("here", subsetRouteExists)
				}
			}
		}

		// if a subset does not exists in the current VirtualService, create it from scratch
		if ! subsetRouteExists {
			// create it
			fmt.Println(vs.Name)
			newRoute, err := CreateNewVirtualServiceHttpRoute(cid, labels, "hostname", "subset", 8080)
			if err != nil {
				utils.Fatal(fmt.Sprintf("Could not create local httpRoute object for virtualservice '%s' due to error '%s'", vs.Name, err), cid)
			}
			fmt.Println(vs.Spec.Http)
			vs.Spec.Http = append(vs.Spec.Http, newRoute)
			err = UpdateVirtualService(cid, v.Namespace, &vs)
			if err != nil {
				utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", vs.Name, err), cid)
			}
		}
	}

	return nil
}

func (v IstioValues) SetLabelsDestinationRule(cid string, name string, labels map[string]string) error {

	return nil
}

func (v IstioValues) SetLabelsVirtualService(cid string, name string, labels map[string]string) error {

	return nil
}
