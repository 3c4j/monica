package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func initDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&User{})
	return db, func() {
		db.Migrator().DropTable(&User{})
		os.Remove("test.db")
	}
}

func TestRepository_Create(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	db.Raw("SELECT * FROM users").Scan(&user)
	assert.Equal(t, user.Name, "test")
	assert.Equal(t, user.Email, "test@test.com")
}

func TestRepository_Update(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user1 := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	user2 := &User{
		ID:    2,
		Name:  "test2",
		Email: "test2@test.com",
	}

	err := repo.Create(context.Background(), user1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user2)
	assert.NoError(t, err)

	user1.Name = "test1"
	err = repo.Update(context.Background(), User{ID: user1.ID}, user1)
	assert.NoError(t, err)

	_user1 := &User{}
	db.Raw("SELECT * FROM users WHERE id = ?", user1.ID).Scan(&_user1)
	assert.Equal(t, _user1.Name, "test1")
	assert.Equal(t, _user1.Email, "test@test.com")

	_user2 := &User{}
	db.Raw("SELECT * FROM users WHERE id = ?", user2.ID).Scan(&_user2)
	assert.Equal(t, _user2.Name, "test2")
	assert.Equal(t, _user2.Email, "test2@test.com")
}

func TestRepository_FindOne(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user1 := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	user2 := &User{
		ID:    2,
		Name:  "test2",
		Email: "test2@test.com",
	}

	err := repo.Create(context.Background(), user1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user2)
	assert.NoError(t, err)

	user, err := repo.FindOne(context.Background(), User{ID: user1.ID})
	assert.NoError(t, err)
	assert.Equal(t, user.Name, "test")
	assert.Equal(t, user.Email, "test@test.com")

	user, err = repo.FindOne(context.Background(), User{ID: 3})
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestRepository_Delete(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user1 := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	user2 := &User{
		ID:    2,
		Name:  "test2",
		Email: "test2@test.com",
	}

	err := repo.Create(context.Background(), user1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user2)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), User{ID: user1.ID})
	assert.NoError(t, err)

	user, err := repo.FindOne(context.Background(), User{ID: user1.ID})
	assert.Error(t, err)
	assert.Nil(t, user)

	user, err = repo.FindOne(context.Background(), User{ID: user2.ID})
	assert.NoError(t, err)
	assert.Equal(t, user.Name, "test2")
	assert.Equal(t, user.Email, "test2@test.com")
}

func TestRepository_Find(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user1 := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	user2 := &User{
		ID:    2,
		Name:  "test2",
		Email: "test2@test.com",
	}

	err := repo.Create(context.Background(), user1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user2)
	assert.NoError(t, err)

	users, err := repo.Find(context.Background(), User{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, user1.ID, users[0].ID)
	assert.Equal(t, user2.ID, users[1].ID)
	assert.Equal(t, user1.Name, users[0].Name)
	assert.Equal(t, user2.Name, users[1].Name)
	assert.Equal(t, user1.Email, users[0].Email)
	assert.Equal(t, user2.Email, users[1].Email)
}

func TestRepository_Count(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	user1 := &User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	user2 := &User{
		ID:    2,
		Name:  "test2",
		Email: "test2@test.com",
	}

	err := repo.Create(context.Background(), user1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user2)
	assert.NoError(t, err)

	count, err := repo.Count(context.Background(), User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestRepository_Paginate(t *testing.T) {
	db, cleanup := initDB(t)
	defer cleanup()

	repo := NewRepository[User](db)

	for i := 1; i <= 10; i++ {
		user := &User{
			ID:    i,
			Name:  fmt.Sprintf("test%d", i),
			Email: fmt.Sprintf("test%d@test.com", i),
		}
		err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	users, err := repo.Paginate(context.Background(), User{}, 1, 8)
	assert.NoError(t, err)
	assert.Equal(t, 8, len(users))
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, 8, users[7].ID)

	users, err = repo.Paginate(context.Background(), User{}, 2, 8)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, 9, users[0].ID)
	assert.Equal(t, 10, users[1].ID)
}
