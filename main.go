package main

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen-config.yaml schema.yaml

import "merchshop/cmd"

func main() {
	cmd.Execute()
}
