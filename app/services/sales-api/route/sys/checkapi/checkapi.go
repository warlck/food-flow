package checkapi

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"runtime"

	"github.com/warlck/food-flow/business/web/response"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

type api struct {
	build string
	log   *logger.Logger
}

func newAPI(build string, log *logger.Logger) *api {
	return &api{
		build: build,
		log:   log,
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
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	// TODO: LOG WHEN NOT OK ONLY to avoid spamming the logs
	// a.log.Info(ctx, "readiness", "status", status)

	return web.Respond(ctx, w, status, http.StatusOK)
}

func (a *api) testError(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return response.NewError(errors.New("this message is trusted"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
