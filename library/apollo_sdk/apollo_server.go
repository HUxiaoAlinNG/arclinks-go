package apollo_sdk

import (
	"github.com/zouyx/agollo"
)

const (
	DEFAULT_CLUSTER = "default"
)

var ApolloServer *apolloServer

// apolloServer warp up agollo sdk.
type apolloServer struct {
	appConfig *agollo.AppConfig
	logger    agollo.LoggerInterface
}

// init apolloServer
func NewApolloServer(apolloServerUrl string) {
	appConfig := &agollo.AppConfig{
		Cluster: DEFAULT_CLUSTER,
		Ip:      apolloServerUrl,
	}

	ApolloServer = &apolloServer{appConfig: appConfig}
}

// Start connect to apollo server.
func (a *apolloServer) Start(appId, namespaceName string) error {
	a.appConfig.AppId = appId
	a.appConfig.NamespaceName = namespaceName

	agollo.InitCustomConfig(func() (*agollo.AppConfig, error) {
		return a.appConfig, nil
	})

	if a.logger != nil {
		agollo.SetLogger(a.logger)
	}
	return agollo.Start()
}

// set logger
func (a *apolloServer) SetLog(logger agollo.LoggerInterface) *apolloServer {
	a.logger = logger
	return a
}

// ListenChangeEvent returns ChangeEvent.
func (a *apolloServer) ListenChangeEvent() <-chan *agollo.ChangeEvent {
	return agollo.ListenChangeEvent()
}

// GetStringValue gets string value from apollo.
func (a *apolloServer) GetStringValue(key, defaultValue string) string {
	return agollo.GetStringValue(key, defaultValue)
}

// GetIntValue gets int value from apollo.
func (a *apolloServer) GetIntValue(key string, defaultValue int) int {
	return agollo.GetIntValue(key, defaultValue)
}

// GetFloatValue gets float value from apollo.
func (a *apolloServer) GetFloatValue(key string, defaultValue float64) float64 {
	return agollo.GetFloatValue(key, defaultValue)
}

// GetBoolValue gets bool value from apollo.
func (a *apolloServer) GetBoolValue(key string, defaultValue bool) bool {
	return agollo.GetBoolValue(key, defaultValue)
}
