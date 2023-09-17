package teste

import (
	"fmt"
	"time"

	"github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const ELEMENT = "//*[@id='app']/div/div/div[3]/div[1]/div/div/div[2]/div"

const DRIVER_PATH = "src/chromedriver/chromedriver"

func StartService() {
	if service, err := selenium.NewChromeDriverService(DRIVER_PATH, 4444); err != nil {
		panic(err)
	} else {
		// defer service.Stop()
		fmt.Println(service)
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

		if driver, err := selenium.NewRemote(caps, ""); err == nil {
			driver.Get("https://web.whatsapp.com")
			if err := driver.Wait(conditions.ElementIsLocated(selenium.ByID, "side")); err != nil {
				fmt.Println("not founded")
			}
			fmt.Println("side element founded")
			if err = driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, "//*[@id='app']/div/span[2]/div/span/div/div/div/div")); err == nil {
				fmt.Println("Not foundend element Founded")
				time.Sleep(3 * time.Second)
				element, err := driver.FindElement(selenium.ByXPATH, "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[2]/div/button/div/div");
				if err != nil {
					fmt.Println("Button ok not found")
				}
				fmt.Println("Button ok found")
				element.Click()
				
				// if err = driver.Wait(conditions.ElementTextContains())
			}
			fmt.Print("Not found element not founded")

		} else {
			panic(err)
		}
	}
}