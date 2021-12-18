// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"

	"github.com/flamego/flamego"
)

// Context represents context of a request.
type Context struct {
	flamego.Context
}

// Contexter initializes a classic context for a request.
func Contexter() flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}

		ctx.Map(c)
	}
}

// String writes the plain text response body with the given status code.
func (ctx Context) String(statusCode int, body string, v ...interface{}) {
	ctx.ResponseWriter().Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.ResponseWriter().WriteHeader(statusCode)
	msg := fmt.Sprintf(body, v...)
	_, _ = ctx.ResponseWriter().Write([]byte(msg))
}
