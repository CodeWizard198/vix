package vix

import (
	"fmt"
	"github.com/valyala/fastjson"
	"testing"
)

type rec struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Pram   []int  `json:"pram"`
}

func TestTool(t *testing.T) {
	path := "/a/d"
	checkPath("GET", path)
}

// func Login(v *vix.Context, form *LoginForm)
// func (l *LoginForm) GetType() *LoginForm {return l}
func TestRec(t *testing.T) {
	j := "{\"name\": \"along\",\"method\": \"get\",\"pram\":[1,2,3,4,5,6]}"
	re := &rec{}
	RJSON(re, j)
	fmt.Println(re)
}

func RJSON(res any, jsons string) {
	var parser fastjson.Parser
	parse, err := parser.Parse(jsons)
	if err != nil {
		fmt.Printf("fastjson error, %s\n", err.Error())
		return
	}
	fmt.Println(parse)
}
