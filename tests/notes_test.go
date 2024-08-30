package tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/countenum404/Veksel/internal/api"
	"github.com/countenum404/Veksel/internal/types"
)

const baseUrl = "http://localhost:4567"

func TestGetNotesWithoutAuth(t *testing.T) {
	url := baseUrl + "/api/notes"

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var apiErr types.ApiError
	body := resp.Body
	json.NewDecoder(body).Decode(&apiErr)
	if apiErr.Error != "Authentication failed: username and password are required" {
		t.Errorf("The user service passed request")
	}
}

func TestCreateNotesWithoutAuth(t *testing.T) {
	url := baseUrl + "/api/notes"

	noteReq := &types.CreateNoteRequest{
		Header:  "Free Software, Free Society",
		Content: "1. The GNU Project and Free Software\n2. What's in a Name?\n3. Copyright and Injustice Software Patents\n and etc",
	}

	reqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", api.JSON)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Status code is %d", resp.StatusCode)
	}

	var apiErr types.ApiError
	body := resp.Body
	json.NewDecoder(body).Decode(&apiErr)
	if apiErr.Error != "Authentication failed: username and password are required" {
		t.Errorf("The user service passed request")
	}
}

func TestGetNotesUserHasnt(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rgosling"
	user.Password = "drive"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var empty []types.Note
	body := resp.Body
	json.NewDecoder(body).Decode(&empty)
	if len(empty) > 0 {
		t.Error("The user has notes")
	}
}

func TestCreateNote(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rstallman"
	user.Password = "gnulinux"

	noteReq := &types.CreateNoteRequest{
		Header:  "Free Software, Free Society",
		Content: "1. The GNU Project and Free Software 2. What's in a Name? 3. Copyright and Injustice Software Patents",
	}

	reqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", api.JSON)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var spellResponse types.SpellResponse
	respBody := resp.Body
	json.NewDecoder(respBody).Decode(&spellResponse)
	if len(spellResponse.Spells) > 0 {
		t.Errorf("Spells integration is working for en language")
	}
	if spellResponse.NoteRequest.Header != noteReq.Header {
		t.Errorf("Wrong header")
	}
	if spellResponse.NoteRequest.Content != noteReq.Content {
		t.Errorf("Wrong content")
	}
}

func TestOver10kCreateNote(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rstallman"
	user.Password = "gnulinux"

	file, err := os.Open("large.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	large, err := io.ReadAll(file)

	noteReq := &types.CreateNoteRequest{
		Header:  "A LARGE",
		Content: string(large),
	}
	reqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", api.JSON)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var apiError types.ApiError
	respBody := resp.Body
	json.NewDecoder(respBody).Decode(&apiError)
	if apiError.Error != "note is too large" {
		t.Errorf("Note over 10k is accepted")
	}
}

func TestGetNotesWithWrongPassword(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rgosling"
	user.Password = "lalaland"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var apiErr types.ApiError
	body := resp.Body
	json.NewDecoder(body).Decode(&apiErr)
	if apiErr.Error != "Invalid username or password" {
		t.Errorf("The user service passed request")
	}
}

func TestCreateNoteWithWrongPassword(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rstallman"
	user.Password = "linux"

	noteReq := &types.CreateNoteRequest{
		Header:  "Free Software, Free Society",
		Content: "1. The GNU Project and Free Software 2. What's in a Name? 3. Copyright and Injustice Software Patents",
	}

	reqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", api.JSON)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var apiErr types.ApiError
	body := resp.Body
	json.NewDecoder(body).Decode(&apiErr)
	if apiErr.Error != "Invalid username or password" {
		t.Errorf("The user service passed request")
	}
}

func TestCreateNoteWithSpellingErrors(t *testing.T) {
	url := baseUrl + "/api/notes"

	user := new(types.User)
	user.Username = "rstallman"
	user.Password = "gnulinux"

	preparedSpellRes := types.SpellResponse{
		NoteRequest: types.CreateNoteRequest{Header: "Places to visit", Content: "1. Mascow"},
		Spells: types.SpellResult{
			types.SpellingError{
				Code: 1,
				Pos:  0,
				Row:  0,
				Col:  0,
				Len:  6,
				Word: "Mascow",
				S: []string{
					"Moscow",
					"Mascot",
					"Madcow",
				},
			},
		},
	}

	noteReq := &types.CreateNoteRequest{
		Header:  "Places to visit",
		Content: "1. Mascow",
	}

	reqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	auth := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", api.JSON)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Some errors occured:%s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %d", resp.StatusCode)
	}
	var spellResp types.SpellResponse
	body := resp.Body
	json.NewDecoder(body).Decode(&spellResp)

	if spellResp.NoteRequest != preparedSpellRes.NoteRequest {
		t.Error("Note Request struct not equals")
	}
	if len(spellResp.Spells) != len(preparedSpellRes.Spells) {
		t.Error("Not equals")
	}
}
