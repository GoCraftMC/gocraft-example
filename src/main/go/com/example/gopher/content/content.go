// Package content registers the mod's items and blocks.
package content

import (
	"github.com/GoCraftMC/gocraft-sdk/mc"
	"github.com/GoCraftMC/gocraft-sdk/registry"
)

// Registered content, available to the rest of the mod after Register.
var (
	Wand  mc.Item
	Block registry.RegisteredBlock
)

// Register adds everything to the game's registries.
func Register() {
	Wand = registry.RegisterItem("gopher_wand", registry.Item{
		MaxStack: 1,
		OnUse: func(ctx registry.UseItem) mc.InteractionResult {
			if !ctx.Level.IsClient() {
				below := ctx.Player.Entity().BlockPos()
				ctx.Level.SetBlock(below.Above().Above(), Block.DefaultState())
				ctx.Player.Send("A gopher block appears!")
			}
			return mc.Success()
		},
	})

	Block = registry.RegisterBlock("gopher_block", registry.Block{
		Hardness:   1.5,
		Resistance: 6,
		Light:      12,
		OnUse: func(ctx registry.UseBlock) mc.InteractionResult {
			if !ctx.Level.IsClient() {
				ctx.Player.Send("The gopher block hums quietly.")
			}
			return mc.Success()
		},
	})

	registry.AddToTab(registry.TabTools, Wand)
	registry.AddToTab(registry.TabBuildingBlocks, Block.Item)
}
