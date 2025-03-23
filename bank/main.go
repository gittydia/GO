package main

import (
	"bank/function" // Changed the import path
	"fmt"
)


func main() {
    var u1 function.SigninU
    var u2 function.LoginU
    var account function.UserAccount

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
            if u2.Email == u1.Email && u2.Password == u1.Password {

                fmt.Println("Welcome to the bank" + " " + u1.Name)
                
                
                for {
                    fmt.Println("\nChoose an option: ")
                    fmt.Println("1. Deposit")
                    fmt.Println("2. Withdraw")
                    fmt.Println("3. Check balance")
                    fmt.Println("4. Exit")
                    fmt.Scanln(&choice)
                    if choice == 1 {
                        function.Deposit(&account)
                    } else if choice == 2 {
                        function.Withdraw(&account)
                    } else if choice == 3 {
                        function.CheckBalance(&account)
                    } else if choice == 4 {
                        fmt.Println("Exiting...")
                        break
                    } else {
                        fmt.Println("Invalid choice")
                    }
                
            }
        } else if choice == 3 {
            fmt.Println("Exiting...")
            break
        } else {
            fmt.Println("Invalid choice")
        }

        
    }
    fmt.Println("\nChoose an option: ")
    fmt.Println("1. Signin")
    fmt.Println("2. Login")
    fmt.Println("3. Exit")
    fmt.Scanln(&choice)
}
}