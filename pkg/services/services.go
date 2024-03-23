package services

import legendasdivx "github.com/dronept/go-stremio-legendasdivx/pkg/services/legendas_divx"

type Services struct {
	LegendasDivx *legendasdivx.LegendasDivx
}

func NewServices() *Services {
	return &Services{
		LegendasDivx: legendasdivx.NewLegendasDivx(),
	}
}
