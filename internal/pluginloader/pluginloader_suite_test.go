package pluginloader_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPluginloader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pluginloader Suite")
}
