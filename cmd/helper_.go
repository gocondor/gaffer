package cmd

import (
	"os"
	"strings"
	"unicode"
)

func getEventNameStatement(constName string, eventName string) string {
	t := `const {constName} = "{eventName}"`
	r := strings.Replace(t, "{constName}", constName, 1)
	r = strings.Replace(r, "{eventName}", eventName, 1)
	return r
}

func prepareEventNameConst(eventName string) string {
	var res string
	words := strings.Split(eventName, "-")
	for k, v := range words {
		if k == 0 {
			res = strings.ToUpper(v) + "_"
		} else if k < (len(words) - 1) {
			res = res + strings.ToUpper(v) + "_"
		} else {
			res = res + strings.ToUpper(v)
		}
	}
	return res
}

func camelCaseToSnake(name string) string {
	var res string
	namesB := []byte(name)
	for i, v := range namesB {
		if i == 0 {
			res = res + strings.ToLower(string(v))
		} else {
			if !unicode.IsUpper(rune(v)) {
				res = res + string(v)
			} else {
				res = res + "-"
				res = res + strings.ToLower(string(v))
			}
		}
	}
	return res
}

func prepareJobContent(jobName string) string {
	t := `package eventjobs

import (
	"github.com/gocondor/core"
)

var {JobName} core.EventJob = func(event *core.Event, c *core.Context) {
	// logic implementation goes here...
}
`
	res := strings.Replace(t, "{JobName}", jobName, 1)
	return res
}

func prepareHandlerContent(HandlerName string) string {
	t := `
func {HandlerName}(c *core.Context) *core.Response {
	// logic implementation goes here...

	return nil
}
`
	res := strings.Replace(t, "{HandlerName}", HandlerName, 1)
	return res
}

func createHandlerFile(handlersFilePath string) (*os.File, error) {
	file, err := os.Create(handlersFilePath)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(`package handlers

import (
	"github.com/gocondor/core"
)
`)

	if err != nil {
		return nil, err
	}
	return file, nil
}

func prepareMiddlewareContent(middlewareName string) string {
	t := `package middlewares

import (
	"github.com/gocondor/core"
)

var {middlewareName} core.Middleware = func (c *core.Context) {
	c.Next()
}
`
	res := strings.Replace(t, "{middlewareName}", middlewareName, 1)
	return res
}
