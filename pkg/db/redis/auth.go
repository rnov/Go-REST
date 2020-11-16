package redis

import "github.com/rnov/Go-REST/pkg/errors"

const tokenPattern = "TOKEN_"

func (rProxy *Proxy) CheckAuth(auth string) error {
	exist, err := rProxy.master.Exists(tokenPattern + auth).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		return errors.NewFailedAuthErr()
	}

	return nil
}
