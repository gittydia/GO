package main

import (
	"bank/function" // Changed the import path
	"fmt"
)

func main() {
    var u1 function.SigninU
    var u2 function.LoginU
    var choice int 

    fmt.Println("Choose an option: ")
    fmt.Println("1. Signin")
    fmt.Println("2. Login")
    fmt.Scanln(&choice)

        
    for {
        if choice == 1 {
            function.Signin(&u1)
        } else if choice == 2 {
            function.Login(&u2, &u1)
        } else if choice == 3 {
            fmt.Println("Exiting...")
            break
        } else {
            fmt.Println("Invalid choice")
        }

        fmt.Println("\nChoose an option: ")
        fmt.Println("1. Signin")
        fmt.Println("2. Login")
        fmt.Println("3. Exit")
        fmt.Scanln(&choice)
    }

    
   
}