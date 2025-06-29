package Models

import (
	"context"
	"elucid503/SproutLoggingRelay/Misc"
	"time"

	"github.com/gin-gonic/gin"
)

type Subscription struct {

	LastNotifSent int64 `json:"LastNotifSent"` // timestamp of the last notification sent

	Details gin.H `json:"Details"`
	
	Key string `json:"Key"`
	PublicKey string `json:"PublicKey"`

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

	// Save the subscription to the database

	_, Err := Misc.Database.Collection("Subscriptions").InsertOne(context.TODO(), Associated)

	if Err != nil {

		return Err
		
	}

	return nil

}

func (Associated *Subscription) SendNotification() error {

	// now, we must request the Details.endpoint with 

	// we should also the last notification time

	Associated.LastNotifSent = time.Now().Unix()
	return Associated.Save()

}