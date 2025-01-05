package candycache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

// Структура виде ключ-значение для возвращения списка элементов кэша с их ключами.
type KeyItemPair struct {
	Key  string
	Item Item
}

// Элемент в кэше - это данные и время их жизни.
type Item struct {
	destroyTimestamp int64       // Момент в Unix-секундах, когда элемент становится устаревшим
	data             interface{} // Данные
}

// Кэш - это хранилище элементов и инервал его очистки (ну и мьютекс на всякий случай).
// Интервал очистки хранилища укахывается в НАНОСЕКУНДАХ (используй множители для преобразования во что-то другое).
type Cache struct {
	sync.RWMutex                    // Мьютекс ждя реализации безопасного доступа к общим данным
	storage         map[string]Item // Хранилище элементов
	cleanupInterval time.Duration   // Интервал очистки хранилища в наносекундах
}

// Создает новый экземпляр Cache с интервалом очистки cleanupInterval.
// Если cleanupInterval < 0, то кэш не будет очищаться автоматически.
func Cacher(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		storage:         make(map[string]Item),
		cleanupInterval: cleanupInterval,
	}

	// Запускаем Garbage Collector если интервал очистки больше 0
	// Иначе (если он отрицательный) кэш будет жить до ручного вызова Cleanup
	if cleanupInterval > 0 {
		go cache.gc(cleanupInterval)
	}

	return cache
}

// gc = Garbage Collector.
func (c *Cache) gc(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.Cleanup()
	}
}

// Перебирает все элементы в кэше, удаляет устаревшие.
func (c *Cache) Cleanup() {
	c.Lock()
	defer c.Unlock()

	for key, item := range c.storage {
		if item.destroyTimestamp <= time.Now().Unix() {
			delete(c.storage, key)
		}
	}
}

// Удаление всех элементов из кэша.
func (c *Cache) Flush() {
	c.Lock()
	defer c.Unlock()

	for key := range c.storage {
		delete(c.storage, key)
	}
}

// Получение элемента из кэша по ключу.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.storage[key]

	// Элемент не найден в кэше
	if !found {
		return nil, false
	}

	return item.data, true
}

// Удаление элемента по ключу.
func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.storage[key]; !found {
		return errors.New("key not found")
	}

	delete(c.storage, key)

	return nil
}

// Добавление элемента в кэш.
// key - ключ.
// data - данные.
// ttl - время жизни элемента (time to life) в наносекундах.
func (c *Cache) Add(key string, data interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	c.storage[key] = Item{
		destroyTimestamp: time.Now().Unix() + int64(ttl.Seconds()),
		data:             data,
	}
}

// Вернет количество элементов в кэше.
func (c *Cache) Count() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.storage)
}

// Печать всех элементов кэша (ключ и время уничтожения).
func (c *Cache) List() []KeyItemPair {
	c.RLock()
	defer c.RUnlock()

	// Создаем срез для хранения пар ключ-значение
	items := make([]KeyItemPair, 0, len(c.storage))

	// Заполняем срез парами ключ-значение
	for key, item := range c.storage {
		items = append(items, KeyItemPair{Key: key, Item: item})
	}

	return items
}

// Вернет размер всего кэша в байтах.
func (c *Cache) Size() int {
	c.RLock()
	defer c.RUnlock()

	size := 0
	for key, item := range c.storage {
		size += isize(key) + isize(item.data) + isize(item.destroyTimestamp)
	}

	return size
}

// ПОДДЕРЖИВАЕМЫЕ ТИПЫ:
// Bool +
// Int +
// Int8 +
// Int16 +
// Int32 +
// Int64 +
// Uint +
// Uint8 +
// Uint16 +
// Uint32 +
// Uint64 +
// Uintptr +
// Float32 +
// Float64 +
// Complex64 +
// Complex128 +
// Array +
// Func +
// Map +
// Slice +
// String +
// Struct
func isize(i interface{}) int {
	if i == nil {
		return 0
	}
	val := reflect.ValueOf(i)
	kind := val.Kind()
	size := 0
	switch kind {
	case reflect.Slice, reflect.Array, reflect.String:
		len := val.Len()
		for i := 0; i < len; i++ {
			size += isize(val.Index(i).Interface())
		}
		return size
	case reflect.Map:
		for _, key := range val.MapKeys() {
			size += isize(key.Interface()) + isize(val.MapIndex(key).Interface())
		}
		return size
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			size += isize(val.Field(i).Interface())
		}
		return size
	default:
		return int(reflect.TypeOf(i).Size())
	}
}

// Возвращает данные элемента кэша.
func (i *Item) Data() interface{} {
	return i.data
}

// Возвращает момент смерти элемента кэша.
func (i *Item) DestroyTimestamp() int64 {
	return i.destroyTimestamp
}
