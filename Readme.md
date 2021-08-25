# What is Commandblock_go?

Commandblock_go is a port of [Commandblock_py](https://github.com/skandabhairava/Datapack_generator) for Go

# How to use?

## Register a new datapack
`mypack := datapack.New("<name of datapack>", "<namespace>", 7, "<load function name>", "<tick function name>", "<description>")`

## Register a new command
`mypack.RegisterFunction("<function name>", "<content>")`

## Generate the datapack
`mypack.Generate()`