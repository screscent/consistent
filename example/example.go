package main

import (
    "github.com/nokivan/consistent"
)

func main() {
    c := consistent.New()
    c.Add("Server1")
    c.Add("Server2")
    c.Add("Server3")

    users := []string{ "foo", "bar", "foobar" }
    for _, u := range users {
        ret, err := c.Get(u)
        if err != nil {
            println(err.Error())
        }
        println(ret)
    }
}
