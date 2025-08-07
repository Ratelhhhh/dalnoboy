package app

import "fmt"

// App представляет основное приложение
type App struct {
	Name string
}

// New создает новый экземпляр приложения
func New(name string) *App {
	return &App{
		Name: name,
	}
}

// Run запускает приложение
func (a *App) Run() error {
	fmt.Printf("Приложение %s запущено\n", a.Name)
	return nil
}
