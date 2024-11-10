package user

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	collection *mongo.Collection
}

func NewService(collection *mongo.Collection) *Service {
	return &Service{collection: collection}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	log.Printf("Creating user with email: %s", req.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}

	user := &User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Successfully created user with ID: %s", user.ID.Hex())
	return user, nil
}

func (s *Service) List(ctx context.Context, filter FilterOptions) ([]User, error) {
	findOptions := options.Find()
	findOptions.SetSkip(int64((filter.Page - 1) * filter.PageSize))
	findOptions.SetLimit(int64(filter.PageSize))

	if filter.Sort != "" {
		findOptions.SetSort(bson.D{{Key: filter.Sort, Value: 1}})
	}

	query := bson.M{}
	if filter.Search != "" {
		query = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
				{"email": bson.M{"$regex": filter.Search, "$options": "i"}},
			},
		}
	}

	cursor, err := s.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name":       req.Name,
			"role":       req.Role,
			"updated_at": time.Now(),
		},
	}

	var user User
	err = s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
} 