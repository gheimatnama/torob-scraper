package torob

import (
	"github.com/jinzhu/gorm"
	"time"
	"torobSpider/rotator"
)
type (
	SearchItem struct {
		ImageUrl string `gorm:"primary_key" json:"image_url" gorm:"size:1024"`
		RandomKey string `json:"random_key"`
		PersianName string `json:"name1" gorm:"size:1024"`
		EnglishName string `json:"name2" gorm:"size:1024"`
		MoreInfoUrl string `json:"more_info_url" gorm:"size:1024"`
		WebClientAbsoluteUrl string `json:"web_client_absolute_url" gorm:"size:1024"`
		PriceText string `json:"price_text"`
		PriceTextMode string `json:"price_text_mode"`
		ShopText string `json:"shop_text" gorm:"size:1024"`
	}

	RuntimeInfo struct {
		MaxRunningWorkers int
		WorkerPool chan int
		MaxParallelProductPerSearch chan int
		MaxParallelSearch chan int
		DB *gorm.DB
		QueriesFile string
		SearchResultLimit int
		OnlyRepairDownloadedSources bool
		ProxyRotator *rotator.ProxyRotator
	}

	SearchQuery struct {
		ProductID int `gorm:"primary_key" json:"id"`
		Title string `json:"title" gorm:"size:1024"`
		CreatedAt time.Time
	}

	Pagination struct {
		CurrentPage int `json:"current_page"`
		Data []interface{} `json:"data"`
		FirstPageUrl string `json:"first_page_url"`
		NextPageUrl string `json:"next_page_url"`
		PrevPageUrl string `json:"prev_page_url"`
		LastPageUrl string `json:"last_page_url"`
		From int `json:"from"`
		LastPage int `json:"last_page"`
		PerPage int `json:"per_page"`
		To int `json:"to"`
		Total int `json:"total"`
	}

	ProductSource struct {
		gorm.Model
		ProductID string `json:"-"`
		PrimaryName string `json:"name1"`
		SecondaryName string `json:"name2"`
		ShopName string `json:"shop_name"`
		ShopLocation string `json:"shop_name2"`
		ShopId int64 `json:"shop_id"`
		PageUrl string `json:"page_url" gorm:"size:1024"`
		DirectPageUrl string `json:"-" gorm:"size:1024"`
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
