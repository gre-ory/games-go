package main

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	yaml "gopkg.in/yaml.v2"

	"github.com/gre-ory/games-go/internal/util/list"

	share_api "github.com/gre-ory/games-go/internal/game/share/api"

	ttt_api "github.com/gre-ory/games-go/internal/game/tictactoe/api"
	ttt_service "github.com/gre-ory/games-go/internal/game/tictactoe/service"
	ttt_store "github.com/gre-ory/games-go/internal/game/tictactoe/store"

	skj_api "github.com/gre-ory/games-go/internal/game/skyjo/api"
	skj_service "github.com/gre-ory/games-go/internal/game/skyjo/service"
	skj_store "github.com/gre-ory/games-go/internal/game/skyjo/store"
)

// //////////////////////////////////////////////////
// main

func main() {

	// exit process immediately upon sigterm
	handleSigTerms()

	//
	// context
	//

	ctx := context.Background()

	//
	// random
	//

	rand.Seed(time.Now().UnixNano())

	//
	// config
	//

	config := readConfig()
	secret := readSecrets()

	// logger
	logger := NewLogger(config.Log)
	logger.Info("")
	logger.Info(" -------------------------------------------------- ")
	logger.Info("")
	logger.Info("starting app...", zap.String("env", config.Env), zap.String("app", config.App), zap.String("version", config.Version), zap.Any("config", config))

	//
	// store
	//

	ttt_gameStore := ttt_store.NewGameStore()
	ttt_playerStore := ttt_store.NewPlayerStore()
	skj_gameStore := skj_store.NewGameStore()
	skj_playerStore := skj_store.NewPlayerStore()

	//
	// service
	//

	ttt_service := ttt_service.NewGameService(logger, ttt_gameStore, ttt_playerStore)
	skj_service := skj_service.NewGameService(logger, skj_gameStore, skj_playerStore)

	//
	// api
	//

	cookie_server := share_api.NewCookieServer(logger, config.Cookie.Key, config.Cookie.MaxAge, secret.CookieSecret)
	ttt_server := ttt_api.NewGameServer(logger, cookie_server, ttt_service)
	skj_server := skj_api.NewGameServer(logger, cookie_server, skj_service)

	//
	// router
	//

	router := httprouter.New()
	cookie_server.RegisterRoutes(router)
	ttt_server.RegisterRoutes(router)
	skj_server.RegisterRoutes(router)
	router.NotFound = http.FileServer(http.FS(staticFS))

	//
	// server
	//

	server := http.Server{
		Addr: config.Server.Address,
		Handler: AllowCORS(logger, config.Server.WhiteListOrigins)(
			WithRequestLogging(logger)(
				router,
			),
		),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	logger.Info(fmt.Sprintf("starting backend server on %s", server.Addr))
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal("backend server failed", zap.Error(err))
	}
}

// //////////////////////////////////////////////////
// logger

func NewLogger(config LogConfig) *zap.Logger {

	fmt.Printf("config.Env: [%s] \n", config.Env)
	var cfg zap.Config
	switch config.Env {
	case "prd":
		cfg = zap.NewProductionConfig()
	case "dev":
		cfg = zap.NewDevelopmentConfig()
		// cfg.Development = false
	default:
	}

	logger, err := cfg.Build()
	if err == nil {
		logger.Info("logger has been initialized", zap.String("env", config.Env))
		return logger
	} else {
		fmt.Printf("error: %s \n", err.Error())
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.File,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})

	fmt.Printf("config.Encoder: [%s] \n", config.Encoder)
	var zapEncoderCfg zapcore.EncoderConfig
	switch config.Encoder {
	case "prd":
		zapEncoderCfg = zap.NewProductionEncoderConfig()
	case "dev":
		zapEncoderCfg = zap.NewDevelopmentEncoderConfig()
	default:
		fmt.Printf("invalid LOG_CONFIG >>> FALLBACK to 'dev'!\n")
		zapEncoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	fmt.Printf("config.Level: [%s] \n", config.Level)
	var zapLevel zapcore.Level
	switch config.Level {
	case "err":
		zapLevel = zap.ErrorLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "debug":
		zapLevel = zap.DebugLevel
	default:
		fmt.Printf("invalid LOG_LEVEL >>> FALLBACK to 'info'!\n")
		zapLevel = zap.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapEncoderCfg),
		writer,
		zapLevel,
	)

	logger = zap.New(core)
	logger.Info("logger has been initialized", zap.String("encoder", config.Encoder), zap.String("level", config.Level))
	return logger
}

// //////////////////////////////////////////////////
// request logging

const (
	DebugStaticResource = false
	DebugApiCall        = true
)

func WithRequestLogging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static/") {
				if DebugStaticResource {
					logger.Info(fmt.Sprintf("[%s] %s %s", r.Method, r.URL.Path, r.UserAgent()))
				}
			} else {
				if DebugApiCall {
					now := time.Now()
					logger.Info(fmt.Sprintf("[%s] ------------------------- %s ------------------------- %s", r.Method, r.URL.Path, r.UserAgent()))
					defer func() {
						logger.Info(fmt.Sprintf("[%s] ------------------------- %s ( %s ) ------------------------- %s", r.Method, r.URL.Path, time.Since(now), r.UserAgent()))
					}()
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// //////////////////////////////////////////////////
// cors

// var whitelistOrigins []string = []string{
// 	"http://localhost:3000",
// 	"http://localhost:9090",
// 	"http://158.178.206.68:8080",
// 	"http://158.178.206.68:8081",
// 	"http://158.178.206.68:8082",
// }

func AllowCORS(logger *zap.Logger, whitelistOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := list.Contains(whitelistOrigins, origin)
			if allowed {
				// logger.Info(fmt.Sprintf("[COR] OK - Origin: %s", origin))
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Headers", "content-type,authorization")
			} else {
				logger.Info(fmt.Sprintf("[COR] BLOCKED - Origin: %s", origin))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// //////////////////////////////////////////////////
// static

var (
	//go:embed static/*
	staticFS embed.FS
	//go:embed *
	rootFS embed.FS
)

func registerStaticFiles(router *httprouter.Router) {
	// router.ServeFiles("/static/*filepath", http.FS(rootFS))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	// router.Handle("/static/", http.FileServer(http.FS(staticFS)))
}

// //////////////////////////////////////////////////
// cors

// cors := func(h http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// in development, the Origin is the the Hugo server, i.e. http://localhost:1313
// 		// but in production, it is the domain name where one's site is deployed
// 		//
// 		// CHANGE THIS: You likely do not want to allow any origin (*) in production. The value should be the base URL of
// 		// where your static content is served
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, hx-target, hx-current-url, hx-request")
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		h.ServeHTTP(w, r)
// 	}
// }

// //////////////////////////////////////////////////
// sigterms

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}

// //////////////////////////////////////////////////
// config

type Config struct {
	Env     string       `yaml:"env"`
	App     string       `yaml:"app"`
	Version string       `yaml:"version"`
	Log     LogConfig    `yaml:"log"`
	Cookie  CookieConfig `yaml:"cookie"`
	Server  ServerConfig `yaml:"server"`
}

type LogConfig struct {
	Env     string `yaml:"env"`
	Encoder string `yaml:"encoder"`
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
}

type CookieConfig struct {
	Key    string `yaml:"key"`
	MaxAge int    `yaml:"max-age"`
}

type ServerConfig struct {
	Address          string   `yaml:"address"`
	WhiteListOrigins []string `yaml:"white-list-origins"`
}

func readConfig() *Config {

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		panic(fmt.Errorf("missing CONFIG_FILE env variable"))
	}

	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("missing config file: %s", err.Error()))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read config file %s: %s", path, err.Error()))
	}
	fmt.Printf("\n\n ----- %s ----- \n%s\n\n", path, string(file))

	config := Config{}
	err = yaml.UnmarshalStrict(file, &config)
	if err != nil {
		panic(fmt.Errorf("unable to decode config file %s: %s", path, err.Error()))
	}

	// replace env variables
	config.Log.File = replaceEnvVariables(config.Log.File)

	return &config
}

// //////////////////////////////////////////////////
// secrets

type Secrets struct {
	SessionSecretKey string `yaml:"session-secret-key"`
	CookieSecret     string `yaml:"cookie-secret"`
}

func readSecrets() *Secrets {

	path := os.Getenv("SECRET_FILE")
	if path == "" {
		panic(fmt.Errorf("missing SECRET_FILE env variable"))
	}

	stats, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("missing secret file: %s", err.Error()))
	}

	permissions := stats.Mode().Perm()
	if permissions != 0o600 {
		panic(fmt.Errorf("incorrect permissions for secret file %s (0%o), must be 0600 for '%s'", permissions, permissions, path))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read secret file %s: %s", path, err.Error()))
	}

	secrets := &Secrets{}
	err = yaml.UnmarshalStrict(file, secrets)
	if err != nil {
		panic(fmt.Errorf("unable to decode secret file %s: %s", path, err.Error()))
	}

	return secrets
}

// //////////////////////////////////////////////////
// env variable

func replaceEnvVariables(value string) string {
	regexp := regexp.MustCompile(`\$([A-Z_]+)`)
	matches := regexp.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		envVariable := match[1]
		envValue := os.Getenv(envVariable)
		value = strings.ReplaceAll(value, match[0], envValue)
	}
	return value
}
