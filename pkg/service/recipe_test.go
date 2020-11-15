package service

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
	"reflect"
	"testing"
)

type recipeDBMock struct {
	getRecipeById func(recipeId string) (*recipe.Recipe, error)
	getAllRecipes func() ([]*recipe.Recipe, error)
	createRecipe  func(recipe *recipe.Recipe) error
	updateRecipe  func(recipe *recipe.Recipe) error
	deleteRecipe  func(recipeId string) error
}

func (rm *recipeDBMock) GetRecipeByID(recipeId string) (*recipe.Recipe, error) {
	if rm.getRecipeById != nil {
		return rm.getRecipeById(recipeId)
	}
	panic("Not implemented")
}

func (rm *recipeDBMock) GetAllRecipes() ([]*recipe.Recipe, error) {
	if rm.getAllRecipes != nil {
		return rm.getAllRecipes()
	}
	panic("Not implemented")
}

func (rm *recipeDBMock) CreateRecipe(recipe *recipe.Recipe) error {
	if rm.createRecipe != nil {
		return rm.createRecipe(recipe)
	}
	panic("Not implemented")
}

func (rm *recipeDBMock) UpdateRecipe(recipe *recipe.Recipe) error {
	if rm.updateRecipe != nil {
		return rm.updateRecipe(recipe)
	}
	panic("Not implemented")
}

func (rm *recipeDBMock) DeleteRecipe(recipeId string) error {
	if rm.deleteRecipe != nil {
		return rm.deleteRecipe(recipeId)
	}
	panic("Not implemented")
}

func TestRcp_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		rcpDB       recipeDBMock
		inputRcpID  string
		expectedRcp *recipe.Recipe
		expectedErr error
	}{
		{
			name: "successful retrieval",
			rcpDB: recipeDBMock{
				getRecipeById: func(recipeId string) (*recipe.Recipe, error) {
					return &recipe.Recipe{
						ID:         "654321",
						Name:       "qwerty",
						PrepTime:   20,
						Difficulty: 3,
						Vegetarian: false,
					}, nil
				},
			},
			inputRcpID: "654321",
			expectedRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			expectedErr: nil,
		},
		{
			name: "error recipe not found",
			rcpDB: recipeDBMock{
				getRecipeById: func(recipeId string) (*recipe.Recipe, error) {
					return nil, errors.NewExistErr(false)
				},
			},
			inputRcpID:  "654321",
			expectedRcp: nil,
			expectedErr: errors.NewExistErr(false),
		},
		{
			name:        "error recipe ID does not match regex",
			inputRcpID:  "0987654321qwerty",
			expectedRcp: nil,
			expectedErr: errors.NewInputError("Invalid ID format", nil),
		},
		{
			name: "error DB issue",
			rcpDB: recipeDBMock{
				getRecipeById: func(recipeId string) (*recipe.Recipe, error) {
					return nil, errors.NewDBErr(fmt.Sprint("error parsing recipe from DB"))
				},
			},
			inputRcpID:  "654321",
			expectedRcp: nil,
			expectedErr: errors.NewDBErr(fmt.Sprint("error parsing recipe from DB")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcpSvr := NewRecipe(&test.rcpDB)
			rcp, err := rcpSvr.GetByID(test.inputRcpID)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
			if rcp != nil {
				if !reflect.DeepEqual(rcp, test.expectedRcp) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedRcp, rcp)
				}
			}
		})
	}
}

func TestRcp_ListAll(t *testing.T) {
	tests := []struct {
		name         string
		rcpDB        recipeDBMock
		expectedRcps []*recipe.Recipe
		expectedErr  error
	}{
		{
			name: "successful retrieval",
			rcpDB: recipeDBMock{
				getAllRecipes: func() ([]*recipe.Recipe, error) {
					return []*recipe.Recipe{
						{
							ID:         "654321",
							Name:       "qwerty",
							PrepTime:   20,
							Difficulty: 3,
							Vegetarian: false,
						}, {
							ID:         "987654",
							Name:       "zxcv",
							PrepTime:   35,
							Difficulty: 3,
							Vegetarian: false,
						},
					}, nil
				},
			},
			expectedRcps: []*recipe.Recipe{
				{
					ID:         "654321",
					Name:       "qwerty",
					PrepTime:   20,
					Difficulty: 3,
					Vegetarian: false,
				}, {
					ID:         "987654",
					Name:       "zxcv",
					PrepTime:   35,
					Difficulty: 3,
					Vegetarian: false,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Empty result - no recipes in DB",
			rcpDB: recipeDBMock{
				getAllRecipes: func() ([]*recipe.Recipe, error) {
					return nil, nil
				},
			},
			expectedRcps: nil,
		},
		{
			name: "error DB issue",
			rcpDB: recipeDBMock{
				getAllRecipes: func() ([]*recipe.Recipe, error) {
					return nil, errors.NewDBErr(fmt.Sprint("error parsing recipe from DB"))
				},
			},
			expectedRcps: nil,
			expectedErr:  errors.NewDBErr(fmt.Sprint("error parsing recipe from DB")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcpSvr := NewRecipe(&test.rcpDB)
			rcps, err := rcpSvr.ListAll()
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
			if rcps != nil {
				if !reflect.DeepEqual(rcps, test.expectedRcps) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedRcps, rcps)
				}
			}
		})
	}
}

func TestRcp_Create(t *testing.T) {
	tests := []struct {
		name        string
		rcpDB       recipeDBMock
		inputRcp    *recipe.Recipe
		expectedErr error
	}{
		{
			name: "successful create",
			rcpDB: recipeDBMock{
				createRecipe: func(recipe *recipe.Recipe) error {
					return nil
				},
			},
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error input validation - Difficulty out of range",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			expectedErr: errors.NewInputError("Invalid input parameters", nil),
		},
		{
			name: "error DB issue - recipe already exists",
			rcpDB: recipeDBMock{
				createRecipe: func(recipe *recipe.Recipe) error {
					return errors.NewExistErr(true)
				},
			},
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			expectedErr: errors.NewExistErr(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcpSvr := NewRecipe(&test.rcpDB)
			err := rcpSvr.Create(test.inputRcp)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}

func TestRcp_Update(t *testing.T) {
	tests := []struct {
		name        string
		ID          string
		rcpDB       recipeDBMock
		inputRcp    *recipe.Recipe
		expectedErr error
	}{
		{
			name: "successful update",
			ID:   "654321",
			rcpDB: recipeDBMock{
				updateRecipe: func(recipe *recipe.Recipe) error {
					return nil
				},
			},
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			expectedErr: nil,
		},
		{
			name: "error input validation - Difficulty out of range",
			ID:   "654321",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			expectedErr: errors.NewInputError("Invalid input parameters", nil),
		},
		{
			name: "error input validation - IDs do not match",
			ID:   "123456",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			expectedErr: errors.NewInputError("ID param and recipe ID do not match", nil),
		},
		{
			name: "error DB issue - recipe does not exist",
			ID:   "654321",
			rcpDB: recipeDBMock{
				updateRecipe: func(recipe *recipe.Recipe) error {
					return errors.NewExistErr(false)
				},
			},
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			expectedErr: errors.NewExistErr(false),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcpSvr := NewRecipe(&test.rcpDB)
			err := rcpSvr.Update(test.ID, test.inputRcp)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}

func TestRcp_Delete(t *testing.T) {
	tests := []struct {
		name        string
		rcpDB       recipeDBMock
		inputRcpID  string
		expectedErr error
	}{
		{
			name: "successful delete",
			rcpDB: recipeDBMock{
				deleteRecipe: func(recipeId string) error {
					return nil
				},
			},
			inputRcpID:  "654321",
			expectedErr: nil,
		},
		{
			name: "error DB issue - DB connection error ",
			rcpDB: recipeDBMock{
				deleteRecipe: func(recipeId string) error {
					return errors.NewDBErr("error DB connection")
				},
			},
			inputRcpID:  "654321",
			expectedErr: errors.NewDBErr("error DB connection"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcpSvr := NewRecipe(&test.rcpDB)
			err := rcpSvr.Delete(test.inputRcpID)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}
