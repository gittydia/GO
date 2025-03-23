CREATE DATABASE bank;

USE bank;

-- create costumer table

CREATE TABLE IF NOT EXISTS Customers (
	CustomerID INT AUTO_INCREMENT PRIMARY KEY,
    FirstName VARCHAR(50) NOT NULL,
    LastName VARCHAR(50) NOT NULL,
    Email VARCHAR(100) NOT NULL,
    Password VARCHAR(100) NOT NULL,
    PhoneNumber VARCHAR(11),
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE IF NOT EXISTS Accounts(
	AccountID INT AUTO_INCREMENT PRIMARY KEY,
    CustomerID INT NOT NULL,
    AccountNumber VARCHAR(20) UNIQUE NOT NULL,
    Balance DECIMAL(15,2) DEFAULT 0.00,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(CustomerID) REFERENCES Customers(CustomerID)
);

CREATE TABLE IF NOT EXISTS TRANSACTIONS(
	TransactionID INT AUTO_INCREMENT PRIMARY KEY,
    AccountID INT NOT NULL,
    TransactionType ENUM('DEPOSIT', 'WITHDRAW', 'TRANSFER') NOT NULL,
    Amount DECIMAL(15,2) NOT NULL,
    TransactionsDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    RelatedAccountID INT,
    FOREIGN KEY(AccountID) REFERENCES Accounts(AccountID),
    FOREIGN KEY(RelatedAccountID) REFERENCES Accounts(AccountID)
);

-- store procedure for DEPOSIT
USE bank; -- Ensure the correct database is selected
DELIMITER //
DROP PROCEDURE IF EXISTS DepositMoney; -- Ensure no duplicate procedure
CREATE PROCEDURE DepositMoney(
	IN p_AccountID INT,
    IN p_Amount DECIMAL(15, 2)
)
BEGIN
	START TRANSACTION;
    
    -- Update account balance
    UPDATE Accounts
    SET Balance = Balance + p_Amount
    WHERE AccountID = p_AccountID;
    
    -- Record the transaction
    INSERT INTO Transactions (AccountID, TransactionType, Amount)
    VALUES (p_AccountID, 'DEPOSIT', p_Amount);
    
    COMMIT;
END //
DELIMITER ;

-- stored procedure for WITHDRAW
DELIMITER //
CREATE PROCEDURE WithdrawMoney(
	IN p_AccountID INT,
    IN p_Amount DECIMAL(15, 2)
)
BEGIN
	DECLARE v_Balance DECIMAL(15, 2);
    START TRANSACTION;
    
    -- Check if the account has sufficient balance
    SELECT Balance INTO v_Balance
    FROM Accounts
    WHERE AccountID = p_AccountID;

    IF v_Balance >= p_Amount THEN
        -- Update the account balance
        UPDATE Accounts
        SET Balance = Balance - p_Amount
        WHERE AccountID = p_AccountID;

        -- Record the transaction
        INSERT INTO Transactions (AccountID, TransactionType, Amount)
        VALUES (p_AccountID, 'WITHDRAW', p_Amount);

        COMMIT;
    ELSE
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient balance';
    END IF;
END //
DELIMITER ;

-- Stored Procedure for Transferring Money
DELIMITER //
DROP PROCEDURE IF EXISTS TransferMoney; -- Ensure no duplicate procedure
CREATE PROCEDURE TransferMoney(
    IN p_FromAccountID INT,
    IN p_ToAccountID INT,
    IN p_Amount DECIMAL(15, 2)
)
BEGIN
    DECLARE v_FromBalance DECIMAL(15, 2);

    START TRANSACTION;

    -- Check if the source account has sufficient balance
    SELECT Balance INTO v_FromBalance
    FROM Accounts
    WHERE AccountID = p_FromAccountID;

    IF v_FromBalance >= p_Amount THEN
        -- Deduct from the source account
        UPDATE Accounts
        SET Balance = Balance - p_Amount
        WHERE AccountID = p_FromAccountID;

        -- Add to the destination account
        UPDATE Accounts
        SET Balance = Balance + p_Amount
        WHERE AccountID = p_ToAccountID;

        -- Record the transaction for the source account
        INSERT INTO Transactions (AccountID, TransactionType, Amount, RelatedAccountID)
        VALUES (p_FromAccountID, 'TRANSFER', p_Amount, p_ToAccountID);

        -- Record the transaction for the destination account
        INSERT INTO Transactions (AccountID, TransactionType, Amount, RelatedAccountID)
        VALUES (p_ToAccountID, 'TRANSFER', p_FromAccountID);

        COMMIT;
    ELSE
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient balance';
    END IF;
END //
DELIMITER ;

ALTER TABLE Customers CHANGE COLUMN Passwords Password VARCHAR(100) NOT NULL;
ALTER TABLE Customers ADD COLUMN Password VARCHAR(100) NOT NULL;