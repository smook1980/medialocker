package medialocker

type Server interface {
	Start(listen ...string) error
	Stop(timeout int) error
}
