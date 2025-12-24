package services

import (
	"iss-model-backend/internal/models"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(
	title, content, excerpt, author, publishDate, readTime, image string,
	images, tags []string,
) (*models.Post, error) {
	// Ensure arrays are not nil
	if images == nil {
		images = []string{}
	}
	if tags == nil {
		tags = []string{}
	}

	post := &models.Post{
		Title:       title,
		Content:     content,
		Excerpt:     excerpt,
		Author:      author,
		PublishDate: publishDate,
		ReadTime:    readTime,
		Image:       image,
		Images:      images,
		Tags:        tags,
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := s.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := s.db.Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) UpdatePost(
	id uint,
	title, content, excerpt, author, publishDate, readTime, image string,
	images, tags []string,
) (*models.Post, error) {
	post, err := s.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure arrays are not nil
	if images == nil {
		images = []string{}
	}
	if tags == nil {
		tags = []string{}
	}

	post.Title = title
	post.Content = content
	post.Excerpt = excerpt
	post.Author = author
	post.PublishDate = publishDate
	post.ReadTime = readTime
	post.Image = image
	post.Images = images
	post.Tags = tags

	if err := s.db.Save(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) DeletePost(id uint) error {
	if err := s.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}
	return nil
}
