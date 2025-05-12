package models

type Books struct {
	ID            int     `json:db:"id"`
	Title         string  `json:db:"title"`
	AuthorID      *int    `json:db:"author_id"`
	GenreID       *int    `json:db:"genre_id"`
	Price         float64 `json:db:"price"`
	StockQuantity int     `json:db:"stock_quantity"`
}

type Genres struct {
	ID   int    `json:db:"id`
	Name string `json:db:"name`
}

type Authors struct {
	ID   int    `json:db:"id"`
	Name string `json:db:"name"`
	Bio  string `json:db:"bio"`
}
