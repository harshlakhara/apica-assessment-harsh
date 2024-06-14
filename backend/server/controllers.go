package server

import (
	"lru-cache/cache"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	lruCache cache.LRUCache = *cache.NewLRUCache(5, time.Second*10)
)

type CacheBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Ttl   int    `json:"ttl"`
}

func PutCache(c *fiber.Ctx) error {
	body := new(CacheBody)
	err := c.BodyParser(body)
	lruCache.Set(body.Key, body.Value, time.Duration(body.Ttl)*time.Second)
	res := make(map[string]interface{})
	if err != nil {
		res["ok"] = false
		res["message"] = err
		return c.JSON(res)
	}
	res["ok"] = true
	res["message"] = "updated cache succesfully"
	return c.JSON(res)
}

func GetCache(c *fiber.Ctx) error {
	key := c.Params("key")
	value, bool := lruCache.Get(key)
	res := make(map[string]interface{})
	res["key"] = key
	res["value"] = value
	res["ok"] = true
	if bool {
		return c.JSON(res)
	}
	res = make(map[string]interface{})
	res["message"] = "Key not found"
	res["ok"] = false
	return c.JSON(res)
}

func DeleteCache(c *fiber.Ctx) error {
	key := c.Params("key")
	deleted := lruCache.Delete(key)
	res := make(map[string]interface{})
	if deleted {
		res["message"] = "Deleted Succesfully"
		res["ok"] = true
		return c.JSON(res)
	}
	res["ok"] = false
	res["message"] = "Key not found"
	return c.JSON(res)
}

func GetSnapshot(c *fiber.Ctx) error {
	return c.JSON(lruCache.Snapshot())
}

func ClearCache(c *fiber.Ctx) error {
	lruCache.Clear()
	res := make(map[string]interface{})
	res["message"] = "Cache cleared succesfully"
	res["ok"] = true
	return c.JSON(res)
}
