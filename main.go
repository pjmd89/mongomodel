package main

import (
	"fmt"

	"github.com/pjmd89/mongomodel/mongomodel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trabajo struct{
	Nombre string `bson:"nombre"`
	Direccion string `bson:"direccion"`
}
type Usuarios struct {
	mongomodel.Model					`bson:"-"`
	Id				primitive.ObjectID	`bson:"_id,omitempty" gql:"name=id,id=true,objectID=true"`
	Nombre 			string				`bson:"nombre" gql:"name=nombre"`
	Apellido 		string				`bson:"apellido" gql:"name=apellido"`
	Edad 			int					`bson:"edad" gql:"name=edad"`
	Trabajos		[]Trabajo			`bson:"trabajos" gql:"name=trabajos"`
}
type UsuariosUpdate struct{
	mongomodel.Model					`bson:"-"`
	Id				primitive.ObjectID	`bson:"_id,omitempty" gql:"name=id,id=true,objectID=true"`
	Nombre 			string				`bson:"nombre,omitempty" gql:"name=nombre"`
	Apellido 		string				`bson:"apellido,omitempty" gql:"name=apellido"`
	Edad 			int					`bson:"edad,omitempty" gql:"name=edad"`
	Trabajos		[]Trabajo			`bson:"trabajos,omitempty" gql:"name=trabajos"`
}
func OnDB(currentDB string, currentCollection string) ( r string ){
	r = currentDB;
	if currentCollection == "esta coleccion cambia la BD"{
		r = "otra base de datos"
	}
	if currentCollection == "Usuarios"{
		r = "Pruebita"
	}
	return r;
}
func main(){
	db := mongomodel.NewConn(nil);
	db.(*mongomodel.MongoDBConn).OnDatabase = OnDB
	Users := Usuarios{};
	Users.Init(Usuarios{}, db);
	userData := map[string]interface{}{
		"nombre": "Pablo",
		"apellido": "Munoz",
		"edad": 30,
	}
	
	updateData := map[string]interface{}{
		"nombre":"Pablo Jose",
	}
	result,_ := Users.Create(userData, nil);
	fmt.Println(result.(*Usuarios).Id)

	count,_ := Users.Count(nil,nil)
	fmt.Println(count);

	readResults,_ := Users.Read(nil,nil);
	fmt.Println(readResults.([]Usuarios)[0].Nombre);

	where := map[string]interface{}{
		"_id":readResults.([]Usuarios)[0].Id,
	}
	update,_ := Users.Update(updateData,where,nil);
	fmt.Println(update.([]Usuarios)[0].Nombre, len(update.([]Usuarios)));

	whereDelete := map[string]interface{}{
		"_id":result.(*Usuarios).Id,
	}
	deleteModel,_ := Users.Delete(whereDelete,nil);
	fmt.Println(deleteModel.([]Usuarios)[0].Nombre);
}