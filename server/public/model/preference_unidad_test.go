package model

import (
	"strconv"
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func preferenciaValidaUnidad() *Preference {
	return &Preference{
		UserId:   NewId(),
		Category: PreferenceCategoryFlaggedPost,
		Name:     NewId(),
		Value:    "true",
	}
}

func TestPreferenceIsValid_Unidad(t *testing.T) {
	casos := []struct {
		nombre  string
		mutar   func(p *Preference)
		errorID string
	}{
		{"P1 guardar un mensaje (flagged_post)", func(p *Preference) {}, ""},
		{"P2 UserId inválido", func(p *Preference) { p.UserId = "abc" }, "model.preference.is_valid.id.app_error"},
		{"P3 categoría vacía", func(p *Preference) { p.Category = "" }, "model.preference.is_valid.category.app_error"},
		{"P4 categoría de 33 caracteres", func(p *Preference) { p.Category = strings.Repeat("c", 33) }, "model.preference.is_valid.category.app_error"},
		{"P5 nombre de 33 caracteres", func(p *Preference) { p.Name = strings.Repeat("n", 33) }, "model.preference.is_valid.name.app_error"},
		{"P6 valor con una runa de más", func(p *Preference) { p.Value = strings.Repeat("v", MaxPreferenceValueLength+1) }, "model.preference.is_valid.value.app_error"},
		{"P7 tema con JSON inválido", func(p *Preference) { p.Category = PreferenceCategoryTheme; p.Value = "no-json" }, "model.preference.is_valid.theme.app_error"},
		{"P8 tema con JSON válido", func(p *Preference) { p.Category = PreferenceCategoryTheme; p.Value = "{}" }, ""},
		{"P9 barra lateral con otro nombre", func(p *Preference) { p.Category = PreferenceCategorySidebarSettings; p.Name = "otro" }, ""},
		{"P10 límite de DMs no numérico", func(p *Preference) { limiteDMsUnidad(p, "abc") }, "model.preference.is_valid.limit_visible_dms_gms.app_error"},
		{"P11 límite de DMs menor que 1", func(p *Preference) { limiteDMsUnidad(p, "0") }, "model.preference.is_valid.limit_visible_dms_gms.app_error"},
		{"P12 límite de DMs por encima del máximo", func(p *Preference) {
			limiteDMsUnidad(p, strconv.Itoa(PreferenceMaxLimitVisibleDmsGmsValue+1))
		}, "model.preference.is_valid.limit_visible_dms_gms.app_error"},
		{"P13 límite de DMs válido", func(p *Preference) { limiteDMsUnidad(p, "10") }, ""},
		{"límite: categoría de 32 caracteres es válida", func(p *Preference) { p.Category = strings.Repeat("c", 32) }, ""},
		{"límite: nombre de 32 caracteres es válido", func(p *Preference) { p.Name = strings.Repeat("n", 32) }, ""},
		{"límite: valor de longitud máxima es válido", func(p *Preference) { p.Value = strings.Repeat("v", MaxPreferenceValueLength) }, ""},
		{"límite: límite de DMs igual a 1 es válido", func(p *Preference) { limiteDMsUnidad(p, "1") }, ""},
		{"límite: límite de DMs igual al máximo es válido", func(p *Preference) {
			limiteDMsUnidad(p, strconv.Itoa(PreferenceMaxLimitVisibleDmsGmsValue))
		}, ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)
			preferencia := preferenciaValidaUnidad()
			c.mutar(preferencia)

			appErr := preferencia.IsValid()

			if c.errorID == "" {
				g.Expect(appErr).To(gomega.BeNil())
				return
			}
			esperarAppErrorUnidad(g, appErr, c.errorID)
		})
	}
}

func limiteDMsUnidad(p *Preference, valor string) {
	p.Category = PreferenceCategorySidebarSettings
	p.Name = PreferenceLimitVisibleDmsGms
	p.Value = valor
}
