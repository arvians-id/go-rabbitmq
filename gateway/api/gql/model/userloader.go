package model

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"net/http"
	"time"
)

func NewUserLoaderConfig(userService services.UserServiceContract, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := UserLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 100,
			fetch: func(keys []int64) ([]*User, []error) {
				users, _, err := userService.FindByIDs(r.Context(), &pb.GetUserByIDsRequest{
					Ids: keys,
				})
				if err != nil {
					return nil, []error{err}
				}

				u := make(map[int64]*User, len(users.Users))
				for _, user := range users.Users {
					u[user.Id] = &User{
						Id:        user.Id,
						Name:      user.Name,
						Email:     user.Email,
						CreatedAt: user.CreatedAt,
						UpdatedAt: user.UpdatedAt,
					}
				}

				result := make([]*User, len(keys))

				for i, key := range keys {
					result[i] = u[key]
				}

				return result, nil
			},
		}

		ctx := context.WithValue(r.Context(), "userLoader", &userLoader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Loader(ctx context.Context, userService services.UserServiceContract) UserLoader {
	userLoader := UserLoader{
		wait:     2 * time.Millisecond,
		maxBatch: 100,
		fetch: func(keys []int64) ([]*User, []error) {
			users, _, err := userService.FindByIDs(ctx, &pb.GetUserByIDsRequest{
				Ids: keys,
			})
			if err != nil {
				return nil, []error{err}
			}

			u := make(map[int64]*User, len(users.Users))
			for _, userData := range users.Users {
				u[userData.Id] = &User{
					Id:        userData.Id,
					Name:      userData.Name,
					Email:     userData.Email,
					CreatedAt: userData.CreatedAt,
					UpdatedAt: userData.UpdatedAt,
				}
			}

			result := make([]*User, len(keys))

			for i, key := range keys {
				result[i] = u[key]
			}

			return result, nil
		},
	}

	return userLoader
}

func GetUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value("userLoader").(*UserLoader)
}
