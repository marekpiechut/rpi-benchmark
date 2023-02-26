package main

import (
	"bytes"
	"github.com/h2non/bimg"
	image_size "github.com/marekpiechut/rpi-benchmark/image_size"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"runtime"

	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberJWT "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
)

type ImageData struct {
	Height  int    `json:"height"`
	Width   int    `json:"width"`
	ImgType string `json:"type"`
	User    string `json:"user"`
}

func main() {
	log.Printf("Starting server... (CPU: %d)", runtime.NumCPU())
	requests := 0
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(withJWT())
	app.Use(func(c *fiber.Ctx) error {
		requests++
		if requests%1000 == 0 {
			log.Printf("Requests: %d", requests)
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/identify", func(c *fiber.Ctx) error {
		user := getUser(c)

		body := c.Request().Body()
		img := bimg.NewImage(body)
		size, err := img.Size()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		res := ImageData{
			Height:  size.Height,
			Width:   size.Width,
			ImgType: "png",
			User:    user,
		}
		return c.JSON(res)
	})

	app.Post("/identify2", func(c *fiber.Ctx) error {
		user := getUser(c)

		body := c.Request().Body()
		reader := bytes.NewReader(body)
		m, _, err := image.Decode(reader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		bounds := m.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()

		res := ImageData{
			Height:  h,
			Width:   w,
			ImgType: "png",
			User:    user,
		}
		return c.JSON(res)
	})

	app.Post("/identify3", func(c *fiber.Ctx) error {
		user := getUser(c)

		body := c.Request().Body()
		reader := bytes.NewReader(body)
		size, err := image_size.DetectSize(reader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		res := ImageData{
			Height:  size.Height,
			Width:   size.Width,
			ImgType: "png",
			User:    user,
		}
		return c.JSON(res)
	})

	log.Fatal(app.Listen(":3000"))
}

func getUser(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user, ok := claims["sub"].(string)
	if !ok {
		user = ""
	}
	return user
}

func withJWT() fiber.Handler {
	return fiberJWT.New(fiberJWT.Config{
		SigningKey: []byte("super-secret-key"),
	})
}
