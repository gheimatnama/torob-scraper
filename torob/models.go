package torob

import (
	"github.com/jinzhu/gorm"
	"torobSpider/rotator"
)
type (
	SearchItem struct {
		ImageUrl string `gorm:"primary_key" json:"image_url"`
		RandomKey string `json:"random_key"`
		PersianName string `json:"name1"`
		EnglishName string `json:"name2"`
		MoreInfoUrl string `json:"more_info_url"`
		WebClientAbsoluteUrl string `json:"web_client_absolute_url"`
		PriceText string `json:"price_text"`
		PriceTextMode string `json:"price_text_mode"`
		ShopText string `json:"shop_text"`
	}

	RuntimeInfo struct {
		MaxRunningWorkers int
		WorkerPool chan int
		DB *gorm.DB
		QueriesFile string
		SearchResultLimit int
		ProxyRotator *rotator.ProxyRotator
	}

	ProductSource struct {
		gorm.Model
		ProductID string `json:"-"`
		PrimaryName string `json:"name1"`
		SecondaryName string `json:"name2"`
		ShopName string `json:"shop_name"`
		ShopLocation string `json:"shop_name2"`
		ShopId int64 `json:"shop_id"`
		PageUrl string `json:"page_url"`
		DirectPageUrl string `json:"-"`
		SourceID string `json:"-"`
		IDInSource string `json:"-"`
		ShopScorePercentile float64 `json:"shop_score_percentile"`
		ShopVotesCount int64 `json:"shop_votes_count"`
		Price int64 `json:"price"`
		PriceText string `json:"price_text"`
		PriceTextStriked string `json:"price_text_striked"`
		PriceTextMode string `json:"price_text_mode"`
	}

	ProductsInfo struct {
		Result []ProductSource
	}

	Product struct {
		SearchItem
		ProductsInfo ProductsInfo `json:"products_info" gorm:"-"`
	}
	
	SearchResult struct {
		Results []SearchItem `json:"results"`
	}
)