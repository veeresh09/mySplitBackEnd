package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Expense represents an expense in a group
type Expense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	GroupID     primitive.ObjectID `bson:"groupId"`     // ID of the group this expense belongs to
	PaidBy      primitive.ObjectID `bson:"paidBy"`      // ID of the user who paid the expense
	Amount      float64            `bson:"amount"`      // Total amount of the expense
	Description string             `bson:"description"` // Description of the expense
	Split       []ExpenseSplit     `bson:"split"`       // Information on how the expense is split among users
	CreatedAt   time.Time          `bson:"createdAt"`   // Timestamp of when the expense was created
	ModifiedAt  time.Time          `bson:"modifiedAt"`  // Timestamp of last modification
	CreatedBy   primitive.ObjectID `bson:"userId"`      // ID of the user who created the expense
}

// ExpenseSplit represents how an individual expense is split among the users
type ExpenseSplit struct {
	UserID primitive.ObjectID `bson:"userId"` // ID of the user
	Amount float64            `bson:"amount"` // Amount attributed to this user
}
