package command

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/evergreen-ci/evergreen/agent/internal/client"
	agentutil "github.com/evergreen-ci/evergreen/agent/internal/testutil"
	"github.com/evergreen-ci/evergreen/db"
	"github.com/evergreen-ci/evergreen/model"
	modelutil "github.com/evergreen-ci/evergreen/model/testutil"
	"github.com/evergreen-ci/evergreen/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
)

func TestIncKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Convey("With keyval plugin installed", t, func() {
		err := db.Clear(model.KeyValCollection)
		require.NoError(t, err)
		require.NoError(t, err)

		testConfig := testutil.TestConfig()
		// These test don't actually need the integration test settings, but
		// MakeTaskConfigFromModelData needs it to create an (unused) GitHub app
		// token.
		testutil.ConfigureIntegrationTest(t, testConfig, t.Name())

		configPath := filepath.Join(testutil.GetDirectoryOfFile(), "testdata", "plugin_keyval.yml")

		comm := client.NewMock("http://localhost.com")

		modelData, err := modelutil.SetupAPITestData(testConfig, "testinc", "rhel55", configPath, modelutil.NoPatch)
		require.NoError(t, err)
		conf, err := agentutil.MakeTaskConfigFromModelData(ctx, testConfig, modelData)
		require.NoError(t, err)

		Convey("Inc command should increment a key successfully", func() {
			logger, err := comm.GetLoggerProducer(ctx, &conf.Task, nil)
			require.NoError(t, err)

			for _, task := range conf.Project.Tasks {
				So(len(task.Commands), ShouldNotEqual, 0)
				for _, command := range task.Commands {
					pluginCmds, err := Render(command, &model.Project{}, BlockInfo{})
					require.NoError(t, err)
					So(pluginCmds, ShouldNotBeNil)
					So(err, ShouldBeNil)
					for _, cmd := range pluginCmds {
						err = cmd.Execute(ctx, comm, logger, conf)
						So(err, ShouldBeNil)
					}
				}
				So(conf.Expansions.Get("testkey"), ShouldEqual, "2")
				So(conf.Expansions.Get("testkey_x"), ShouldEqual, "1")
			}
		})
	})
}
