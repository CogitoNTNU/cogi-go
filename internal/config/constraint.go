package config

type SQLConfig interface {
	GetSQLHost() string
	GetSQLUser() string
	GetSQLPassword() string
	GetSQLPort() int
	GetSQLDatabase() string
}

func (c LocalDbConfig) GetSQLHost() string {
	return c.SQLHost
}

func (c LocalDbConfig) GetSQLUser() string {
	return c.SQLUser
}

func (c LocalDbConfig) GetSQLPassword() string {
	return c.SQLPassword
}

func (c LocalDbConfig) GetSQLPort() int {
	return c.SQLPort
}

func (c LocalDbConfig) GetSQLDatabase() string {
	return c.SQLDatabase
}
