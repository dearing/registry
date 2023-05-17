package registry

import (
	"encoding/gob"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Account represents a user account
type Account struct {
	Password []byte // hashed password for client
}

// Registry is a database of accounts
type Registry struct {
	Accounts map[string]Account
}

// NewRegistry creates a new account database
func NewRegistry() *Registry {
	return &Registry{
		Accounts: make(map[string]Account),
	}
}

// Register creates a new account and saves it to the database
func (r *Registry) Register(username, password string) error {

	if _, ok := r.Accounts[username]; ok {
		return fmt.Errorf("account already exists: %v", username)
	}

	// generate the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	r.Accounts[username] = Account{
		Password: hashedPassword,
	}

	return err
}

// Unregister deletes an existing account from the database
func (r *Registry) Unregister(username string) error {

	if _, ok := r.Accounts[username]; !ok {
		return fmt.Errorf("account not found: %v", username)
	}

	delete(r.Accounts, username)

	return nil
}

// Verify checks the username and password against the database
func (r *Registry) Verify(username, password string) error {

	account, ok := r.Accounts[username]
	if !ok {
		return fmt.Errorf("account not found: %v", username)
	}

	err := bcrypt.CompareHashAndPassword(account.Password, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

// Save takes a writer and saves the account database to it
func (r *Registry) Save(writer io.Writer) error {

	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(r.Accounts)
	if err != nil {
		return err
	}

	return nil
}

// Load takes a reader and loads the account database from it
func (r *Registry) Load(reader io.Reader) error {

	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&r.Accounts)
	if err != nil {
		return err
	}

	return nil
}
