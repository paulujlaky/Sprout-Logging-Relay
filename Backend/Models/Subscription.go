package Models

import (
	"context"
	"elucid503/SproutLoggingRelay/Misc"
	"encoding/json"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Subscription struct {

	LastNotifSent int64 `json:"LastNotifSent"` // timestamp of the last notification sent

	Subscription webpush.Subscription `json:"Subscription"`
	
}

// Retrievers 

func GetAllSubscriptions() ([]*Subscription, error) {

	// Get all subscriptions from the database

	Cursor, Err := Misc.Database.Collection("Subscriptions").Find(context.TODO(), gin.H{})

	if Err != nil {

		return nil, Err
		
	}

	var Subscriptions []*Subscription

	for Cursor.Next(context.TODO()) {

		var Subscription Subscription

		Err := Cursor.Decode(&Subscription)

		if Err != nil {

			return nil, Err
			
		}

		Subscriptions = append(Subscriptions, &Subscription)

	}

	return Subscriptions, nil

}

// Attached Functions

func (Associated *Subscription) Save() error {

	// Saves the subscription to the database

	Filter := gin.H{"subscription": Associated.Subscription}
	ToSet := gin.H{"$set": Associated}

	Options := options.UpdateOne().SetUpsert(true) 

	_, Err := Misc.Database.Collection("Subscriptions").UpdateOne(context.TODO(), Filter, ToSet, Options)

	if Err != nil {

		return Err
		
	}

	return nil

}

func (Associated *Subscription) Delete() error {

	// Deletes the subscription from the database

	_, Err := Misc.Database.Collection("Subscriptions").DeleteOne(context.TODO(), gin.H{
		"subscription": Associated.Subscription,
	})

	if Err != nil {

		return Err
		
	}

	return nil

}

func (Associated *Subscription) SendNotification(Title string, Body string, Silent bool, Navigate string) error {

	payload := map[string]interface{}{

		"web_push": 8030,

		"notification": map[string]interface{}{

			"title":   Title,
			"body":    Body,
			"silent":  Silent,
			"navigate": Navigate,

		},

	}

	SerializedNotificationBody, ErrorSerializing := json.Marshal(payload)

	if ErrorSerializing != nil { return ErrorSerializing }

	Response, NotificationSendError := webpush.SendNotification(SerializedNotificationBody, &Associated.Subscription, &webpush.Options{

		Topic: "Logs",
		TTL: 60, // 1 minute

		VAPIDPublicKey: Misc.Config.WebPush.PublicKey,
		VAPIDPrivateKey: Misc.Config.WebPush.PrivateKey,
		Subscriber: "https://relay.sprout.software",

	});

	if Response != nil && (Response.StatusCode == 404 || Response.StatusCode == 410) {

		return Associated.Delete() // subscription is no longer valid
		
	}

	if NotificationSendError != nil { return NotificationSendError }

	// updating the time

	Associated.LastNotifSent = time.Now().Unix()
	return Associated.Save()

}