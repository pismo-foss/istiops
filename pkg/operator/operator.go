package operator

type Shift struct {
	Port     uint32
	Hostname string
	Selector *Selector
	Headers  map[string]string
	Weight   int32
}

type Operator interface {
	Create(s *Shift) error
	Delete(s *Shift) error
	Update(s *Shift) error
	Clear(map[string]string) error
}
