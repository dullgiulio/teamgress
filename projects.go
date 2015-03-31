package teamgress

import (
	"fmt"
)

type Project struct {
	Name    string `json:"name"`
	Service string `json:"type,omitempty"`
	Stage   string `json:"stage,omitempty"`
	Server  string `json:"server,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
}

func MakeProject(name, service, stage string) Project {
	return Project{
		Name:    name,
		Service: service,
		Stage:   stage,
	}
}

func (p *Project) String() string {
	return fmt.Sprintf("%s[%s] %s/%s", p.Name, p.Server, p.Service, p.Stage)
}
