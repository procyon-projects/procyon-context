package context

import (
	"github.com/codnect/goo"
	peas "github.com/procyon-projects/procyon-peas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestScannedPeaDefinition_GetComponentName(t *testing.T) {
	definition := NewScannedPeaDefinition("testComponent", goo.GetType(NewConfigurationPropertiesBindingProcessor))
	assert.Equal(t, peas.SharedScope, definition.GetScope())
	assert.Equal(t, "testComponent", definition.GetComponentName())
	assert.Equal(t, "github.com.procyon.projects.procyon.context.ConfigurationPropertiesBindingProcessor", definition.GetTypeName())
	assert.NotNil(t, definition.GetPeaType())
}

func TestScannedPeaNameGenerator_GenerateName(t *testing.T) {
	definition := NewScannedPeaDefinition("testComponent", goo.GetType(NewConfigurationPropertiesBindingProcessor))
	scannedPeaNameGenerator := NewScannedPeaNameGenerator()
	peaName := scannedPeaNameGenerator.GenerateName(definition)
	assert.Equal(t, "configurationPropertiesBindingProcessor", peaName)
}

type mockPeaRegistry struct {
	mock.Mock
}

func (registry mockPeaRegistry) RegisterPeaDefinition(peaName string, definition peas.PeaDefinition) {
	registry.Called(peaName, definition)
}

func (registry mockPeaRegistry) RemovePeaDefinition(peaName string) {
	registry.Called(peaName)
}

func (registry mockPeaRegistry) ContainsPeaDefinition(peaName string) bool {
	results := registry.Called(peaName)
	return results.Bool(0)
}

func (registry mockPeaRegistry) GetPeaDefinition(peaName string) peas.PeaDefinition {
	results := registry.Called(peaName)
	return results.Get(0).(peas.PeaDefinition)
}

func (registry mockPeaRegistry) GetPeaDefinitionNames() []string {
	results := registry.Called()
	if results == nil {
		return nil
	}
	return results.Get(0).([]string)
}

func (registry mockPeaRegistry) GetPeaDefinitionCount() int {
	results := registry.Called()
	return results.Int(0)
}

func (registry mockPeaRegistry) GetPeaNamesForType(typ goo.Type) []string {
	results := registry.Called(typ)
	if results == nil {
		return nil
	}
	return results.Get(0).([]string)
}

func TestComponentPeaDefinitionScanner_DoScan(t *testing.T) {
	mockRegistry := &mockPeaRegistry{}
	mockRegistry.On("ContainsPeaDefinition", mock.AnythingOfType("string")).Return(false)
	mockRegistry.On("RegisterPeaDefinition", mock.AnythingOfType("string"), mock.AnythingOfType("context.ScannedPeaDefinition")).Return(false)

	scanner := NewComponentPeaDefinitionScanner(mockRegistry)
	scanner.DoScan()
}
