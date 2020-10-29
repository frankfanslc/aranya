package conf

type VirtualnodeNetworkConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	Mesh    VirtualnodeNetworkMeshConfig    `json:"mesh" yaml:"mesh"`
	Backend VirtualnodeNetworkBackendConfig `json:"backend" yaml:"backend"`

	NetworkService struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
		IP   string `json:"ip" yaml:"ip"`
		Port int    `json:"port" yaml:"port"`
	} `json:"networkService" yaml:"networkService"`

	AbbotService struct {
		Name     string `json:"name" yaml:"name"`
		PortName string `json:"portName" yaml:"portName"`
	} `json:"abbotService" yaml:"abbotService"`
}

type VirtualnodeNetworkBackendConfig struct {
	Driver    string `json:"driver" yaml:"driver"`
	Wireguard struct {
		Name       string `json:"name" yaml:"name"`
		MTU        int    `json:"mtu" yaml:"mtu"`
		ListenPort int    `json:"listenPort" yaml:"listenPort"`
		PrivateKey string `json:"privateKey" yaml:"privateKey"`
	} `json:"wireguard" yaml:"wireguard"`
}

type VirtualnodeNetworkMeshConfig struct {
	// allocate ip addresses to mesh network device
	IPv4Blocks []AddressBlock `json:"ipv4Blocks" yaml:"ipv4Blocks"`
	IPv6Blocks []AddressBlock `json:"ipv6Blocks" yaml:"ipv6Blocks"`
}

type AddressBlock struct {
	CIDR  string `json:"cidr"`
	Start string `json:"start"`
	End   string `json:"end"`
}
