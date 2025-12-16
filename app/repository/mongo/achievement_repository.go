package repositorymongo

import (
	"context"
	"errors"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/mongo"
	"github.com/Ahlul-Mufi/data-prestasi/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const achievementsCollection = "achievements"

type AchievementRepository interface {
	Create(achievement m.Achievement) (m.Achievement, error)
	GetByID(id primitive.ObjectID) (m.Achievement, error)
	GetByStudentID(studentID string) ([]m.Achievement, error)
	GetMultipleByIDs(ids []primitive.ObjectID) ([]m.Achievement, error)
	Update(id primitive.ObjectID, achievement m.Achievement) (m.Achievement, error)
	SoftDelete(id primitive.ObjectID) error
	GetAll(filter bson.M, skip, limit int64) ([]m.Achievement, int64, error)
}

type achievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository() AchievementRepository {
	return &achievementRepository{
		collection: database.GetCollection(achievementsCollection),
	}
}

func (r *achievementRepository) Create(achievement m.Achievement) (m.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	achievement.CreatedAt = now
	achievement.UpdatedAt = now
	achievement.IsDeleted = false

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return m.Achievement{}, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)

	return achievement, nil
}

func (r *achievementRepository) GetByID(id primitive.ObjectID) (m.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var achievement m.Achievement
	filter := bson.M{
		"_id":       id,
		"isDeleted": false,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&achievement)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return m.Achievement{}, errors.New("achievement not found")
		}
		return m.Achievement{}, err
	}

	return achievement, nil
}

func (r *achievementRepository) GetByStudentID(studentID string) ([]m.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"studentId": studentID,
		"isDeleted": false,
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []m.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *achievementRepository) GetMultipleByIDs(ids []primitive.ObjectID) ([]m.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(ids) == 0 {
		return []m.Achievement{}, nil
	}

	filter := bson.M{
		"_id":       bson.M{"$in": ids},
		"isDeleted": false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []m.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *achievementRepository) Update(id primitive.ObjectID, achievement m.Achievement) (m.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":       id,
		"isDeleted": false,
	}

	update := bson.M{
		"$set": bson.M{
			"achievementType": achievement.AchievementType,
			"title":           achievement.Title,
			"description":     achievement.Description,
			"details":         achievement.Details,
			"tags":            achievement.Tags,
			"points":          achievement.Points,
			"updatedAt":       achievement.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated m.Achievement

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return m.Achievement{}, errors.New("achievement not found")
		}
		return m.Achievement{}, err
	}

	return updated, nil
}

func (r *achievementRepository) SoftDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":       id,
		"isDeleted": false,
	}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"updatedAt": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("achievement not found or already deleted")
	}

	return nil
}

func (r *achievementRepository) GetAll(filter bson.M, skip, limit int64) ([]m.Achievement, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}
	filter["isDeleted"] = false

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var achievements []m.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, 0, err
	}

	return achievements, total, nil
}

func (r *achievementRepository) AddAttachment(id primitive.ObjectID, attachment m.Attachment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":       id,
		"isDeleted": false,
	}

	update := bson.M{
		"$push": bson.M{"attachments": attachment},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}
