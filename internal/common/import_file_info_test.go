package common

import (
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes/internal/utils"
	"testing"
)

func Test_NewImportFileInfoFromLocalFile(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		originalPath string
		want         *ImportFileInfo
		wantErr      bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "import file info",
			prepare: func(f *fields) {
				f.originalPath = "./testfiles/test.bww"
				f.want = &ImportFileInfo{
					OriginalPath: "./testfiles/test.bww",
					Name:         "test",
					Hash:         "534b1d50f10ee4ea30604ce01660e2429682fe6e53a4ef6a9d01c835ef73b866",
					Data:         []byte(`Bagpipe Reader:1.0`),
					FileFormat:   fileformat.Format_BWW,
				}
				f.wantErr = false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			fileInfo, err := NewImportFileInfoFromLocalFile(
				afero.NewOsFs(),
				f.originalPath,
				fileformat.Format_BWW,
			)
			if f.wantErr {
				g.Expect(err).Should(HaveOccurred())
			} else {
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(fileInfo).To(Equal(f.want))
			}
		})
	}
}

func Test_NewImportFileInfo(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		fileName string
		fileData []byte
		want     *ImportFileInfo
		wantErr  bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "import file info",
			prepare: func(f *fields) {
				f.fileName = "test.bww"
				f.fileData = []byte(`Bagpipe Reader:1.0`)
				f.want = &ImportFileInfo{
					OriginalPath: f.fileName,
					Name:         "test",
					Hash:         "534b1d50f10ee4ea30604ce01660e2429682fe6e53a4ef6a9d01c835ef73b866",
					Data:         f.fileData,
					FileFormat:   fileformat.Format_BWW,
				}
				f.wantErr = false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			fileInfo, err := NewImportFileInfo(f.fileName, fileformat.Format_BWW, f.fileData)
			if f.wantErr {
				g.Expect(err).Should(HaveOccurred())
			} else {
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(fileInfo).To(Equal(f.want))
			}
		})
	}
}
