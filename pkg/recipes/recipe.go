package recipes

import (
	"errors"

	"github.com/goRest/pkg/db"
	e "github.com/goRest/pkg/errors"
	log "github.com/goRest/pkg/logger"
)

type Recipe struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PrepTime   int    `json:"prepTime"`
	Difficulty int    `json:"difficulty"`
	Vegetarian bool   `json:"vegetarian"`
}

// this is a must, struct can not implement interface from different package.
type RecipeService struct {
	rcpDb  db.RecipesDbCalls
	logger log.Loggers
	// add more func fields
}

func NewRecipeSrv(rcpDb db.RecipesDbCalls, logger log.Loggers) *RecipeService {
	recipeSrv := &RecipeService{
		rcpDb:  rcpDb,
		logger: logger,
	}
	return recipeSrv
}

func (rcp *RecipeService) GetById(recipeId string) (*Recipe, error) {

	recipe, err := rcp.rcpDb.GetRecipeById(recipeId)

	if recipe == nil {
		rcp.logger.Notice(err)
		return nil, &e.NotFoundErr{}
	}

	if err != nil {
		rcp.logger.Error(err)
		return nil, &e.DBErr{}
	}

	return recipe, nil

}

func (rcp *RecipeService) ListAll() ([]*Recipe, error) {

	recipes, err := rcp.rcpDb.GetAllRecipes()

	if err != nil {
		rcp.logger.Error(err)
		return nil, &e.DBErr{}
	}

	return recipes, nil
}

func (rcp *RecipeService) Create(recipe *Recipe, auth string) (map[string]string, error) {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return nil, err
	}

	valid := validateRecipeDataRange(recipe)
	if len(valid) > 0 {
		return valid, &e.InvalidParamsErr{}
	}

	err = rcp.rcpDb.CreateRecipe(recipe)

	if err != nil {

		if errors.Is(err, &e.DBErr{}) {
			rcp.logger.Error(err)
		}

		return nil, err
	}

	return nil, nil
}

func (rcp *RecipeService) Update(recipe *Recipe, urlId string, auth string) (map[string]string, error) {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return nil, err
	}

	valid := validateRecipeDataRange(recipe)
	if len(valid) == 0 && urlId != recipe.Id {
		return valid, &e.InvalidParamsErr{}
	}

	err = rcp.rcpDb.UpdateRecipe(recipe)
	if err != nil {
		if errors.Is(err, &e.DBErr{}) {
			rcp.logger.Error(err)
		}
		return nil, err
	}

	return nil, nil
}

func (rcp *RecipeService) Delete(recipeId string, auth string) error {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return err
	}

	err = rcp.rcpDb.DeleteRecipe(recipeId)
	if err != nil {
		if errors.Is(err, &e.DBErr{}) {
			rcp.logger.Error(err)
		}
		return err
	}
	return nil
}

// fixme : error msg
func validateRecipeDataRange(recipe *Recipe) map[string]string {

	valid := make(map[string]string)

	if recipe.Difficulty <= 1 || recipe.Difficulty > 3 {
		valid[e.Difficulty] = e.OutOfRange
	}
	if len(recipe.Name) > 100 {
		valid[e.Name] = e.TooLong
	}
	if recipe.PrepTime <= 1 || recipe.PrepTime > 1000 {
		valid[e.Preptime] = e.OutOfRange
	}
	return valid
}
