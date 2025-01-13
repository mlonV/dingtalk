package aws

// import (
// 	awssns "github.com/aws/aws-sdk-go/service/sns"

// 	"github.com/gin-gonic/gin"
// 	"log"
// 	"net/http"
// )

// {
// 	"Type" : "Notification",
// 	"MessageId" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
// 	"TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
// 	"Subject" : "My First Message",
// 	"Message" : "Hello world!",
// 	"Timestamp" : "2012-05-02T00:54:06.655Z",
// 	"SignatureVersion" : "1",
// 	"Signature" : "EXAMPLEw6JRN...",
// 	"SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
// 	"UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96"
// 	}
// type SNSMessage struct {
// 	Type             string `json:"Type"`
// 	MessageId        string `json:"MessageId"`
// 	Token            string `json:"Token"`
// 	TopicArn         string `json:"TopicArn"`
// 	Subject          string `json:"Subject"`
// 	Message          string `json:"Message"`
// 	SubscribeURL     string `json:"SubscribeURL"`
// 	Timestamp        string `json:"Timestamp"`
// 	SignatureVersion string `json:"SignatureVersion"`
// 	Signature        string `json:"Signature"`
// 	SigningCertURL   string `json:"SigningCertURL"`
// 	UnsubscribeURL   string `json:"UnsubscribeURL"`
// }

// // Handler to process SNS messages
// func SnsHandler(c *gin.Context) {

// 	var snsMessage SNSMessage

// 	// Parse the JSON body into the SNSMessage struct
// 	if err := c.BindJSON(&snsMessage); err != nil {
// 		log.Printf("Error parsing SNS message: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SNS message format"})
// 		return
// 	}

// 	// Handle subscription confirmation
// 	if snsMessage.Type == "SubscriptionConfirmation" {
// 		log.Println("Subscription confirmation received.")
// 		log.Printf("Visit this URL to confirm the subscription: %s", snsMessage.SubscribeURL)

// 		// Optionally confirm the subscription automatically
// 		resp, err := http.Get(snsMessage.SubscribeURL)
// 		if err != nil {
// 			log.Printf("Error confirming subscription: %v", err)
// 		} else {
// 			defer resp.Body.Close()
// 			log.Println("Subscription confirmed.")
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "Subscription confirmed"})
// 		return
// 	}
// 	awssns
// 	// Handle notification message
// 	if snsMessage.Type == "Notification" {
// 		log.Println("Notification received.")
// 		log.Printf("Message: %s", snsMessage.Message)

// 		// Here you can add logic to process the SNS message as needed
// 		c.JSON(http.StatusOK, gin.H{"status": "Notification processed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "Message received but unhandled"})
// }
