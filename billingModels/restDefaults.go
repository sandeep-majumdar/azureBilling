package billingModels

type restDefaults struct {
	InitialHTTPSPort int    // InitialHTTPSPort is the start port offset for all the servers we will start
	ServerPem        string // ServerPem certification location
	ServerKey        string // ServerKey certification location
	ClientPem        string // ClientPem certification location
	ClientKey        string // ClientKey certification location
}

func (def *restDefaults) Init() {
	def.InitialHTTPSPort = 443
	def.ServerPem = "certs/server.pem"
	def.ServerKey = "certs/server.key"
	def.ClientPem = "certs/client.pem"
	def.ClientKey = "certs/client.key"
}
