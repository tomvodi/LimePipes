package test

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func DataFromFile(filePath string) []byte {
	bwwFile, err := os.Open(filePath)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	var data []byte
	data, err = io.ReadAll(bwwFile)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	return data
}

func ExportToYaml(muMo music_model.MusicModel, filePath string) {
	data, err := yaml.Marshal(muMo)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = os.WriteFile(filePath, data, 0664)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func ImportFromYaml(filePath string, embExpander interfaces.EmbellishmentExpander) music_model.MusicModel {
	muMo := make(music_model.MusicModel, 0)
	fileData, err := os.ReadFile(filePath)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = yaml.Unmarshal(fileData, &muMo)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	if embExpander != nil {
		embExpander.ExpandModel(muMo)
	}

	return muMo
}
