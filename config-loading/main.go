package main

import (
	"fmt"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/environment"
)

func main() {
	fmt.Println("=== Config Loading Example ===")
	fmt.Println()

	fmt.Println("Example 1: Default config loading")
	app1, err := boot.NewApplication(
		boot.WithAppName("config-example"),
		boot.WithProfiles("dev"),
	)
	if err != nil {
		panic(err)
	}

	if err := app1.Start(); err != nil {
		panic(err)
	}
	defer app1.Stop()

	env := app1.Environment()
	fmt.Printf("App Name: %s\n", env.GetString("app.name", "unknown"))
	fmt.Printf("App Port: %d\n", env.GetInt("app.port", 0))
	fmt.Printf("Debug Mode: %v\n", env.GetBool("app.debug", false))
	fmt.Println()

	fmt.Println("Example 2: Custom config location")
	app2, err := boot.NewApplication(
		boot.WithAppName("custom-location-example"),
		boot.WithConfigLocation("./application-prod.json"),
	)
	if err != nil {
		panic(err)
	}

	if err := app2.Start(); err != nil {
		panic(err)
	}
	defer app2.Stop()

	env2 := app2.Environment()
	fmt.Printf("App Name: %s\n", env2.GetString("app.name", "unknown"))
	fmt.Printf("App Port: %d\n", env2.GetInt("app.port", 0))
	fmt.Printf("Debug Mode: %v\n", env2.GetBool("app.debug", false))
	fmt.Println()

	fmt.Println("Example 3: Custom property source")
	customSource := environment.NewMapPropertySource(
		"custom",
		environment.PriorityHigh,
		map[string]any{
			"app.name": "custom-app",
			"app.port": 9999,
		},
	)

	app3, err := boot.NewApplication(
		boot.WithAppName("custom-example"),
		boot.WithPropertySource(customSource),
	)
	if err != nil {
		panic(err)
	}

	if err := app3.Start(); err != nil {
		panic(err)
	}
	defer app3.Stop()

	env3 := app3.Environment()
	fmt.Printf("App Name: %s\n", env3.GetString("app.name", "unknown"))
	fmt.Printf("App Port: %d\n", env3.GetInt("app.port", 0))
	fmt.Println()

	fmt.Println("Config loading example completed!")
}
