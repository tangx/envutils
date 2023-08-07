package envutils

import (
	"log"
	"os"
)

// Export config struct to default.yml, variables with prefix
func Export(prefix string, config interface{}) error {

	err := CallSetDefaults(config)
	if err != nil {
		return err
	}

	b, err := Marshal(config, prefix)
	if err != nil {
		return err
	}

	return os.WriteFile("default.yml", b, os.ModePerm)
}

// MustExport export config struct to default.yml, variables with prefix
// panic if error
func MustExport(prefix string, config interface{}) {
	err := Export(prefix, config)
	if err != nil {
		panic(err)
	}
}

// Import variables from config.yml and additional config files to config struct
// then import variable from env to config struct
// overwrite if variable already exists
func Import(prefix string, config interface{}, cfgs ...string) error {

	// initial defualt config
	CallSetDefaults(config)

	// read variables from files
	files := append([]string{"config.yml"}, cfgs...)
	for _, file := range files {

		err := UnmarshalFile(config, prefix, file)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("WARN: skip, %v\n", err)
				continue
			}

			return err
		}
	}

	// read variables from environment
	return UnmarshalEnv(config, prefix)
}

// MustImport import variable from config.yml and additional config files and environment, panic if error
// Get more detail, see Import API
func MustImport(prefix string, config interface{}, cfgs ...string) {
	err := Import(prefix, config, cfgs...)
	if err != nil {
		panic(err)
	}
}
