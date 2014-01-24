package main

import (
    . "github.com/screscent/consistent"
)

func main() {
	AddKey("Server1")
	AddKey("Server2")
	AddKey("Server3")
	AddKey("Server4")
	AddKey("Server5")
	Update()

	users := []string{"foo", "bar", "foobar", "i", "enjoy", "golang"}
	for _, u := range users {
		ret, err := GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}

	println("---")
	Remove("Server5")
	Remove("Server3")
	Update()
	for _, u := range users {
		ret, err := GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}

	println("---")
	AddKey("Server3")
	AddKey("Server5")
	Update()
	for _, u := range users {
		ret, err := GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}
}
