package types

// UserDetailsConfig - User Details config
type UserDetailsConfig struct {
	UserDetails map[string]UserDetails `yaml:"userDetails"`
}

// UserDetails - User Details
type UserDetails struct {
	UserName string   `yaml:"userName"`
	Password string   `yaml:"password"`
	Email    string   `yaml:"email"`
	UID      string   `yaml:"uid"`
	Groups   []string `yaml:"groups"`
}
