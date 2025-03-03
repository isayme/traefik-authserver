package conf

type SessionConfig struct {
	Name     string `json:"name" yaml:"name"`
	Secret   string `json:"secret" yaml:"secret"`
	Domain   string `json:"domain" yaml:"domain"`
	MaxAge   int    `json:"maxAge" yaml:"maxAge"`
	HttpOnly bool   `json:"httpOnly" yaml:"httpOnly"`
	Secure   bool   `json:"secure" yaml:"secure"`
	LoginUrl string `json:"loginUrl" yaml:"loginUrl"`
}

type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Github   string `json:"github" yaml:"github"`
}

type Github struct {
	ClientId      string   `json:"clientId" yaml:"clientId"`
	ClientSecret  string   `json:"clientSecret" yaml:"clientSecret"`
	RedirectUrl   string   `json:"redirectUrl" yaml:"redirectUrl"`
	WhitelistUser []string `json:"whitelistUser" yaml:"whitelistUser"`
}

type Config struct {
	Session SessionConfig `json:"session" yaml:"session"`
	Users   []User        `json:"users" yaml:"users"`
	Github  Github        `json:"github" yaml:"github"`
}
