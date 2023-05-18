package main

import (
	"fmt"
	"time"

	"github.com/dearing/registry"
)

func main() {

	// register sample account
	accounts := registry.NewRegistry()
	accounts.Register("admin", "admin")

	fmt.Println("Secure Area; please login.")

	var username, password string

	// get username
	fmt.Print("Enter username:")
	_, err := fmt.Scanln(&username)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// get password
	fmt.Print("Enter password:")
	_, err = fmt.Scanln(&password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// verify username and password
	err = accounts.Verify(username, password)
	if err != nil {
		fmt.Println("Invalid username or password")
		return
	}

	// create a new JWT token, set expiry to +1 hour
	jwt := registry.NewJWT(username, time.Now().Add(time.Hour).Unix())

	fmt.Printf("Welcome back %s, your login will expire at %s\n\n", jwt.Payload.Subject, jwt.Payload.ExpiresAt.HumanReadable())

}
