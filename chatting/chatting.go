package chatting

import "fmt"

type SigninU struct {
    Name     string
    Email    string
    Password string
}

type LoginU struct {
    Name     string
    Email    string
    Password string
}

func Signin(u1 *SigninU) {
    fmt.Println("Enter your name: ")
    fmt.Scanln(&u1.Name)
    fmt.Println("Enter your email: ")
    fmt.Scanln(&u1.Email)
    fmt.Println("Enter your password: ")
    fmt.Scanln(&u1.Password)
    if u1.Name != "" && u1.Email != "" && u1.Password != "" {
        fmt.Println("Signin successful")
    } else {
        fmt.Println("Signin failed")
    }
}

func Login(u2 *LoginU, u1 *SigninU) {
    fmt.Println("Enter your email: ")
    fmt.Scanln(&u2.Email)
    fmt.Println("Enter your password: ")
    fmt.Scanln(&u2.Password)
    if u2.Email == u1.Email && u2.Password == u1.Password {
        fmt.Println("Login successful")
    } else {
        fmt.Println("Login failed")
    }
}

type UserAccount struct {
    balance float64
    accountType string
    amount float64
}
//bank account
func deposit(account *UserAccount) {
    fmt.Println("Enter the amount you want to deposit: ")
    fmt.Scanln(&account.amount)
    account.balance += account.amount
    fmt.Println("Deposit successful")
}

func withdraw(account *UserAccount) {
    fmt.Println("Enter the amount you want to withdraw: ")
    fmt.Scanln(&account.amount)
    if account.balance >= account.amount {
        account.balance -= account.amount
        fmt.Println("Withdrawal successful")
    } else {
        fmt.Println("Insufficient balance")
    }
}