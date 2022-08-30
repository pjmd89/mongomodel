package tests

import (
	"github.com/pjmd89/mongomodel/mongomodel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestSubTypes struct {
}
type TestTypes struct {
	mongomodel.Model `bson:"-"`
	Id               primitive.ObjectID  `bson:"_id,omitempty" gql:"name=_id,id=true,objectID=true"`
	IdWithVal        primitive.ObjectID  `bson:"idWithVal" gql:"name=idWithVal,objectID=true"`
	IdPtrWithVal     *primitive.ObjectID `bson:"idPtrWithVal" gql:"name=idPtrWithVal,objectID=true"`
	IdWithIDVal      primitive.ObjectID  `bson:"idWithIDVal" gql:"name=idWithIDVal,objectID=true"`
	IdPtrWithIDVal   *primitive.ObjectID `bson:"idPtrWithIDVal" gql:"name=idPtrWithIDVal,objectID=true"`
	IdOutVal         primitive.ObjectID  `bson:"idOutVal" gql:"name=idOutVal,objectID=true"`
	IdPtrOutVal      *primitive.ObjectID `bson:"idPtrOutVal" gql:"name=idPtrOutVal,objectID=true"`
	IntWithVal       int                 `bson:"intWithVal" gql:"name=intWithVal"`
	IntPtrWithVal    *int                `bson:"intPtrWithVal" gql:"name=intPtrWithVal"`
	Int              int                 `bson:"int" gql:"name=int"`
	IntPtr           *int                `bson:"intPtr" gql:"name=intPtr"`
	IntDef           int                 `bson:"intDef" gql:"name=intDef,default=120"`
	IntPtrDef        *int                `bson:"intPtrDef" gql:"name=intPtrDef,default=15"`
	StringWithVal    string              `bson:"stringWithVal" gql:"name=stringWithVal"`
	StringPtrWithVal *string             `bson:"stringPtrWithVal" gql:"name=stringPtrWithVal"`
	String           string              `bson:"string" gql:"name=string"`
	StringPtr        *string             `bson:"stringPtr" gql:"name=stringPtr"`
	StringDef        string              `bson:"stringDef" gql:"name=stringDef,default=test default"`
	StringPtrDef     *string             `bson:"stringPtrDef" gql:"name=stringPtrDef,default=test ptr default"`
	Arr              []string            `bson:"arr" gql:"name=arr"`
	ArrPtr           *[]string           `bson:"arrPtr" gql:"name=arrPtr"`
	ArrWithVal       []string            `bson:"arrWithVal" gql:"name=arrWithVal"`
	ArrPtrWithVal    *[]string           `bson:"arrPtrWithVal" gql:"name=arrPtrWithVal"`
	Created          int64               `bson:"created" gql:"name=created,createdDate=true"`
	Updated          int64               `bson:"updated" gql:"name=updated,updatedDate=true"`
	CreatedPtr       *int64              `bson:"createdPtr" gql:"name=createdPtr,createdDate=true"`
	UpdatedPrt       *int64              `bson:"updatedPtr" gql:"name=updatedPtr,updatedDate=true"`
}
