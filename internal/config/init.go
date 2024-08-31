package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// initViperConfigPaths initializes the viper config paths.
// When withParents is true, it will search for the limepipes config
// file in the parent directories up to 3 levels.
func initViperConfigPaths(withParents bool) error {
	viper.AddConfigPath(".")

	if !withParents {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	addParentDirsUntilSourceRoot(wd, 3)

	return nil
}

// addParentDirsUntilSourceRoot goes maxDirUp directories up from the given dir
// until it finds the source root directory called "limepipes" which then
// it adds to the viper config paths.
// If it doesn't find the source root directory, it will not add any paths.
func addParentDirsUntilSourceRoot(dir string, maxDirUp int) {
	if filepath.Base(dir) != "limepipes" {
		for i := 0; i < maxDirUp; i++ {
			dir = filepath.Dir(dir)

			if filepath.Base(dir) == "limepipes" {
				viper.AddConfigPath(dir)
				break
			}
		}
	}
}

func initViperConfig() {
	viper.SetConfigName("limepipes")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
}

func readViperConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

func Init() (*Config, error) {
	err := initViperConfigPaths(false)
	if err != nil {
		return nil, err
	}

	initViperConfig()

	return readViperConfig()
}

func InitTest() (*Config, error) {
	err := initViperConfigPaths(true)
	if err != nil {
		return nil, err
	}

	initViperConfig()

	return readViperConfig()
}
