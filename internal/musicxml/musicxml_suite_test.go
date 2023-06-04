package musicxml_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMusicxml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Musicxml Suite")
}
