// Copyright 2020 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	promgrpc "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pelletier/go-toml/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	"github.com/scionproto/scion/daemon"
	"github.com/scionproto/scion/daemon/config"
	sd_drkey "github.com/scionproto/scion/daemon/drkey"
	sd_grpc "github.com/scionproto/scion/daemon/drkey/grpc"
	"github.com/scionproto/scion/daemon/fetcher"
	api "github.com/scionproto/scion/daemon/mgmtapi"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/experimental/hiddenpath"
	hpgrpc "github.com/scionproto/scion/pkg/experimental/hiddenpath/grpc"
	libgrpc "github.com/scionproto/scion/pkg/grpc"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/metrics"
	"github.com/scionproto/scion/pkg/private/prom"
	"github.com/scionproto/scion/pkg/private/serrors"
	cryptopb "github.com/scionproto/scion/pkg/proto/crypto"
	sdpb "github.com/scionproto/scion/pkg/proto/daemon"
	"github.com/scionproto/scion/pkg/scrypto/cppki"
	"github.com/scionproto/scion/pkg/scrypto/signed"
	"github.com/scionproto/scion/private/app"
	cppkiapi "github.com/scionproto/scion/private/mgmtapi/cppki/api"
	segapi "github.com/scionproto/scion/private/mgmtapi/segments/api"
	"github.com/scionproto/scion/private/pathdb"
	"github.com/scionproto/scion/private/periodic"
	"github.com/scionproto/scion/private/revcache"
	"github.com/scionproto/scion/private/segment/segfetcher"
	segfetchergrpc "github.com/scionproto/scion/private/segment/segfetcher/grpc"
	infra "github.com/scionproto/scion/private/segment/verifier"
	"github.com/scionproto/scion/private/service"
	"github.com/scionproto/scion/private/storage"
	"github.com/scionproto/scion/private/storage/drkey/level2"
	pathstoragemetrics "github.com/scionproto/scion/private/storage/path/metrics"
	truststoragemetrics "github.com/scionproto/scion/private/storage/trust/metrics"
	"github.com/scionproto/scion/private/topology"
	"github.com/scionproto/scion/private/trust"
	"github.com/scionproto/scion/private/trust/compat"
	trustmetrics "github.com/scionproto/scion/private/trust/metrics"
)

var daemonConfig config.Config

//func main() {
//	application := launcher.Application{
//		TOMLConfig: &daemonConfig,
//		ShortName:  "SCION Daemon",
//		Main:       realMain,
//	}
//	application.Run()
//}

func runDaemon(ctx context.Context, configFile string) error {

	/*if err := libconfig.LoadFile(configFile, a.TOMLConfig); err != nil {
		return serrors.WrapStr("loading config from file", err,
			"file", a.config.GetString(cfgConfigFile))
	}
	a.TOMLConfig.InitDefaults()*/

	// Load the configuration from the file into daemonConfig.
	// Load configFIle into []byte
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return serrors.WrapStr("reading config file", err)
	}

	err = toml.Unmarshal(configFileBytes, &daemonConfig)
	if err != nil {
		return serrors.WrapStr("unmarshalling config file", err)
	}

	topo, err := topology.NewLoader(topology.LoaderCfg{
		File:      daemonConfig.General.Topology(),
		Reload:    app.SIGHUPChannel(ctx),
		Validator: &topology.DefaultValidator{},
		Metrics:   loaderMetrics(),
	})
	if err != nil {
		return serrors.WrapStr("creating topology loader", err)
	}
	g, errCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer log.HandlePanic()
		return topo.Run(errCtx)
	})

	closer, err := daemon.InitTracer(daemonConfig.Tracing, daemonConfig.General.ID)
	if err != nil {
		return serrors.WrapStr("initializing tracer", err)
	}
	defer closer.Close()

	revCache := storage.NewRevocationStorage()
	pathDB, err := storage.NewPathStorage(daemonConfig.PathDB)
	if err != nil {
		return serrors.WrapStr("initializing path storage", err)
	}
	pathDB = pathstoragemetrics.WrapDB(pathDB, pathstoragemetrics.Config{
		Driver: string(storage.BackendSqlite),
	})
	defer pathDB.Close()
	defer revCache.Close()
	cleaner := periodic.Start(pathdb.NewCleaner(pathDB, "sd_segments"),
		300*time.Second, 295*time.Second)
	defer cleaner.Stop()
	rcCleaner := periodic.Start(revcache.NewCleaner(revCache, "sd_revocation"),
		10*time.Second, 10*time.Second)
	defer rcCleaner.Stop()

	dialer := &libgrpc.TCPDialer{
		SvcResolver: func(dst addr.SVC) []resolver.Address {
			if base := dst.Base(); base != addr.SvcCS {
				panic("Unsupported address type, implementation error?")
			}
			targets := []resolver.Address{}
			for _, entry := range topo.ControlServiceAddresses() {
				targets = append(targets, resolver.Address{Addr: entry.String()})
			}
			return targets
		},
	}

	trustDB, err := storage.NewTrustStorage(daemonConfig.TrustDB)
	if err != nil {
		return serrors.WrapStr("initializing trust database", err)
	}
	defer trustDB.Close()
	trustDB = truststoragemetrics.WrapDB(trustDB, truststoragemetrics.Config{
		Driver: string(storage.BackendSqlite),
		QueriesTotal: metrics.NewPromCounterFrom(
			prometheus.CounterOpts{
				Name: "trustengine_db_queries_total",
				Help: "Total queries to the database",
			},
			[]string{"driver", "operation", prom.LabelResult},
		),
	})
	engine, err := daemon.TrustEngine(daemonConfig.General.ConfigDir, topo.IA(), trustDB, dialer)
	if err != nil {
		return serrors.WrapStr("creating trust engine", err)
	}
	engine.Inspector = trust.CachingInspector{
		Inspector:          engine.Inspector,
		Cache:              daemonConfig.TrustEngine.Cache.New(),
		CacheHits:          metrics.NewPromCounter(trustmetrics.CacheHitsTotal),
		MaxCacheExpiration: daemonConfig.TrustEngine.Cache.Expiration.Duration,
	}
	trcLoader := periodic.Start(periodic.Func{
		Task: func(ctx context.Context) {
			trcDirs := filepath.Join(daemonConfig.General.ConfigDir, "certs")
			res, err := trust.LoadTRCs(ctx, trcDirs, trustDB)
			if err != nil {
				log.SafeInfo(log.FromCtx(ctx), "TRC loading failed", "err", err)
			}
			if len(res.Loaded) > 0 {
				log.SafeInfo(log.FromCtx(ctx), "Loaded TRCs from disk", "trcs", res.Loaded)
			}
		},
		TaskName: "daemon_trc_loader",
	}, 10*time.Second, 10*time.Second)
	defer trcLoader.Stop()

	var drkeyClientEngine *sd_drkey.ClientEngine
	if daemonConfig.DRKeyLevel2DB.Connection != "" {
		backend, err := storage.NewDRKeyLevel2Storage(daemonConfig.DRKeyLevel2DB)
		if err != nil {
			return serrors.WrapStr("creating level2 DRKey DB", err)
		}
		counter := metrics.NewPromCounter(
			promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "drkey_level2db_queries_total",
					Help: "Total queries to the database",
				},
				[]string{"operation", prom.LabelResult},
			),
		)
		level2DB := &level2.Database{
			Backend: backend,
			Metrics: &level2.Metrics{
				QueriesTotal: func(op, label string) metrics.Counter {
					return metrics.CounterWith(
						counter,
						"operation", op,
						prom.LabelResult, label,
					)
				},
			},
		}
		defer level2DB.Close()

		drkeyFetcher := &sd_grpc.Fetcher{
			Dialer: dialer,
		}
		drkeyClientEngine = &sd_drkey.ClientEngine{
			IA:      topo.IA(),
			DB:      level2DB,
			Fetcher: drkeyFetcher,
		}
		cleaners := drkeyClientEngine.CreateStorageCleaners()
		for _, cleaner := range cleaners {
			cleaner_task := periodic.Start(cleaner,
				5*time.Minute, 5*time.Minute)
			defer cleaner_task.Stop()
		}
	}

	listen := daemon.APIAddress(daemonConfig.SD.Address)
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return serrors.WrapStr("listening", err)
	}

	hpGroups, err := hiddenpath.LoadHiddenPathGroups(daemonConfig.SD.HiddenPathGroups)
	if err != nil {
		return serrors.WrapStr("loading hidden path groups", err)
	}
	var requester segfetcher.RPC = &segfetchergrpc.Requester{
		Dialer: dialer,
	}
	if len(hpGroups) > 0 {
		requester = &hpgrpc.Requester{
			RegularLookup: requester,
			HPGroups:      hpGroups,
			Dialer:        dialer,
		}
	}

	createVerifier := func() infra.Verifier {
		if daemonConfig.SD.DisableSegVerification {
			return acceptAllVerifier{}
		}
		return compat.Verifier{Verifier: trust.Verifier{
			Engine:             engine,
			Cache:              daemonConfig.TrustEngine.Cache.New(),
			CacheHits:          metrics.NewPromCounter(trustmetrics.CacheHitsTotal),
			MaxCacheExpiration: daemonConfig.TrustEngine.Cache.Expiration.Duration,
		}}
	}

	server := grpc.NewServer(
		libgrpc.UnaryServerInterceptor(),
		libgrpc.DefaultMaxConcurrentStreams(),
	)
	sdpb.RegisterDaemonServiceServer(server, daemon.NewServer(
		daemon.ServerConfig{
			IA:       topo.IA(),
			MTU:      topo.MTU(),
			Topology: topo,
			Fetcher: fetcher.NewFetcher(
				fetcher.FetcherConfig{
					IA:         topo.IA(),
					MTU:        topo.MTU(),
					Core:       topo.Core(),
					NextHopper: topo,
					RPC:        requester,
					PathDB:     pathDB,
					Inspector:  engine,
					Verifier:   createVerifier(),
					RevCache:   revCache,
					Cfg:        daemonConfig.SD,
				},
			),
			Engine:      engine,
			RevCache:    revCache,
			DRKeyClient: drkeyClientEngine,
		},
	))

	promgrpc.Register(server)

	var cleanup app.Cleanup
	g.Go(func() error {
		defer log.HandlePanic()
		if err := server.Serve(listener); err != nil {
			return serrors.WrapStr("serving gRPC API", err, "addr", listen)
		}
		return nil
	})
	cleanup.Add(func() error { server.GracefulStop(); return nil })

	if daemonConfig.API.Addr != "" {
		r := chi.NewRouter()
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
		}))
		r.Get("/", api.ServeSpecInteractive)
		r.Get("/openapi.json", api.ServeSpecJSON)
		server := api.Server{
			SegmentsServer: segapi.Server{
				Segments: pathDB,
			},
			CPPKIServer: cppkiapi.Server{
				TrustDB: trustDB,
			},
			Config:   service.NewConfigStatusPage(daemonConfig).Handler,
			Info:     service.NewInfoStatusPage().Handler,
			LogLevel: service.NewLogLevelStatusPage().Handler,
		}
		log.Info("Exposing API", "addr", daemonConfig.API.Addr)
		h := api.HandlerFromMuxWithBaseURL(&server, r, "/api/v1")
		mgmtServer := &http.Server{
			Addr:    daemonConfig.API.Addr,
			Handler: h,
		}
		g.Go(func() error {
			defer log.HandlePanic()
			err := mgmtServer.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return serrors.WrapStr("serving service management API", err)
			}
			return nil
		})
		cleanup.Add(mgmtServer.Close)
	}

	// Start HTTP endpoints.
	statusPages := service.StatusPages{
		"info":      service.NewInfoStatusPage(),
		"config":    service.NewConfigStatusPage(daemonConfig),
		"log/level": service.NewLogLevelStatusPage(),
		"topology":  service.NewTopologyStatusPage(topo),
	}
	if err := statusPages.Register(http.DefaultServeMux, daemonConfig.General.ID); err != nil {
		return serrors.WrapStr("registering status pages", err)
	}

	g.Go(func() error {
		defer log.HandlePanic()
		return daemonConfig.Metrics.ServePrometheus(errCtx)
	})

	g.Go(func() error {
		defer log.HandlePanic()
		<-errCtx.Done()
		return cleanup.Do()
	})

	return g.Wait()
}

type acceptAllVerifier struct{}

func (acceptAllVerifier) Verify(ctx context.Context, signedMsg *cryptopb.SignedMessage,
	associatedData ...[]byte) (*signed.Message, error) {

	return nil, nil
}

func (v acceptAllVerifier) WithValidity(cppki.Validity) infra.Verifier {
	return v
}

func (v acceptAllVerifier) WithServer(net.Addr) infra.Verifier {
	return v
}

func (v acceptAllVerifier) WithIA(addr.IA) infra.Verifier {
	return v
}

func loaderMetrics() topology.LoaderMetrics {
	updates := prom.NewCounterVec("", "",
		"topology_updates_total",
		"The total number of updates.",
		[]string{prom.LabelResult},
	)
	return topology.LoaderMetrics{
		ValidationErrors: metrics.NewPromCounter(updates).With(prom.LabelResult, "err_validate"),
		ReadErrors:       metrics.NewPromCounter(updates).With(prom.LabelResult, "err_read"),
		LastUpdate: metrics.NewPromGauge(
			prom.NewGaugeVec("", "",
				"topology_last_update_time",
				"Timestamp of the last successful update.",
				[]string{},
			),
		),
		Updates: metrics.NewPromCounter(updates).With(prom.LabelResult, prom.Success),
	}
}
