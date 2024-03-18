package main

import (
	"fmt"
	"sort"
	"strings"
)

var game *Game
var locations map[string]*Location
var items map[string]*Item
var msg *Message

func initLocations() {
	locations = map[string]*Location{
		"кухня": {
			id:            0,
			name:          "кухня",
			enter:         "кухня, ничего интересного",
			description:   "ты находишься на кухне, ",
			locationsToGo: map[string]*Location{},
		},
		"комната": {
			id:            1,
			name:          "комната",
			enter:         "ты в своей комнате",
			description:   "",
			locationsToGo: map[string]*Location{},
		},
		"коридор": {
			id:            2,
			name:          "коридор",
			enter:         "ничего интересного",
			description:   "",
			locationsToGo: map[string]*Location{},
		},
		"улица": {
			id:            3,
			name:          "улица",
			enter:         "на улице весна",
			description:   "",
			locationsToGo: map[string]*Location{},
			locked:        true,
			warning:       "дверь закрыта",
		},
	}
}

func initLocationsToGo() {
	for _, location := range locations {
		switch location.name {
		case "кухня":
			location.locationsToGo["коридор"] = locations["коридор"]
		case "комната":
			location.locationsToGo["коридор"] = locations["коридор"]
		case "улица":
			location.locationsToGo["домой"] = locations["домой"]
		case "коридор":
			location.locationsToGo["кухня"] = locations["кухня"]
			location.locationsToGo["комната"] = locations["комната"]
			location.locationsToGo["улица"] = locations["улица"]

		}

	}
}

func initItems() {
	items = make(map[string]*Item)

	items["ключи"] = &Item{
		id:       0,
		name:     "ключи",
		location: "комната",
		place:    "на столе",
	}
	items["конспекты"] = &Item{
		id:       1,
		name:     "конспекты",
		location: "комната",
		place:    "на столе",
	}
	items["чай"] = &Item{
		id:       3,
		name:     "чай",
		location: "кухня",
		place:    "на столе",
	}

	items["рюкзак"] = &Item{
		id:       4,
		name:     "рюкзак",
		location: "комната",
		place:    "на стуле",
	}

	items["дверь"] = &Item{
		id:       5,
		name:     "дверь",
		location: "",
		place:    "",
	}
}

func initGame() {
	game = &Game{}
	msg = &Message{}
	initLocations()
	initLocationsToGo()
	initItems()
	initQuests()
	game.currentLocation = locations["кухня"]

}

func initQuests() {
	game.quests = make(map[string]string)
	game.quests["кухня"] = ", надо собрать рюкзак и идти в универ"
	game.quests["рюкзак"] = "не выполнен"
	game.quests["ключ"] = "не выполнен"
	game.quests["дверь"] = "не выполнен"
}

func checkQuests() {
	if items["ключи"].location == "инвентарь" {
		game.quests["ключ"] = "выполнен"
	}
	if items["ключи"].location == "инвентарь" &&
		items["конспекты"].location == "инвентарь" {
		game.quests["кухня"] = ", надо идти в универ"
	}
	if game.quests["дверь"] == "выполнен" {
		locations["улица"].locked = false
	}
}

func itemsMsg(loc string) {
	var onTableItems []*Item
	var onChairItems []*Item
	message := ""

	// Фильтрация предметов на столе и на стуле
	for _, item := range items {
		if item.location == loc {
			if item.place == "на столе" {
				onTableItems = append(onTableItems, item)
			} else if item.place == "на стуле" {
				onChairItems = append(onChairItems, item)
			}
		}
	}

	// Сортировка предметов по порядковому номеру
	sort.Slice(onTableItems, func(i, j int) bool {
		return onTableItems[i].id < onTableItems[j].id
	})
	sort.Slice(onChairItems, func(i, j int) bool {
		return onChairItems[i].id < onChairItems[j].id
	})

	// Формирование строки сообщения
	if len(onTableItems) > 0 {
		tableItemsNames := make([]string, len(onTableItems))
		for i, item := range onTableItems {
			tableItemsNames[i] = item.name
		}
		message += "на столе: " + strings.Join(tableItemsNames, ", ")
	}

	if len(onChairItems) > 0 {
		chairItemsNames := make([]string, len(onChairItems))
		for i, item := range onChairItems {
			chairItemsNames[i] = item.name
		}
		message += ", на стуле: " + strings.Join(chairItemsNames, ", ")
	}

	// Установка значения сообщения
	if len(onTableItems)+len(onChairItems) == 0 {
		msg.items = "пустая комната"
		return
	}

	msg.items = message
}

func isItemInLocation(itemName, locationName string) bool {
	for _, item := range items {
		if item.name == itemName && item.location == locationName {
			return true
		}
	}
	return false
}

func LocationsToGoMsg(loc *Location) {
	var locationNames []string
	for key := range loc.locationsToGo {
		locationNames = append(locationNames, key)
	}

	// Сортировка имен локаций по их ID
	sort.Slice(locationNames, func(i, j int) bool {
		loc1 := loc.locationsToGo[locationNames[i]]
		loc2 := loc.locationsToGo[locationNames[j]]
		return loc1.id < loc2.id
	})

	msg.locationsToGo = "можно пройти - " + strings.Join(locationNames, ", ")
}

func (msg Message) Print() string {
	if msg.err != "" {
		return msg.err
	}

	if game.quests["универ"] != "" && game.currentLocation == locations["кухня"] {
		msg.items += "," + game.quests["универ"]
	}

	return fmt.Sprintf("%s. %s %s %s",
		msg.location,
		msg.items,
		msg.locationsToGo,
		msg.action)

}

func (msg *Message) Clear() {
	msg.location = ""
	msg.items = ""
	msg.locationsToGo = ""
	msg.action = ""
	msg.err = ""
}

type Message struct {
	location      string
	items         string
	locationsToGo string
	action        string
	err           string
}

type Game struct {
	quests          map[string]string
	currentLocation *Location
}

type Location struct {
	id            int
	name          string
	enter         string
	description   string
	locationsToGo map[string]*Location
	locked        bool
	warning       string
}

type Item struct {
	id       int
	name     string
	location string
	place    string
}

type Player struct {
}

func parceCommand(command string) string {
	var answer string
	splittedCmd := strings.Split(command, " ")

	switch splittedCmd[0] {

	case "идти":
		if _, ok := game.currentLocation.locationsToGo[splittedCmd[1]]; ok {

			if locations[splittedCmd[1]].locked {
				return locations[splittedCmd[1]].warning
			}
			game.currentLocation = locations[splittedCmd[1]]
			LocationsToGoMsg(game.currentLocation)

			answer = fmt.Sprintf("%s. %s",
				game.currentLocation.enter,
				msg.locationsToGo)
		} else {
			return "нет пути в " + splittedCmd[1]
		}

	case "осмотреться":
		msg.location = game.currentLocation.description
		itemsMsg(game.currentLocation.name)
		LocationsToGoMsg(game.currentLocation)

		answer = fmt.Sprintf("%s%s. %s",
			game.currentLocation.description,
			msg.items+game.quests[game.currentLocation.name],
			msg.locationsToGo)

	case "надеть":
		if splittedCmd[1] == "рюкзак" {
			game.quests[splittedCmd[1]] = "выполнен"
			items[splittedCmd[1]].place = "инвентарь"
			return "вы надели: рюкзак"
		} else {
			return "нет такого"
		}

	case "взять":
		if game.quests["рюкзак"] == "не выполнен" {
			return "некуда класть"
		}

		item_to_take := splittedCmd[1]

		if isItemInLocation(item_to_take, game.currentLocation.name) {
			items[item_to_take].location = "инвентарь"
			return "предмет добавлен в инвентарь: " + item_to_take

		} else {
			return "нет такого"
		}
	case "применить":
		item_to_use := splittedCmd[1]
		item_to_use_to := splittedCmd[2]

		if !isItemInLocation(item_to_use, "инвентарь") {
			return "нет предмета в инвентаре - " + item_to_use
		}

		if _, ok := items[item_to_use_to]; !ok {
			return "не к чему применить"
		}

		if item_to_use+item_to_use_to == "ключидверь" {
			game.quests["дверь"] = "выполнен"
			return "дверь открыта"
		}

	default:
		return "неизвестная команда"
	}

	return answer
}

func handleCommand(command string) string {

	msg.Clear()
	checkQuests()

	answer := parceCommand(command)
	return answer

}

func main() {

}
