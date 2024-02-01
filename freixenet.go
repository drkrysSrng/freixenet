package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new Chrome instance
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-web-security", true), // Removing this line
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
		cookieBytes, err := os.ReadFile(cookieFileName)
		if err != nil {
			log.Fatal(err)
		}
		cookieValue = string(cookieBytes)

		// Set the session cookie using JavaScript
		err = chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`document.cookie = "%s"`, cookieValue), nil),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Using existing session cookie: %s\n", cookieValue)
	}

	// Navigate to a webpage
	url := "https://deportesweb.madrid.es/DeportesWeb/login"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the 'h4' element with text "Usuario y contraseña" is visible
	var isH4Visible bool
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`//h4[contains(text(), 'Usuario y contraseña')]`, chromedp.BySearch),
		chromedp.Evaluate(`true`, &isH4Visible),
	)
	if err != nil {
		log.Fatal(err)
	}

	if isH4Visible {

		// Click the login button
		err = chromedp.Run(ctx,
			chromedp.Click(`//h4[contains(text(), 'Usuario y contraseña')]`),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Wait for some time to see the result
		time.Sleep(3 * time.Second)

		// Click the login button
		err = chromedp.Run(ctx,
			chromedp.Click(`#acceso_pass`),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Wait for some time to see the result
		time.Sleep(3 * time.Second)

		// Fill in the username and password fields
		var username, password string
		fmt.Print("Enter username: ")
		fmt.Scan(&username)
		fmt.Print("Enter password: ")
		fmt.Scan(&password)

		err = chromedp.Run(ctx,
			chromedp.SendKeys("#correoelectronico", username, chromedp.ByID),
			chromedp.SendKeys("#contrasenia", password, chromedp.ByID),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Click the login button
		err = chromedp.Run(ctx,
			chromedp.Click(`//button[contains(text(), 'Acceder')]`),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Wait for some time to see the result
		time.Sleep(5 * time.Second)

		// Retrieve and print the session cookie
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.cookie`, &cookieValue),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Save the session cookie to a file
		err = os.WriteFile(cookieFileName, []byte(cookieValue), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Session cookie saved to %s\n", cookieFileName)
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

	// Click the login button
	err = chromedp.Run(ctx,
		chromedp.Click(`//div[contains(text(), 'Sala multitrabajo')]`),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(50 * time.Second)
}
