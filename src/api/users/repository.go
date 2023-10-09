package users

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/grpcgateway/src/db"
	"github.com/myrachanto/grpcgateway/src/pasetos"
	"github.com/myrachanto/grpcgateway/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Userrepository repository
var (
	Userrepository UserrepoInterface = &userrepository{}
	ctx                              = context.TODO()
	Userrepo                         = userrepository{}
	locker                           = sync.Mutex{}
	collectionName                   = "user"
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

type UserrepoInterface interface {
	Create(user *User) (*User, httperrors.HttpErr)
	Login(user *LoginUser) (*Auth, httperrors.HttpErr)
	Logout(token string) (string, httperrors.HttpErr)
	GetOne(code string) (user *User, errors httperrors.HttpErr)
	GetAll(search string) ([]*User, httperrors.HttpErr)
	Update(code string, user *User) (*User, httperrors.HttpErr)
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
}
type userrepository struct {
}

func NewUserRepo() UserrepoInterface {
	return &userrepository{}
}

func (r *userrepository) Create(user *User) (*User, httperrors.HttpErr) {
	if err1 := user.Validate(); err1 != nil {
		return nil, err1
	}
	var err httperrors.HttpErr
	ok := r.emailexist(user.Email)
	if ok {
		return nil, httperrors.NewBadRequestError("that email exist in the our system!")
	}
	if err != nil {
		return nil, err
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return nil, err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return nil, httperrors.NewNotFoundError("Your email format is wrong!")
	}
	code, errs := r.genecode()
	if errs != nil {
		return nil, errs
	}
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword
	user.Usercode = code
	user.Base.Updated_At = time.Now()
	user.Base.Created_At = time.Now()
	collection := db.Mongodb.Collection(collectionName)
	result1, errd := collection.InsertOne(ctx, &user)
	if errd != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Create user Failed, %d", errd))
	}
	user.ID = result1.InsertedID.(primitive.ObjectID)
	return user, nil

}

func (r *userrepository) Login(user *LoginUser) (*Auth, httperrors.HttpErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	var auser User
	filter := bson.M{"email": user.Email}
	collection := db.Mongodb.Collection(collectionName)
	errs := collection.FindOne(ctx, filter).Decode(&auser)
	if errs != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("User with this email does exist @ - , %d", errs))
	}
	ok := user.Compare(user.Password, auser.Password)
	if !ok {
		return nil, httperrors.NewNotFoundError("wrong email password combo!")
	}
	maker, err := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, err
	}
	data := &pasetos.Data{
		Usercode: auser.Usercode,
		Username: auser.Usercode,
		Email:    auser.Email,
	}
	tokenString, payload, err := maker.CreateToken(data, time.Hour*3)
	if err != nil {
		return nil, err
	}
	auths := &Auth{Usercode: auser.Usercode, Role: auser.Role, Picture: auser.Picture, UserName: auser.Username, Token: tokenString, TokenExpires: payload.ExpiredAt}

	// locker.Lock()
	// r.Sessions = append(r.Sessions, *auths)
	// locker.Unlock()
	return auths, nil
}
func (r *userrepository) Logout(token string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(token)
	if stringresults.Noerror() {
		return "", stringresults
	}
	collection := db.Mongodb.Collection("auth")
	filter1 := bson.M{"token": token}
	_, err3 := collection.DeleteOne(ctx, filter1)
	if err3 != nil {
		return "", httperrors.NewBadRequestError("something went wrong login out!")
	}
	return "something went wrong login out!", nil
}
func (r *userrepository) GetOne(code string) (user *User, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection(collectionName)
	filter := bson.M{"usercode": code}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return user, nil
}

func (r *userrepository) GetAll(search string) ([]*User, httperrors.HttpErr) {
	collection := db.Mongodb.Collection(collectionName)
	results := []*User{}
	fmt.Println(search)
	if search != "" {
		// 	filter := bson.D{
		// 		{"name", primitive.Regex{Pattern: search, Options: "i"}},
		// }
		filter := bson.D{
			{"$or", bson.A{
				bson.D{{"firstname", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"lastname", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"username", primitive.Regex{Pattern: search, Options: "i"}}},
				bson.D{{"email", primitive.Regex{Pattern: search, Options: "i"}}},
			}},
		}
		// fmt.Println(filter)
		cursor, err := collection.Find(ctx, filter)
		fmt.Println(cursor)
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		fmt.Println(results)
		return results, nil
	} else {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, httperrors.NewNotFoundError("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, httperrors.NewNotFoundError("Error decoding!")
		}
		return results, nil
	}

}

func (r *userrepository) Update(code string, user *User) (*User, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	uuser := &User{}
	collection := db.Mongodb.Collection(collectionName)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&uuser)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this id, %d", err))
	}

	if user.Firstname == "" {
		user.Firstname = uuser.Firstname
	}
	if user.Lastname == "" {
		user.Lastname = uuser.Lastname
	}
	if user.Username == "" {
		user.Username = uuser.Username
	}
	if user.Phone == "" {
		user.Phone = uuser.Phone
	}
	if user.Address == "" {
		user.Address = uuser.Address
	}
	if user.Picture == "" {
		user.Picture = uuser.Picture
	}
	if user.Email == "" {
		user.Email = uuser.Email
	}
	if user.Usercode == "" {
		user.Usercode = uuser.Usercode
	}
	update := bson.M{"$set": user}
	_, errs := collection.UpdateOne(ctx, filter, update)
	if errs != nil {
		return nil, httperrors.NewNotFoundError("Error updating!")
	}
	return uuser, nil
}

func (r *userrepository) PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(oldpassword)
	if stringresults.Noerror() {
		return "", "", stringresults
	}
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", "", stringresults2
	}
	stringresults3 := httperrors.ValidStringNotEmpty(newpassword)
	if stringresults3.Noerror() {
		return "", "", stringresults3
	}
	upay := &User{}
	collection := db.Mongodb.Collection(collectionName)
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return "", "", httperrors.NewBadRequestError(fmt.Sprintf("Could not find resource with this email, %d", err))
	}
	ok := upay.Compare(oldpassword, upay.Password)
	if !ok {
		return "", "", httperrors.NewNotFoundError("wrong password combo!")
	}
	newhashpassword, err2 := upay.HashPassword(newpassword)
	if err2 != nil {
		return "", "", err2
	}
	upay.Password = newhashpassword
	_, errs := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"password", newhashpassword}}},
		},
	)
	if errs != nil {
		return "", "", httperrors.NewNotFoundError("Error updating!")
	}
	return email, newpassword, nil
}
func (r userrepository) Delete(code string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	bl, errs := r.GetOne(code)
	if errs != nil {
		return "", errs
	}
	go support.Cleaner(bl.Picture)
	collection := db.Mongodb.Collection(collectionName)

	filter := bson.M{"usercode": code}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", httperrors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r userrepository) genecode() (string, httperrors.HttpErr) {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := db.Mongodb.Collection(collectionName)
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "UserCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r userrepository) getuno(code string) (result *User, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	collection := db.Mongodb.Collection(collectionName)
	filter := bson.M{"usercode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r userrepository) emailexist(email string) bool {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return stringresults.Noerror()
	}
	collection := db.Mongodb.Collection(collectionName)
	result := &User{}
	filter := bson.M{"email": email}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	return err1 == nil
}
