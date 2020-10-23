package recipe

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	log "github.com/rnov/Go-REST/pkg/logger"
)

type Recipe struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PrepTime   int    `json:"prepTime"`
	Difficulty int    `json:"difficulty"`
	Vegetarian bool   `json:"vegetarian"`
}

type RcpSrv interface {
	GetById(recipeId string) (*Recipe, error)
	ListAll() ([]*Recipe, error)
	Create(recipe *Recipe, auth string) (*Recipe, error)
	Update(recipe *Recipe, urlId string, auth string) error
	Delete(recipeId string, auth string) error
}

// this is a must, struct can not implement interface from different package.
type rcp struct {
	rcpDb  db.RecipesDbCalls
	logger log.Loggers
	// add more func fields
}

func NewRecipeSrv(rcpDb db.RecipesDbCalls, logger log.Loggers) *rcp {
	recipeSrv := &rcp{
		rcpDb:  rcpDb,
		logger: logger,
	}
	return recipeSrv
}

func (r *rcp) GetById(recipeId string) (*Recipe, error) {

	recipe, err := r.rcpDb.GetRecipeById(recipeId)

	if recipe == nil {
		r.logger.Notice(err)
		return nil, errors.NewNotFoundErr(fmt.Sprintf("recipe with id: %s was not found", recipe.ID))
	}

	if err != nil {
		r.logger.Error(err)
		return nil, errors.NewDBErr(err.Error())
	}

	return recipe, nil

}

func (r *rcp) ListAll() ([]*Recipe, error) {

	recipes, err := r.rcpDb.GetAllRecipes()
	if err != nil {
		r.logger.Error(err)
		return nil, errors.NewDBErr(err.Error())
	}

	return recipes, nil
}

func (r *rcp) Create(recipe *Recipe, auth string) (*Recipe, error) {

	// todo move CheckAuthToken in the middleware
	err := r.rcpDb.CheckAuthToken(auth)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	if v := validateRecipeInput(recipe); len(v) > 0 {
		return nil, errors.NewInvalidParamsErr(v)
	}

	if err = r.rcpDb.CreateRecipe(recipe); err != nil {
		return nil, err
	}

	return recipe, nil
}

func (r *rcp) Update(recipe *Recipe, urlId string, auth string) error {

	// todo move CheckAuthToken in the middleware
	err := r.rcpDb.CheckAuthToken(auth)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	if v := validateRecipeInput(recipe); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}

	err = r.rcpDb.UpdateRecipe(recipe)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *rcp) Delete(recipeId string, auth string) error {

	// todo move CheckAuthToken in the middleware
	err := r.rcpDb.CheckAuthToken(auth)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	if err = r.rcpDb.DeleteRecipe(recipeId); err != nil {
		return err
	}

	return nil
}

func validateRecipeInput(recipe *Recipe) map[string]string {

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
