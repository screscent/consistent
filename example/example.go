package main

import (
    "github.com/nokivan/consistent"
)

func main() {
    c := consistent.New()
    c.Add("Server1")
    c.Add("Server2")
    c.Add("Server3")
    c.Add("Server4")
    c.Add("Server5")

    users := []string{ "foo", "bar", "foobar", "i", "enjoy", "golang" }
    for _, u := range users {
        ret, err := c.Get(u)
        if err != nil {
            println(err.Error())
        }
        println(ret)
    }

    println("---")
    c.Remove("Server5")
    c.Remove("Server3")
    for _, u := range users {
        ret, err := c.Get(u)
        if err != nil {
            println(err.Error())
        }
        println(ret)
    }

    println("---")
    c.Add("Server3")
    c.Add("Server5")
    for _, u := range users {
        ret, err := c.Get(u)
        if err != nil {
            println(err.Error())
        }
        println(ret)
    }
}
