package create_super_admin

import (
	"errors"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	fmt_color "go-api-docker/cmd/cli/commands/enums/fmt_color"
	password "go-api-docker/internal/common/security/password"
	users_entities "go-api-docker/internal/go_crm/users/domains/entities"
)

func СreateSuperAdmin(DB *gorm.DB, accessControlModel *casbin.Enforcer) *cobra.Command {
	return &cobra.Command{
		Use:   "create-super-admin",
		Short: "Create super admin command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt_color.PrintMessage(fmt_color.Yellow, "Start create super admin!")

			err := DB.Transaction(func(tx *gorm.DB) error {
				err := createSuperAdmin(DB, accessControlModel)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				fmt_color.PrintError("Faild create super admin!" + err.Error())
			}
		},
	}
}

func createSuperAdmin(DB *gorm.DB, accessControlModel *casbin.Enforcer) error {
	userByEmail, err := findUserByEmail(DB)
	if err != nil {
		fmt_color.PrintError("Error when findUserByEmail:" + err.Error())
		return err
	}
	if userByEmail != nil {
		fmt_color.PrintError("Super Admin alredy created.")
		return err
	}

	user, err := factoryUser()
	if err != nil {
		fmt_color.PrintError("Error when creating user:" + err.Error())
		return err
	}

	result := DB.Create(user)
	if result.Error != nil {
		fmt_color.PrintError("Error inserting user:" + result.Error.Error())
		return err
	}

	fmt_color.PrintMessage(fmt_color.Green, "User created with ID: "+user.ID.String())
	accessControlModel.AddGroupingPolicy(user.ID.String(), "super_admin")

	return err
}

func findUserByEmail(DB *gorm.DB) (*users_entities.User, error) {
	var user users_entities.User
	result := DB.Where("email = ?", os.Getenv("SUPER_ADMIN_EMAIL")).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
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
		LastName:       "Admin",
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
