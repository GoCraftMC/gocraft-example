// Package gopher is the mod's main entrypoint. go.jar links in every package under src/main/go, and
// each one registers what it needs from init.
package gopher

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/GoCraftMC/gocraft-example/src/main/go/com/example/gopher/content"
	"github.com/GoCraftMC/gocraft-sdk/command"
	"github.com/GoCraftMC/gocraft-sdk/event"
	"github.com/GoCraftMC/gocraft-sdk/gocraft"
	"github.com/GoCraftMC/gocraft-sdk/mc"
)

// Config is stored in config/gopher.json.
type Config struct {
	Greeting string `json:"greeting"`
}

func init() {
	// A hook is a mixin declared in gomod.json. This one appends to the server's message of the day.
	gocraft.Hook("motd", func(h *gocraft.HookContext) {
		h.SetReturn(h.ReturnValue().String() + " (powered by Go)")
	})

	gocraft.OnInitialize(func() {
		cfg := Config{Greeting: "Welcome, %s. This server speaks Go."}
		if err := gocraft.LoadConfig(&cfg); err != nil {
			gocraft.Log().Error("could not read config", "err", err)
		}

		content.Register()
		registerCommands()
		watchUptime()

		event.PlayerJoin(func(player mc.Player, _ mc.Server) {
			player.Send(fmt.Sprintf(cfg.Greeting, player.Name()))
		})

		gocraft.Log().Info("Gopher is ready")
	})
}

func registerCommands() {
	command.Register(
		command.Literal("gopher").
			Executes(func(c *command.Context) int {
				c.Reply("Hello from Go, " + c.SenderName() + "!")
				return 1
			}).
			Then(
				command.Literal("say").Then(
					command.Arg("message", command.GreedyString()).Executes(func(c *command.Context) int {
						c.Server().Broadcast("[Gopher] " + c.String("message"))
						return 1
					}),
				),
				command.Literal("fib").Then(
					command.Arg("n", command.IntRange(1, 40)).Executes(func(c *command.Context) int {
						n := c.Int("n")
						c.Reply(fmt.Sprintf("fib(%d) = %d", n, fibonacci(n)))
						return 1
					}),
				),
				// Work that takes real time belongs in a goroutine, not on the server thread. Results go
				// back through Server.Execute, because game state belongs to the server thread.
				command.Literal("later").Executes(func(c *command.Context) int {
					server := c.Server()
					gocraft.Go(func() {
						time.Sleep(time.Second)
						server.Execute(func() { server.Broadcast("[Gopher] One second has passed, and the game never stalled.") })
					})
					return 1
				}),
			),
	)
}

func fibonacci(n int) int {
	a, b := 0, 1
	for range n {
		a, b = b, a+b
	}
	return a
}

// watchUptime shows the ordinary Go way to do background work: a goroutine, a ticker and a context
// that ends with the server.
func watchUptime() {
	var ticks atomic.Int64
	event.ServerTick(func(mc.Server) { ticks.Add(1) })

	event.ServerStarted(func(server mc.Server) {
		ctx, cancel := context.WithCancel(context.Background())
		event.ServerStopping(func(mc.Server) { cancel() })

		gocraft.Go(func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					gocraft.Log().Info("still running", "ticks", ticks.Load())
				}
			}
		})
	})
}
