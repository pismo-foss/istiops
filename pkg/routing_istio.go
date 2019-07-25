package pkg

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SanitizeVersionString returns a non-special character string given a semantic version. Ex. 3.0.0 -> 300
func SanitizeVersionString(version string) (sanitizedVersion string, error error) {
	replacer := strings.NewReplacer(
		".", "",
		"-", "",
		"/", "",
		"_", "",
	)
	sanitizedVersion = replacer.Replace(version)
	sanitizedVersion = strings.ToLower(sanitizedVersion)

	return sanitizedVersion, nil
}

// IstioOperationsInterface set IstiOps interface for handling routing
type IstioOperationsInterface interface {
	SetLabelsVirtualService(cid string, name string, labels map[string]string) error
	SetLabelsDestinationRule(cid string, name string, labels map[string]string) error
	SetHeaders(cid string, labels map[string]string, headers map[string]string) (subset string, error error)
	SetPercentage(cid string, virtualServiceName string, subset string, percentage int32) error
	ClearDestinationRules(cid string, labels map[string]string) error
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
		return nil, nil, err
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
func (v IstioValues) SetPercentage(cid string, virtualServiceName string, subset string, percentage int32) error {
	vs, err := GetVirtualService(cid, virtualServiceName, v.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for _, httpRules := range vs.Spec.Http {
		for _, httpRoute := range httpRules.Route {
			if httpRoute.Destination.Subset == subset {
				utils.Info(fmt.Sprintf("Setting %d of traffic routing to subset '%s'", percentage, subset), cid)
				httpRoute.Weight = percentage
			}
		}
	}

	err = UpdateVirtualService(cid, v.Namespace, vs)
	if err != nil {
		return err
	}

	return nil
}

// Headers set headers as routing-match strategy for istio resources
func (v IstioValues) SetHeaders(cid string, labels map[string]string, headers map[string]string) (subset string, error error) {
	var subsetRouteExists bool

	sanitizedVersion, err := SanitizeVersionString(v.Version)
	if err != nil {
		return "", err
	}

	subsetRuleName := fmt.Sprintf("%s-%d", sanitizedVersion, v.Build)

	vss, drs, err := GetResourcesToUpdate(cid, v, labels)
	if err != nil {
		return "", err
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
					utils.Warn(fmt.Sprintf("Subset '%s' already created for vs '%s", subsetRuleName, vs.Name), cid)
					subsetRouteExists = true
				}
			}
		}

		// if a subset does not exists in the current VirtualService, create it from scratch
		if !subsetRouteExists {
			// create it
			newRoute, err := CreateNewVirtualServiceHttpRoute(cid, labels, "hostname", subsetRuleName, 8080)
			if err != nil {
				utils.Fatal(fmt.Sprintf("Could not create local httpRoute object for virtualservice '%s' due to error '%s'", vs.Name, err), cid)
			}
			vs.Spec.Http = append(vs.Spec.Http, newRoute)
			err = UpdateVirtualService(cid, v.Namespace, &vs)
			if err != nil {
				utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", vs.Name, err), cid)
				return "", err
			}
		}
	}

	return subsetRuleName, nil
}

// SetLabelsDestinationRule set
func (v IstioValues) SetLabelsDestinationRule(cid string, name string, labels map[string]string) error {
	dr, err := GetDestinationRule(cid, name, v.Namespace, metav1.GetOptions{})
	if err != nil {
		utils.Fatal(fmt.Sprintf("Could not find destination rule '%s", name), cid)
		return err
	}

	if dr.Labels == nil {
		dr.Labels = map[string]string{}
	}

	for labelKey, labelValue := range labels {
		dr.Labels[labelKey] = labelValue
	}

	utils.Info(fmt.Sprintf("Setting labels '%s' to destination rule '%s'...", labels, dr.Name), cid)

	err = UpdateDestinationRule(cid, v.Namespace, dr)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Could not update destination rule '%s", dr.Name), cid)
		return err
	}

	return nil
}

func (v IstioValues) SetLabelsVirtualService(cid string, name string, labels map[string]string) error {

	vs, err := GetVirtualService(cid, name, v.Namespace, metav1.GetOptions{})
	if err != nil {
		utils.Fatal(fmt.Sprintf("Could not find virtualService '%s' due to error '%s'", name, err), cid)
		return err
	}

	if vs.Labels == nil {
		vs.Labels = map[string]string{}
	}

	for labelKey, labelValue := range labels {
		vs.Labels[labelKey] = labelValue
	}

	utils.Info(fmt.Sprintf("Setting labels '%s' to virtualService '%s'...", labels, vs.Name), cid)

	err = UpdateVirtualService(cid, v.Namespace, vs)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Could not update virtualService '%s', due to error '%s'", vs.Name, err), cid)
		return err
	}

	return nil
}
