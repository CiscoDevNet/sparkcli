package api

import (
	"errors"
	"log"
	"strings"

	"github.com/tdeckers/sparkcli/util"
)

type MessageService struct {
	Client *util.Client
}

type Message struct {
	Id            string `json:"id,omitempty"`
	RoomId        string `json:"roomId,omitempty"`
	Text          string `json:"text,omitempty"`
	Files         string `json:"files,omitempty"`
	ToPersonId    string `json:"toPersonId,omitempty"`
	ToPersonEmail string `json:"toPersonEmail,omitempty"`
	PersonId      string `json:"personId,omitempty"`
	PersonEmail   string `json:"personEmail,omitempty"`
	Created       string `json:"created,omitempty"`
}

type MessageItems struct {
	Items []Message `json:"items"`
}

func (m MessageService) list() (*[]Message, error) {
	log.Fatal("Not implemented")
	return nil, nil
}

func (m MessageService) List(roomId string) (*[]Message, error) {
	req, err := m.Client.NewGetRequest("/messages?roomId=" + roomId)
	if err != nil {
		return nil, err
	}
	var result MessageItems
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Items, nil
}

func (m MessageService) Create(roomId string, txt string) (*Message, error) {
	// Check for default roomId
	if roomId == "-" {
		config := util.GetConfiguration()
		if config.DefaultRoomId != "" {
			roomId = config.DefaultRoomId
		} else {
			return nil, errors.New("No DefaultRoomId configured.")
		}
	}

	recipientType := "room"
	msg := Message{Text: txt}

	// roomId can either be a UUID for a specific room, or an email address
	// prefixed with 'email:' to send to an individual. (Currently users
	// can't be addressed directly by their UUID.)
	splitRoom := strings.SplitN(roomId, ":", 2)
	if len(splitRoom) > 1 {
		recipientType = splitRoom[0]
		switch recipientType {
		case "email":
			msg.ToPersonEmail = splitRoom[1]
		default:
			msg.RoomId = splitRoom[1]
		}
	} else {
		msg.RoomId = roomId
	}

	req, err := m.Client.NewPostRequest("/messages", msg)
	if err != nil {
		return nil, err
	}
	var result Message
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MessageService) CreateFile(roomId string, file string) (*Message, error) {
	// Check for default roomId
	if roomId == "-" {
		config := util.GetConfiguration()
		if config.DefaultRoomId != "" {
			roomId = config.DefaultRoomId
		} else {
			return nil, errors.New("No DefaultRoomId configured.")
		}
	}

	req, err := m.Client.NewFilePostRequest("/messages", roomId, file)
	if err != nil {
		return nil, err
	}
	var result Message
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MessageService) Get(id string) (*Message, error) {
	if id == "" {
		return nil, errors.New("id can't be empty when getting message")
	}
	req, err := m.Client.NewGetRequest("/messages/" + id)
	if err != nil {
		return nil, err
	}
	var result Message
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MessageService) Delete(id string) error {
	if id == "" {
		return errors.New("id can't be empty when deleting a message")
	}
	req, err := m.Client.NewDeleteRequest("/messages/" + id)
	if err != nil {
		return err
	}
	_, err = m.Client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil //success
}
