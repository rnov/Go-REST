package service

import (
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	r "github.com/rnov/Go-REST/pkg/recipe"
	"regexp"
)

type RcpSrv interface {
	GetByID(recipeId string) (*r.Recipe, error)
	ListAll() ([]*r.Recipe, error)
	Create(recipe *r.Recipe) error
	Update(ID string, recipe *r.Recipe) error
	Delete(recipeId string) error
}

// this is a must, struct can not implement interface from different package.
type rcp struct {
	rcpDb db.Recipe
	//logger log.Loggers
	// add more func fields
}

func NewRecipe(rcpDb db.Recipe) *rcp {
	recipeSrv := &rcp{
		rcpDb: rcpDb,
	}
	return recipeSrv
}

func (r *rcp) GetByID(ID string) (*r.Recipe, error) {
	if !validateRcpID(ID) {
		return nil, errors.NewUserErr("invalid ID format")
	}
	rcp, err := r.rcpDb.GetRecipeById(ID)
	if err != nil {
		return nil, err
	}

	return rcp, nil
}

func (r *rcp) ListAll() ([]*r.Recipe, error) {
	recipes, err := r.rcpDb.GetAllRecipes()
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *rcp) Create(recipe *r.Recipe) error {
	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}
	if err := r.rcpDb.CreateRecipe(recipe); err != nil {
		return err
	}

	return nil
}

func (r *rcp) Update(ID string, recipe *r.Recipe) error {
	if !validateRcpID(ID) {
		return errors.NewUserErr("invalid ID format")
	}
	if ID != recipe.ID {
		return errors.NewUserErr("ID param and recipe ID do not match")
	}

	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}
	err := r.rcpDb.UpdateRecipe(recipe)
	if err != nil {
		return err
	}

	return nil
}

func (r *rcp) Delete(recipeID string) error {
	if !validateRcpID(recipeID) {
		return errors.NewUserErr("invalid ID format")
	}
	if err := r.rcpDb.DeleteRecipe(recipeID); err != nil {
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
