package main

// 产品
type DBConnection interface {
	Connect() (string, error)
	Close() (string, error)
}

type DBCommand interface {
	Execute(cmd string) (string, error)
}

// 工厂
type DBFactory interface {
	CreateConnection() DBConnection
	CreateCommand() DBCommand
}

func NewFactory(dbType string) DBFactory {
	switch dbType {
	case "mysql":
		return &MySQLFactory{}
	case "postgresql":
		return &PostgreSQLFactory{}
	default:
		return nil
	}
}

//// 具体工厂
//
type MySQLFactory struct{}

func (f *MySQLFactory) CreateConnection() DBConnection {
	return &MySQLConnection{}
}

func (f *MySQLFactory) CreateCommand() DBCommand {
	return &MySQLCommand{}
}

type MySQLConnection struct{}

func (c *MySQLConnection) Connect() (string, error) {
	return "MySQL connected", nil
}

func (c *MySQLConnection) Close() (string, error) {
	return "MySQL connection closed", nil
}

type MySQLCommand struct{}

func (cmd *MySQLCommand) Execute(command string) (string, error) {
	return "MySQL executed: " + command, nil
}

//
type PostgreSQLFactory struct{}

func (f *PostgreSQLFactory) CreateConnection() DBConnection {
	return &PostgreSQLConnection{}
}

func (f *PostgreSQLFactory) CreateCommand() DBCommand {
	return &PostgreSQLCommand{}
}

type PostgreSQLConnection struct{}

func (c *PostgreSQLConnection) Connect() (string, error) {
	return "PostgreSQL connected", nil
}

func (c *PostgreSQLConnection) Close() (string, error) {
	return "PostgreSQL connection closed", nil
}

type PostgreSQLCommand struct{}

func (cmd *PostgreSQLCommand) Execute(command string) (string, error) {
	return "PostgreSQL executed: " + command, nil
}

//