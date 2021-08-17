package registry

// Service 服务抽象
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

// Node 服务节点的抽象
type Node struct {
	Id     string `json:"id"`
	IP     string `json:"ip"`
	Port   uint   `json:"port"`
	Weight int    `json:"weight"`
}
