package api

import "github.com/imroc/req/v3"

var ReqClint *req.Client

func init() {
	ReqClint = req.DefaultClient()
}
