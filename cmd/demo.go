package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func main() {
	server := NewServer()
	server.Register("admin", "password")
	go http.ListenAndServe(":8080", server)

	fmt.Println("server online")

	// Create a cookie jar to store cookies across requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// register a new user
	err = clientPost(client, "http://localhost:8080/register")
	if err != nil {
		log.Fatal(err)
	}

	// login as the new user
	err = clientPost(client, "http://localhost:8080/login")
	if err != nil {
		log.Fatal(err)
	}

	// check the session
	err = clientSession(client)
	if err != nil {
		log.Fatal(err)
	}

}

func clientPost(c *http.Client, endpoint string) error {
	data := url.Values{}
	data.Set("username", "newuser")
	data.Set("password", "password")
	resp, err := c.PostForm(endpoint, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response status code and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	return nil
}

func clientSession(c *http.Client) error {
	resp, err := c.Get("http://localhost:8080/session")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response status code and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	return nil
}
