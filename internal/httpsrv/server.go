package httpsrv

import (
	"go.infratographer.com/x/ginx"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"
)

func NewServer(logger *zap.SugaredLogger, cfg ginx.Config) *ginx.Server {
	// TODO - Add a storage/controller implementation here.
	router := NewRouter(nil)
	server := ginx.NewServer(logger.Desugar(), cfg, versionx.BuildDetails())
	server = server.AddHandler(router)

	return &server
}
