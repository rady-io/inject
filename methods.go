package rady

import "reflect"

type (
	/*
	GET is a tag to mark a method with path and http GET method

	Usage:

		type UserController struct {
			Controller 	`prefix:"api/v1"`
			GET 		`path:":id" method:"GetUserInfo"`
		}

		func (u *UserController) GetUserInfo(ctx echo.Context) error {
			// do something
		}

		type Root struct {
			*UserController
		}

		func main() {
			CreateApplication(new(Root)).Run()
		}
	 */
	GET struct {
	}

	// POST is a tag to mark a method with path and http POST method
	POST struct {
	}

	// PUT is a tag to mark a method with path and http PUT method
	PUT struct {
	}

	// HEAD is a tag to mark a method with path and http HEAD method
	HEAD struct {
	}

	// DELETE is a tag to mark a method with path and http DELETE method
	DELETE struct {
	}

	// CONNECT is a tag to mark a method with path and http CONNECT method
	CONNECT struct {
	}

	// OPTIONS is a tag to mark a method with path and http OPTIONS method
	OPTIONS struct {
	}

	// TRACE is a tag to mark a method with path and http TRACE method
	TRACE struct {
	}

	// PATCH is a tag to mark a method with path and http PATCH method
	PATCH struct {
	}
)

const (
	GetStr     = "Get"
	PostStr    = "Post"
	PutStr     = "Put"
	HeadStr    = "Head"
	DeleteStr  = "Delete"
	ConnectStr = "Connect"
	OptionsStr = "Options"
	TraceStr   = "Trace"
	PatchStr   = "Patch"
)

var (
	StrToMethod = map[string]interface{}{
		GetStr:     GET{},
		PostStr:    POST{},
		PutStr:     PUT{},
		HeadStr:    HEAD{},
		DeleteStr:  DELETE{},
		ConnectStr: CONNECT{},
		OptionsStr: OPTIONS{},
		TraceStr:   TRACE{},
		PatchStr:   PATCH{},
	}
	MethodToStr = map[interface{}]string{
		GET{}:     GetStr,
		POST{}:    PostStr,
		PUT{}:     PutStr,
		HEAD{}:    HeadStr,
		DELETE{}:  DeleteStr,
		CONNECT{}: ConnectStr,
		OPTIONS{}: OptionsStr,
		TRACE{}:   TraceStr,
		PATCH{}:   PatchStr,
	}

	MethodsTypeSet = make(map[reflect.Type]string)
)
