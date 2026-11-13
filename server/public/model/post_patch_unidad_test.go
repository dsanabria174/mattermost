package model

import (
	"testing"

	"github.com/onsi/gomega"
)

func punteroUnidad[T any](v T) *T { return &v }

func propsQuitarVistaPreviaUnidad(valor any) *StringInterface {
	props := StringInterface{PostPropsRemoveLinkPreview: valor}
	return &props
}

func TestPostPatchIsOnlyRemoveLinkPreview_Unidad(t *testing.T) {
	archivos := StringArray{"abc"}

	casos := []struct {
		nombre   string
		parche   *PostPatch
		esperado bool
	}{
		{"P1 sin Props", &PostPatch{}, false},
		{"P2 con Message", &PostPatch{Message: punteroUnidad("hola"), Props: propsQuitarVistaPreviaUnidad("true")}, false},
		{"P3 con FileIds", &PostPatch{FileIds: &archivos, Props: propsQuitarVistaPreviaUnidad("true")}, false},
		{"P4 que fija el mensaje (IsPinned)", &PostPatch{IsPinned: punteroUnidad(true), Props: propsQuitarVistaPreviaUnidad("true")}, false},
		{"P5 con HasReactions", &PostPatch{HasReactions: punteroUnidad(true), Props: propsQuitarVistaPreviaUnidad("true")}, false},
		{"P6 solo quita la vista previa", &PostPatch{Props: propsQuitarVistaPreviaUnidad("true")}, true},
		{"P6 variante: valor distinto de true", &PostPatch{Props: propsQuitarVistaPreviaUnidad("false")}, false},
		{"P6 variante: valor que no es cadena", &PostPatch{Props: propsQuitarVistaPreviaUnidad(true)}, false},
		{"desfijar tampoco cuenta como quitar la vista previa", &PostPatch{IsPinned: punteroUnidad(false), Props: propsQuitarVistaPreviaUnidad("true")}, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)

			resultado := c.parche.IsOnlyRemoveLinkPreview()

			g.Expect(resultado).To(gomega.Equal(c.esperado))
		})
	}
}
