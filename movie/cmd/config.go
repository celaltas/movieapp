package main

type apiConfig struct {
	Port string `yaml:"port"`
}
type jaegerConfig struct {
	URL string `yaml:"url"`
}

type serverConfig struct {
	API    apiConfig    `yaml:"api"`
	Jaeger jaegerConfig `yaml:"jaeger"`
}
