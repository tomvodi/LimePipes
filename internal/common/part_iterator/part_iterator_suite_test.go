package part_iterator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPartIterator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PartIterator Suite")
}
