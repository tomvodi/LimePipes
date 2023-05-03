package bww_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBww(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bww Suite")
}
