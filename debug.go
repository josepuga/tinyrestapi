//go:build debug
package main

import "fmt"

//debugPrint Only to print in debug mode: `go run -tags debug main.go`
func debugPrint(v ...any) {
	if len(v) == 0 {
		return
	}
	var msg string
	if format, ok := v[0].(string); ok {
		msg = fmt.Sprintf(format, v[1:]...)
	} else {
		msg = fmt.Sprint(v...)
	}    
	fmt.Printf("[DEBUG] %s\n",msg)
}