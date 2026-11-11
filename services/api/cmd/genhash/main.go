package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "owner123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hash for '%s':\n%s\n\n", password, string(hash))

	// Also update DB directly
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/linkchat?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE owners SET password_hash=? WHERE username='owner'", string(hash))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Database updated successfully!")

	// Verify
	var storedHash string
	db.QueryRow("SELECT password_hash FROM owners WHERE username='owner'").Scan(&storedHash)
	fmt.Printf("Stored hash: %s\n", storedHash)

	// Test verify
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Printf("❌ Verification FAILED: %v\n", err)
	} else {
		fmt.Println("✅ Verification PASSED — login will work!")
	}
}
