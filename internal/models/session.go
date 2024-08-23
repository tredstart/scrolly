package models

import (
	"log"

	"github.com/tredstart/scrolly/internal/database"
)

type Text struct {
	Id      string
	Session string
	Text    string
}

func FetchTexts(session string) ([]Text, error) {
	var texts []Text
	if err := database.DB.Select(&texts, `select * from text where session = ?`, session); err != nil {
		log.Println("cannot load texts by sesh ", err)
		return []Text{}, err
	}
	return texts, nil
}
func FetchTextById(id string) (Text, error) {
	var text Text
	if err := database.DB.Get(&text, `select * from text where id = ?`, id); err != nil {
		log.Println("cannot get by id ", err)
		return Text{}, err
	}

	return text, nil
}
func CreateText(text Text) error {
	if _, err := database.DB.Exec(`
    insert into text
    values (?, ?, ?)
    `, text.Id, text.Session, text.Text); err != nil {
		log.Println("cannot create text ", err)
		return err
	}
	return nil
}
func UpdateText(id, text string) error {
	if _, err := database.DB.Exec(`
    update text
    set text = ?
    where id = ?
    `, text, id); err != nil {
		log.Println("cannot update text ", err)
		return err
	}
	return nil
}

func CreateSession(id string) error {
	if _, err := database.DB.Exec("insert into session values (?)", id); err != nil {
		log.Println("cannot create session ", err)
		return err
	}

	return nil
}
