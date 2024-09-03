package pluginloader

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	pimocks "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces/mocks"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/interfaces/mocks"
	"os"
)

var _ = Describe("Pluginloader", func() {
	var err error
	var loader *Loader
	var fs afero.Fs
	var pluginsDir string
	var pluginID string
	var lpPlugin *pimocks.LimePipesPlugin
	var processHandler *mocks.PluginProcessHandler

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		pluginsDir = "/testdir"
		pluginID = "testplugin"
		err := fs.MkdirAll(pluginsDir, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		processHandler = mocks.NewPluginProcessHandler(GinkgoT())
		lpPlugin = pimocks.NewLimePipesPlugin(GinkgoT())

		loader = &Loader{
			fs:               fs,
			processHandler:   processHandler,
			supportedPlugins: []string{pluginID},
			pluginInfos:      make(map[string]*messages.PluginInfoResponse),
		}
	})

	Context("LoadPluginsFromDir", func() {
		JustBeforeEach(func() {
			err = loader.LoadPluginsFromDir(pluginsDir)
		})

		When("loading a plugin which does not have an executable file", func() {
			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("having a plugin file to start", func() {
			BeforeEach(func() {
				_, err = fs.Create(
					fmt.Sprintf("%s/limepipes-plugin-%s", pluginsDir, pluginID),
				)
				Expect(err).NotTo(HaveOccurred())
			})

			When("getting LimePipes plugin fails", func() {
				BeforeEach(func() {
					processHandler.EXPECT().RunPlugin(
						pluginID,
						fmt.Sprintf("%s/limepipes-plugin-%s", pluginsDir, pluginID)).
						Return(nil)
					processHandler.EXPECT().GetPlugin(
						pluginID,
					).Return(nil, fmt.Errorf("no plugin with this id"))
				})

				It("should return an error", func() {
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("getting LimePipes plugin succeeds", func() {
				BeforeEach(func() {
					processHandler.EXPECT().RunPlugin(
						pluginID,
						fmt.Sprintf("%s/limepipes-plugin-%s", pluginsDir, pluginID)).
						Return(nil)
					processHandler.EXPECT().GetPlugin(
						pluginID,
					).Return(lpPlugin, nil)
				})

				When("getting plugin info fails", func() {
					BeforeEach(func() {
						lpPlugin.EXPECT().PluginInfo().
							Return(nil, fmt.Errorf("no plugin info"))
					})

					It("should return an error", func() {
						Expect(err).Should(HaveOccurred())
					})
				})

				Context("getting plugin info succeeds", func() {
					BeforeEach(func() {
						lpPlugin.EXPECT().PluginInfo().
							Return(&messages.PluginInfoResponse{
								Name:           pluginID,
								Type:           messages.PluginType_INOUT,
								FileFormat:     fileformat.Format_BWW,
								FileExtensions: []string{".bww", ".bmw"},
							}, nil)
					})

					It("should not return an error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should have the plugin info in the pluginInfos map", func() {
						Expect(loader.LoadedPlugins()).To(HaveLen(1))
					})

					When("getting the file type for an unhandled file extension", func() {
						var ff fileformat.Format
						JustBeforeEach(func() {
							ff, err = loader.FileTypeForFileExtension(".xxx")
						})

						It("should return an error", func() {
							Expect(err).Should(HaveOccurred())
							Expect(ff).To(Equal(fileformat.Format_Unknown))
						})
					})

					When("getting the file type for a valid extension", func() {
						var ff fileformat.Format
						JustBeforeEach(func() {
							ff, err = loader.FileTypeForFileExtension(".bww")
						})

						It("should return the correct file type", func() {
							Expect(err).ShouldNot(HaveOccurred())
							Expect(ff).To(Equal(fileformat.Format_BWW))
						})
					})

					When("unloading plugins fails", func() {
						BeforeEach(func() {
							processHandler.EXPECT().KillPlugins().
								Return(fmt.Errorf("failed killing plugins"))
						})

						JustBeforeEach(func() {
							err = loader.UnloadPlugins()
						})

						It("should return an error", func() {
							Expect(err).Should(HaveOccurred())
						})
					})

					When("unloading plugins succeeds", func() {
						BeforeEach(func() {
							processHandler.EXPECT().KillPlugins().
								Return(nil)
						})

						JustBeforeEach(func() {
							err = loader.UnloadPlugins()
						})

						It("should return an error", func() {
							Expect(err).ShouldNot(HaveOccurred())
						})
					})

					When("getting the plugin for an unhandled file extension", func() {
						JustBeforeEach(func() {
							_, err = loader.PluginForFileExtension(".xxx")
						})

						It("should return an error", func() {
							Expect(err).Should(HaveOccurred())
						})
					})

					When("getting the plugin for a valid file extension", func() {
						var plug interfaces.LimePipesPlugin
						JustBeforeEach(func() {
							plug, err = loader.PluginForFileExtension(".bww")
						})

						It("should return an error", func() {
							Expect(err).ShouldNot(HaveOccurred())
							Expect(plug).To(Equal(lpPlugin))
						})
					})
				})
			})
		})
	})
})
