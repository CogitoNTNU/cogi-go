package config

type LocalDbOption func(l *LocalDbConfig)

func WithDbName(name string) LocalDbOption {
	return func(l *LocalDbConfig) {
		l.SQLDatabase = name
	}
}

func WithDbHost(host string) LocalDbOption {
	return func(l *LocalDbConfig) {
		l.SQLHost = host
	}
}

func WithDbPassword(password string) LocalDbOption {
	return func(l *LocalDbConfig) {
		l.SQLPassword = password
	}
}

func WithDbUser(user string) LocalDbOption {
	return func(l *LocalDbConfig) {
		l.SQLUser = user
	}
}

func WithDbPort(port int) LocalDbOption {
	return func(l *LocalDbConfig) {
		l.SQLPort = port
	}
}

