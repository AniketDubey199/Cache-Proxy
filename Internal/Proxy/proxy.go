package proxy

import (
	"bytes"
	"io"
	"net/http"

	cache "github.com/AniketDubey199/Cache-Proxy/Internal/Cache"
	"github.com/gofiber/fiber/v3"
)

type ProxyHandler struct {
	Origin string
	Cache  *cache.Cache
}

func (p *ProxyHandler) Caching(c fiber.Ctx) error {
	cacheKey := c.Method() + "|" + c.OriginalURL()

	if c.Method() == fiber.MethodGet {
		if cached, ok := p.Cache.Get(cacheKey); ok {
			for k, v := range cached.Headers {
				c.Set(k, v[0])
			}
			c.Set("X-Cache", "HIT")
			return c.Status(cached.StatusCode).Send(cached.Body)
		}
	}

	req, err := http.NewRequest(
		c.Method(),
		p.Origin+c.OriginalURL(),
		bytes.NewReader(c.Body()),
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	for k, values := range c.GetReqHeaders() {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(502).SendString(err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	for k, v := range resp.Header {
		c.Set(k, v[0])
	}
	c.Set("X-Cache", "MISS")

	if c.Method() == fiber.MethodGet && resp.StatusCode == http.StatusOK {
		p.Cache.Set(cacheKey, cache.CacheReponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       body,
		})
	}

	return c.Status(resp.StatusCode).Send(body)
}
