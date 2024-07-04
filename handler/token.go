package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gcp-access-token/initial"
	"gcp-access-token/utils"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2/google"
)

type token struct {
	ServiceName string `json:"serviceName"`
	File        string `json:"file"`
	Token       string `json:"token"`
	Expire      int64  `json:"expire"`
}

func Healthcheck(c *fiber.Ctx) error {

	powerBy := make(map[string]string)

	powerBy["Application"] = "Google OAuth 2.0 for Server to Server Applications"
	powerBy["CreatedBy"] = "Mr.Navamin Sawasdee"
	powerBy["Email"] = "navaminsawasdee@gmail.com"

	return c.Status(fiber.StatusOK).JSON(powerBy)
}

func UseCommand(c *fiber.Ctx) error {
	var data []token

	for _, config := range initial.LoadConfig().ServiceAccount {
		log.Printf("%s", "-----------------------------------------------")
		log.Printf("%s to %s", "Project", config.Project)
		log.Printf("%s", "-----------------------------------------------")

		for _, service := range config.Service {
			if _, err := exec.Command("gcloud", "auth", "activate-service-account", "--key-file="+service.File+"").Output(); err != nil {
				log.Printf("[%s] %s - %s", "X", service.File, err)
			} else {
				resp, err := exec.Command("gcloud", "auth", "print-access-token").Output()
				if err != nil {
					log.Printf("[%s] %s - %s", "X", service.File, err)
				}

				now := time.Now()

				log.Printf("%v", token{
					ServiceName: service.Name,
					File:        service.File,
					Token:       string(resp),
					Expire:      now.Add(time.Hour).Unix(),
				})

				data = append(data, token{
					ServiceName: service.Name,
					File:        service.File,
					Token:       string(resp),
					Expire:      now.Add(time.Hour).Unix(),
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func UseLib(c *fiber.Ctx) error {

	var (
		data  []token
		scope = "https://www.googleapis.com/auth/cloud-platform"
	)

	for _, config := range initial.LoadConfig().ServiceAccount {
		log.Printf("%s", "-----------------------------------------------")
		log.Printf("%s to %s", "Project", config.Project)
		log.Printf("%s", "-----------------------------------------------")

		for _, service := range config.Service {
			// Load the service account key JSON file
			sa, err := os.ReadFile(service.File)
			if err != nil {
				log.Printf("[%s] %s - %s", "X", service.File, fmt.Sprintf("Error reading service account file: %v", err))
			} else {
				// Parse the service account key JSON file to get a token source
				config, err := google.JWTConfigFromJSON(sa, scope)
				if err != nil {
					log.Printf("[%s] %s - %s", "X", service.File, fmt.Sprintf("Error creating JWT config: %v", err))
				} else {
					tokenSource := config.TokenSource(context.Background())

					// Get an OAuth2 token
					t, err := tokenSource.Token()
					if err != nil {
						log.Printf("[%s] %s - %s", "X", service.File, fmt.Sprintf("Error getting token: %v", err))
					}

					now := time.Now()

					log.Printf("%v", token{
						ServiceName: service.Name,
						File:        service.File,
						Token:       t.AccessToken,
						Expire:      now.Add(time.Hour).Unix(),
					})

					data = append(data, token{
						ServiceName: service.Name,
						File:        service.File,
						Token:       t.AccessToken,
						Expire:      now.Add(time.Hour).Unix(),
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func NotUseLib(c *fiber.Ctx) error {

	var (
		data  []token
		sa    utils.ServiceAccount
		scope = "https://www.googleapis.com/auth/cloud-platform"
	)

	for _, config := range initial.LoadConfig().ServiceAccount {
		log.Printf("%s", "-----------------------------------------------")
		log.Printf("%s to %s", "Project", config.Project)
		log.Printf("%s", "-----------------------------------------------")

		for _, service := range config.Service {
			// Load the service account key JSON file
			result, err := os.ReadFile(service.File)
			if err != nil {
				log.Printf("[%s] %s - %s", "X", service.File, fmt.Sprintf("Error reading service account file: %v", err))
			} else {
				err = json.Unmarshal(result, &sa)
				if err != nil {
					log.Printf("[%s] %s - %s", "X", service.File, fmt.Sprintf("Error parsing service account file: %v", err))
				} else {
					pk, err := utils.GetPrivateKey(sa.PrivateKey) // GetPrivateKey
					if err != nil {
						log.Printf("[%s] %s - %s", "X", service.File, err)
					} else {
						jwt, err := utils.GenerateJWT(sa, scope, pk) // GenerateJWT
						if err != nil {
							log.Printf("[%s] %s - %s", "X", service.File, err)
						} else {
							act, err := utils.GetAccessToken(sa, jwt) // GetAccessToken
							if err != nil {
								log.Printf("[%s] %s - %s", "X", service.File, err)
							} else {
								now := time.Now()

								log.Printf("%v", token{
									ServiceName: service.Name,
									File:        service.File,
									Token:       act,
									Expire:      now.Add(time.Hour).Unix(),
								})

								data = append(data, token{
									ServiceName: service.Name,
									File:        service.File,
									Token:       act,
									Expire:      now.Add(time.Hour).Unix(),
								})
							}
						}
					}
				}
			}
		}

	}

	return c.Status(fiber.StatusOK).JSON(data)
}
