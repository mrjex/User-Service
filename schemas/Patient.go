package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Username string
    Password string
}

