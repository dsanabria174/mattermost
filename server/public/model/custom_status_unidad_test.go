package model

import (
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

func TestCustomStatusDuracionYExpiracion_Unidad(t *testing.T) {
	futuro := time.Now().Add(1 * time.Hour)
	pasado := time.Now().Add(-1 * time.Hour)

	casos := []struct {
		nombre   string
		estado   CustomStatus
		esperado bool
	}{
		{"P1 sin duración ni expiración", CustomStatus{}, true},
		{"P2 sin duración, expiración futura", CustomStatus{ExpiresAt: futuro}, true},
		{"P3 sin duración, expiración pasada", CustomStatus{ExpiresAt: pasado}, false},
		{"P4 duración válida, expiración futura", CustomStatus{Duration: "one_hour", ExpiresAt: futuro}, true},
		{"P5 duración válida, expiración pasada", CustomStatus{Duration: "one_hour", ExpiresAt: pasado}, false},
		{"P6 duración inexistente", CustomStatus{Duration: "una_semana", ExpiresAt: futuro}, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)

			valido := c.estado.AreDurationAndExpirationTimeValid()

			g.Expect(valido).To(gomega.Equal(c.esperado))
		})
	}

	t.Run("todas las duraciones permitidas son válidas con expiración futura", func(t *testing.T) {
		for _, duracion := range []string{"thirty_minutes", "one_hour", "four_hours", "today", "this_week", "date_and_time"} {
			g := gomega.NewWithT(t)
			estado := CustomStatus{Duration: duracion, ExpiresAt: futuro}

			valido := estado.AreDurationAndExpirationTimeValid()

			g.Expect(valido).To(gomega.BeTrue(), "la duración %q debería ser válida", duracion)
		}
	})
}

func TestCustomStatusPreSave_Unidad(t *testing.T) {
	futuro := time.Now().Add(1 * time.Hour)
	pasado := time.Now().Add(-1 * time.Hour)

	t.Run("sin duración y con expiración futura pasa a date_and_time", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{ExpiresAt: futuro}

		estado.PreSave()

		g.Expect(estado.Duration).To(gomega.Equal("date_and_time"))
	})

	t.Run("sin duración y con expiración pasada no cambia", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{ExpiresAt: pasado}

		estado.PreSave()

		g.Expect(estado.Duration).To(gomega.BeEmpty())
	})

	t.Run("sin duración y sin expiración no cambia", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{}

		estado.PreSave()

		g.Expect(estado.Duration).To(gomega.BeEmpty())
	})

	t.Run("con duración definida no cambia", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{Duration: "one_hour", ExpiresAt: futuro}

		estado.PreSave()

		g.Expect(estado.Duration).To(gomega.Equal("one_hour"))
	})

	t.Run("límite: texto de 100 caracteres se conserva aunque tenga más de 100 bytes", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{Text: strings.Repeat("ñ", 100)}

		estado.PreSave()

		g.Expect([]rune(estado.Text)).To(gomega.HaveLen(100))
	})

	t.Run("texto de 101 caracteres se recorta a 100", func(t *testing.T) {
		g := gomega.NewWithT(t)
		estado := CustomStatus{Text: strings.Repeat("ñ", 101)}

		estado.PreSave()

		g.Expect([]rune(estado.Text)).To(gomega.HaveLen(100))
	})
}

func estadoUnidad(texto string) CustomStatus {
	return CustomStatus{Emoji: "smile", Text: texto}
}

func textosUnidad(rcs RecentCustomStatuses) []string {
	textos := make([]string, 0, len(rcs))
	for _, s := range rcs {
		textos = append(textos, s.Text)
	}
	return textos
}

func TestRecentCustomStatuses_Unidad(t *testing.T) {
	t.Run("Add agrega el nuevo estado al frente", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("b"), estadoUnidad("c")}
		nuevo := estadoUnidad("a")

		resultado := recientes.Add(&nuevo)

		g.Expect(textosUnidad(resultado)).To(gomega.Equal([]string{"a", "b", "c"}))
	})

	t.Run("Add mueve al frente el estado con el mismo texto sin duplicarlo", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b"), estadoUnidad("c")}
		nuevo := CustomStatus{Emoji: "tada", Text: "b"}

		resultado := recientes.Add(&nuevo)

		g.Expect(textosUnidad(resultado)).To(gomega.Equal([]string{"b", "a", "c"}))
		g.Expect(resultado[0].Emoji).To(gomega.Equal("tada"))
	})

	t.Run("Add descarta el más antiguo al superar el máximo de 5", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{
			estadoUnidad("1"), estadoUnidad("2"), estadoUnidad("3"), estadoUnidad("4"), estadoUnidad("5"),
		}
		nuevo := estadoUnidad("6")

		resultado := recientes.Add(&nuevo)

		g.Expect(resultado).To(gomega.HaveLen(MaxRecentCustomStatuses))
		g.Expect(textosUnidad(resultado)).To(gomega.Equal([]string{"6", "1", "2", "3", "4"}))
	})

	t.Run("Contains devuelve false para un estado nulo", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a")}

		contiene, err := recientes.Contains(nil)

		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(contiene).To(gomega.BeFalse())
	})

	t.Run("Contains devuelve false para un estado vacío", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a")}
		vacio := CustomStatus{}

		contiene, err := recientes.Contains(&vacio)

		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(contiene).To(gomega.BeFalse())
	})

	t.Run("Contains encuentra un estado presente y no uno ausente", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b")}
		presente := estadoUnidad("b")
		ausente := estadoUnidad("z")

		estaPresente, err1 := recientes.Contains(&presente)
		estaAusente, err2 := recientes.Contains(&ausente)

		g.Expect(err1).NotTo(gomega.HaveOccurred())
		g.Expect(err2).NotTo(gomega.HaveOccurred())
		g.Expect(estaPresente).To(gomega.BeTrue())
		g.Expect(estaAusente).To(gomega.BeFalse())
	})

	t.Run("Remove quita el estado indicado", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b"), estadoUnidad("c")}
		aQuitar := estadoUnidad("b")

		resultado, err := recientes.Remove(&aQuitar)

		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(textosUnidad(resultado)).To(gomega.Equal([]string{"a", "c"}))
	})

	t.Run("Remove con un estado ausente no cambia la lista", func(t *testing.T) {
		g := gomega.NewWithT(t)
		recientes := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b")}
		ausente := estadoUnidad("z")

		resultado, err := recientes.Remove(&ausente)

		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(textosUnidad(resultado)).To(gomega.Equal([]string{"a", "b"}))
	})
}
