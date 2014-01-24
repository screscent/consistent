package main

import (
    "github.com/screscent/consistent"
)

func main() {
   c := consistent.New()
	c.AddKey("Server1")
	c.AddKey("Server2")
	c.AddKey("Server3")
	c.AddKey("Server4")
	c.AddKey("Server5")
	c.Update()

	users := []string{"foo", "bar", "foobar", "i", "enjoy", "golang"}
	for _, u := range users {
		ret, err := c.GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}

	println("---")
	c.Remove("Server5")
	c.Remove("Server3")
	c.Update()
	for _, u := range users {
		ret, err := c.GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}

	println("---")
	c.AddKey("Server3")
	c.AddKey("Server5")
	c.Update()
	for _, u := range users {
		ret, err := c.GetKey(u)
		if err != nil {
			println(err.Error())
		}
		println(ret)
	}
}
