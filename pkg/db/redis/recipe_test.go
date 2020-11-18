package redis

import (
	e "errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
)

type redisAccessorMock struct {
	getAllAccessor func(key string) (map[string]string, error)
	keysAccessor   func(pattern string) ([]string, error)
	existsAccessor func(key string) (int64, error)
	setAccessor    func(key string, fields map[string]interface{}) (string, error)
	setErrAccessor func(key string, fields map[string]interface{}) error
	delAccessor    func(key string) (int64, error)
}

func (rm *redisAccessorMock) getAll(key string) (map[string]string, error) {
	if rm.getAllAccessor != nil {
		return rm.getAllAccessor(key)
	}
	panic("Not implemented")
}

func (rm *redisAccessorMock) keys(pattern string) ([]string, error) {
	if rm.keysAccessor != nil {
		return rm.keysAccessor(pattern)
	}
	panic("Not implemented")
}

func (rm *redisAccessorMock) exists(key string) (int64, error) {
	if rm.existsAccessor != nil {
		return rm.existsAccessor(key)
	}
	panic("Not implemented")
}

func (rm *redisAccessorMock) set(key string, fields map[string]interface{}) (string, error) {
	if rm.setAccessor != nil {
		return rm.setAccessor(key, fields)
	}
	panic("Not implemented")
}

func (rm *redisAccessorMock) setErr(key string, fields map[string]interface{}) error {
	if rm.setErrAccessor != nil {
		return rm.setErrAccessor(key, fields)
	}
	panic("Not implemented")
}

func (rm *redisAccessorMock) del(key string) (int64, error) {
	if rm.delAccessor != nil {
		return rm.delAccessor(key)
	}
	panic("Not implemented")
}

func TestProxy_GetRecipeByID(t *testing.T) {
	tests := []struct {
		name        string
		ID          string
		accessor    *redisAccessorMock
		expectedRcp *recipe.Recipe
		expectedErr error
	}{
		{
			name: "successful retrieve",
			ID:   "654321",
			accessor: &redisAccessorMock{
				getAllAccessor: func(key string) (map[string]string, error) {
					rcp := map[string]string{rcpID: "654321", name: "qwerty", prepTime: "20", difficulty: "3", vegetarian: "False"}
					return rcp, nil
				},
			},
			expectedRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error - get all DB call",
			ID:   "654321",
			accessor: &redisAccessorMock{
				getAllAccessor: func(key string) (map[string]string, error) {
					return nil, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - recipe does not exists",
			ID:   "654321",
			accessor: &redisAccessorMock{
				getAllAccessor: func(key string) (map[string]string, error) {
					return nil, nil
				},
			},
			expectedErr: errors.NewExistErr(false),
		},
		{
			name: "error - mapping from redis to recipe struct",
			ID:   "654321",
			accessor: &redisAccessorMock{
				getAllAccessor: func(key string) (map[string]string, error) {
					rcp := map[string]string{rcpID: "654321", name: "qwerty", prepTime: "isNotAnInt", difficulty: "3", vegetarian: "False"}
					return rcp, nil
				},
			},
			expectedErr: errors.NewDBErr(fmt.Sprintf("error parsing recipe from redis: %s", "strconv.Atoi: parsing \"isNotAnInt\": invalid syntax")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			rcp, err := proxy.GetRecipeByID(test.ID)
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

func TestProxy_GetAllRecipes(t *testing.T) {
	// getAll could be called twice, this way we could mock both calls with different results.
	var nCalls = 2
	tests := []struct {
		name         string
		accessor     *redisAccessorMock
		expectedRcps []*recipe.Recipe
		expectedErr  error
	}{
		{
			name: "successful retrieve - empty",
			accessor: &redisAccessorMock{
				keysAccessor: func(pattern string) ([]string, error) {
					return make([]string, 0), nil
				},
			},
			expectedRcps: nil,
			expectedErr:  nil,
		},
		{
			name: "successful retrieve - multiple results",
			accessor: &redisAccessorMock{
				keysAccessor: func(pattern string) ([]string, error) {
					rcpKeys := []string{"654321", "98765"}
					return rcpKeys, nil
				},
				getAllAccessor: func(key string) (map[string]string, error) {
					nCalls--
					if nCalls == 1 {
						return map[string]string{rcpID: "654321", name: "qwerty", prepTime: "20", difficulty: "3", vegetarian: "False"}, nil
					}
					return map[string]string{rcpID: "98765", name: "zxcvb", prepTime: "60", difficulty: "4", vegetarian: "False"}, nil
				},
			},
			expectedRcps: []*recipe.Recipe{
				{
					ID:         "654321",
					Name:       "qwerty",
					PrepTime:   20,
					Difficulty: 3,
					Vegetarian: false,
				},
				{
					ID:         "98765",
					Name:       "zxcvb",
					PrepTime:   60,
					Difficulty: 4,
					Vegetarian: false,
				},
			},
		},
		{
			name: "error - DB retrieving keys",
			accessor: &redisAccessorMock{
				keysAccessor: func(pattern string) ([]string, error) {
					return nil, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - get all DB call",
			accessor: &redisAccessorMock{
				keysAccessor: func(pattern string) ([]string, error) {
					rcpKeys := []string{"654321", "98765"}
					return rcpKeys, nil
				},
				getAllAccessor: func(key string) (map[string]string, error) {
					return nil, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - mapping from redis to recipe struct",
			accessor: &redisAccessorMock{
				keysAccessor: func(pattern string) ([]string, error) {
					rcpKeys := []string{"654321", "98765"}
					return rcpKeys, nil
				},
				getAllAccessor: func(key string) (map[string]string, error) {
					return map[string]string{rcpID: "654321", name: "qwerty", prepTime: "20", difficulty: "isNotAnInt", vegetarian: "False"}, nil
				},
			},
			expectedErr: errors.NewDBErr(fmt.Sprintf("error parsing recipe from redis: %s", "strconv.Atoi: parsing \"isNotAnInt\": invalid syntax")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			rcps, err := proxy.GetAllRecipes()
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
			if rcps != nil {
				if len(rcps) != 0 && len(test.expectedRcps) != 0 && !reflect.DeepEqual(rcps, test.expectedRcps) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedRcps, rcps)
				}
			}
		})
	}
}

func TestProxy_CreateRecipe(t *testing.T) {
	tests := []struct {
		name        string
		inputRcp    *recipe.Recipe
		accessor    *redisAccessorMock
		expectedErr error
	}{
		{
			name: "successful create",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, nil
				},
				setAccessor: func(key string, fields map[string]interface{}) (string, error) {
					return "", nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "error - DB exist call",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - recipe already exists",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
			},
			expectedErr: errors.NewExistErr(true),
		},
		{
			name: "error - DB inserting recipe",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, nil
				},
				setAccessor: func(key string, fields map[string]interface{}) (string, error) {
					return "", e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			err := proxy.CreateRecipe(test.inputRcp)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}

func TestProxy_UpdateRecipe(t *testing.T) {
	tests := []struct {
		name        string
		inputRcp    *recipe.Recipe
		accessor    *redisAccessorMock
		expectedErr error
	}{
		{
			name: "successful update",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				setErrAccessor: func(key string, fields map[string]interface{}) error {
					return nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "successful update",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				setErrAccessor: func(key string, fields map[string]interface{}) error {
					return nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "error - DB exist call",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - recipe does not exist",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, nil
				},
			},
			expectedErr: errors.NewExistErr(false),
		},
		{
			name: "error - DB setting new values",
			inputRcp: &recipe.Recipe{
				ID:         "654321",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				setErrAccessor: func(key string, fields map[string]interface{}) error {
					return e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			err := proxy.UpdateRecipe(test.inputRcp)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}

func TestProxy_DeleteRecipe(t *testing.T) {
	// del to redis could be called twice, this way we could mock both calls with different results.
	var nCalls = 2
	tests := []struct {
		name        string
		ID          string
		accessor    *redisAccessorMock
		expectedErr error
		delCalls    int
	}{
		{
			name: "successful delete",
			ID:   "654321",
			accessor: &redisAccessorMock{
				delAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "error - DB exist call",
			ID:   "654321",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - DB del call",
			ID:   "654321",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				delAccessor: func(key string) (int64, error) {
					return 0, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
		{
			name: "error - recipe does not exist",
			ID:   "654321",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				delAccessor: func(key string) (int64, error) {
					return 0, nil
				},
			},
			expectedErr: errors.NewExistErr(false),
		},
		{
			name: "error - DB deleting recipe",
			ID:   "654321",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				delAccessor: func(key string) (int64, error) {
					nCalls--
					if nCalls == 1 {
						return 1, nil
					}
					return 0, e.New("DB issue")
				},
			},
			expectedErr: errors.NewDBErr("DB issue"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			err := proxy.DeleteRecipe(test.ID)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}
