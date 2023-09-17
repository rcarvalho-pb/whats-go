package web

import (
	"fmt"
	"path/filepath"
	"time"
	"whats/internal/excel"

	"github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const DRIVER_PATH = "src/chromedriver/chromedriver"
const WHATSAPP_URL = "https://web.whatsapp.com"
const SIDE_ELEMENT = "side"
const NOT_FOUND_ELEMENT = "//*[@id='app']/div/span[2]/div/span/div/div/div/div"
const INPUT_ELEMENT = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[2]"
const CLICK = "MOUSE"
const ENTER = "ENTER"
const ATTACH_MENU = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[1]/div[2]/div/div"
const INPUT_FILE = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[1]/div[2]/div/span/div/ul/div/div[2]/li/div/input"
const SEND_IMAGE_BUTTON = "//*[@id='app']/div/div/div[3]/div[2]/span/div/span/div/div/div[2]/div/div[2]/div[2]/div/div"
const BUTTON_OK = "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[2]/div/button/div/div"

func StartService() {
	if service, err := selenium.NewChromeDriverService(DRIVER_PATH, 4444); err != nil {
		panic(err)
	} else {
		fmt.Println("Starting application...")
		driver := sendMessages(service)
		defer driver.Close() 
		defer service.Stop()
		fmt.Println("Ending application!")
	}

}

func sendMessages(service *selenium.Service) selenium.WebDriver {
	fmt.Println("Starting SendMessages...")
	var logs []excel.Log

	driver := getWebDriver()
	driver.Get("https://www.google.com")
	driver.Get(WHATSAPP_URL)

	links := excel.GetLinks()

	if isElementLoaded(driver, selenium.ByID, SIDE_ELEMENT) {
		for i, link := range links {
			fmt.Printf("Starting %dº message...\n", i + 1)
			driver.Get(link.Link)

			if isElementLoaded(driver, selenium.ByXPATH, NOT_FOUND_ELEMENT) {
				logs = append(logs, excel.Log{OS: link.OS, Contact: "NÃO"})
				continue
			} else {
				fmt.Println("Element not found not fount")
			}

			if isElementLoaded(driver, selenium.ByXPATH, INPUT_ELEMENT) {
				fmt.Println("The main element was found")
				sendContent(driver)
				logs = append(logs, excel.Log{OS: link.OS, Contact: "SIM"})
			}
			
			fmt.Printf("Ending %d message...\n", i + 1)
			time.Sleep(5 * time.Second)
		} 
	} else {
		return driver
	}

	excel.SaveLog(logs)

	fmt.Println("Ending message service")
	return driver
}

func sendContent(driver selenium.WebDriver) {
	fmt.Println("Starting send content")
	send(driver, INPUT_ELEMENT, ENTER)
	if attachMenu, err := driver.FindElement(selenium.ByXPATH, ATTACH_MENU); err == nil {
		attachMenu.Click()
		time.Sleep(10 * time.Second)
		if inputElement, err := driver.FindElement(selenium.ByXPATH, INPUT_FILE); err == nil {
			if absPath, err := filepath.Abs("./src/images/foto.jpeg"); err == nil {
				inputElement.SendKeys(absPath)
				if sendImage, err := driver.FindElement(selenium.ByXPATH, SEND_IMAGE_BUTTON); err == nil {
					sendImage.SendKeys(selenium.EnterKey)
					fmt.Println("Ending send content")
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}

	} else {
		panic(err)
	}
}

func send(driver selenium.WebDriver, elementXPath, typeOfSending string) {
	fmt.Println("Starting send")
	if typeOfSending == CLICK {
		if isElementLoaded(driver, selenium.ByXPATH, elementXPath) {
			if elementSearched, err := driver.FindElement(selenium.ByXPATH, elementXPath); err == nil {
				elementSearched.Click()
			} else {
				panic(err)
			}
		}
	}

	if typeOfSending == ENTER {
		if isElementLoaded(driver, selenium.ByXPATH, elementXPath) {
			if elementSearched, err := driver.FindElement(selenium.ByXPATH, elementXPath); err == nil {
				elementSearched.SendKeys(selenium.EnterKey)
			} else {
				panic(err)
			}
		}
	}
	fmt.Println("Ending send")
}

func isElementLoaded(driver selenium.WebDriver, typeOfSearch, elementSearched string) bool {
	if elementSearched == NOT_FOUND_ELEMENT {
		if err := driver.Wait(conditions.ElementIsLocated(typeOfSearch, elementSearched)); err != nil {
			fmt.Printf("Element NotFound not found\n")
		}
		time.Sleep(5 * time.Second)
		_, err := driver.FindElement(selenium.ByXPATH, "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[2]/div/button/div/div");
		if err != nil {
			return false
		} else {
			return true
		}
	}
	else {
		
	} 
}

func getWebDriver() selenium.WebDriver {
	fmt.Println("Starting get WebDriver")
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
		panic(err)
	} else {
		fmt.Println("Ending get WebDriver")
		return driver
	}
}

