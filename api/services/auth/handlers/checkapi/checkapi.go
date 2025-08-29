package checkapi

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

type api struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func newAPI(build string, log *logger.Logger, db *sqlx.DB) *api {
	return &api{
		build: build,
		log:   log,
		db:    db,
	}
}

func (a *api) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := Info{
		Status:     "up",
		Build:      a.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// TODO: LOG WHEN NOT OK ONLY to avoid spamming the logs
	// a.log.Info(ctx, "liveness", "info", info, "status", "ok")

	return web.Respond(ctx, w, info, http.StatusOK)
}

func (a *api) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	// if err := sqldb.StatusCheck(ctx, a.db); err != nil {
	// 	status = "db not ready"
	// 	statusCode = http.StatusInternalServerError
	// 	a.log.Info(ctx, "readiness failure", "status", status)
	// }

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	a.log.Info(ctx, "readiness", "status", status)

	return web.Respond(ctx, w, data, statusCode)
}
