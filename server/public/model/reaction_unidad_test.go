package model

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func reaccionValidaUnidad() *Reaction {
	return &Reaction{
		UserId:    NewId(),
		PostId:    NewId(),
		EmojiName: "+1",
		CreateAt:  GetMillis(),
		UpdateAt:  GetMillis(),
	}
}

func TestReactionIsValid_Unidad(t *testing.T) {
	casos := []struct {
		nombre  string
		mutar   func(r *Reaction)
		errorID string
	}{
		{"P1 UserId inválido", func(r *Reaction) { r.UserId = "abc" }, "model.reaction.is_valid.user_id.app_error"},
		{"P2 PostId inválido", func(r *Reaction) { r.PostId = "abc" }, "model.reaction.is_valid.post_id.app_error"},
		{"P3 emoji vacío", func(r *Reaction) { r.EmojiName = "" }, "model.reaction.is_valid.emoji_name.app_error"},
		{"P4 emoji de 65 caracteres", func(r *Reaction) { r.EmojiName = strings.Repeat("a", 65) }, "model.reaction.is_valid.emoji_name.app_error"},
		{"P5 emoji con carácter no permitido", func(r *Reaction) { r.EmojiName = "mi emoji!" }, "model.reaction.is_valid.emoji_name.app_error"},
		{"P6 sin fecha de creación", func(r *Reaction) { r.CreateAt = 0 }, "model.reaction.is_valid.create_at.app_error"},
		{"P7 sin fecha de actualización", func(r *Reaction) { r.UpdateAt = 0 }, "model.reaction.is_valid.update_at.app_error"},
		{"P8 reacción válida", func(r *Reaction) {}, ""},
		{"límite: emoji de 64 caracteres es válido", func(r *Reaction) { r.EmojiName = strings.Repeat("a", 64) }, ""},
		{"acepta mayúsculas, guion y guion bajo", func(r *Reaction) { r.EmojiName = "Mi-Emoji_1" }, ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			g := gomega.NewWithT(t)
			reaccion := reaccionValidaUnidad()
			c.mutar(reaccion)

			appErr := reaccion.IsValid()

			if c.errorID == "" {
				g.Expect(appErr).To(gomega.BeNil())
				return
			}
			esperarAppErrorUnidad(g, appErr, c.errorID)
		})
	}
}

func TestReactionPreSave_Unidad(t *testing.T) {
	t.Run("asigna CreateAt cuando es cero", func(t *testing.T) {
		g := gomega.NewWithT(t)
		reaccion := &Reaction{}

		reaccion.PreSave()

		g.Expect(reaccion.CreateAt).To(gomega.BeNumerically(">", 0))
	})

	t.Run("conserva CreateAt cuando ya existe", func(t *testing.T) {
		g := gomega.NewWithT(t)
		reaccion := &Reaction{CreateAt: 123}

		reaccion.PreSave()

		g.Expect(reaccion.CreateAt).To(gomega.Equal(int64(123)))
		g.Expect(reaccion.UpdateAt).To(gomega.BeNumerically(">", 123))
	})

	t.Run("reinicia DeleteAt e inicializa RemoteId", func(t *testing.T) {
		g := gomega.NewWithT(t)
		reaccion := &Reaction{DeleteAt: 999}

		reaccion.PreSave()

		g.Expect(reaccion.DeleteAt).To(gomega.BeZero())
		g.Expect(reaccion.RemoteId).NotTo(gomega.BeNil())
		g.Expect(reaccion.GetRemoteID()).To(gomega.BeEmpty())
	})
}
