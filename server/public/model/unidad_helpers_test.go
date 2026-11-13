package model

import (
	"net/http"

	"github.com/onsi/gomega"
)

func esperarAppErrorUnidad(g *gomega.WithT, appErr *AppError, errorID string) {
	g.Expect(appErr).NotTo(gomega.BeNil())
	g.Expect(appErr).To(gomega.HaveField("Id", errorID))
	g.Expect(appErr).To(gomega.HaveField("StatusCode", http.StatusBadRequest))
}
