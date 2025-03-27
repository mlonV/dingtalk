package supervisor

type Mysql struct {
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
}
