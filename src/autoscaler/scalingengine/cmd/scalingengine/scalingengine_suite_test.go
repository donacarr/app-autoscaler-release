package main_test

import (
	"path/filepath"

	. "code.cloudfoundry.org/app-autoscaler/src/autoscaler/testhelpers"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/cf"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/db"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/scalingengine/config"

	"code.cloudfoundry.org/cfhttp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	yaml "gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestScalingengine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scalingengine Suite")
}

var (
	enginePath       string
	conf             config.Config
	port             int
	healthport       int
	configFile       *os.File
	ccUAA            *MockServer
	appId            string
	httpClient       *http.Client
	healthHttpClient *http.Client
)

var _ = SynchronizedBeforeSuite(
	func() []byte {
		compiledPath, err := gexec.Build("code.cloudfoundry.org/app-autoscaler/src/autoscaler/scalingengine/cmd/scalingengine", "-race")
		Expect(err).NotTo(HaveOccurred())
		return []byte(compiledPath)
	},
	func(pathBytes []byte) {
		enginePath = string(pathBytes)

		ccUAA = NewMockServer()
		ccUAA.RouteToHandler("GET", "/v2/info", ghttp.RespondWithJSONEncoded(http.StatusOK,
			cf.Endpoints{
				TokenEndpoint:   ccUAA.URL(),
				DopplerEndpoint: strings.Replace(ccUAA.URL(), "http", "ws", 1),
			}))

		ccUAA.RouteToHandler("POST", "/oauth/token", ghttp.RespondWithJSONEncoded(http.StatusOK, cf.Tokens{}))

		appId = fmt.Sprintf("app-id-%d", GinkgoParallelProcess())

		ccUAA.Add().GetApp(models.AppStatusStarted)
		ccUAA.Add().GetAppProcesses(2)

		ccUAA.RouteToHandler("PUT", "/v2/apps/"+appId, ghttp.RespondWith(http.StatusCreated, ""))

		conf.CF = cf.Config{
			API:      ccUAA.URL(),
			ClientID: "autoscaler_client_id",
			Secret:   "autoscaler_client_secret",
		}

		port = 7000 + GinkgoParallelProcess()
		healthport = 8000 + GinkgoParallelProcess()
		testCertDir := "../../../../../test-certs"

		verifyCertExistence(testCertDir)

		conf.Server.Port = port
		conf.Server.TLS.KeyFile = filepath.Join(testCertDir, "scalingengine.key")
		conf.Server.TLS.CertFile = filepath.Join(testCertDir, "scalingengine.crt")
		conf.Server.TLS.CACertFile = filepath.Join(testCertDir, "autoscaler-ca.crt")
		conf.Health.Port = healthport
		conf.Logging.Level = "debug"

		dbUrl := GetDbUrl()
		conf.DB.PolicyDB = db.DatabaseConfig{
			URL:                   dbUrl,
			MaxOpenConnections:    10,
			MaxIdleConnections:    5,
			ConnectionMaxLifetime: 10 * time.Second,
		}
		conf.DB.ScalingEngineDB = db.DatabaseConfig{
			URL:                   dbUrl,
			MaxOpenConnections:    10,
			MaxIdleConnections:    5,
			ConnectionMaxLifetime: 10 * time.Second,
		}
		conf.DB.SchedulerDB = db.DatabaseConfig{
			URL:                   dbUrl,
			MaxOpenConnections:    10,
			MaxIdleConnections:    5,
			ConnectionMaxLifetime: 10 * time.Second,
		}

		conf.DefaultCoolDownSecs = 300
		conf.LockSize = 32
		conf.HttpClientTimeout = 10 * time.Second

		conf.Health.HealthCheckUsername = "scalingenginehealthcheckuser"
		conf.Health.HealthCheckPassword = "scalingenginehealthcheckpassword"

		configFile = writeConfig(&conf)

		database, err := db.GetConnection(dbUrl)
		Expect(err).NotTo(HaveOccurred())

		testDB, err := sqlx.Open(database.DriverName, database.DSN)
		FailOnError("open db failed", err)
		defer func() { _ = testDB.Close() }()

		_, err = testDB.Exec(testDB.Rebind("DELETE FROM scalinghistory WHERE appid = ?"), appId)
		FailOnError("delete from scalinghistory", err)

		_, err = testDB.Exec(testDB.Rebind("DELETE from policy_json WHERE app_id = ?"), appId)
		FailOnError("delete from policy_json", err)

		_, err = testDB.Exec(testDB.Rebind("DELETE from activeschedule WHERE appid = ?"), appId)
		FailOnError("delete from activeschedule", err)

		_, err = testDB.Exec(testDB.Rebind("DELETE from app_scaling_active_schedule WHERE app_id = ?"), appId)
		FailOnError("delete from app_scaling_active_schedule", err)

		policy := `
		{
 			"instance_min_count": 1,
  			"instance_max_count": 5
		}`
		_, err = testDB.Exec(testDB.Rebind("INSERT INTO policy_json(app_id, policy_json, guid) values(?, ?, ?)"), appId, policy, "1234")
		FailOnError("insert failed", err)

		//nolint:staticcheck  // SA1019 TODO: https://github.com/cloudfoundry/app-autoscaler-release/issues/548
		tlsConfig, err := cfhttp.NewTLSConfig(
			filepath.Join(testCertDir, "eventgenerator.crt"),
			filepath.Join(testCertDir, "eventgenerator.key"),
			filepath.Join(testCertDir, "autoscaler-ca.crt"))
		Expect(err).NotTo(HaveOccurred())
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
		healthHttpClient = &http.Client{}

	})

func verifyCertExistence(testCertDir string) {
	_, err := ioutil.ReadFile(filepath.Join(testCertDir, "scalingengine.key"))
	Expect(err).NotTo(HaveOccurred())
	_, err = ioutil.ReadFile(filepath.Join(testCertDir, "scalingengine.crt"))
	Expect(err).NotTo(HaveOccurred())
	_, err = ioutil.ReadFile(filepath.Join(testCertDir, "autoscaler-ca.crt"))
	Expect(err).NotTo(HaveOccurred())
}

var _ = SynchronizedAfterSuite(
	func() {
		ccUAA.Close()
		_ = os.Remove(configFile.Name())
	},
	func() {
		gexec.CleanupBuildArtifacts()
	})

func writeConfig(c *config.Config) *os.File {
	cfg, err := ioutil.TempFile("", "engine")
	Expect(err).NotTo(HaveOccurred())

	defer func() { _ = cfg.Close() }()

	bytes, err := yaml.Marshal(c)
	Expect(err).NotTo(HaveOccurred())

	_, err = cfg.Write(bytes)
	Expect(err).NotTo(HaveOccurred())

	return cfg
}

type ScalingEngineRunner struct {
	configPath string
	startCheck string
	Session    *gexec.Session
}

func NewScalingEngineRunner() *ScalingEngineRunner {
	return &ScalingEngineRunner{
		configPath: configFile.Name(),
		startCheck: "scalingengine.started",
	}
}

func (engine *ScalingEngineRunner) Start() {
	// #nosec G204
	engineSession, err := gexec.Start(
		exec.Command(
			enginePath,
			"-c",
			engine.configPath,
		),
		gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[32m[engine]\x1b[0m ", GinkgoWriter),
		gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[32m[engine]\x1b[0m ", GinkgoWriter),
	)
	Expect(err).NotTo(HaveOccurred())
	engine.Session = engineSession
}

func (engine *ScalingEngineRunner) Interrupt() {
	if engine.Session != nil {
		engine.Session.Interrupt().Wait(5 * time.Second)
	}
}

func (engine *ScalingEngineRunner) KillWithFire() {
	if engine.Session != nil {
		engine.Session.Kill().Wait(5 * time.Second)
	}
}
