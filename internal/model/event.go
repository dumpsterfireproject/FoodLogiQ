package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// {
// 	"id": "ljadfj", // derived. Internal ID of  event
// 	"createdAt": "2020-12-29T14:41:31.123Z", // derived. datetime the  event is created
// 	"createdBy": "<userid>", // derived. id of the user  who created the event
// 	"isDeleted": false, // derived. False when  created. True when deleted.
// 	"type": "shipping", // required. valid entries  are shipping and receiving
// 	"contents": [
// 	{
// 	"gtin": "1234", // required. Global Trade  Item Number. 14-digit number.
// 	"lot": "adffda", // required. any value. GTIN  + Lot are a compound identifier
// 	"bestByDate": "2021-01-13", // optional. date value
//  "expirationDate": "2021-01-17", // optional. date value
//  },
// 	...
// 	]
//    }

const ShippingType = "shipping"
const ReceivingType = "receiving"

type Contents struct {
	Gtin           string     `json:"gtin" bson:"gtin"`
	Lot            string     `json:"lot" bson:"lot"`
	BestByDate     *time.Time `json:"bestByDate,omitempty" bson:"bestByDate,omitempty"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
}

type Event struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt *time.Time         `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	CreatedBy *string            `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
	Type      string             `json:"type" bson:"type"`
	Contents  []Contents         `json:"contents" bson:"contents"`
}
