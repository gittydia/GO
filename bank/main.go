package main

import (
	"bank/function" // Changed the import path
)

func main() {
    var u1 function.SigninU
    var u2 function.LoginU
    
    function.Signin(&u1)
    function.Login(&u2, &u1)
}