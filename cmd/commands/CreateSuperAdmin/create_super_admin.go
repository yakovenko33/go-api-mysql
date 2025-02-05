package commands

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	password "go-api-docker/GoCrm/Common/Security/Password"
	users_entities "go-api-docker/GoCrm/Users/Domains/Entities"
	database "go-api-docker/database"
)

var СreateSuperAdmin = &cobra.Command{
	Use:   "create-super-admin",
	Short: "Create super admin command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create super admin!")

		userByEmail, err := findUserByEmail()
		if err != nil {
			fmt.Println("Error when findUserByEmail:", err)
			return
		}
		if userByEmail != nil {
			fmt.Println("Super Admin alredy created.", err)
			return
		}

		user, err := factoryUser()
		if err != nil {
			fmt.Println("Error when factory user:", err)
			return
		}

		result := database.DB.Create(user)

		if result.Error != nil {
			fmt.Println("Error inserting user:", result.Error)
		} else {
			fmt.Println("User created with ID:", user.ID)
		}
	},
}

func findUserByEmail() (*users_entities.User, error) {
	fmt.Println("Find record by email")
	var user users_entities.User
	result := database.DB.Where("email = ?", os.Getenv("SUPER_ADMIN_EMAIL")).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, nil
}

func factoryUser() (users_entities.User, error) {
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	hashedPassword, err := password.HashPassword(os.Getenv("SUPER_ADMIN_PASSWORD"))
	if err != nil {
		return users_entities.User{}, err
	}

	user := users_entities.User{
		ID:             uuid.New(),
		FirstName:      "Super",
		LastName:       "Adminr",
		Email:          os.Getenv("SUPER_ADMIN_EMAIL"),
		Password:       hashedPassword,
		CreatedAt:      formattedTime,
		UpdatedAt:      formattedTime,
		Status:         1,
		CreatedBy:      formattedTime,
		ModifiedUserId: formattedTime,
	}
	return user, nil
}
