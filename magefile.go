//+build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Test() error {
	return sh.Run("go", "test", "./...", "-v")
}

func Benchmark() error {
	return sh.Run("go", "test", "./...", "-v", "-bench", ".", "-test.run=filterouttests")
}
