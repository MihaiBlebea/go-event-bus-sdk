package eventbus

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const base = "https://go-event-bus.cap-rover.purpletreetech.com/api"

type Service interface {
	Subscribe(eventName, handler string) error
	Unsubscribe(eventName string) error
	Trigger(eventName string, payload interface{}) error
	EventHistory(page, perPage int) ([]EventResponse, error)
}

type service struct {
	apiVersion int
	token      string
}

func New(token string) Service {
	return &service{
		apiVersion: 1,
		token:      token,
	}
}

// Project - creates a new project on the event bus server with supplied name.
// If request is successful, sets the api token on the client struct
func (s *service) Project(name string) error {
	resp := ProjectResponse{}
	data := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	if err := post(s.getServerURL("/project"), &data, &resp); err != nil {
		return err
	}

	if !resp.Success {
		return errors.New("failed to obtain token")
	}

	s.token = resp.Token

	return nil
}

// Subscribe - Subscribe to an event providing the event name and the handler url
// to be notified when the event is fired.
func (s *service) Subscribe(eventName, handler string) error {
	data := struct {
		Token string `json:"token"`
		Event string `json:"event"`
		URL   string `json:"url"`
	}{
		Token: s.token,
		Event: eventName,
		URL:   handler,
	}
	if err := post(s.getServerURL("/subscribe"), &data, nil); err != nil {
		return err
	}

	return nil
}

// Unsubscribe - removes a subscriber from the event listeners.
// Must provide a valid event name
func (s *service) Unsubscribe(eventName string) error {
	data := struct {
		Token string `json:"token"`
		Event string `json:"event"`
	}{
		Token: s.token,
		Event: eventName,
	}
	if err := post(s.getServerURL("/unsubscribe"), &data, nil); err != nil {
		return err
	}

	return nil
}

// Trigger - trigger an event. provide payload to be sent as event data.
func (s *service) Trigger(eventName string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data := struct {
		Token   string `json:"token"`
		Event   string `json:"event"`
		Payload string `json:"payload"`
	}{
		Token:   s.token,
		Event:   eventName,
		Payload: string(b),
	}
	if err := post(s.getServerURL("/handle"), &data, nil); err != nil {
		return err
	}

	return nil
}

// EventHistory - returns all the events for the particular project in order.
func (s *service) EventHistory(page, perPage int) ([]EventResponse, error) {
	resp := EventHistoryResponse{}

	url := s.getServerURL(
		fmt.Sprintf("/events&token=%s&per_page=%d&page=%d", s.token, perPage, page),
	)
	if err := get(url, &resp); err != nil {
		return []EventResponse{}, err
	}

	return resp.Events, nil
}

// GetToken - returns the api token set on the client
func (s *service) GetToken() string {
	return s.token
}

func (s *service) getServerURL(endpoint string) string {
	if string(endpoint[0]) == "/" {
		endpoint = strings.Replace(endpoint, "/", "", 1)
	}

	return fmt.Sprintf("%s/v%d/%s", base, s.apiVersion, endpoint)
}
