package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Username string  `bson:"username,omitempty" json:"username,omitempty"`
    Password string `bson:"password,omitempty" json:"password,omitempty"`
}

