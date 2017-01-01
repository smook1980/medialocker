package medialocker

type Cli interface {
	//ConfigFromCli()
	Exec([]string) error
}
