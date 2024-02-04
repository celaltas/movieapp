package main


type apiConfig struct {
	Port string `yaml:"port"`
}


type serverConfig struct {
	API apiConfig `yaml:"api"`
}