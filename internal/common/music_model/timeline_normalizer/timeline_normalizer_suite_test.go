package timeline_normalizer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTimelineNormalizer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TimelineNormalizer Suite")
}
