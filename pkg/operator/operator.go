package operator

type Operator interface {
	Create(ir *IstioRoute) error
	Delete(ir *IstioRoute) error
	Update(ir *IstioRoute) error
	Clear(map[string]string) error
}
