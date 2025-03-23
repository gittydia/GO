package main

import (
	"bank/function"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Establish database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	var u1 function.SigninU
	var u2 function.LoginU
	var account function.UserAccount

	var choice int

	for {
		fmt.Println("\nChoose an option: ")
		fmt.Println("1. Signin")
		fmt.Println("2. Login")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}
		if choice == 1 {
			function.Signin(&u1, db)
		} else if choice == 2 {
			function.Login(&u2, db)
			if u2.Email != "" && u2.Name != "" { // Ensure login was successful
				fmt.Println("Welcome to the bank" + ", " + u2.Email + "!")

				// Retrieve AccountID for the logged-in user
				err := db.QueryRow("SELECT AccountID FROM Accounts WHERE CustomerID = (SELECT CustomerID FROM Customers WHERE Email = ?)", u2.Email).Scan(&account.AccountID)
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("Error: No account found for the user. Creating a new account...")
						_, err = db.Exec("INSERT INTO Accounts (CustomerID, Balance) VALUES ((SELECT CustomerID FROM Customers WHERE Email = ?), 0)", u2.Email)
						if err != nil {
							fmt.Println("Error creating account for the user:", err)
							continue
						}
						// Retrieve the newly created AccountID
						err = db.QueryRow("SELECT AccountID FROM Accounts WHERE CustomerID = (SELECT CustomerID FROM Customers WHERE Email = ?)", u2.Email).Scan(&account.AccountID)
						if err != nil {
							fmt.Println("Error retrieving newly created account information:", err)
							continue
						}
						fmt.Println("New account created successfully.")
					} else {
						fmt.Println("Error retrieving account information:", err)
						continue
					}
				}

				for {
					fmt.Println("\nChoose an option: ")
					fmt.Println("1. Deposit")
					fmt.Println("2. Withdraw")
					fmt.Println("3. Check balance")
					fmt.Println("4. Exit")
					fmt.Print("Enter your choice: ")
					_, err := fmt.Scanln(&choice)
					if err != nil {
						fmt.Println("Invalid input. Please try again.")
						continue
					}
					if choice == 1 {
						function.Deposit(&account, db)
					} else if choice == 2 {
						function.Withdraw(&account, db)
					} else if choice == 3 {
						function.CheckBalance(&account, db)
					} else if choice == 4 {
						fmt.Println("Exiting...")
						break
					} else {
						fmt.Println("Invalid choice")
					}
				}
			} else {
				fmt.Println("Login failed: User not recognized.")
			}
		} else if choice == 3 {
			fmt.Println("Exiting...")
			break
		} else {
			fmt.Println("Invalid choice")
		}
	}
}