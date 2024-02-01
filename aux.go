package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// Create a new Chrome instance
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-web-security", true),
	)
	ctx, cancel = chromedp.NewExecAllocator(ctx, options...)
	defer cancel()

	// Create a new Chrome browser context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Check if the cookie file exists
	cookieFileName := "cookie.txt"
	var cookieValue string
	if _, err := os.Stat(cookieFileName); err == nil {
		// Read the cookie value from the file
		cookieBytes, err := ioutil.ReadFile(cookieFileName)
		if err != nil {
			log.Fatal(err)
		}
		cookieValue = string(cookieBytes)

		fmt.Printf("Using existing session cookie: %s\n", cookieValue)
	}

	// Check if the 'h4' element with text "Usuario y contraseña" is visible
	var isH4Visible bool
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://deportesweb.madrid.es/DeportesWeb/login"),
		chromedp.WaitVisible(`//h4[contains(text(), 'Usuario y contraseña')]`, chromedp.BySearch),
		chromedp.Evaluate(`true`, &isH4Visible),
	)
	if err != nil {
		log.Fatal(err)
	}

	if isH4Visible {
		// Perform login via HTTP request
		var username, password string
		fmt.Print("Enter username: ")
		fmt.Scan(&username)
		fmt.Print("Enter password: ")
		fmt.Scan(&password)

		// Prepare login request
		loginURL := "https://deportesweb.madrid.es/DeportesWeb/login"
		formData := url.Values{
			"correoelectronico": {username},
			"contrasenia":       {password},
		}

		// Send login request
		response, err := http.PostForm(loginURL, formData)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// Check response status
		if response.StatusCode != http.StatusOK {
			log.Fatalf("Login request failed with status: %s", response.Status)
		}

		// Extract and print cookies from the response
		cookies := response.Cookies()
		for _, cookie := range cookies {
			fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
		}

		// Save cookies to a file or use them as needed
		cookieValue = response.Header.Get("Set-Cookie")
		err = os.WriteFile(cookieFileName, []byte(cookieValue), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Session cookie saved to %s\n", cookieFileName)

		// Rest of your Chromedp interactions here...
	}

	// Get the title of the logged-in page
	var title string
	err = chromedp.Run(ctx,
		chromedp.Title(&title),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print the title
	fmt.Printf("Title: %s\n", title)

	time.Sleep(50 * time.Second)
}
