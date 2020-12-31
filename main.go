package confucius

import "os"

type confucius struct {
	Container Container
}

func NewConfuciusServer(services map[string]Service) *confucius {
	bootstrap()
	server := confucius{Container: NewContainer(&Config{})}

	for name, service := range services {
		server.Container.Set(name, service)
	}

	return &server
}

func bootstrap() {
	// Set working dir
	workDir, err := os.Getwd()
	err = os.Setenv("WORKING_PATH", workDir)

	if err != nil {
		panic("No working dir")
	}

	// Set config path
	err = os.Setenv("CONFIG_PATH", workDir+"/config")

	if err != nil {
		panic("Cannot set config path")
	}

	if _, err := os.Stat(os.Getenv("CONFIG_PATH")); os.IsNotExist(err) {
		err = os.Mkdir(os.Getenv("CONFIG_PATH"), 0755)

		if err != nil {
			panic("Config path doesn't exists and cannot be created")
		}
	}

	// Set certificates path
	err = os.Setenv("CERTS_PATH", workDir+"/certs")

	if err != nil {
		panic("Cannot set certs path")
	}

	if _, err := os.Stat(os.Getenv("CERTS_PATH")); os.IsNotExist(err) {
		err = os.Mkdir(os.Getenv("CERTS_PATH"), 0755)

		if err != nil {
			panic("Certs path doesn't exists and cannot be created")
		}
	}
}
