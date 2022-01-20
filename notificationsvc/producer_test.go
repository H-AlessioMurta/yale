package main

import (
	"os"
    "testing"
	"encoding/json"
	"log"
	"bytes"
	"net/http"
	"github.com/gofiber/fiber/v2"
)

func TestMain(m *testing.M) {
    srv := fiber.New()
	api := srv.Group("/api/v1") // /api
	api.Post("/comments", createComment)
    code := m.Run()
    os.Exit(code)
}

func TestPostComment(t *testing.T){
	values := map[string]string{"Service": "Borrowingsvc", "fatal error": "Missing connection to mongo db"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("https://notificationsvc:3000", "application/json",
		bytes.NewBuffer(json_data))

	CheckErr(err)

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["success"] != "true"{
		t.Errorf("Find this error %v", res["message"])
	}	
}
