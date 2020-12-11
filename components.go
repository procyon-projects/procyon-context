package context

import (
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	"strings"
)

type ScannedPeaDefinition struct {
	*peas.SimplePeaDefinition
	componentName string
}

func NewScannedPeaDefinition(componentName string, peaType goo.Type) ScannedPeaDefinition {
	return ScannedPeaDefinition{
		peas.NewSimplePeaDefinition(peaType),
		componentName,
	}
}

func (definition ScannedPeaDefinition) GetComponentName() string {
	return definition.componentName
}

type ScannedPeaNameGenerator struct {
}

func NewScannedPeaNameGenerator() ScannedPeaNameGenerator {
	return ScannedPeaNameGenerator{}
}

func (generator ScannedPeaNameGenerator) GenerateName(peaDefinition peas.PeaDefinition) string {
	peaTypeName := peaDefinition.GetTypeName()
	lastDotIndex := strings.LastIndex(peaTypeName, ".")
	shortName := ""
	if lastDotIndex != -1 {
		shortName = peaTypeName[lastDotIndex+1:]
	} else {
		shortName = peaTypeName
	}
	shortName = strings.ToLower(shortName[:1]) + shortName[1:]
	return shortName
}

type ComponentPeaDefinitionScanner struct {
	peaNameGenerator peas.PeaNameGenerator
	peaRegistry      peas.PeaDefinitionRegistry
}

func NewComponentPeaDefinitionScanner(registry peas.PeaDefinitionRegistry) ComponentPeaDefinitionScanner {
	return ComponentPeaDefinitionScanner{
		NewScannedPeaNameGenerator(),
		registry,
	}
}

func (scanner ComponentPeaDefinitionScanner) DoScan() {
	scannedPeaDefinitions := make([]ScannedPeaDefinition, 0)
	err := core.ForEachComponentType(func(componentName string, componentType goo.Type) error {
		scannedPeaDefinition := NewScannedPeaDefinition(componentName, componentType)
		scannedPeaDefinitions = append(scannedPeaDefinitions, scannedPeaDefinition)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, peaDefinition := range scannedPeaDefinitions {
		peaName := scanner.peaNameGenerator.GenerateName(peaDefinition)
		if !scanner.peaRegistry.ContainsPeaDefinition(peaName) {
			peaDefinitionHolder := peas.NewPeaDefinitionHolder(peaName, peaDefinition)
			scanner.peaRegistry.RegisterPeaDefinition(peaName, peaDefinitionHolder.GetPeaDefinition())
		}
	}
}
