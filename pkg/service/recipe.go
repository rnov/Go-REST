package service

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/db"
	r "github.com/rnov/Go-REST/pkg/recipe"
	"github.com/rnov/Go-REST/pkg/errors"
)


type RcpSrv interface {
	GetById(recipeId string) (*r.Recipe, error)
	ListAll() ([]*r.Recipe, error)
	Create(recipe *r.Recipe) (*r.Recipe, error)
	Update(recipe *r.Recipe) error
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

func (r *rcp) GetById(recipeId string) (*r.Recipe, error) {

	rcp, err := r.rcpDb.GetRecipeById(recipeId)
	if rcp == nil {
		return nil, errors.NewNotFoundErr(fmt.Sprintf("rcp with id: %s was not found", rcp.ID))
	}
	if err != nil {
		return nil, errors.NewDBErr(err.Error())
	}

	return rcp, nil

}

func (r *rcp) ListAll() ([]*r.Recipe, error) {
	recipes, err := r.rcpDb.GetAllRecipes()
	if err != nil {
		return nil, errors.NewDBErr(err.Error())
	}

	return recipes, nil
}

func (r *rcp) Create(recipe *r.Recipe) (*r.Recipe, error) {
	if v := validateRecipeInput(recipe); len(v) > 0 {
		return nil, errors.NewInvalidParamsErr(v)
	}
	if err := r.rcpDb.CreateRecipe(recipe); err != nil {
		return nil, err
	}

	return recipe, nil
}

func (r *rcp) Update(recipe *r.Recipe) error {
	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}
	err := r.rcpDb.UpdateRecipe(recipe)
	if err != nil {
		return err
	}

	return nil
}

func (r *rcp) Delete(recipeId string) error {
	if err := r.rcpDb.DeleteRecipe(recipeId); err != nil {
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
	if recipe.PrepTime <= 1 || recipe.PrepTime > 1000 {
		valid[errors.Preptime] = errors.OutOfRange
	}
	return valid
}
