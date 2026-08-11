package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/warlck/food-flow/api/services/sales/build/all"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/debug"
	"github.com/warlck/food-flow/app/sdk/mux"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/addonbus/stores/addondb"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/categorybus/stores/categorydb"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/menuitembus/stores/menuitemdb"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/orderbus/stores/orderdb"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/promobus/stores/promodb"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus/stores/restaurantdb"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/otel"
	"github.com/warlck/food-flow/foundation/web"
)

var build = "develop"

func main() {

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	// -------------------------------------------------------------------------

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}

		Auth struct {
			Host string `conf:"default:http://auth-service:6000"`
		}

		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:database-service"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}

		Stripe struct {
			SecretKey     string `conf:"mask"`
			WebhookSecret string `conf:"mask"`
		}

		Otel struct {
			TraceEndpoint string `conf:"default:tempo-service:4317"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "SALES",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Initialize OpenTelemetry

	log.Info(ctx, "startup", "status", "initializing OpenTelemetry support", "traceEndpoint", cfg.Otel.TraceEndpoint)

	shutdownOtel, err := otel.Init(ctx, "sales", cfg.Otel.TraceEndpoint)
	if err != nil {
		return fmt.Errorf("initializing otel: %w", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownOtel(ctx); err != nil {
			log.Error(ctx, "shutdown", "status", "otel shutdown error", "err", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer db.Close()

	// -------------------------------------------------------------------------
	// Create Business Packages
	userstore := userdb.NewStore(log, db)
	userBus := userbus.NewBusiness(log, userstore)

	restaurantstore := restaurantdb.NewStore(log, db)
	restaurantBus := restaurantbus.NewBusiness(log, restaurantstore)

	categorystore := categorydb.NewStore(log, db)
	categoryBus := categorybus.NewBusiness(log, categorystore)

	menuitemstore := menuitemdb.NewStore(log, db)
	menuitemBus := menuitembus.NewBusiness(log, menuitemstore)

	addonstore := addondb.NewStore(log, db)
	addonBus := addonbus.NewBusiness(log, addonstore)

	promostore := promodb.NewStore(log, db)
	promoBus := promobus.NewBusiness(log, promostore)

	orderstore := orderdb.NewStore(log, db)
	orderBus := orderbus.NewBusiness(log, orderstore, menuitemBus, restaurantBus, addonBus, promoBus)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	authClient := authclient.New(log, cfg.Auth.Host)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr: cfg.Web.APIHost,
		Handler: mux.WebAPI(mux.Config{
			Build:      cfg.Version.Build,
			Log:        log,
			AuthClient: authClient,
			DB:         db,
			BusConfig: mux.BusConfig{
				UserBus:             userBus,
				RestaurantBus:       restaurantBus,
				CategoryBus:         categoryBus,
				MenuItemBus:         menuitemBus,
				OrderBus:            orderBus,
				AddonBus:            addonBus,
				PromoBus:            promoBus,
				StripeSecretKey:     cfg.Stripe.SecretKey,
				StripeWebhookSecret: cfg.Stripe.WebhookSecret,
			},
		}, all.Routes()),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil

}
