package bel

import (
	"os"
    "testing"
)

func TestGenerateStuff(t *testing.T) {
    extr, err := Extract(StructOfAllKind{})
    if err != nil {
        t.Error(err)
        return
    }

    err = extr.RenderInterface(os.Stdout)
    if err != nil {
        t.Error(err)
        return
    }
}