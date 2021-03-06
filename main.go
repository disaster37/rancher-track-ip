package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/disaster37/rancher-track-ip/trackip"
	repository "github.com/disaster37/rancher-track-ip/trackip/repository"
	usecase "github.com/disaster37/rancher-track-ip/trackip/usecase"
	elastic "github.com/elastic/go-elasticsearch/v7"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Read config file
	configHandler := viper.New()
	configHandler.SetConfigFile(`config/config.yaml`)
	err := configHandler.ReadInConfig()
	if err != nil {
		panic(err)
	}

	level, err := log.ParseLevel(configHandler.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	// Connect on repositories
	cfg := elastic.Config{
		Addresses: configHandler.GetStringSlice("elasticsearch.urls"),
		Username:  configHandler.GetString("elasticsearch.username"),
		Password:  configHandler.GetString("elasticsearch.password"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	elasticClient, err := elastic.NewClient(cfg)
	if err != nil {
		log.Errorf("failed to connect on elasticsearch: %s", err.Error())
		panic("failed to connect on elasticsearch")
	}
	elasticRepo := repository.NewElasticsearchRepository(elasticClient, configHandler.GetString("elasticsearch.index"))

	var containerRepo trackip.RancherRepository
	switch configHandler.GetString("platform") {
	case "cattle":
		opts := &rancher.ClientOpts{
			Url:       configHandler.GetString("rancher.url"),
			AccessKey: configHandler.GetString("rancher.key"),
			SecretKey: configHandler.GetString("rancher.secret"),
			Timeout:   time.Second * 5,
		}
		rancherClient, err := rancher.NewRancherClient(opts)
		if err != nil {
			log.Errorf("failed to connect on Rancher: %s", err.Error())
			panic("failed to connect on Rancher")
		}

		containerRepo = repository.NewRancherRepository(rancherClient)
		break
	case "kubernetes":
		config, err := clientcmd.BuildConfigFromFlags("", "config/kubeconfig")
		if err != nil {
			log.Errorf("Failed to load kubernetes config config/kubeconfig: %s", err.Error())
			panic("Failed to load kubernetes config")
		}
		kubeClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Errorf("Failed to connect on kubernetes: %s", err.Error())
			panic("Failed to connect on kubernetes")
		}
		containerRepo = repository.NewKubernetesRepository(kubeClient)
		break
	default:
		panic("Only cattle or kubernetes as been supported")
	}

	// Init usecase
	trackIPusecase := usecase.NewTrackIPUsecase(elasticRepo, containerRepo)

	err = trackIPusecase.TrackContainers(context.Background(), configHandler.GetInt64("loop"))
	if err != nil {
		log.Error(err)
		panic("Track IP crash")
	}
}
