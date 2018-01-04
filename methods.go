package rhapsody

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
type GET struct {
}

// POST is a tag to mark a method with path and http POST method
type POST struct {
}


// PUT is a tag to mark a method with path and http PUT method
type PUT struct {
}

// HEAD is a tag to mark a method with path and http HEAD method
type HEAD struct{
}

// DELETE is a tag to mark a method with path and http DELETE method
type DELETE struct {
}

// CONNECT is a tag to mark a method with path and http CONNECT method
type CONNECT struct {
}

// OPTIONS is a tag to mark a method with path and http OPTIONS method
type OPTIONS struct {
}

// TRACE is a tag to mark a method with path and http TRACE method
type TRACE struct {
}

// PATCH is a tag to mark a method with path and http PATCH method
type PATCH struct {
}
