//go:generate genny -in=handlers.tpl -out=posthandler.go gen "EventType=*data.Post name=post NAME=Post"
//go:generate genny -in=handlers.tpl -out=commenthandler.go gen "EventType=*data.Comment name=comment NAME=Comment"
//go:generate genny -in=handlers.tpl -out=messagehandler.go gen "EventType=*data.Message name=message NAME=Message"
// Package handlers provides handler interfaces for Reddit types.
package handlers
