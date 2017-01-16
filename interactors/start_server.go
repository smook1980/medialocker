package interactors

func ServerInteractor(app *medialocker.App) error {
	server := app.GetServer()
	&ServerManager{server:server}
}

struct ServerManager {
	medialocker.Server
}
