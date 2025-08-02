package testdata

import (
	"github.com/pandahawk/blog-api/internal/shared/model"
)

var (
	Alice = &model.User{
		ID:       UserIDs[0],
		Username: "alice",
		Email:    "alice@example.com",
	}
	Bob = &model.User{
		ID:       UserIDs[1],
		Username: "bob",
		Email:    "bob@example.com",
	}
	Caren = &model.User{
		ID:       UserIDs[2],
		Username: "caren",
		Email:    "caren@example.com",
	}
	Dave = &model.User{
		ID:       UserIDs[3],
		Username: "dave",
		Email:    "dave@example.com",
	}
)
var SampleUsers = []*model.User{Alice, Bob, Caren, Dave}
