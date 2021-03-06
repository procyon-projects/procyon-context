package context

import (
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testProperties struct {
	DriverName   string `yaml:"driver-name" json:"driver-name"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	DatabaseName string `yaml:"database-name" json:"database-name"`
	Port         uint16 `yaml:"port" json:"port"`
	Timeout      uint32 `yaml:"timeout" json:"timeout"`
	DefaultValue string `yaml:"default" json:"default" default:"test"`
}

func (testProperties) GetConfigurationPrefix() string {
	return "test"
}

func TestConfigurationPropertiesBinder(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	tempArgs := os.Args
	os.Args = append(os.Args, "--test.driver-name=test-driver")
	os.Args = append(os.Args, "--test.username=test-user")
	os.Args = append(os.Args, "--test.password=test-pass")
	os.Args = append(os.Args, "--test.database-name=test-db")
	os.Args = append(os.Args, "--test.port=3000")
	os.Args = append(os.Args, "--test.timeout=1000")

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := newConfigurationPropertiesBinder(standardEnvironment, core.NewDefaultTypeConverterService())

	properties := &testProperties{}
	err := propertiesBindingProcessor.Bind(properties)
	assert.Nil(t, err)

	assert.Equal(t, "test-driver", properties.DriverName)
	assert.Equal(t, "test-user", properties.Username)
	assert.Equal(t, "test-pass", properties.Password)
	assert.Equal(t, "test-db", properties.DatabaseName)
	assert.Equal(t, uint16(3000), properties.Port)
	assert.Equal(t, uint32(1000), properties.Timeout)
	assert.Equal(t, "test", properties.DefaultValue)

	os.Args = tempArgs
}

type testYamlProperties struct {
	DriverName   string `yaml:"driver-name"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database-name"`
	Port         uint16 `yaml:"port"`
	Timeout      uint32 `yaml:"timeout"`
}

func (testYamlProperties) GetConfigurationPrefix() string {
	return "test"
}

func TestConfigurationPropertiesBinderForYaml(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	tempArgs := os.Args
	os.Args = append(os.Args, "--test.driver-name=test-driver")
	os.Args = append(os.Args, "--test.username=test-user")
	os.Args = append(os.Args, "--test.password=test-pass")
	os.Args = append(os.Args, "--test.database-name=test-db")
	os.Args = append(os.Args, "--test.port=3000")
	os.Args = append(os.Args, "--test.timeout=1000")

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := newConfigurationPropertiesBinder(standardEnvironment, core.NewDefaultTypeConverterService())

	properties := &testProperties{}
	err := propertiesBindingProcessor.Bind(properties)
	assert.Nil(t, err)

	assert.Equal(t, "test-driver", properties.DriverName)
	assert.Equal(t, "test-user", properties.Username)
	assert.Equal(t, "test-pass", properties.Password)
	assert.Equal(t, "test-db", properties.DatabaseName)
	assert.Equal(t, uint16(3000), properties.Port)
	assert.Equal(t, uint32(1000), properties.Timeout)
	assert.Equal(t, "test", properties.DefaultValue)

	os.Args = tempArgs
}

func TestConfigurationPropertiesBinder_BindWithEmptyInstance(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := newConfigurationPropertiesBinder(standardEnvironment, core.NewDefaultTypeConverterService())

	err := propertiesBindingProcessor.Bind(nil)
	assert.Nil(t, err)
}

func TestConfigurationPropertiesBinder_BindWithNonConfigurationPropertiesInstance(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := newConfigurationPropertiesBinder(standardEnvironment, core.NewDefaultTypeConverterService())

	err := propertiesBindingProcessor.Bind(testContext{})
	assert.Nil(t, err)
}
