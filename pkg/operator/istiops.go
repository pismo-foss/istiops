package operator

import (
	"fmt"

	"github.com/pismo/istiops/pkg/router"
	"github.com/pismo/istiops/utils"
)

type Istiops struct {
	DrRouter router.Router
	VsRouter router.Router
}

func (ips *Istiops) Get(r *router.Shift) error {
	return nil
}

func (ips *Istiops) Create(r *router.Shift) error {

	VsRouter := ips.VsRouter
	err := VsRouter.Validate(r)
	if err != nil {
		return err
	}

	DrRouter := ips.DrRouter
	err = DrRouter.Validate(r)
	if err != nil {
		return err
	}

	return nil
}

func (ips *Istiops) Delete(r *router.Shift) error {
	VsRouter := ips.VsRouter
	err := VsRouter.Delete(r)
	if err != nil {
		return err
	}

	DrRouter := ips.DrRouter
	err = DrRouter.Delete(r)
	if err != nil {
		return err
	}

	return nil
}

// Migrar struct Shift para operator e criar struct Route no router, isso vai desacoplar totalmente os pacotes
// toda a dependencia entre os pacotes deve ser apenas de interfaces e suas dependencias
func (ips *Istiops) Update(r *router.Shift) error {
	//nao vale a pena gerar erros mais especificos? fica melhor documentado
	if len(r.Selector.Labels) == 0 || len(r.Traffic.PodSelector) == 0 {
		utils.Fatal(fmt.Sprintf("Selectors must not be empty otherwise istiops won't be able to find any resources."), "")
	}

	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	//pegar os dados do Shift e jogar no Route para fazer a chamada
	err = DrRouter.Validate(r)
	if err != nil {
		return err
	}
	err = DrRouter.Update(r)
	if err != nil {
		return err
	}

	// melhor validar os 2 antes de executar, se der erro na validacao de virtualservice vai ficar destinationrule solta la
	err = VsRouter.Validate(r)
	if err != nil {
		return err
	}

	// se der erro aqui precisa remover o destination rule que foi criado pra nao deixar sujeira
	err = VsRouter.Update(r)
	if err != nil {
		return err
	}

	if r.Traffic.Weight > 0 {
		// update router to serve percentage
		//if err != nil {
		//	utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), "")
		//}
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(s *router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Clear(s)
	if err != nil {
		return err
	}

	err = VsRouter.Clear(s)
	if err != nil {
		return err
	}

	// Clean dr rules ?

	return nil
}
