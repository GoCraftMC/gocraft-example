// Package client holds code that only runs on the physical client. The whole mod compiles into one
// library, so what keeps client classes off a dedicated server is that nothing here runs there.
package client

import (
	"github.com/GoCraftMC/gocraft-sdk/event"
	"github.com/GoCraftMC/gocraft-sdk/gocraft"
	"github.com/GoCraftMC/gocraft-sdk/jvm"
)

func init() {
	gocraft.OnInitializeClient(func() {
		announced := false
		event.ClientTick(func(jvm.Object) {
			if !announced {
				announced = true
				gocraft.Log().Info("the client is ticking, and Go is watching")
			}
		})
	})
}
