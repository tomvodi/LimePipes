package cmd

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	pmocks "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces/mocks"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/cmd/limepipes-cli/importtype"
	"github.com/tomvodi/limepipes/internal/interfaces/mocks"
	"github.com/tomvodi/limepipes/internal/utils"
	"os"
	"path/filepath"
	"testing"
)

var _ = Describe("FileProcessor", func() {
	var err error
	var afs afero.Fs
	var opts *Options
	var pfo *ProcessFilesOptions
	var fp *FileProcessor
	var pl *mocks.PluginLoader
	var filePlug *pmocks.LimePipesPlugin
	var ds *mocks.DataService
	var parsedTunes []*messages.ParsedTune

	BeforeEach(func() {
		afs = afero.NewMemMapFs()
		pfo = &ProcessFilesOptions{}
		opts = &Options{}
		pl = mocks.NewPluginLoader(GinkgoT())
		filePlug = pmocks.NewLimePipesPlugin(GinkgoT())
		ds = mocks.NewDataService(GinkgoT())
		fp = NewFileProcessor(afs, pl, ds)
	})

	JustBeforeEach(func() {
		err = fp.ProcessFiles(pfo, opts)
	})

	Context("when no import types are given", func() {
		BeforeEach(func() {
			opts.ImportTypes = []string{}
		})

		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when no valid import type was given", func() {
		BeforeEach(func() {
			opts.ImportTypes = []string{"not-valid"}
		})

		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when a valid import type was given", func() {
		BeforeEach(func() {
			opts.ImportTypes = []string{
				importtype.FromFileFormat(fileformat.Format_BWW),
			}
			opts.Verbose = true
			pl.EXPECT().FileExtensionsForFileFormat(fileformat.Format_BWW).
				Return([]string{".bww"}, nil).Times(1)
		})

		When("files should be moved to output dir", func() {
			BeforeEach(func() {
				opts.OutputDir = "output"
				pfo.MoveToOutputDir = true
			})

			It("should have created that output dir", func() {
				_, err = afs.Stat(opts.OutputDir)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		When("having one file in a directory", func() {
			BeforeEach(func() {
				pfo.ArgPaths = []string{"testdata"}
				file, err := afs.Create("testdata/tune1.bww")
				Expect(err).NotTo(HaveOccurred())
				_, err = file.Write([]byte("tune1.bww testdata"))
				Expect(err).NotTo(HaveOccurred())
			})

			When("there is no plugin for that file format", func() {
				BeforeEach(func() {
					pl.EXPECT().PluginForFileExtension(".bww").
						Return(nil, fmt.Errorf("no plugin"))
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			When("there is a plugin for that file format", func() {
				BeforeEach(func() {
					pl.EXPECT().PluginForFileExtension(".bww").
						Return(filePlug, nil)
				})

				When("the plugin returns an import error", func() {
					BeforeEach(func() {
						filePlug.EXPECT().Parse([]byte("tune1.bww testdata")).
							Return(nil, fmt.Errorf("import error"))
					})

					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					When("failed files should be skipped", func() {
						BeforeEach(func() {
							opts.SkipFailedFiles = true
						})

						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})
					})
				})

				When("the plugin returns a tune", func() {
					BeforeEach(func() {
						parsedTunes = []*messages.ParsedTune{
							{
								Tune: &tune.Tune{
									Title: "tune1",
								},
								TuneFileData: []byte("tune1.bww testdata"),
							},
						}
						filePlug.EXPECT().Parse([]byte("tune1.bww testdata")).
							Return(parsedTunes, nil)
					})

					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					When("files should be moved to output dir", func() {
						BeforeEach(func() {
							opts.OutputDir = "output"
							pfo.MoveToOutputDir = true
						})

						It("should have moved the file to the output dir", func() {
							_, err = afs.Stat(opts.OutputDir + "/tune1.bww")
							Expect(err).NotTo(HaveOccurred())
						})
					})

					When("tunes should be imported", func() {
						BeforeEach(func() {
							pfo.ImportToDb = true
						})

						When("there is no file format for the file extension", func() {
							BeforeEach(func() {
								pl.EXPECT().FileFormatForFileExtension(".bww").
									Return(fileformat.Format_Unknown, fmt.Errorf("no file format"))
							})

							When("failed files should be skipped", func() {
								BeforeEach(func() {
									opts.SkipFailedFiles = true
								})

								It("should not return an error", func() {
									Expect(err).NotTo(HaveOccurred())
								})
							})

							It("should return an error", func() {
								Expect(err).To(HaveOccurred())
							})
						})

						When("the file format is known", func() {
							BeforeEach(func() {
								pl.EXPECT().FileFormatForFileExtension(".bww").
									Return(fileformat.Format_BWW, nil)
							})

							When("the tune can not be imported to database", func() {
								BeforeEach(func() {
									ds.EXPECT().ImportTunes(parsedTunes, mock.Anything).
										Return(nil, nil, fmt.Errorf("import error"))
								})

								It("should return an error", func() {
									Expect(err).To(HaveOccurred())
								})

								When("failed files should be skipped", func() {
									BeforeEach(func() {
										opts.SkipFailedFiles = true
									})

									It("should not return an error", func() {
										Expect(err).NotTo(HaveOccurred())
									})
								})
							})

							When("the tune can be imported to database", func() {
								BeforeEach(func() {
									ds.EXPECT().ImportTunes(parsedTunes, mock.Anything).
										Return(nil, nil, nil)
								})

								It("should not return an error", func() {
									Expect(err).NotTo(HaveOccurred())
								})
							})
						})
					})
				})
			})
		})
	})
})

func Test_getAllFilesFromPaths(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		files   []string
		paths   []string
		opts    *GetFilesOptions
		want    []string
		wantErr bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "no file extensions given",
			prepare: func(f *fields) {
				f.opts = &GetFilesOptions{
					FileExtensions: nil,
				}
				f.wantErr = true
			},
		},
		{
			name: "get all Bagpipe Music Writer files from a directory",
			prepare: func(f *fields) {
				f.opts = &GetFilesOptions{
					Verbose:        true,
					FileExtensions: []string{".bww", ".bmw"},
				}
				f.files = []string{
					"testdata/tune1.bww", // Bagpipe Player file
					"testdata/tune2.bmw", // Bagpipe Music Writer file
					"testdata/test.txt",
				}
				f.paths = []string{"testdata"}
				f.want = []string{
					"testdata/tune1.bww",
					"testdata/tune2.bmw",
				}
			},
		},
		{
			name: "get all Bagpipe Music Writer files from a directory recursively",
			prepare: func(f *fields) {
				f.opts = &GetFilesOptions{
					Recursive:      true,
					FileExtensions: []string{".bww", ".bmw"},
				}
				f.paths = []string{"testdata"}
				f.files = []string{
					"testdata/tune1.bww",
					"testdata/subdir/tune2.bww",
					"testdata/subdir/tune3.bww",
				}
				f.want = []string{
					"testdata/subdir/tune2.bww",
					"testdata/subdir/tune3.bww",
					"testdata/tune1.bww",
				}
			},
		},
		{
			name: "get single files passed",
			prepare: func(f *fields) {
				f.opts = &GetFilesOptions{
					FileExtensions: []string{".bww", ".bmw"},
				}
				f.files = []string{
					"testdata/tune1.bww",
					"testdata/test.txt",
				}
				f.paths = []string{
					"testdata/tune1.bww",
					"testdata/test.txt",
				}
				f.want = []string{
					"testdata/tune1.bww",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			afs := afero.NewMemMapFs()
			for _, file := range f.files {
				err := afs.MkdirAll(filepath.Dir(file), os.ModePerm)
				g.Expect(err).To(BeNil())
				_, err = afs.Create(file)
				g.Expect(err).To(BeNil())
			}

			got, err := getAllFilesFromPaths(
				afs,
				f.paths,
				f.opts,
			)
			if (err != nil) != f.wantErr {
				t.Errorf("getAllFilesFromPaths() error = %v, wantErr %v", err, f.wantErr)
				return
			}
			g.Expect(got).To(BeComparableTo(f.want))
		})
	}
}
