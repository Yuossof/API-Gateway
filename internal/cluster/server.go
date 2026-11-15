package cluster

type Server struct {
	URL string `yaml:"url" validate:"required,url"`
}

type Cluster struct {
	Name    string   `yaml:"name" validate:"required"`
	Servers []Server `yaml:"servers" validate:"required,min=1,dive"`
}
