package function

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
    Balance float64
    AccountType string
    Amount float64
}
//bank account
func Deposit(account *UserAccount) {
    fmt.Println("Enter the amount you want to deposit: ")
    fmt.Scanln(&account.Amount)
    account.Balance += account.Amount
    fmt.Println("Deposit successful")
}

func Withdraw(account *UserAccount) {
    fmt.Println("Enter the amount you want to withdraw: ")
    fmt.Scanln(&account.Amount)
    if account.Balance >= account.Amount {
        account.Balance -= account.Amount
        fmt.Println("Withdrawal successful")
    } else {
        fmt.Println("Insufficient balance")
    }
}
func CheckBalance(account *UserAccount) {
    fmt.Println("Your balance is: ", account.Balance)
}

    

