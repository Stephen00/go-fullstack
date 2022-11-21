package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null" json:"content"`
	Author    User      `json:"author"`
	AuthorId  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("title required")
	}

	if p.Content == "" {
		return errors.New("content required")
	}

	if p.AuthorId < 1 {
		return errors.New("author required")
	}

	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	if err := db.Debug().Model(&Post{}).Create(&p).Error; err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		if err := db.Debug().Model(&User{}).Where("id = ?", p.AuthorId).Take(&p.Author).Error; err != nil {
			return &Post{}, nil
		}
	}
	return p, nil
}

func (p *Post) FindAllUsers(db *gorm.DB) (*[]Post, error) {
	posts := []Post{}
	if err := db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error; err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			if err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorId).Take(&p.Author).Error; err != nil {
				return &[]Post{}, nil
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostById(db *gorm.DB, pid uint64) (*Post, error) {
	if err := db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error; err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		if err := db.Debug().Model(&User{}).Where("id = ?", p.AuthorId).Take(&p.Author).Error; err != nil {
			return &Post{}, nil
		}
	}

	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	if err := db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error; err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		if err := db.Debug().Model(&User{}).Where("id = ?", p.AuthorId).Take(&p.Author).Error; err != nil {
			return &Post{}, nil
		}
	}
	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
