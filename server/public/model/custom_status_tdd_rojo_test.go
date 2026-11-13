package model

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestRecentCustomStatusesAdd_NoModificaElOriginal_TDD(t *testing.T) {
	g := gomega.NewWithT(t)
	original := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b"), estadoUnidad("c")}
	copia := append(RecentCustomStatuses{}, original...)
	nuevo := CustomStatus{Emoji: "tada", Text: "b"}

	_ = original.Add(&nuevo)

	g.Expect(original).To(gomega.Equal(copia))
}

func TestRecentCustomStatusesRemove_NoModificaElOriginal_TDD(t *testing.T) {
	g := gomega.NewWithT(t)
	original := RecentCustomStatuses{estadoUnidad("a"), estadoUnidad("b"), estadoUnidad("c")}
	copia := append(RecentCustomStatuses{}, original...)
	aQuitar := estadoUnidad("a")

	_, err := original.Remove(&aQuitar)

	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(original).To(gomega.Equal(copia))
}
