package web

import (
	"errors"
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"

	"github.com/gofiber/fiber/v2"
	rec "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars/v2"

	"github.com/luxploit/red"
)

type HttpUtils struct {
}

func New() (*HttpUtils, error) {
	return &HttpUtils{}, nil
}

func (h *HttpUtils) NewFiber() (outapp *fiber.App, outerr error) {

	// Fiber doesn't always return an error so i'd rather have it catch the panic here,
	// and since we panic on error anyways atleast we know where it crashed (roughly)
	defer func() {
		if r := recover(); r != nil {
			outapp = nil
			outerr = fmt.Errorf("web: fiber: %v", r)
		}
	}()

	engine := handlebars.New("./public", ".hbs")
	app := fiber.New(fiber.Config{
		Views:                 engine,
		ErrorHandler:          ErrorHandler,
		DisableStartupMessage: true,
	})

	app.Use(rec.New())
	app.Use(RequestLogging())
	app.Use(XPoweredBy())

	return app, nil
}

func (h HttpUtils) StartFiber(app *fiber.App, port int) (outerr error) {

	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return fmt.Errorf("web: %w", err)
	}

	// Fiber doesn't always return an error so i'd rather have it catch the panic here,
	// and since we panic on error anyways atleast we know where it crashed (roughly)
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("web: fiber: %v", r)
		}
	}()

	app.Use(func(c *fiber.Ctx) error {
		return NotFoundError(c, &Details{
			Message: "seems like you took a wrong turn",
			Err:     errors.New("invalid resource"),
			Context: fiber.Map{
				"method": c.Method(),
			},
		})
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		logger.Info("HTTP Listener", "Listening on 0.0.0.0:%d", port)
		return nil
	})

	if err = app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		return fmt.Errorf("web: fiber: %w", err)
	}

	return nil
}
