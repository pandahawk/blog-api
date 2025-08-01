package post

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/gorm"
	"log"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=post

type Repository interface {
	FindAll() ([]*model.Post, error)
	FindByID(id uuid.UUID) (*model.Post, error)
	Create(post *model.Post) (*model.Post, error)
	Delete(post *model.Post) error
	Update(post *model.Post) (*model.Post, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) FindAll() ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.Preload("User").Find(&posts).Error
	return posts, err
}

func (r repository) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r repository) Create(post *model.Post) (*model.Post, error) {
	err := r.db.Preload("User").Create(post).Error

	if err := r.db.Preload("User").
		First(&post, "id = ?", post.ID).Error; err != nil {
		return nil, err
	}
	return post, err
}

func (r repository) Delete(post *model.Post) error {
	err := r.db.Preload("User").Delete(post).Error
	return err
}

func (r repository) Update(post *model.Post) (*model.Post, error) {
	err := r.db.Preload("User").Save(post).Error
	return post, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func NewDevRepository(db *gorm.DB) Repository {

	var post *model.Post
	if err := db.First(&post).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("no posts found... initializing sample data")
		samplePosts := []*model.Post{
			{
				ID:    uuid.New(),
				Title: "Lorem Ipsum",
				Content: `Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
						Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, 
						when an unknown printer took a galley of type and scrambled it to make a type specimen book.
						
						It has survived not only five centuries, but also the leap into electronic typesetting, 
						remaining essentially unchanged. It was popularised in the 1960s with the release of 
						Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing 
						software like Aldus PageMaker including versions of Lorem Ipsum.`,
				UserID: uuid.MustParse("3d9f18b2-f029-4a44-baf8-7437d51967d7"),
			},
			{
				ID:    uuid.New(),
				Title: "Where does it come from?",
				Content: `Contrary to popular belief, Lorem Ipsum is not simply random text. 
						It has roots in a piece of classical Latin literature from 45 BC, making it over 
						2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College 
						in Virginia, looked up one of the more obscure Latin words, consectetur, from a 
						Lorem Ipsum passage, and going through the cites of the word in classical literature, 
						discovered the undoubtable source.
						
						Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" 
						(The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise 
						on the theory of ethics, very popular during the Renaissance. 
						
						The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in 
						section 1.10.32.
						
						The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those 
						interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by 
						Cicero are also reproduced in their exact original form, accompanied by English 
						versions from the 1914 translation by H. Rackham.`,
				UserID: uuid.MustParse("3d9f18b2-f029-4a44-baf8-7437d51967d7"),
			},
			{
				ID:    uuid.New(),
				Title: "Why do we use it?",
				Content: `It is a long established fact that a reader will bedistracted by the readable 
						content of a page when looking at its layout. The point of using Lorem Ipsum is 
						that it has a more-or-less normal distribution of letters, as opposed to using 
						'Content here, content here', making it look like readable English. 
						
						Many desktop publishing packages and web page editors now use Lorem Ipsum as their 
						default model text, and a search for 'lorem ipsum' will uncover many websites still 
						in their infancy. Various versions have evolved over the years, sometimes by accident, 
						sometimes on purpose (injected humour and the like).`,
				UserID: uuid.MustParse("27e6db8c-3432-456e-a879-e7a0c58c9cc4"),
			},
		}
		if err := db.Save(&samplePosts).Error; err != nil {
			log.Fatal("error creating sample posts", err)
		}
		log.Println("init sample posts successfully")
	}

	return &repository{db: db}
}
