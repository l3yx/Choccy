package util

import (
	"github.com/tidwall/gjson"
	"os"
)

type CodeQLLocation struct {
	Uri              string
	StartLine        int64
	StartColumn      int64
	EndLine          int64
	EndColumn        int64
	ContextStartLine int64
	ContextEndLine   int64
	ContextSnippet   string
	Message          string
}
type CodeQLCodeFlow struct {
	Location []CodeQLLocation
}
type CodeQLRelatedLocation struct {
	Id      int64
	Message string
	CodeQLLocation
}
type CodeQLResult struct {
	Rule             string                  //rule id
	Message          string                  //Newsletter
	RelatedLocations []CodeQLRelatedLocation //The hyperlink in the message points to the location, if not, nil
	Location         CodeQLLocation          //Location of vulnerability
	CodeFlows        []CodeQLCodeFlow        //Code path, a path-problem type vulnerability may have multiple paths (the same endpoint), the problem type field is nil
}
type CodeQLSarif struct {
	SemanticVersion string
	NotificationsId []string       //I don't know what to do, but it can be used to judge language roughly
	Rules           []string       //All rules used id
	Packs           []string       //All package names used
	Results         []CodeQLResult //Scan results, length is the number of vulnerabilities
}

func ParseSarifFile(path string) (*CodeQLSarif, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(bytes)

	codeQLSarif := CodeQLSarif{}
	//CodeQL edition
	codeQLSarif.SemanticVersion = gjson.Get(content, "runs.0.tool.driver.semanticVersion").String()
	//NotificationsId
	gjson.Get(content, "runs.0.tool.driver.notifications.#.id").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.NotificationsId = append(codeQLSarif.NotificationsId, value.String())
		return true
	})
	//All the rules used for scanning, which contains all the attributes of the ql rule, but only takes the id for now
	gjson.Get(content, "runs.0.tool.driver.rules.#.id").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.Rules = append(codeQLSarif.Rules, value.String())
		return true
	})
	//Scan all packs used, which include name, version, etc., but only name is taken for now.
	gjson.Get(content, "runs.0.tool.extensions.#.name").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.Packs = append(codeQLSarif.Packs, value.String())
		return true
	})

	codeQLSarif.Results = make([]CodeQLResult, 0)
	//The length of the scan results is the number of vulnerabilities.
	results := gjson.Get(content, "runs.0.results").Array()
	for _, result := range results {
		codeQLResult := CodeQLResult{}

		codeQLResult.Rule = result.Get("ruleId").String()
		codeQLResult.Message = result.Get("message.text").String()

		result.Get("relatedLocations").ForEach(func(key, value gjson.Result) bool {
			relatedLocation := CodeQLRelatedLocation{}
			relatedLocation.Id = value.Get("id").Int()
			relatedLocation.Message = value.Get("message.text").String()
			relatedLocation.Uri = value.Get("physicalLocation.artifactLocation.uri").String()
			relatedLocation.StartLine = value.Get("physicalLocation.region.startLine").Int()
			relatedLocation.StartColumn = value.Get("physicalLocation.region.startColumn").Int()
			relatedLocation.EndLine = value.Get("physicalLocation.region.endLine").Int()
			relatedLocation.EndColumn = value.Get("physicalLocation.region.endColumn").Int()
			relatedLocation.ContextStartLine = value.Get("physicalLocation.contextRegion.startLine").Int()
			relatedLocation.ContextEndLine = value.Get("physicalLocation.contextRegion.endLine").Int()
			relatedLocation.ContextSnippet = value.Get("physicalLocation.contextRegion.snippet.text").String()
			codeQLResult.RelatedLocations = append(codeQLResult.RelatedLocations, relatedLocation)
			return true
		})

		// locations is an array, but it seems that the length is all 1
		resultCodeQLLocation := CodeQLLocation{}
		resultCodeQLLocation.Uri = result.Get("locations.0.physicalLocation.artifactLocation.uri").String()
		resultCodeQLLocation.StartLine = result.Get("locations.0.physicalLocation.region.startLine").Int()
		resultCodeQLLocation.StartColumn = result.Get("locations.0.physicalLocation.region.startColumn").Int()
		resultCodeQLLocation.EndLine = result.Get("locations.0.physicalLocation.region.endLine").Int()
		resultCodeQLLocation.EndColumn = result.Get("locations.0.physicalLocation.region.endColumn").Int()
		//Context code and location
		resultCodeQLLocation.ContextStartLine = result.Get("locations.0.physicalLocation.contextRegion.startLine").Int()
		resultCodeQLLocation.ContextEndLine = result.Get("locations.0.physicalLocation.contextRegion.endLine").Int()
		resultCodeQLLocation.ContextSnippet = result.Get("locations.0.physicalLocation.contextRegion.snippet.text").String()
		codeQLResult.Location = resultCodeQLLocation

		//All paths of a vulnerability (the end points are the same). If there is no path, it is a problem.
		codeFlows := result.Get("codeFlows").Array()
		for _, codeFlow := range codeFlows {
			codeQLCodeFlow := CodeQLCodeFlow{}
			//Position sequence of a single path
			threadFlowLocations := codeFlow.Get("threadFlows.0.locations").Array()
			for _, threadFlowLocation := range threadFlowLocations {
				codeFlowLocation := CodeQLLocation{}

				codeFlowLocation.Uri = threadFlowLocation.Get("location.physicalLocation.artifactLocation.uri").String()
				codeFlowLocation.StartLine = threadFlowLocation.Get("location.physicalLocation.region.startLine").Int()
				codeFlowLocation.StartColumn = threadFlowLocation.Get("location.physicalLocation.region.startColumn").Int()
				codeFlowLocation.EndLine = threadFlowLocation.Get("location.physicalLocation.region.endLine").Int()
				codeFlowLocation.EndColumn = threadFlowLocation.Get("location.physicalLocation.region.endColumn").Int()
				codeFlowLocation.ContextStartLine = threadFlowLocation.Get("location.physicalLocation.contextRegion.startLine").Int()
				codeFlowLocation.ContextEndLine = threadFlowLocation.Get("location.physicalLocation.contextRegion.endLine").Int()
				codeFlowLocation.ContextSnippet = threadFlowLocation.Get("location.physicalLocation.contextRegion.snippet.text").String()
				codeFlowLocation.Message = threadFlowLocation.Get("location.message.text").String()
				codeQLCodeFlow.Location = append(codeQLCodeFlow.Location, codeFlowLocation)
			}
			codeQLResult.CodeFlows = append(codeQLResult.CodeFlows, codeQLCodeFlow)
		}
		codeQLSarif.Results = append(codeQLSarif.Results, codeQLResult)
	}

	return &codeQLSarif, nil
}
