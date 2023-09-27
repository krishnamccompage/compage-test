package daos

import (
	"context"
	"errors"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/daos/clients/nosqls"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDao struct {
	mongoDBClient *nosqls.MongoDBClient
	collection    *mongo.Collection
}

func NewUserDao() (*UserDao, error) {
	mongoDBClient, err := nosqls.InitMongoDB()
	if err != nil {
		log.Debugf("init mongoDB failed: %v", err)
		return nil, err
	}
	return &UserDao{
		mongoDBClient: mongoDBClient,
		collection:    mongoDBClient.Database.Collection("users"),
	}, nil
}

func (userDao *UserDao) CreateUser(user *models.User) (*models.User, error) {
	// create a document for given user
	insertOneResult, err := userDao.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Debugf("insert failed: %v", err)
		return nil, err
	}
	user.ID = insertOneResult.InsertedID.(primitive.ObjectID).Hex()

	log.Debugf("user created")
	return user, nil
}

func (userDao *UserDao) UpdateUser(id string, user *models.User) (*models.User, error) {
	if id != user.ID {
		log.Debugf("id(%s) and payload(%s) don't match", id, user.ID)
		return nil, errors.New("id and payload don't match")
	}

	existingUser, err := userDao.GetUser(id)
	if err != nil {
		log.Debugf("get user failed: %v", err)
		return nil, err
	}

	// no user retrieved means no user exists to update
	if existingUser == nil || existingUser.ID == "" {
		return nil, nosqls.ErrNotExists
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nosqls.ErrInvalidObjectID
	}
	filter := bson.D{
		{Key: "_id", Value: objectID},
	}
	user.ID = ""
	update := bson.M{
		"$set": user,
	}

	updateResult, err := userDao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Debugf("update user failed: %v", err)
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, nosqls.ErrUpdateFailed
	}

	log.Debugf("user updated")
	return user, nil
}

func (userDao *UserDao) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nosqls.ErrInvalidObjectID
	}
	filter := bson.D{
		{Key: "_id", Value: objectID},
	}
	deleteResult, err := userDao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			return nosqls.ErrNotExists
		}
		log.Debugf("delete user failed: %v", err)
		return nosqls.ErrDeleteFailed
	}
	if deleteResult.DeletedCount == 0 {
		return nosqls.ErrNotExists
	}

	log.Debugf("user deleted")
	return nil
}

func (userDao *UserDao) ListUsers() ([]*models.User, error) {
	filters := bson.D{}
	users, err := userDao.collection.Find(context.TODO(), filters)
	if err != nil {
		return nil, err
	}
	var userList []*models.User
	for users.Next(context.TODO()) {
		var user *models.User
		if err = users.Decode(&user); err != nil {
			log.Debugf("decode user failed: %v", err)
			return nil, err
		}
		userList = append(userList, user)
	}
	if userList == nil {
		return []*models.User{}, nil
	}

	log.Debugf("user listed")
	return userList, nil
}

func (userDao *UserDao) GetUser(id string) (*models.User, error) {
	var user *models.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.User{}, nosqls.ErrInvalidObjectID
	}
	filter := bson.D{
		{Key: "_id", Value: objectID},
	}
	if err = userDao.collection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			return &models.User{}, nosqls.ErrNotExists
		}
		log.Debugf("decode user failed: %v", err)
		return nil, err
	}

	log.Debugf("user retrieved")
	return user, nil
}
