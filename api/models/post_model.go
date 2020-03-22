package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Post struct {
	gorm.Model
	Title    string `gorm:"size:255;not null;unique" json:"title"`
	Content  string `gorm:"size:255;not null;" json:"content"`
	Author   User   `json:"author"`
	AuthorID uint32 `gorm:"not null" json:"author_id"`
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
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Post{}, nil
	}
	return p, nil
}

func (p *Post) FindAllPost(db *gorm.DB) (*[]Post, error) {
	var posts []Post
	err := db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, nil
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err = db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			//if err != nil {
			//	return &[]Post{}, err
			//}
		}
	}
	return &posts, nil
}


func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Post{}, errors.New("data not found")
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) UpdatePostById(db *gorm.DB) (*Post, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) DeletePostById(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
