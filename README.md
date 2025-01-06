# :candy: Candy Cache :candy:

**CandyCache** — это простой и эффективный кэш на языке Go, который позволяет хранить данные с ограниченным временем жизни (**TTL**). 

# Установка

Для использования CandyCache в вашем проекте, установите его, используя ```go get git.hikan.ru/serr/candycache```, далее просто добавьте ```git.hikan.ru/serr/candycache``` в блок импорта.

# Основные возможности

- Автоматическая очистка устаревших элементов и возможность ее отключения.
- Кэшем можно управлять вручную.
- Конкурентный доступ к данным возможен благодаря мьютексам.
- Кэш может хранить данные любых типов.
- Можно создавать и загружать дампы кэша.

# Использование

## Создание кэша

### С автоматической очисткой

Для создания нового экземпляра кэша используйте функцию **Cacher**, передавая интервал очистки в наносекундах.

Если требуется указать интервал в секундах/минутах/часах и т.д. - используйте множители из пакета **time**:
```go
cache := candycache.Cacher(10 * time.Minute) // Очистка каждые 10 минут
```

### Без автоматической очистки

Если автоматичская очистка не нужна - просто передайте параметром любое отрицательное число:

```go
cache := candycache.Cacher(-1) // Кэш не будет очищаться автоматически
```

## Добавление элемента

Для добавления элемента в кэш используйте метод **Set**:
```go
cache.Set("key", "value", 5 * time.Minute) // Элемент будет считаться устаревшим через 5 минут
```
В случае, если по указанном ключу уже что-то хранится, оно будет заменено на новый элемент.

## Получение элемента

Для получения элемента из кэша используйте метод **Get**:

```go
value, err := cache.Get("key") // Получение значения по ключу "key"
```
Если элемент найден, то в переменную **value** будет записано значение, а в **err** - **nil**. Если элемент не найден, то в **err** будет записано **key not found**, а значением вернется **nil**.

## Удаление элемента

Для удаления элемента по ключу используйте метод **Delete**:

```go
err := cache.Delete("key1")
if err != nil {
    fmt.Println("Ошибка:", err) // Не найдено записи по ключу
}
```

Элемент будет удален, не смотря на то, устаревший он или нет.

## Массовое удаление элементов

### Удаление устаревших элементов

Для удаления устаревших элементов используйте метод **Cleanup**:

```go
cache.Cleanup() // Перебирает все элементы кэша, удаляет устаревшие
```

### Удаление всех элементов кэша

Для полной очистки кэша используйте метод **Flush**:

```go
cache.Flush() // Удаляет все элементы кэша, не смотря на то, устаревшие они или нет
```

## Получение информации о кэше

### Получение списка элементов

Для получения списка всех элементов в кэше используйте метод **List**:

```go
items := cache.List() // Список всех элементов кэша
for _, item := range items {
    fmt.Printf("Ключ: %s, Значение: %v, Момент устаревания: %d\n", item.Key, item.Item.Data(), item.Item.DestroyTimestamp())
}
```

### Получение количества элементов

Для получения количества элементов в кэше используйте метод **Count**:

```go
count := cache.Count() // Количество элементов в кэше
```

### Получение размера кэша

Для получения размера всего кэша в байтах используйте метод **Size**:

```go
size := cache.Size() // Размер кэша в байтах
```

Данный метод возвращает корректное значение, если в кэше элементы представлены этими типами данных:

```go
bool
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64, uintptr
float32, float64
complex64, complex128
array, slice, string
map, struct, func
```

**И композициями этих типов**.

В противном случае значение может быть не точным.

# Пример использования №1

```go
cache := candycache.Cacher(10 * time.Minute) // Создаем кэш с интервалом очистки 10 минут

cache.Set("key1", "string", 5*time.Minute)
cache.Set("key2", 2, 10*time.Minute)
cache.Set("key7", -2.1231, 10*time.Minute)
cache.Set("key3", []string{"string1", "string2"}, 10*time.Minute)
cache.Set("key4", map[string]int{"a": 1, "b": 2}, 10*time.Minute)
cache.Set("key5", Person{Name: "Alice", Age: 30, Hobbies: []string{"reading", "swimming"}}, 10*time.Minute)
cache.Set("key6", []Person{
    {Name: "Bob", Age: 25, Hobbies: []string{"coding", "gaming"}},
    {Name: "Charlie", Age: 35, Hobbies: []string{"hiking", "photography"}},
}, 10*time.Minute)

file, err := os.Create("cache_dump.json")
if err != nil {
    log.Fatal("error creating file: ", err)
}

if err := cache.Save(file); err != nil { // Сохранение кэша в файл
    log.Fatal("error saving cache: ", err)
}
file.Close()

cache.Flush() // Удаление всех элементов из кэша

file, err = os.Open("cache_dump.json")
if err != nil {
    log.Fatal("error opening file: ", err)
}

if err := cache.Load(file); err != nil { // Загрузка кэша из файла
    fmt.Println("error load cache:", err)
}

list := cache.List() // Получаю список элементов кэша

for _, i := range list {
    fmt.Println(i.Key, i.Item.Data(), i.Item.DestroyTimestamp())
}
```

# Пример использования №2

```go
cache := candycache.Cacher(10 * time.Minute) // Создаем кэш с интервалом очистки 10 минут

cache.Set("key1", "string", 5*time.Minute)
cache.Set("key2", 2, 10*time.Minute)
cache.Set("key7", -2.1231, 10*time.Minute)
cache.Set("key3", []string{"string1", "string2"}, 10*time.Minute)
cache.Set("key4", map[string]int{"a": 1, "b": 2}, 10*time.Minute)
cache.Set("key5", Person{Name: "Alice", Age: 30, Hobbies: []string{"reading", "swimming"}}, 10*time.Minute)
cache.Set("key6", []Person{
    {Name: "Bob", Age: 25, Hobbies: []string{"coding", "gaming"}},
    {Name: "Charlie", Age: 35, Hobbies: []string{"hiking", "photography"}},
}, 10*time.Minute)

var buffer bytes.Buffer

if err := cache.Save(&buffer); err != nil { // Сохранение бэкапа
    log.Fatal("error saving cache: ", err)
}

cache.Set("key1", "lost", 10*time.Minute)
cache.Set("key2", "lost", 10*time.Minute)
cache.Set("key3", "lost", 10*time.Minute)
cache.Set("key4", "lost", 10*time.Minute)
cache.Set("key5", "lost", 10*time.Minute)
cache.Set("key6", "lost", 10*time.Minute)
cache.Set("key7", "lost", 10*time.Minute)

if err := cache.Load(&buffer); err != nil { // Восстановление бэкапа
    log.Fatal("error loading cache: ", err)
}

list := cache.List() // Получаю список элементов кэша

for _, i := range list {
    fmt.Println(i.Key, i.Item.Data(), i.Item.DestroyTimestamp())
}
```