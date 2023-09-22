package main

import (
	"fmt"
	"time"
	// "time"

	"github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const DRIVER_PATH = "src/chromedriver/chromedriver"
const WHATSAPP_URL = "https://web.whatsapp.com"
const WHATSAPP_URL1 = "https://web.whatsapp.com/send?phone=558393418877&text=%22Ol%C3%A1,%20Ramon.%20Tudo%20bem?%20Gostaria%20de%20confirmar%20se%20o%20seu%20endere%C3%A7o%20%C3%A9%20realmente%20o%20que%20est%C3%A1%20aqui.%22"
const WHATSAPP_URL2 = "https://web.whatsapp.com/send?phone=5583993418572&text=%22Ol%C3%A1,%20Ramon.%20Tudo%20bem?%20Gostaria%20de%20confirmar%20se%20o%20seu%20endere%C3%A7o%20%C3%A9%20realmente%20o%20que%20est%C3%A1%20aqui.%22"

const INICIANDO = "Iniciando conversa"
const SIDE_ELEMENT = "side"
const LOADING_CONVERSATION = "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[1]"
const NOT_FOUND_TEXT = "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[1]"
const NUM_INVALIDO = "O número de telefone compartilhado através de url é inválido."

func teste() {
	
	if service, err := selenium.NewChromeDriverService(DRIVER_PATH, 4444); err != nil {
		panic(err)
	} else {
		fmt.Println("Starting application...")
		defer service.Stop()
		caps := selenium.Capabilities {}
		caps.AddChrome(chrome.Capabilities{
			Args: []string {
				"window-size=1920x1080",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"disable-gpu",
				// "--headless",
			},
		})
		if driver, err := selenium.NewRemote(caps, ""); err != nil {
			fmt.Println("Cant start driver")
		} else {
			driver.Get("https://www.google.com")
			driver.Get(WHATSAPP_URL)

			if !isElementLoaded(driver, selenium.ByID, SIDE_ELEMENT) {
				fmt.Println("not found")
			}

			links := []string {WHATSAPP_URL1, WHATSAPP_URL2}

			for _, link := range links {
				driver.Get(link)

				if isElementLoaded(driver, selenium.ByXPATH, LOADING_CONVERSATION) {
					fmt.Println("Conversa carregada")
					time.Sleep(5*time.Second)
					} else {
					time.Sleep(5*time.Second)
					fmt.Println("Conversa nao carregada")
				}
			}

			
			
		}


		fmt.Println("Ending application!")
	}
}

func isElementLoaded(driver selenium.WebDriver, typeOfSearch, elementSearched string) bool {
	if elementSearched != SIDE_ELEMENT {
		err := driver.Wait(conditions.ElementIsLocated(typeOfSearch, elementSearched))
		if err != nil {
			fmt.Println("Loading element not found.")
			return false
		}
		for {
			element, err := driver.FindElement(typeOfSearch, elementSearched)
			if err != nil {
				fmt.Println("Element text not found.")
				return true
			} 

			text, err := element.Text()
			if err != nil {
				fmt.Println("Text not found")
			}
			if text != INICIANDO {
				if text == NUM_INVALIDO {
					return false
				}
			}
		}

	}

	if err := driver.Wait(conditions.ElementIsLocated(typeOfSearch, elementSearched)); err != nil {
		fmt.Println("Whatsapp loading Timeout. Try again.");
		return false
	} else {
		fmt.Println("Side founded")
		return true
	}
}
