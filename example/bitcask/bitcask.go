package main

import (
	"github.com/adylanrff/bitcask"
)

func main() {
	db, err := bitcask.Open("test.db")
	if err != nil {
		panic(err)
	}
	if val, err := db.Get("a"); err != nil {
		panic(err)
	} else {
		println(val)
	}
}
