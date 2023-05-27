package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

const (
	// WebDriverのパスとSelenium Serverのパスは環境に応じて修正してください。
	seleniumPath     = "resource/selenium-server-standalone-3.141.59.jar"
	chromeDriverPath = "resource/chromedriver"
	targetURL        = "https://redeem.tcg.pokemon.com/en-us/"
	port             = 8080
)

type JsonData struct {
	Values []string `json:"values"`
}

func waitForElement(wd selenium.WebDriver, locator, value string) (selenium.WebElement, error) {
	var element selenium.WebElement
	var err error

	for i := 0; i < 10; i++ { // Adjust the number of attempts as needed.
		element, err = wd.FindElement(locator, value)
		if err == nil {
			break
		}

		fmt.Printf("search count: %v\n", i)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("element not found: %v", err)
	}

	return element, nil
}

func waitForElements(wd selenium.WebDriver, locator, value string, index int) (selenium.WebElement, error) {
	var element selenium.WebElement
	var err error

	for i := 0; i < 10; i++ {
		elements, err := wd.FindElements(locator, value)
		fmt.Printf("elements: %v\n", elements)
		if err == nil && len(elements) > 0 {
			element = elements[index]
			break
		}

		fmt.Printf("search count: %v\n", i)
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
		fmt.Errorf("Error loading .env file")
	}

	os.Setenv("PATH", os.Getenv("PATH")+":"+filepath.Dir(chromeDriverPath))
}

func main() {
	fmt.Printf("program started.\n")
	// .envファイルから環境変数を取得する例
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

	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"args": []string{
				"--headless",
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

	// WebDriverを利用して特定のURLにアクセスします。
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

	fmt.Printf("wrote the source to file.")

	fmt.Printf("loading json.\n")
	// 以下のパスは適宜変更してください。
	err = filepath.Walk("targets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			fileBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}

			jsonData := JsonData{}
			err = json.Unmarshal(fileBytes, &jsonData)
			if err != nil {
				return fmt.Errorf("failed to decode JSON: %v", err)
			}

			for i, val := range jsonData.Values {
				fmt.Printf("search element.\n")
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

				// Submit the code.
				submitButton, err := waitForElements(wd, selenium.ByClassName, "Button_blueButton__1PlZZ", 0)
				if err != nil {
					return fmt.Errorf("failed to find button: %v", err)
				}
				err = submitButton.Click()
				if err != nil {
					return fmt.Errorf("failed to click button: %v", err)
				}

				// Wait for 3 seconds.
				time.Sleep(3 * time.Second)

				if i > 0 && i%10 == 0 {
					// Click the next button every 10 codes.
					nextButton, err := waitForElements(wd, selenium.ByClassName, "Button_blueButton__1PlZZ", 1)
					if err != nil {
						return fmt.Errorf("failed to find next button: %v", err)
					}
					err = nextButton.Click()
					if err != nil {
						return fmt.Errorf("failed to click next button: %v", err)
					}

					// Wait for 3 seconds before continuing to the next page.
					time.Sleep(3 * time.Second)
				}
			}
		}
		return nil
	})
	fmt.Printf("loaded json.\n")
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", "targets", err)
	}

	fmt.Printf("program ended.\n")
}
