package context

import (
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	"strings"
)

type RepositoryMetadata struct {
	typ *core.Type
}

type Repository interface {
	GetRepositoryMetadata() RepositoryMetadata
}

func NewRepositoryMetadata(typ *core.Type) RepositoryMetadata {
	return RepositoryMetadata{
		typ,
	}
}

type ServiceMetadata struct {
}

func NewServiceMetadata() ServiceMetadata {
	return ServiceMetadata{}
}

type Service interface {
	GetServiceMetadata() ServiceMetadata
}

type ScannedPeaDefinition struct {
	peas.SimplePeaDefinition
	componentName string
}

func NewScannedPeaDefinition(componentName string, peaType *core.Type) ScannedPeaDefinition {
	return ScannedPeaDefinition{
		peas.NewSimplePeaDefinition(peaType.String(), peaType, ""),
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
	peaTypeName := peaDefinition.GetName()
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
	err := core.VisitComponentTypes(func(componentName string, componentType *core.Type) error {
		scannedPeaDefinition := NewScannedPeaDefinition(componentName, componentType)
		scannedPeaDefinitions = append(scannedPeaDefinitions, scannedPeaDefinition)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, peaDefinition := range scannedPeaDefinitions {
		peaName := scanner.peaNameGenerator.GenerateName(peaDefinition)
		if scanner.checkPeaDefinition(peaName, peaDefinition) {
			peaDefinitionHolder := peas.NewPeaDefinitionHolder(peaName, peaDefinition)
			scanner.registerPeaDefinition(peaDefinitionHolder)
		}
	}
}

func (scanner ComponentPeaDefinitionScanner) checkPeaDefinition(peaName string, peaDefinition peas.PeaDefinition) bool {
	if !scanner.peaRegistry.ContainsPeaDefinition(peaName) {
		return true
	}
	return false
}

func (scanner ComponentPeaDefinitionScanner) registerPeaDefinition(peaDefinitionHolder *peas.PeaDefinitionHolder) {
	peaName := peaDefinitionHolder.GetPeaName()
	scanner.peaRegistry.RegisterPeaDefinition(peaName, peaDefinitionHolder.GetPeaDefinition())
	// register aliases
}
