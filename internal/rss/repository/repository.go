package repository

import (
	"database/sql"
	"fmt"

	"github.com/Llifintsefv/GoRssFeed/internal/rss/models"
	"github.com/mmcdole/gofeed"
)

type RssRepository interface {
	SaveNews( *gofeed.Feed) error 
	GetNewsById(int) (models.News,error)
	GetAllNews() ([]models.News,error)
}

type rssRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RssRepository{
	return &rssRepository{db: db}
}

func (r *rssRepository) SaveNews(feed *gofeed.Feed) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO news (title,link,published,description) VALUES ($1,$2,$3,$4)")
	if err != nil {
		_ = tx.Rollback() 
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, item := range feed.Items {
		_, err := stmt.Exec(item.Title, item.Link, item.Published, item.Description)
		if err != nil {
			_ = tx.Rollback() 
			return fmt.Errorf("failed to save news item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *rssRepository) GetNewsById(id int) (models.News,error) {
	stmt,err := r.db.Prepare("SELECT id,title, link, published, description FROM news WHERE id = $1")
	if err != nil {
		return models.News{},fmt.Errorf("failte to prepare statement: %w",err)
	}
	defer stmt.Close()
	rows,err :=stmt.Query(id)
	
	if err != nil {
		return models.News{},err
	}
	defer rows.Close()
	var news models.News
	if rows.Next(){
		err = rows.Scan(&news.ID,&news.Title, &news.Link, &news.Published, &news.Description)
		if err != nil {
			return models.News{},err
		}
		return news, nil
	}
	return models.News{},fmt.Errorf("id not found %d",id)
}

func (r *rssRepository) GetAllNews() ([]models.News, error) {
	stmt, err := r.db.Prepare("SELECT title, link, published, description FROM news")
	if err != nil {
		return []models.News{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []models.News{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err = rows.Scan(&n.Title, &n.Link, &n.Published, &n.Description)
		if err != nil {
			return []models.News{}, fmt.Errorf("failed to scan news: %w", err)
		}
		newsList = append(newsList, n)
	}
	if err := rows.Err(); err != nil {
		return []models.News{}, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return newsList, nil
}