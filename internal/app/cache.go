package app

import (
	"log"
	"sync"
)

type Cache struct {
	mx   sync.RWMutex
	Data map[string]string
}

//Создание кеша
func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]string),
	}
}

//Запись данных в кеш
func (c *Cache) PutOrder(id string, o string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.Data[id] = o
}

//Получение данных из кеша
func (c *Cache) GetOrder(id string) (o string, b bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	log.Println("мы тут")
	o, b = c.Data[id]
	log.Println("мы тут2")
	return
}

//Удаление данных из кеша
func (c *Cache) DeleteOrder(id string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.Data, id)
}
