package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/proxy/protocol"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	rec "github.com/gofiber/fiber/v2/middleware/recover"
)

func ListenService(opt *settings.Options, log *log.Logger) (outerr error) {

	//& Fiber doesn't always return an error so i'd rather have it catch the panic here,
	//& and since we panic on error anyways atleast we know where it crashed (roughly)

	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("proxy: fiber: %v", r)
		}
	}()

	if len(opt.Service.Proxies) == 0 {
		return fmt.Errorf("proxy: no targets were defined in env config")
	}

	upstreams := make(map[string]string, len(opt.Service.Proxies))
	for target, host := range opt.Service.Proxies {
		addr, err := protocol.EnsureScheme(host)
		if err != nil {
			return fmt.Errorf("proxy: %w", err)
		}
		upstreams[target] = addr
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:          web.ErrorHandler,
		DisableStartupMessage: true,
	})

	app.Use(rec.New())
	app.Use(protocol.XPoweredBy())
	app.Use(protocol.RequestLogging())

	app.Use(func(c *fiber.Ctx) error {
		host := c.Hostname()
		target, ok := upstreams[host]
		if !ok {
			return web.BadGatewayError(c, &web.Details{
				Message: "invalid upstream",
				Err:     fmt.Errorf("proxy: no target found for host: %s", host),
			})
		}

		c.Request().Header.Add("X-Forwarded-IP", c.IP())
		c.Request().Header.Set("X-Forwarded-Host", host)
		c.Request().Header.Set("X-Forwarded-Proto", "http")

		err := proxy.Do(c, fmt.Sprintf("%s%s", target, c.OriginalURL()))
		if err != nil {
			return web.BadGatewayError(c, &web.Details{
				Message: "invalid proxy to target",
				Err:     fmt.Errorf("proxy: fasthttp: %w", err),
			})
		}

		return nil
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		log.Info("HTTP Listener", "Listening on 0.0.0.0:%d", opt.Service.Ports["proxy"])
		return nil
	})

	if err := app.Listen(fmt.Sprintf(":%d", opt.Service.Ports["proxy"])); err != nil {
		return fmt.Errorf("web: fiber: %w", err)
	}

	return nil
}
