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
const NOT_FOUND_ELEMENT = "//*[@id='app']/div/span[2]/div/span/div/div/div"
const INPUT_ELEMENT = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[2]"
const ATTACH_MENU = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[1]/div[2]/div/div"
const INPUT_FILE = "//*[@id='main']/footer/div[1]/div/span[2]/div/div[1]/div[2]/div/span/div/ul/div/div[2]/li/div/input"
const SEND_IMAGE_BUTTON = "//*[@id='app']/div/div/div[3]/div[2]/span/div/span/div/div/div[2]/div/div[2]/div[2]/div/div"
const INICIANDO = "Iniciando conversa"

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
	time.Sleep(2 * time.Second)
	driver.Get(WHATSAPP_URL)
	
	if isWhatsAppLogged(driver) {
		time.Sleep(2 * time.Second)
		fmt.Println("Whatsapp Logged.")
		links := excel.GetLinks()

		for _, link := range links {
			driver.Get(link.Link)
			time.Sleep(2 * time.Second)
			if isNumberInvalid(driver) {
				continue
			}
			time.Sleep(2 * time.Second)
			sendText(driver)
			time.Sleep(2 * time.Second)
			sendImage(driver)

			time.Sleep(5 * time.Second)
		}
		
	} else {
		fmt.Println("Whatsapp not logged")
		return driver
	}

	CloseWhats(driver)
	excel.SaveLog(logs)

	fmt.Println("Ending message service")
	return driver
}

func sendText(driver selenium.WebDriver) {
	element, err := driver.FindElement(selenium.ByXPATH, "//*[@id='main']/footer/div[1]/div/span[2]/div/div[2]/div[1]/div/div[1]")
	if err != nil {
		fmt.Println("Não consigo enviar")
	}
	element.SendKeys(selenium.EnterKey)
}

func sendImage(driver selenium.WebDriver) {
	image, err := filepath.Abs("./src/images/fotos.jpeg")
	fmt.Println(image)
	if err != nil {
		fmt.Println("Não foi possível encontrar a imagem.")
	}
	
	attachBtn, err := driver.FindElement(selenium.ByXPATH, ATTACH_MENU)
	if err != nil {
		fmt.Println("Não encontrou o menu de arquivos")
	}
	attachBtn.Click()
	time.Sleep(2 * time.Second)
	
	inputFile, err := driver.FindElement(selenium.ByXPATH, "//*[@id='main']/footer/div[1]/div/span[2]/div/div[1]/div[2]/div/span/div/ul/div/div[2]/li/div/input")
	if err != nil {
		fmt.Println("Não encontrou o input de arquivo")
	}

	inputFile.SendKeys(image)

	time.Sleep(2 * time.Second)

	sendButton, err := driver.FindElement(selenium.ByXPATH, "//*[@id='app']/div/div/div[3]/div[2]/span/div/span/div/div/div[2]/div/div[1]/div[3]/div/div/div[2]/div[1]/div[1]")
	if err != nil {
		fmt.Println("Não da pra enviar. Não achei onde envia")
	}

	time.Sleep(2 * time.Second)

	sendButton.SendKeys(selenium.EnterKey)


}

func isWhatsAppLogged(driver selenium.WebDriver) bool {
	fmt.Println("Checking if whatsapp is logged.")
	if err := driver.Wait(conditions.ElementIsLocated(selenium.ByID, SIDE_ELEMENT)); err != nil {
		fmt.Println("Whatsapp loading Timeout. Try again.");
		return false
	} else {
		fmt.Println("Side founded")
		return true
	}
}
func isNumberInvalid(driver selenium.WebDriver) bool {
	fmt.Println("Checking if number is valid")
		err := driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, NOT_FOUND_ELEMENT))
		if err != nil {
			fmt.Println("Carregar conversa não encontrado")
			return false
		}

		fmt.Println("Loading conversation carregado")

		for {
			_, err := driver.FindElement(selenium.ByXPATH, NOT_FOUND_ELEMENT)
			if err != nil {
				fmt.Println("Element text not found.")
				return true
			} 

			time.Sleep(2 * time.Second)

			element, err := driver.FindElement(selenium.ByXPATH, "//*[@id='app']/div/span[2]/div/span/div/div/div/div/div/div[1]")
			if err != nil {
				fmt.Println("Erro ao procurar texto de conversa invalida")
				return false
			}

			text, err := element.Text()
			if err != nil {
				fmt.Println("Text not found")
				return false
			}
			if text != INICIANDO {
				if text == "O número de telefone compartilhado através de url é inválido." {
					fmt.Println("Conversa não carregada")
					return true
				} else {
					fmt.Println("Conversa carregada")
					return false
				}
			}
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

func CloseWhats(browser selenium.WebDriver) {
	fmt.Println("Closing Web Whatsapp...")

	menu, err := browser.FindElement(selenium.ByXPATH, "//*[@id='app']/div/div/div[4]/header/div[2]/div/span/div[5]/div")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	menu.Click()
	time.Sleep(2 * time.Second)
	logout, err := browser.FindElement(selenium.ByXPATH, "//*[@id='app']/div/div/div[4]/header/div[2]/div/span/div[5]/span/div/ul/li[6]/div")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	logout.Click()
	time.Sleep(2 * time.Second)
	getOut, err := browser.FindElement(selenium.ByXPATH, "//*[@id='app']/div/span[2]/div/div/div/div/div/div/div[3]/div/button[2]")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	getOut.Click()

	time.Sleep(5 * time.Second)

	fmt.Println("Web Whatsapp closed")
}

