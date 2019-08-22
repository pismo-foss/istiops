package operator

import (
	"errors"
	"testing"

	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		r    router.Route
		ips  Istiops
		want error
	}{
		"Missing port in route": {
			r: router.Route{
				//Port:     5000,
				Hostname: "api-xpto.domain.io",
				Selector: &router.Selector{
					router.Labels: map[string]string{"environment": "pipeline-go"},
				},
				Headers: map[string]string{
					"x-version": "PR-141",
					"x-cid":     "blau",
				},
				Weight: 0,
			},
			istiops: operator.Istiops{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      26,
				VirtualService: &mockVirtualService{
					validate: func() (VirtualService, error) {
						return nil, errors.New("Missing port for VirtualService spec")
					},
				},
			},
			want: errors.New("Missing port for VirtualService spec"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.istiops.Create(tc.route)
			if !equalError(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

type mockVirtualService struct {
	cool     string
	validate func() (router.Router, error)
	update   func() (router.Router, error)
	delete   func() (router.Router, error)
}

func (mvs *mockVirtualService) Update(s router.Shift) (mockVirtualService, error) { return mvs.update() }
func (mvs *mockVirtualService) Delete(r router.Shift) (mockVirtualService, error) { return mvs.delete() }
func (mvs *mockVirtualService) Validate(r router.Shift) (mockVirtualService, error) {
	return mvs.validate()
}

type mockDestinationRule struct {
	validate func() (DestinationRule, error)
	update   func() (DestinationRule, error)
	delete   func() (DestinationRule, error)
}

func (mvs *mockDestinationRule) Validate(r Route) DestinationRule { return m.validate() }
func (mvs *mockDestinationRule) Update(r Route) DestinationRule   { return m.update() }
func (mvs *mockDestinationRule) Delete(r Route) DestinationRule   { return m.delete() }
