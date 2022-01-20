package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"

)

// Comment struct
type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	srv := fiber.New()
	api := srv.Group("/api/v1") // /api
	api.Post("/comments", createComment)
	log.SetPrefix("\033[36mNotificationsvc\033[0m: ")
	srv.Listen(":3000")

}
// Connection with Kafka inside k8s's cluster
func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	CheckErr(err)
	return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
	brokersUrl := []string{"kafka:9092"}// of mine kafka services
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	LogInfo(fmt.Sprintf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset))
	return nil
}

// createComment handler
func createComment(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(Comment)

	//  Parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	// convert body into bytes and send it to kafka
	cmtInBytes, err := json.Marshal(cmt)
	PushCommentToQueue("comments", cmtInBytes)

	// Return Comment in JSON format
	err = c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		LogError(err.Error())
		return err
	}
	LogResponse(err.Error())
	return err
}