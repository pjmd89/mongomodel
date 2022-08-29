package tests

import (
	"github.com/pjmd89/mongomodel/mongomodel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestSubTypes struct {
}
type TestTypes struct {
	mongomodel.Model `bson:"-"`
	Id               primitive.ObjectID `bson:"_id,omitempty" gql:"name=_id,id=true,objectID=true"`
	IntWithVal       int                `bson:"intWithVal" gql:"name=intWithVal"`
	IntPtrWithVal    *int               `bson:"intPtrWithVal" gql:"name=intPtrWithVal"`
	Int              int                `bson:"int" gql:"name=int"`
	IntPtr           *int               `bson:"intPtr" gql:"name=intPtr"`
	IntDef           int                `bson:"intDef" gql:"name=intDef,default=120"`
	IntPtrDef        *int               `bson:"intPtrDef" gql:"name=intPtrDef,default=15"`
	StringWithVal    string             `bson:"StringWithVal" gql:"name=StringWithVal"`
	StringPtrWithVal *string            `bson:"StringPtrWithVal" gql:"name=StringPtrWithVal"`
	String           string             `bson:"String" gql:"name=String"`
	StringPtr        *string            `bson:"StringPtr" gql:"name=StringPtr"`
	StringDef        string             `bson:"StringDef" gql:"name=StringDef,default=test default"`
	StringPtrDef     *string            `bson:"StringPtrDef" gql:"name=StringPtrDef,default=test ptr default"`
	Created          int64              `bson:"created" gql:"name=created,createdDate=true"`
	Updated          int64              `bson:"updated" gql:"name=updated,updatedDate=true"`
	CreatedPtr       *int64             `bson:"createdPtr" gql:"name=createdPtr,createdDate=true"`
	UpdatedPrt       *int64             `bson:"updatedPtr" gql:"name=updatedPtr,updatedDate=true"`
}
