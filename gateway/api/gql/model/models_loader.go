package model

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"time"
)

func UserServiceFindByIDs(ctx context.Context, userService services.UserServiceContract) UserLoader {
	return UserLoader{
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
}

func TodoServiceFindByUserIDs(ctx context.Context, todoService services.TodoServiceContract) TodoLoader {
	return TodoLoader{
		wait:     2 * time.Millisecond,
		maxBatch: 100,
		fetch: func(keys []int64) ([][]*Todo, []error) {
			todos, _, err := todoService.FindByUserIDs(ctx, &pb.GetTodoByUserIDsRequest{
				Ids: keys,
			})
			if err != nil {
				return nil, []error{err}
			}

			u := make(map[int64][]*Todo, len(todos.Todos))
			for _, todo := range todos.Todos {
				u[todo.UserId] = append(u[todo.UserId], &Todo{
					Id:          todo.Id,
					Title:       todo.Title,
					Description: todo.Description,
					IsDone:      *todo.IsDone,
					UserId:      todo.UserId,
					CreatedAt:   todo.CreatedAt,
					UpdatedAt:   todo.UpdatedAt,
				})
			}

			result := make([][]*Todo, len(keys))
			for i, key := range keys {
				result[i] = u[key]
			}

			return result, nil
		},
	}
}

func CategoryServiceFindByTodoIDs(ctx context.Context, categoryService services.CategoryServiceContract) CategoryLoader {
	return CategoryLoader{
		wait:     2 * time.Millisecond,
		maxBatch: 100,
		fetch: func(keys []int64) ([][]*Category, []error) {
			categories, _, err := categoryService.FindByTodoIDs(ctx, &pb.GetCategoryByTodoIDsRequest{
				Id: keys,
			})
			if err != nil {
				return nil, []error{err}
			}

			u := make(map[int64][]*Category, len(categories.Categories))
			for _, category := range categories.Categories {
				u[category.TodoId] = append(u[category.TodoId], &Category{
					Id:   category.Id,
					Name: category.Name,
				})
			}

			result := make([][]*Category, len(keys))
			for i, key := range keys {
				result[i] = u[key]
			}

			return result, nil
		},
	}
}
