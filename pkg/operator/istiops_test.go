package operator

import (
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		r    Route
		ips  Istiops
		want error
	}{
		"Missing port in route": {
			route: Route{
				//Port:     5000,
				Hostname: "api-xpto.domain.io",
				Selector: operator.Selector{
					Labels: map[string]string{"environment": "pipeline-go"},
				},
				Headers: map[string]string{
					"x-version": "PR-141",
					"x-cid":     "blau",
				},
				Weight: 0,
			},
			istiops: Istiops{
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

type mockVirtualSerivce struct {
	cool     string
	validate func() (VirtualService, error)
	update   func() (VirtualService, error)
	delete   func() (VirtualService, error)
}

func (mvs *mockVirtualService) Validate(r Route) VirtualService { return m.validate() }
func (mvs *mockVirtualService) Update(r Route) VirtualService   { return m.update() }
func (mvs *mockVirtualService) Delete(r Route) VirtualService   { return m.delete() }

type mockDestinationRule struct {
	validate func() (DestinationRule, error)
	update   func() (DestinationRule, error)
	delete   func() (DestinationRule error)
}

func (mvs *mockDestinationRule) Validate(r Route) DestinationRule { return m.validate() }
func (mvs *mockDestinationRule) Update(r Route) DestinationRule   { return m.update() }
func (mvs *mockDestinationRule) Delete(r Route) DestinationRule   { return m.delete() }
