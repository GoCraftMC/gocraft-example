# GoCraft example

A complete [GoCraft](https://github.com/GoCraftMC/gocraft) mod in Go, laid out like a Fabric mod. It has an
item, a light-emitting block with its model, texture and recipe, commands, events, a config file, a hook that
changes the server's message of the day, and goroutines handing work back to the server thread.

```sh
./gow build
./gow runClient
```

In the game, try `/give @s gopher:gopher_wand`, `/gopher fib 30` and `/gopher later`.

`go.mod` points at a sibling checkout of `gocraft-sdk`, and `build.toml` at a sibling checkout of `gocraft`
for the mod jar, so clone all three next to each other and run `./gradlew build` in `gocraft` first.

## License

BSD-3-Clause.
