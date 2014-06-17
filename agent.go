package consulapi

// AgentCheck represents a check known to the agent
type AgentCheck struct {
	Node        string
	CheckID     string
	Name        string
	Status      string
	Notes       string
	Output      string
	ServiceID   string
	ServiceName string
}

// AgentService represents a service known to the agent
type AgentService struct {
	ID      string
	Service string
	Tags    []string
	Port    int
}

// AgentMember represents a cluster member known to the agent
type AgentMember struct {
	Name        string
	Addr        string
	Port        uint16
	Tags        map[string]string
	Status      int
	ProtocolMin uint8
	ProtocolMax uint8
	ProtocolCur uint8
	DelegateMin uint8
	DelegateMax uint8
	DelegateCur uint8
}

// AgentServiceRegistration is used to register a new service
type AgentServiceRegistration struct {
	ID    string   `json:",omitempty"`
	Name  string   `json:",omitempty"`
	Tags  []string `json:",omitempty"`
	Port  int      `json:",omitempty"`
	Check *AgentServiceCheck
}

// AgentServiceCheck is used to create an associated
// check for a service
type AgentServiceCheck struct {
	Script   string `json:",omitempty"`
	Interval string `json:",omitempty"`
	TTL      string `json:",omitempty"`
}

// Agent can be used to query the Agent endpoints
type Agent struct {
	c *Client
}

// Agent returns a handle to the agent endpoints
func (c *Client) Agent() *Agent {
	return &Agent{c}
}

// Self is used to query the agent we are speaking to for
// information about itself
func (a *Agent) Self() (map[string]map[string]interface{}, error) {
	r := a.c.newRequest("GET", "/v1/agent/self")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out map[string]map[string]interface{}
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Checks returns the locally registered checks
func (a *Agent) Checks() (map[string]*AgentCheck, error) {
	r := a.c.newRequest("GET", "/v1/agent/checks")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out map[string]*AgentCheck
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Services returns the locally registered services
func (a *Agent) Services() (map[string]*AgentService, error) {
	r := a.c.newRequest("GET", "/v1/agent/services")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out map[string]*AgentService
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Members returns the known gossip members. The WAN
// flag can be used to query a server for WAN members.
func (a *Agent) Members(wan bool) ([]*AgentMember, error) {
	r := a.c.newRequest("GET", "/v1/agent/members")
	if wan {
		r.params.Set("wan", "1")
	}
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out []*AgentMember
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceRegister is used to register a new service with
// the local agent
func (a *Agent) ServiceRegister(service *AgentServiceRegistration) error {
	r := a.c.newRequest("PUT", "/v1/agent/service/register")
	r.obj = service
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ServiceDeregister is used to deregister a service with
// the local agent
func (a *Agent) ServiceDeregister(serviceID string) error {
	r := a.c.newRequest("PUT", "/v1/agent/service/deregister/"+serviceID)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
