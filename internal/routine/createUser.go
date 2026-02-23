package routine

import (
	"context"
	"log"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/internal/util/crypto"
	"myoptions.info/indigo/backend/model/entity"
)

func RunCreateUser(flags util.CreateUserRuntimeFlags) int {
	database := util.ConnectToDatabase()
	ctx := context.Background()

	userCount, countErr := gorm.G[entity.LocalUser](database).Select("username = ?", flags.Username).Count(ctx, "username")
	if countErr != nil {
		log.Fatalf("Failed to count users: %s", countErr)
	}
	if userCount != 0 {
		log.Fatalf("User entry already exists for username %s", flags.Username)
	}

	hash, hashErr := crypto.Hash(flags.Password)
	if hashErr != nil {
		log.Fatalf("Failed to hash password: %s", hashErr)
	}
	user := &entity.LocalUser{
		Username:     flags.Username,
		PasswordHash: hash,
		ExpiresAt:    nil,
	}

	createErr := gorm.G[entity.LocalUser](database).Create(ctx, user)
	if createErr != nil {
		log.Fatalf("Failed to create user: %s", createErr)
	}

	return 0
}
