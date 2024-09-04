package users

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/myrachanto/grpcgateway/src/pasetos"
	"github.com/myrachanto/grpcgateway/src/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Userrepository repository
const (
	collectionName = "user"
)

var (
	Userrepository UserrepoInterface = &userrepository{}
	ctx                              = context.TODO()
	Userrepo                         = userrepository{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

type UserrepoInterface interface {
	Create(user *User) (*UserDto, error)
	Login(user *LoginUser) (*Auth, error)
	Logout(token string) (string, error)
	GetOne(code string) (user *UserDto, errors error)
	GetAll(search string) ([]*UserDto, error)
	Update(code string, user *User) (*UserDto, error)
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, error)
	Delete(code string) (string, error)
}
type userrepository struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) UserrepoInterface {
	return &userrepository{
		db: db,
	}
}

func (r *userrepository) Create(user *User) (*UserDto, error) {
	if err1 := user.Validate(); err1 != nil {
		return nil, err1
	}
	var err error
	ok := user.ValidateEmail(user.Email)
	if !ok {
		return nil, fmt.Errorf("Your email format is wrong!")
	}
	ok = r.emailexist(user.Email)
	if ok {

		return nil, fmt.Errorf("that email exist in the our system!")
	}
	if err != nil {
		return nil, err
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return nil, err1
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
	collection := r.db.Collection(collectionName)
	result1, errd := collection.InsertOne(ctx, &user)
	if errd != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Create user Failed, %d", errd))
	}
	user.ID = result1.InsertedID.(primitive.ObjectID)
	return user.UserConvter(), nil

}

func (r *userrepository) Login(user *LoginUser) (*Auth, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	var auser User
	filter := bson.M{"email": user.Email}
	collection := r.db.Collection(collectionName)
	errs := collection.FindOne(ctx, filter).Decode(&auser)
	if errs != nil {
		return nil, fmt.Errorf(fmt.Sprintf("User with this email does exist @ - , %d", errs))
	}
	ok := user.Compare(user.Password, auser.Password)
	if !ok {
		return nil, fmt.Errorf("wrong email password combo!")
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
func (r *userrepository) Logout(token string) (string, error) {
	if len(token) == 0 {
		return "", fmt.Errorf("the token is empty")
	}
	collection := r.db.Collection("auth")
	filter1 := bson.M{"token": token}
	_, err3 := collection.DeleteOne(ctx, filter1)
	if err3 != nil {
		return "", fmt.Errorf("something went wrong login out!")
	}
	return "something went wrong login out!", nil
}
func (r *userrepository) GetOne(code string) (user *UserDto, errors error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("the code is empty")
	}
	collection := r.db.Collection(collectionName)
	filter := bson.M{"usercode": code}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Could not find resource with this id, %d", err))
	}
	return user, nil
}

func (r *userrepository) GetAll(search string) ([]*UserDto, error) {
	collection := r.db.Collection(collectionName)
	results := []*User{}
	if search != "" {
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
			return nil, fmt.Errorf("No records found!")
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil, fmt.Errorf("Error decoding!")
		}
		res := []*UserDto{}
		for _, u := range results {
			res = append(res, u.UserConvter())
		}

		return res, nil
	}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("No records found!")
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("Error decoding!")
	}
	res := []*UserDto{}
	for _, u := range results {
		res = append(res, u.UserConvter())
	}
	return res, nil

}

func (r *userrepository) Update(code string, user *User) (*UserDto, error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("the code is empty")
	}
	uuser := &User{}
	collection := r.db.Collection(collectionName)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&uuser)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Could not find resource with this id, %d", err))
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
		return nil, fmt.Errorf("Error updating!")
	}
	return uuser.UserConvter(), nil
}

func (r *userrepository) PasswordUpdate(oldpassword, email, newpassword string) (string, string, error) {
	if len(oldpassword) == 0 {
		return "", "", fmt.Errorf("the oldpassword is empty")
	}
	if len(email) == 0 {
		return "", "", fmt.Errorf("the email is empty")
	}
	if len(newpassword) == 0 {
		return "", "", fmt.Errorf("the newpassword is empty")
	}
	upay := &User{}
	collection := r.db.Collection(collectionName)
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&upay)
	if err != nil {
		return "", "", fmt.Errorf(fmt.Sprintf("Could not find resource with this email, %d", err))
	}
	ok := upay.Compare(oldpassword, upay.Password)
	if !ok {
		return "", "", fmt.Errorf("wrong password combo!")
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
		return "", "", fmt.Errorf("Error updating!")
	}
	return email, newpassword, nil
}
func (r userrepository) Delete(code string) (string, error) {
	if len(code) == 0 {
		return "", fmt.Errorf("the code is empty")
	}
	bl, errs := r.GetOne(code)
	if errs != nil {
		return "", errs
	}
	go support.Cleaner(bl.Picture)
	collection := r.db.Collection(collectionName)

	filter := bson.M{"usercode": code}
	ok, err := collection.DeleteOne(ctx, filter)
	if ok == nil {
		return "", fmt.Errorf(fmt.Sprintf("deletion of %d failed", err))
	}
	return "deleted successfully", nil

}
func (r userrepository) genecode() (string, error) {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	collection := r.db.Collection(collectionName)
	filter := bson.M{}
	count, err := collection.CountDocuments(ctx, filter)
	co := count + 1
	if err != nil {
		return "", fmt.Errorf("no results found")
	}
	cod := "UserCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", fmt.Errorf("THe string is empty")
	}
	return code, nil
}
func (r userrepository) getuno(code string) (result *User, err error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("the code is empty")
	}
	collection := r.db.Collection(collectionName)
	filter := bson.M{"usercode": code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		return nil, fmt.Errorf("no results found")
	}
	return result, nil
}
func (r userrepository) emailexist(email string) bool {
	if len(email) == 0 {
		return false
	}
	collection := r.db.Collection(collectionName)
	result := &User{}
	filter := bson.M{"email": email}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	return err1 == nil
}
