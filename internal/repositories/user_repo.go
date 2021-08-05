package repositories

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shurcooL/graphql"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
)

type GraphError struct {
	Message string
	Path    []string
}

// UserRepo connects to minerva owl GraphQL server to manage users
// Implements ports.UserRepo interface
type UserRepo struct {
	config *domain.Config
	client *graphql.Client
}

// NewUserRepo creates an instance of UserRepo
func NewUserRepo(config *domain.Config) *UserRepo {
	client := graphql.NewClient(config.UserRepo.Url, nil)
	return &UserRepo{
		config: config,
		client: client,
	}
}

// TODO: Handle repository connection errors, currently they will generate internal server error

func (repo *UserRepo) Create(user domain.Register) (domain.User, error) {
	var m struct {
		CreateUser struct {
			Id       graphql.String
			Name     graphql.String
			Username graphql.String
			Picture  graphql.String
		} `graphql:"createUser(input:{name: $name, username: $username, role: $role, tokenID: $tokenID, provider: $provider, picture: $picture, status: \"active\"})"`
	}

	vars := map[string]interface{}{
		"name":     graphql.String(user.Name),
		"username": graphql.String(user.Username),
		"picture":  graphql.String(user.Picture),
		"role":     graphql.String(user.Role),
		"provider": graphql.String(user.Provider),
		"tokenID":  graphql.String(user.TokenID),
	}

	err := repo.client.Mutate(context.Background(), &m, vars)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       string(m.CreateUser.Id),
		Name:     string(m.CreateUser.Name),
		Username: string(m.CreateUser.Username),
		Picture:  string(m.CreateUser.Picture),
	}, nil
}

func (repo *UserRepo) GetById(id string) (domain.User, error) {
	var query struct {
		User struct {
			Id       graphql.String
			Name     graphql.String
			Username graphql.String
			Picture  graphql.String
		} `graphql:"user(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": id,
	}

	err := repo.client.Query(context.Background(), &query, vars)
	if err != nil {
		log.Debug().Err(err).Msgf("Repo GetById Error:")
		return domain.User{}, err
	}

	return domain.User{
		Id:       string(query.User.Id),
		Name:     string(query.User.Name),
		Username: string(query.User.Username),
		Picture:  string(query.User.Picture),
	}, nil
}

func (repo *UserRepo) GetByUsername(username string) (domain.User, error) {
	var query struct {
		User struct {
			Id       graphql.String
			Name     graphql.String
			Username graphql.String
			Picture  graphql.String
		} `graphql:"userByUsername(username: $username)"`
	}

	vars := map[string]interface{}{
		"username": graphql.String(username),
	}

	err := repo.client.Query(context.Background(), &query, vars)
	if err != nil {
		log.Debug().Err(err).Msgf("Repo GetById Error:")
		return domain.User{}, err
	}

	return domain.User{
		Id:       string(query.User.Id),
		Name:     string(query.User.Name),
		Username: string(query.User.Username),
		Picture:  string(query.User.Picture),
	}, nil
}
