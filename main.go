package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/achsanit/go-quiz/model"
	"github.com/gin-gonic/gin"
)

var users []model.User

func main() {
	users = []model.User{}
	g := gin.Default()

	v1 := g.Group("api/v1")
	{
		usersGroup := v1.Group("/users")
		{
			// [GET] --> get all users
			usersGroup.GET("", func(ctx *gin.Context) {
				// return all users
				ctx.JSON(http.StatusOK, gin.H{
					"status": "ok",
					"users":  users,
				})
			})

			// [POST] --> create user
			usersGroup.POST("", func(ctx *gin.Context) {
				// binding payload
				user := model.User{}
				if err := ctx.Bind(&user); err != nil {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "failed to bind body",
					})
					return
				}
				user.ID = uint(len(users) + 1)
				users = append(users, user)
				ctx.JSON(http.StatusAccepted, map[string]any{
					"message": "user created",
				})
			})

			//[GET] --> get user by ID
			usersGroup.GET("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id")) // get id
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "id not found",
					})
					return
				}

				for _, user := range users {
					// find user with id
					if user.ID == uint(id) { // return when id match
						ctx.JSON(http.StatusOK, user)
						return
					}
				}

				// return json when id not found
				ctx.JSON(http.StatusNotFound, map[string]any{
					"message": "user not found",
				})
			})

			// [PUT]--> delete User by ID
			usersGroup.PUT("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id")) // get id
				if err != nil {
					// return id not found when path/param not number
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "id not found",
					})
					return
				}

				// binding body to model User
				user := model.User{}
				if err := ctx.Bind(&user); err != nil {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "failed to bind body",
					})
					return
				}

				// update user with data that already set with bind
				updatedUsers, err := updateUserByID(users, uint(id), user)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"message": "user not found",
					})
					return
				}
				users = updatedUsers // set list users with updated list

				//return json success update
				ctx.JSON(http.StatusOK, gin.H{
					"status":  "ok",
					"message": "user updated successfully",
				})
			})

			// [DELETE]--> delete User by ID
			usersGroup.DELETE("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id")) // get id
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "id not found",
					})
					return
				}

				// lop to search user with id
				for idx, user := range users {
					if user.ID == uint(id) {
						users = deleteUser(users, idx) // delete element with index item and update list user
						ctx.JSON(http.StatusOK, gin.H{
							"status":  "ok",
							"message": fmt.Sprintf("deleted on index %d", idx),
						})
						return
					}
				}

				// return when user not found
				ctx.JSON(http.StatusNotFound, map[string]any{
					"message": "user not found",
				})
			})
		}
	}

	g.Run(":8008")
}

func deleteUser(slice []model.User, index int) []model.User {
	// deleting element in slice with index
	// append slice item from first item until index - 1
	// and append slice item from index +1 until last item
	return append(slice[:index], slice[index+1:]...)
}

func updateUserByID(users []model.User, id uint, updatedUser model.User) ([]model.User, error) {
	// loop for check eact item in list
	for idx, user := range users {
		if user.ID == id { // when id item was match
			users[idx].Username = updatedUser.Username // update username
			users[idx].Email = updatedUser.Email       // update email
			return users, nil
		}
	}
	return nil, errors.New("user not found") // return error when id item not found
}
