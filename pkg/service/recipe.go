package service

import (
	"regexp"

	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	r "github.com/rnov/Go-REST/pkg/recipe"
)

type RecipeMng interface {
	GetByID(recipeID string) (*r.Recipe, error)
	ListAll() ([]*r.Recipe, error)
	Create(recipe *r.Recipe) error
	Update(ID string, recipe *r.Recipe) error
	Delete(recipeID string) error
}

// this is a must, struct can not implement interface from different package.
type Recipe struct {
	rcpDB db.Recipe
	//logger log.Loggers
	// add more func fields
}

func NewRecipe(rcpDB db.Recipe) *Recipe {
	recipeSrv := &Recipe{
		rcpDB: rcpDB,
	}
	return recipeSrv
}

func (r *Recipe) GetByID(ID string) (*r.Recipe, error) {
	if !validateRcpID(ID) {
		return nil, errors.NewInputError("Invalid ID format", nil)
	}
	rcp, err := r.rcpDB.GetRecipeByID(ID)
	if err != nil {
		return nil, err
	}

	return rcp, nil
}

func (r *Recipe) ListAll() ([]*r.Recipe, error) {
	recipes, err := r.rcpDB.GetAllRecipes()
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *Recipe) Create(recipe *r.Recipe) error {
	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInputError("Invalid input parameters", v)
	}
	if err := r.rcpDB.CreateRecipe(recipe); err != nil {
		return err
	}

	return nil
}

func (r *Recipe) Update(ID string, recipe *r.Recipe) error {
	if !validateRcpID(ID) {
		return errors.NewInputError("Invalid ID format", nil)
	}
	if ID != recipe.ID {
		return errors.NewInputError("ID param and recipe ID do not match", nil)
	}

	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInputError("Invalid input parameters", v)
	}
	err := r.rcpDB.UpdateRecipe(recipe)
	if err != nil {
		return err
	}

	return nil
}

func (r *Recipe) Delete(recipeID string) error {
	if !validateRcpID(recipeID) {
		return errors.NewInputError("Invalid ID format", nil)
	}
	if err := r.rcpDB.DeleteRecipe(recipeID); err != nil {
		return err
	}
	return nil
}

func validateRecipeInput(recipe *r.Recipe) map[string]string {
	valid := make(map[string]string)

	if recipe.Difficulty <= 1 || recipe.Difficulty > 3 {
		valid[errors.Difficulty] = errors.OutOfRange
	}
	if len(recipe.Name) > 100 {
		valid[errors.Name] = errors.TooLong
	}
	if len(recipe.Name) == 0 {
		valid[errors.Name] = errors.MissingName
	}
	if recipe.PrepTime <= 1 || recipe.PrepTime > 1000 {
		valid[errors.Preptime] = errors.OutOfRange
	}
	return valid
}

func validateRcpID(ID string) bool {
	regex, _ := regexp.Compile("^[a-zA-Z0-9]{1,12}$")
	return regex.MatchString(ID)
}
