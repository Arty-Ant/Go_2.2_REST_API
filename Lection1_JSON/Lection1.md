## Лекция 1. Работа с JSON файлами

**JSON** - формат файлов (расширение файлов), которое повсеместно используется для реализации передачи данных
между серверами на уровне API.

**JSON** == **JavaScript Object Notation** (Object - аналог map в Go только для мира JS)

**JSON** - это простейшее файловое расширение, поддерживающее элементарную структуризацию (выглядит как набор пар ключ: значение)

### Шаг 1. Создадим простой .json файл
Для этого определим файл ```users.json```.
```
{"users" : [{"name": "Vasya"}, {"name" : "Vitya"}]}
```
Обратите внимание на то, что ***ПО СТАНДАРТУ В JSON ИСПОЛЬЗУЮТСЯ ДВОИНЫЕ КАВЫЧКИ***.

### Шаг 2. Создадим чуть более сложный .json
Создадим сразу читаемым
```
{
    "users" : [
        {
            "name" : "Alex",
            "type" : "Admin",
            "age" : 32,
            "social" : {
                "vkontakte" : "https://vk.com/id=123512",
                "telegram": "https://t.me/id=172835"
            }
        },
        {
            "name" : "Bob",
            "type" : "Regular",
            "age" : 12,
            "social" : {
                "vkontakte" : "https://vk.com/id=123561235",
                "telegram": "https://t.me/id=19283712"
            }
        },
        {
            "name" : "Alice",
            "type" : "Regular",
            "age" : 19,
            "social" : {
                "vkontakte" : "https://vk.com/id=12123123",
                "telegram": "https://t.me/id=172123123"
            }
        },
        {
            "name" : "George",
            "type" : "Regular",
            "age" : 42,
            "social" : {
                "vkontakte" : "https://vk.com/id=999999",
                "telegram": "https://t.me/id=98888888"
            }
        }
    ]
}
```

***ВАЖНО*** - .json не гарантирует соблюдения упорядоченности при выдаче ключей!


### Шаг 3. Как в принципе читают из таких файлов?
* Для начала нужно создать файловый экземпляр (файловый дескриптор)
```
jsonFile, err := os.Open("users.json")
```
* Сразу же обрабатываем ошибки!
```
if err != nil {....}
```

* Не забываем файл закрывать!
```
defer jsonFile.Close()
```

* Затем нам нужно из файл-дескриптора забрать данные и куда-то их поместить!
```
json.Unmarshall(byteArr, &куда_помещаем)
```

### Шаг 4. Теперь более конкретно.

В Go существует 2 способа работы с JSON файлами:
* структуризованная сериализация/десериализация
* неструктуризованная -//-


#### Шаг 4.1 Структуризация
***Сериализация*** - процесс конвертации объекта в последовательность байтов.  
***Десериализация*** - процесс конвертации последовательности байтов в объект.

Идея структуризованного подхода состоит в том, что мы заранее подготавливаем набор структур, с ***ЯВНО ПРОПИСАННЫМИ ПРАВИЛАМИ*** сериализации/десериализации полей объектов.

```
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Struct for representation total slice
// First Level of JSON object parsing
type Users struct {
	Users []User `json:"users"`
}

// Internal user representation
// Second level of JSON object parsing
type User struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"age"`
	Social Social `json:"social"`
}

// Social block representation
// Third level of JSON object parsing
type Social struct {
	Vkontakte string `json:"vkontakte"`
	Telegram  string `json:"telegram"`
}

// Функция для распечатывания User
func PrintUser(u *User) {
	fmt.Printf("Name: %s\n", u.Name)
	fmt.Printf("Type: %s\n", u.Type)
	fmt.Printf("Age: %d\n", u.Age)
	fmt.Printf("Social. VK: %s and Telegram: %s\n", u.Social.Vkontakte, u.Social.Telegram)
}

// 1. Рассмотрим процесс десериализации (то есть когда из последовательности в объект)
func main() {
	//1. Создадим файл дескриптор
	jsonFile, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	fmt.Println("File descriptor successfully created!")

	// 2. Теперь десериализуем содержимое jsonFile в экземпляр Go
	// Инициализируем экземпляр Users
	var users Users

	// Вычитываем содержимое jsonFile в ВИДЕ ПОСЛЕДОВАТЕЛЬНОСТИ БАЙТ!
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	// Теперь задача - перенести все из byteValue в users - это и есть десериализация!
	json.Unmarshal(byteValue, &users)
	for _, u := range users.Users {
		fmt.Println("================================")
		PrintUser(&u)
	}
}

```

Идея структурированной сериализации/десериализации состоит в том, чтобы общаться с JSON объектами напрямую, через стыковку полей.

Для того, чтобы настроить стыковку полей нужно:
* Определить необходимые уровни объектности JSON (в нашем случае их 3)
* Для каждого уровня объектности подготовить свою структуру, учитывающую набор полей объекта.
```
type Person struct {
    Name string `json:"name"`
}
```
* И все! Больше ничего делать не нужно. Остается только считать из файла и поместить в экземпляр!


#### Шаг 4.2 Неструктуризованный подход
В этом подходе читаемость кода стремится к нулю, но его можно использовать на этапе отладки, чтобы быстро
посмотреть, что вообще в принципе находится в json.

```
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//1. Создадим файл дескриптор
	jsonFile, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	fmt.Println("File descriptor successfully created!")

	//2. Вычитываю набор байт из файл-дескриптора
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	fmt.Println(result["users"])

}
```


### Шаг 5. Сериализация
* Только один способ - структуризованный.
```
//1. Превратим профессора в последовательность байтов
	byteArr, err := json.MarshalIndent(prof1, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(byteArr))
	err = os.WriteFile("output.json", byteArr, 0666) //-rw-rw-rw-
	if err != nil {
		log.Fatal(err)
	}
```

Сериализация - процесс перегона в байты. Поэтому на это этапе у нас на руках будет ```[]byte``` , который в последствии будет помещен в файл ```output.json```

***ВАЖНО*** : 0664/0666 - права доступа к файлу (в нашем случае это ```rw```)