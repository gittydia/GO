package function

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SigninU struct {
	Name     string
	LastName string
	Email    string
	Password string
}

type LoginU struct {
	Name     string
	Email    string
	Password string
}

// GenerateAccountNumber generates a unique account number.
func GenerateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("ACC%09d", rand.Intn(1000000000)) // Example: ACC123456789
}

func Signin(u1 *SigninU, db *sql.DB) {
	fmt.Println("Enter your name: ")
	fmt.Scanln(&u1.Name)
	u1.Name = strings.TrimSpace(u1.Name)
	fmt.Println("Enter your last name: ")
	fmt.Scanln(&u1.LastName)
	u1.LastName = strings.TrimSpace(u1.LastName)
	fmt.Println("Enter your email: ")
	fmt.Scanln(&u1.Email)
	u1.Email = strings.TrimSpace(u1.Email)
	fmt.Println("Enter your password: ")
	fmt.Scanln(&u1.Password)
	u1.Password = strings.TrimSpace(u1.Password)

	if u1.Name != "" && u1.LastName != "" && u1.Email != "" && u1.Password != "" {
		_, err := db.Exec("INSERT INTO Customers (FirstName, LastName, Email, Password) VALUES (?, ?, ?, ?)", u1.Name, u1.LastName, u1.Email, u1.Password)
		if err != nil {
			fmt.Println("Error during Signin:", err)
			return
		}

		// Automatically create an account for the user
		accountNumber := GenerateAccountNumber()
		_, err = db.Exec("INSERT INTO Accounts (CustomerID, AccountNumber, Balance) VALUES ((SELECT CustomerID FROM Customers WHERE Email = ?), ?, 0)", u1.Email, accountNumber)
		if err != nil {
			fmt.Println("Error creating account for the user:", err)
			return
		}

		fmt.Println("Signin successful and account created.")
	} else {
		fmt.Println("Signin failed: All fields are required.")
	}
}

func Login(u2 *LoginU, db *sql.DB) {
	fmt.Println("Enter your email: ")
	fmt.Scanln(&u2.Email)
	u2.Email = strings.TrimSpace(u2.Email)
	fmt.Println("Enter your password: ")
	fmt.Scanln(&u2.Password)
	u2.Password = strings.TrimSpace(u2.Password)

	var dbEmail, dbPassword string
	err := db.QueryRow("SELECT Email, Password FROM Customers WHERE LOWER(Email) = LOWER(?)", u2.Email).Scan(&dbEmail, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Login failed: Email not found in the database")
		} else {
			fmt.Printf("Login failed: Database query error: %v\n", err)
		}
		return
	}

	if strings.EqualFold(u2.Email, dbEmail) && u2.Password == dbPassword {
		// Retrieve and set the user's name
		err := db.QueryRow("SELECT FirstName FROM Customers WHERE LOWER(Email) = LOWER(?)", u2.Email).Scan(&u2.Name)
		if err != nil {
			fmt.Printf("Login failed: Unable to retrieve user name: %v\n", err)
			return
		}
		fmt.Println("Login successful")
	} else if strings.EqualFold(u2.Email, dbEmail) {
		fmt.Println("Login failed: Invalid password.")
	} else {
		fmt.Println("Login failed: User not recognized.")
	}

	// Automatically create a new account if none exists
	err = db.QueryRow("SELECT AccountID FROM Accounts WHERE CustomerID = (SELECT CustomerID FROM Customers WHERE Email = ?)", u2.Email).Scan(&u2.Name)
	if err == sql.ErrNoRows {
		fmt.Println("Error: No account found for the user. Creating a new account...")
		accountNumber := GenerateAccountNumber()
		_, err = db.Exec("INSERT INTO Accounts (CustomerID, AccountNumber, Balance) VALUES ((SELECT CustomerID FROM Customers WHERE Email = ?), ?, 0)", u2.Email, accountNumber)
		if err != nil {
			fmt.Println("Error creating account for the user:", err)
			return
		}
		fmt.Println("New account created successfully.")
	} else if err != nil {
		fmt.Println("Error retrieving account information:", err)
		return
	}
}

type UserAccount struct {
	AccountID     int
	AccountNumber string
	Balance       float64
	Amount        float64
}

func Deposit(account *UserAccount, db *sql.DB) {
	// Verify account exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Accounts WHERE AccountID = ?)", account.AccountID).Scan(&exists)
	if err != nil || !exists {
		fmt.Println("Error: Account not found")
		return
	}

	fmt.Println("Enter the amount you want to deposit: ")
	fmt.Scanln(&account.Amount)

	if account.Amount <= 0 {
		fmt.Println("Error: Deposit amount must be greater than zero.")
		return
	}

	// Call the DepositMoney stored procedure
	_, err = db.Exec("CALL DepositMoney(?, ?)", account.AccountID, account.Amount)
	if err != nil {
		fmt.Println("Error during deposit:", err)
		return
	}

	// Show updated balance
	err = db.QueryRow("SELECT Balance FROM Accounts WHERE AccountID = ?", account.AccountID).Scan(&account.Balance)
	if err != nil {
		fmt.Println("Deposit successful, but unable to retrieve new balance")
		return
	}
	fmt.Printf("Deposit successful. New balance: %.2f\n", account.Balance)
}

func Withdraw(account *UserAccount, db *sql.DB) {
	fmt.Println("Enter the amount you want to withdraw: ")
	fmt.Scanln(&account.Amount)

	if account.Amount <= 0 {
		fmt.Println("Error: Withdrawal amount must be greater than zero.")
		return
	}

	_, err := db.Exec("CALL WithdrawMoney(?, ?)", account.AccountID, account.Amount)
	if err != nil {
		fmt.Println("Error during withdrawal:", err)
		return
	}
	fmt.Println("Withdrawal successful")
}

func CheckBalance(account *UserAccount, db *sql.DB) {
	err := db.QueryRow("SELECT Balance FROM Accounts WHERE AccountID = ?", account.AccountID).Scan(&account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Error: Account not found")
		} else {
			fmt.Println("Error retrieving balance:", err)
		}
		return
	}
	fmt.Println("Your balance is: ", account.Balance)
}

func TransferMoney(account *UserAccount, db *sql.DB) {
	fmt.Println("Enter the account number of the receiver: ")
	var receiverAccountNumber string
	fmt.Scanln(&receiverAccountNumber)

	if receiverAccountNumber == account.AccountNumber {
		fmt.Println("Error: Cannot transfer to the same account.")
		return
	}
	if strings.TrimSpace(receiverAccountNumber) == "" {
		fmt.Println("Error: Account number cannot be empty.")
		return
	}

	// Resolve the receiver's AccountID from the AccountNumber
	var receiverAccountID int
	err := db.QueryRow("SELECT AccountID FROM Accounts WHERE AccountNumber = ?", receiverAccountNumber).Scan(&receiverAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Error: Receiver's account not found")
		} else {
			fmt.Println("Error retrieving receiver's account information:", err)
		}
		return
	}

	fmt.Println("Enter the amount you want to transfer: ")
	fmt.Scanln(&account.Amount)

	if account.Amount <= 0 {
		fmt.Println("Error: Transfer amount must be greater than zero.")
		return
	}

	// Call the TransferMoney stored procedure with correct parameters
	_, err = db.Exec("CALL TransferMoney(?, ?, ?)", account.AccountID, receiverAccountID, account.Amount)
	if err != nil {
		fmt.Println("Error during transfer:", err)
		return
	}
	fmt.Println("Transfer successful!")

	// Show updated balance
	err = db.QueryRow("SELECT Balance FROM Accounts WHERE AccountID = ?", account.AccountID).Scan(&account.Balance)
	if err != nil {
		fmt.Println("Unable to retrieve new balance")
		return
	}
	fmt.Printf("Your updated balance is: %.2f\n", account.Balance)
}

