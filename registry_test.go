package registry

import (
	"testing"

	"github.com/spf13/afero"
)

var username, password = "admin", "password"

// TestRegistry tests the account registry from end to end
func TestRegistry(t *testing.T) {

	// Setup a new account registry
	registy := NewRegistry()

	// Setup a new memory filesystem
	fs := afero.NewMemMapFs()

	err := registy.Register(username, password)
	if err != nil {
		t.Error("Registry should not error when adding account")
	}

	// Get a writer to the account database file
	writer, err := fs.Create("accounts.db")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = registy.Save(writer)
	if err != nil {
		t.Error("Registry should not error when saving")
	}

	// Blank out the account registry
	registy = NewRegistry()

	// Get a reader to the account database file
	reader, err := fs.Open("accounts.db")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = registy.Load(reader)
	if err != nil {
		t.Error("Registry should not error when loading")
	}

	if len(registy.Accounts) != 1 {
		t.Error("Registry should have exactly one account")
	}

	if _, ok := registy.Accounts[username]; !ok {
		t.Error("Registry should have created account with username:", username)
		for k, v := range registy.Accounts {
			t.Logf("%s => %s", k, v.Password)
		}
	}

	if err := registy.Verify(username, password); err != nil {
		t.Error("Validate should not error for valid password")
	}

	if err := registy.Verify(username, "!password"); err == nil {
		t.Error("Validate should error for invalid password")
	}

	if err := registy.Unregister(username); err != nil {
		t.Error("Registry should not error when deleting account")
	}

	if len(registy.Accounts) != 0 {
		t.Error("Registry should not have any accounts")
	}

	err = registy.Register(username, password)
	if err == nil {
		t.Error("Registry should error when adding extant account")
	}

}

// BenchmarkRegistry benchmarks the account registry
func BenchmarkRegistryAdd(b *testing.B) {

	// Setup a new account registry
	registy := NewRegistry()

	// Benchmark creating a new account
	for i := 0; i < b.N; i++ {
		registy.Register(username, password)
	}
}

// BenchmarkRegistryValidate benchmarks the account registry
func BenchmarkRegistryValidate(b *testing.B) {

	// Setup a new account registry
	registy := NewRegistry()

	// Create a new account
	registy.Register(username, password)

	// Benchmark validating the account password
	for i := 0; i < b.N; i++ {
		registy.Verify(username, password)
	}
}
