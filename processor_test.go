package context

import (
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestBootstrapProcessor(t *testing.T) {
	mockRegistry := &mockPeaRegistry{}
	mockRegistry.On("ContainsPeaDefinition", mock.AnythingOfType("string")).Return(false)
	mockRegistry.On("RegisterPeaDefinition", mock.AnythingOfType("string"), mock.AnythingOfType("context.ScannedPeaDefinition")).Return(false)

	bootstrapProcessor := NewBootstrapProcessor()
	bootstrapProcessor.processPeaDefinitions(mockRegistry)

	mockRegistry.AssertExpectations(t)
}

func TestConfigurationPropertiesBindingProcessor_WithNonPointer(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()
	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := NewConfigurationPropertiesBindingProcessor(standardEnvironment, core.NewDefaultTypeConverterService())

	properties := testProperties{}
	pea, err := propertiesBindingProcessor.BeforePeaInitialization("test", properties)
	assert.NotNil(t, err)
	assert.Equal(t, "configuration properties cannot be bound as it is not a type of pointer", err.Error())
	assert.Nil(t, pea)

	pea, err = propertiesBindingProcessor.AfterPeaInitialization("test", properties)
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}

type testPropertiesWithEmptyPrefix struct {
}

func (testPropertiesWithEmptyPrefix) GetConfigurationPrefix() string {
	return ""
}

func TestConfigurationPropertiesBindingProcessor_WithEmptyPrefix(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()
	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := NewConfigurationPropertiesBindingProcessor(standardEnvironment, core.NewDefaultTypeConverterService())

	properties := testPropertiesWithEmptyPrefix{}
	pea, err := propertiesBindingProcessor.BeforePeaInitialization("test", properties)
	assert.NotNil(t, err)
	assert.Equal(t, "prefix must not be null", err.Error())
	assert.Nil(t, pea)

	pea, err = propertiesBindingProcessor.AfterPeaInitialization("test", properties)
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}

func TestConfigurationPropertiesBindingProcessor(t *testing.T) {
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
	propertiesBindingProcessor := NewConfigurationPropertiesBindingProcessor(standardEnvironment, core.NewDefaultTypeConverterService())

	properties := &testProperties{}
	pea, err := propertiesBindingProcessor.BeforePeaInitialization("test", properties)
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	pea, err = propertiesBindingProcessor.AfterPeaInitialization("test", properties)
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	assert.Equal(t, "test-driver", properties.DriverName)
	assert.Equal(t, "test-user", properties.Username)
	assert.Equal(t, "test-pass", properties.Password)
	assert.Equal(t, "test-db", properties.DatabaseName)
	assert.Equal(t, uint16(3000), properties.Port)
	assert.Equal(t, uint32(1000), properties.Timeout)

	os.Args = tempArgs
}

func TestConfigurationPropertiesBindingProcessor_WithEmptyInstance(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := NewConfigurationPropertiesBindingProcessor(standardEnvironment, core.NewDefaultTypeConverterService())
	pea, err := propertiesBindingProcessor.BeforePeaInitialization("test", nil)
	assert.Nil(t, err)
	assert.Nil(t, pea)
}

func TestConfigurationPropertiesBindingProcessor_WithNonConfigurationPropertiesInstance(t *testing.T) {
	standardEnvironment := core.NewStandardEnvironment()

	propertySources := core.NewSimpleCommandLinePropertySource(os.Args)
	standardEnvironment.GetPropertySources().Add(propertySources)
	propertiesBindingProcessor := NewConfigurationPropertiesBindingProcessor(standardEnvironment, core.NewDefaultTypeConverterService())
	pea, err := propertiesBindingProcessor.BeforePeaInitialization("test", testContext{})
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}
