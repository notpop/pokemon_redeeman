package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

const (
	seleniumPath     = "resource/selenium-server-standalone-3.141.59.jar"
	chromeDriverPath = "resource/chromedriver"
	targetURL        = "https://redeem.tcg.pokemon.com/en-us/"
	port             = 8080
	disRedeemableSrc = "/static/media/invalid.33aae293.svg"
)

type JsonData struct {
	Values []string `json:"values"`
}

func waitForElement(wd selenium.WebDriver, locator, value string) (selenium.WebElement, error) {
	var element selenium.WebElement
	var err error

	for i := 0; i < 10; i++ {
		element, err = wd.FindElement(locator, value)
		if err == nil {
			break
		}

		fmt.Printf("search %v count: %v\n", value, i + 1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("element not found: %v", err)
	}

	return element, nil
}

func waitForElementToBeClickable(wd selenium.WebDriver, locator, value string, index int) (selenium.WebElement, error) {
	return waitGetElementForElements(wd, locator, value, index, false)
}

func waitForLastElementToBeClickable(wd selenium.WebDriver, locator, value string) (selenium.WebElement, error) {
	return waitGetElementForElements(wd, locator, value, -1, true)
}

func waitGetElementForElements(wd selenium.WebDriver, locator, value string, index int, last bool) (selenium.WebElement, error) {
	var (
		element selenium.WebElement
		err     error
	)

	for i := 0; i < 10; i++ {
		elements, err := wd.FindElements(locator, value)
		if err != nil {
			return nil, fmt.Errorf("element not found: %v", err)
		}

		if len(elements) > 0 {
			if last {
				index = len(elements) - 1
			}

			if len(elements) > index {
				element = elements[index]
				enabled, err := element.IsEnabled()
				if err == nil && enabled {
					break
				}
			}
		}

		fmt.Printf("search %v count: %v\n", value, i + 1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("element not found: %v", err)
	}

	if element == nil {
		return nil, fmt.Errorf("element is empty")
	}

	return element, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to loading .env file: %v", err)
	}

	os.Setenv("PATH", os.Getenv("PATH")+":"+filepath.Dir(chromeDriverPath))
}

func main() {
	fmt.Printf("program started.\n")

	pokemonId := os.Getenv("ACCESS_POKEMON_ID")
	pokemonPassword := os.Getenv("ACCESS_POKEMON_PASSWORD")

	opts := []selenium.ServiceOption{}
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		fmt.Printf("Failed to start selenium service: %v", err)
		return
	}
	defer service.Stop()
	fmt.Printf("selenium setting...\n")

	// Headless option does not work properly
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"args": []string{
				// "--headless",
				"--disable-gpu",
				"--no-sandbox",
				"--window-size=1280x800",
			},
		},
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Printf("Failed to open session: %v", err)
		return
	}
	defer wd.Quit()
	fmt.Printf("selenium started.\n")

	err = wd.Get(targetURL)
	if err != nil {
		fmt.Printf("Failed to access website: %v", err)
		return
	}
	fmt.Printf("target accessed.\n")

	fmt.Printf("search email button.\n")
	emailInput, err := waitForElement(wd, selenium.ByID, "email")
	if err != nil {
		fmt.Printf("failed to find element: %v", err)
	}
	fmt.Printf("search password button.\n")
	passwordInput, err := waitForElement(wd, selenium.ByID, "password")
	if err != nil {
		fmt.Printf("failed to find element: %v", err)
	}

	fmt.Printf("set email.\n")
	err = emailInput.SendKeys(pokemonId)
	if err != nil {
		fmt.Printf("failed to send key: %v", err)
	}
	fmt.Printf("set password.\n")
	err = passwordInput.SendKeys(pokemonPassword)
	if err != nil {
		fmt.Printf("failed to send key: %v", err)
	}

	fmt.Printf("search accept button.\n")
	loginButton, err := waitForElement(wd, selenium.ByID, "accept")
	if err != nil {
		fmt.Printf("failed to find element: %v", err)
	}
	fmt.Printf("click accept button.\n")
	err = loginButton.Click()
	if err != nil {
		fmt.Printf("failed to click element: %v", err)
	}

	url, err := wd.CurrentURL()
	if err != nil {
		fmt.Printf("failed to get current url: %v", err)
	}
	fmt.Printf("redirected to target page => %s\n", url)

	err = os.MkdirAll("response", os.ModePerm)
	if err != nil {
		fmt.Println("failed to make directory:", err)
	}

	source, err := wd.PageSource()
	if err != nil {
		fmt.Printf("failed to get source: %v", err)
	}

	err = ioutil.WriteFile("response/response.html", []byte(source), 0644)
	if err != nil {
		fmt.Println("failed to write html file:", err)
	}
	fmt.Printf("wrote the source to file.\n")

	cookieAcceptButton, err := waitForElement(wd, selenium.ByID, "onetrust-accept-btn-handler")
	if err != nil {
		fmt.Printf("failed to find the cookie accept button: %v", err)
	} else {
		err = cookieAcceptButton.Click()
		if err != nil {
			fmt.Printf("failed to click the cookie accept button: %v", err)
		}
	}
	fmt.Printf("cookie accepted.\n")

	fmt.Printf("loading txt.\n")
	err = filepath.Walk("targets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			lines := []string{}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			for i, val := range lines {
				fmt.Printf("count: %v\n", i + 1)
				fmt.Printf("search id=code element.\n")
				elem, err := waitForElement(wd, selenium.ByID, "code")
				if err != nil {
					return fmt.Errorf("failed to find element: %v", err)
				}

				err = elem.Clear()
				if err != nil {
					return fmt.Errorf("failed to clear element: %v", err)
				}

				err = elem.SendKeys(val)
				if err != nil {
					return fmt.Errorf("failed to send keys: %v", err)
				}

				submitButton, err := waitForElementToBeClickable(wd, selenium.ByXPATH, "//button[contains(@class, 'Button_blueButton__1PlZZ')]", 0)
				if err != nil {
					return fmt.Errorf("failed to find button: %v", err)
				}
				err = submitButton.Click()
				if err != nil {
					return fmt.Errorf("failed to click button: %v", err)
				}
				fmt.Printf("submitted. code: %s\n", val)

				time.Sleep(2 * time.Second)

				deleteButton, _ := waitForLastElementToBeClickable(wd, selenium.ByXPATH, "//*[contains(@class, 'RedeemModule_tdDelete__2-YLO')]")
				if deleteButton != nil {
					img, _ := deleteButton.FindElement(selenium.ByTagName, "img")
					if img != nil {
						src, _ := img.GetAttribute("src")
						if src == disRedeemableSrc {
							if err := deleteButton.Click(); err != nil {
								return fmt.Errorf("found but failed to click delete button: %v", err)
							}
							fmt.Printf("delete button clicked.\n")
							time.Sleep(2 * time.Second)
						}
					}
				}

				if i == len(lines)-1 || (i > 0 && i%10 == 0) {
					nextButton, err := waitForElementToBeClickable(wd, selenium.ByXPATH, "//button[contains(@class, 'Button_blueButton__1PlZZ')]", 1)
					if err != nil {
						return fmt.Errorf("failed to find next button: %v", err)
					}
					err = nextButton.Click()
					if err != nil {
						return fmt.Errorf("failed to click next button: %v", err)
					}
					fmt.Printf("redeemed.\n")

					time.Sleep(3 * time.Second)
				}
			}
		}
		return nil
	})
	fmt.Printf("loaded txt.\n")
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", "targets", err)
	}

	fmt.Printf("program ended.\n")
}
