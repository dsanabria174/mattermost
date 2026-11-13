package model

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func TestIsValidEmojiName_Unidad(t *testing.T) {
	casos := []struct {
		nombre  string
		entrada string
		errorID string 
	}{
		{"P1 nombre vacío", "", "model.emoji.name.app_error"},
		{"P2 nombre de 65 caracteres", strings.Repeat("a", 65), "model.emoji.name.app_error"},
		{"P3 carácter no permitido", "mi emoji!", "model.emoji.name.app_error"},
		{"P4 nombre de emoji del sistema", "smile", "model.emoji.system_emoji_name.app_error"},
		{"P5 nombre válido", "mi_emoji-1", ""},
		{"límite: 64 caracteres es válido", strings.Repeat("a", 64), ""},
		{"acepta el signo +", "mi+emoji", ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)

			appErr := IsValidEmojiName(c.entrada)

			if c.errorID == "" {
				g.Expect(appErr).To(gomega.BeNil())
				return
			}
			esperarAppErrorUnidad(g, appErr, c.errorID)
		})
	}
}

func emojiValidoUnidad() *Emoji {
	return &Emoji{
		Id:        NewId(),
		CreateAt:  GetMillis(),
		UpdateAt:  GetMillis(),
		CreatorId: NewId(),
		Name:      "mi_emoji",
	}
}

func TestEmojiIsValid_Unidad(t *testing.T) {
	casos := []struct {
		nombre  string
		mutar   func(e *Emoji)
		errorID string
	}{
		{"emoji válido", func(e *Emoji) {}, ""},
		{"Id inválido", func(e *Emoji) { e.Id = "abc" }, "model.emoji.id.app_error"},
		{"sin fecha de creación", func(e *Emoji) { e.CreateAt = 0 }, "model.emoji.create_at.app_error"},
		{"sin fecha de actualización", func(e *Emoji) { e.UpdateAt = 0 }, "model.emoji.update_at.app_error"},
		{"CreatorId de 27 caracteres", func(e *Emoji) { e.CreatorId = strings.Repeat("a", 27) }, "model.emoji.user_id.app_error"},
		{"límite: CreatorId de 26 caracteres", func(e *Emoji) { e.CreatorId = strings.Repeat("a", 26) }, ""},
		{"CreatorId vacío (comportamiento observado)", func(e *Emoji) { e.CreatorId = "" }, ""},
		{"nombre inválido se delega en IsValidEmojiName", func(e *Emoji) { e.Name = "" }, "model.emoji.name.app_error"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)
			emoji := emojiValidoUnidad()
			c.mutar(emoji)

			appErr := emoji.IsValid()

			if c.errorID == "" {
				g.Expect(appErr).To(gomega.BeNil())
				return
			}
			esperarAppErrorUnidad(g, appErr, c.errorID)
		})
	}
}

func TestEmojiPreSave_Unidad(t *testing.T) {
	t.Run("asigna un Id cuando está vacío", func(t *testing.T) {
		g := gomega.NewWithT(t)
		emoji := &Emoji{Name: "mi_emoji"}

		emoji.PreSave()

		g.Expect(IsValidId(emoji.Id)).To(gomega.BeTrue())
	})

	t.Run("conserva el Id cuando ya existe", func(t *testing.T) {
		g := gomega.NewWithT(t)
		id := NewId()
		emoji := &Emoji{Id: id, Name: "mi_emoji"}

		emoji.PreSave()

		g.Expect(emoji.Id).To(gomega.Equal(id))
	})

	t.Run("pasa el nombre a minúsculas", func(t *testing.T) {
		g := gomega.NewWithT(t)
		emoji := &Emoji{Name: "MiEmoji"}

		emoji.PreSave()

		g.Expect(emoji.Name).To(gomega.Equal("miemoji"))
	})

	t.Run("fechas de creación y actualización iguales y mayores que cero", func(t *testing.T) {
		g := gomega.NewWithT(t)
		emoji := &Emoji{Name: "mi_emoji"}

		emoji.PreSave()

		g.Expect(emoji.CreateAt).To(gomega.BeNumerically(">", 0))
		g.Expect(emoji.UpdateAt).To(gomega.Equal(emoji.CreateAt))
	})
}
