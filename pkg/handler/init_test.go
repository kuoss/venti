package handler

import (
	"fmt"
	"os"
)

func init() {
	err := os.Chdir("../..")
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("working directory:", wd)
}
