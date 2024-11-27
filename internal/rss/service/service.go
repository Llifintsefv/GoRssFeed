package service

import (
	"fmt"

	"github.com/Llifintsefv/GoRssFeed/internal/rss/models"
	"github.com/Llifintsefv/GoRssFeed/internal/rss/repository"
	"github.com/mmcdole/gofeed"
)

type RssService interface {
	GetNewsById(int) (models.News,error)
	FetchNewsFromSource(string) (error)
	GetAllNews() ([]models.News,error)
}

type rssService struct{
	repo repository.RssRepository
}

func NewRssService(repo repository.RssRepository) RssService {
	return &rssService{repo: repo}
} 

func (s *rssService) GetNewsById(id int) (models.News,error){
	result,err := s.repo.GetNewsById(id)
	if err != nil {
		return models.News{},err
	}
	return result,nil
}

func (s *rssService) FetchNewsFromSource(source string) (error){
	fp := gofeed.NewParser()
	feed,err := fp.ParseURL(source)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	err = s.repo.SaveNews(feed)
	if err != nil {
		return fmt.Errorf("failed to save news: %w", err)
	}
	return nil
}

func (s *rssService) GetAllNews() ([]models.News,error){
	return s.repo.GetAllNews()
}