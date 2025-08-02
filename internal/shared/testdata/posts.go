package testdata

import (
	"github.com/pandahawk/blog-api/internal/shared/model"
)

var (
	Post1 = &model.Post{
		ID:      PostIDs[0],
		Title:   "Exploring the Cosmos",
		Content: "Today, I pondered the vastness of space and the mysteries it holds beyond our imagination.",
		UserID:  Alice.ID,
		User:    Alice,
	}
	Post2 = &model.Post{
		ID:      PostIDs[1],
		Title:   "Why Cats Rule the Internet",
		Content: "From memes to videos, cats have conquered our hearts and the digital world alike.",
		UserID:  Bob.ID,
		User:    Bob,
	}
	Post3 = &model.Post{
		ID:      PostIDs[2],
		Title:   "The Art of Making Pizza",
		Content: "Nothing brings people together like the smell of a freshly baked pizza in the kitchen.",
		UserID:  Caren.ID,
		User:    Caren,
	}
	Post4 = &model.Post{
		ID:      PostIDs[3],
		Title:   "Running in the Rain",
		Content: "Despite the gloomy weather, today's run was refreshing and oddly peaceful.",
		UserID:  Alice.ID,
		User:    Alice,
	}
	Post5 = &model.Post{
		ID:      PostIDs[4],
		Title:   "Tech Trends in 2025",
		Content: "Artificial intelligence and quantum computing are shaping our future in surprising ways.",
		UserID:  Bob.ID,
		User:    Bob,
	}
	Post6 = &model.Post{
		ID:      PostIDs[5],
		Title:   "A Quiet Morning",
		Content: "The world seems to pause at sunrise, offering a moment of calm before the day begins.",
		UserID:  Caren.ID,
		User:    Caren,
	}
)

var SamplePosts = []*model.Post{Post1, Post2,
	Post3, Post4, Post5, Post6}
