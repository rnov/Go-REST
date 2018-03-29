package controller

import (
	"Go-REST/application/common"
	"Go-REST/application/dbInterface"
	"Go-REST/application/model"
	"errors"
)

// this is a must, struct can not implement interface from different package.
type RecipeController struct {
	rcpDb  dbInterface.RecipesDbCalls
	logger common.LogInterface
	// add more func fields
}

func NewRecipeController(rcpDb dbInterface.RecipesDbCalls, logger common.LogInterface) *RecipeController {
	recipeController := &RecipeController{
		rcpDb:  rcpDb,
		logger: logger,
	}
	return recipeController
}

func (rcp *RecipeController) GetById(recipeId string) (*model.Recipe, error) {

	recipe, err := rcp.rcpDb.GetRecipeById(recipeId)

	if recipe == nil {
		rcp.logger.Notice(err)
		return nil, errors.New(common.NotFound)
	}

	if err != nil {
		rcp.logger.Error(err)
		return nil, errors.New(common.DbError)
	}

	return recipe, nil

}

func (rcp *RecipeController) ListAll() ([]*model.Recipe, error) {

	recipes, err := rcp.rcpDb.GetAllRecipes()

	if err != nil {
		rcp.logger.Error(err)
		return nil, errors.New(common.DbError)
	}

	return recipes, nil
}

func (rcp *RecipeController) Create(recipe *model.Recipe, auth string) (map[string]string, error) {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return nil, err
	}

	valid := validateRecipeDataRange(recipe)
	if len(valid) > 0 {
		return valid, errors.New(common.InvalidParams)
	}

	err = rcp.rcpDb.CreateRecipe(recipe)

	if err != nil {
		if err.Error() == common.DbError {
			rcp.logger.Error(err)
		}
		return nil, err
	}

	return nil, nil
}

func (rcp *RecipeController) Update(recipe *model.Recipe, urlId string, auth string) (map[string]string, error) {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return nil, err
	}

	valid := validateRecipeDataRange(recipe)
	if len(valid) == 0 && urlId != recipe.Id {
		return valid, errors.New(common.InvalidParams)
	}

	err = rcp.rcpDb.UpdateRecipe(recipe)
	if err != nil {
		if err.Error() == common.DbError {
			rcp.logger.Error(err)
		}
		return nil, err
	}

	return nil, nil
}

func (rcp *RecipeController) Delete(recipeId string, auth string) error {

	err := rcp.rcpDb.CheckAuthToken(auth)
	if err != nil {
		rcp.logger.Error(err)
		return err
	}

	err = rcp.rcpDb.DeleteRecipe(recipeId)
	if err != nil {
		if err.Error() == common.DbError {
			rcp.logger.Error(err)
		}
		return err
	}
	return nil
}

func validateRecipeDataRange(recipe *model.Recipe) map[string]string {

	valid := make(map[string]string)

	if recipe.Difficulty <= 1 || recipe.Difficulty > 3 {
		valid[common.Difficulty] = common.OutOfRange
	}
	if len(recipe.Name) > 100 {
		valid[common.Name] = common.TooLong
	}
	if recipe.PrepTime <= 1 || recipe.PrepTime > 1000 {
		valid[common.Preptime] = common.OutOfRange
	}
	return valid
}
