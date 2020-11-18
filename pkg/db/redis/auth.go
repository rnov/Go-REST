package redis

import "github.com/rnov/Go-REST/pkg/errors"

const tokenPattern = "TOKEN_"

func (p *Proxy) CheckAuth(auth string) error {
	//exist, err := p.main.Exists(tokenPattern + auth).Result()
	exist, err := p.exists(tokenPattern + auth)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exist == 0 {
		return errors.NewFailedAuthErr()
	}

	return nil
}
