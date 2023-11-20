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
	Rule             string                  //规则id
	Message          string                  //提示信息
	RelatedLocations []CodeQLRelatedLocation //信息中的超链接指向的位置，没有的话为nil
	Location         CodeQLLocation          //漏洞所在位置
	CodeFlows        []CodeQLCodeFlow        //代码路径，一个path-problem类型的漏洞可能有多条路径（终点一样），problem类型该字段为nil
}
type CodeQLSarif struct {
	SemanticVersion string
	NotificationsId []string       //不知道干啥的，但是可以用来大致判断语言
	Rules           []string       //所用全部规则id
	Packs           []string       //所用全部pack名称
	Results         []CodeQLResult //扫描结果，长度就是漏洞数量
}

func ParseSarifFile(path string) (*CodeQLSarif, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(bytes)

	codeQLSarif := CodeQLSarif{}
	//CodeQL版本
	codeQLSarif.SemanticVersion = gjson.Get(content, "runs.0.tool.driver.semanticVersion").String()
	//NotificationsId
	gjson.Get(content, "runs.0.tool.driver.notifications.#.id").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.NotificationsId = append(codeQLSarif.NotificationsId, value.String())
		return true
	})
	//扫描所用的全部规则，里面包含了ql规则的所有属性，但暂只取id
	gjson.Get(content, "runs.0.tool.driver.rules.#.id").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.Rules = append(codeQLSarif.Rules, value.String())
		return true
	})
	//扫描所用的全部pack，里面包含了名称版本等，但暂只取name
	gjson.Get(content, "runs.0.tool.extensions.#.name").ForEach(func(key, value gjson.Result) bool {
		codeQLSarif.Packs = append(codeQLSarif.Packs, value.String())
		return true
	})

	codeQLSarif.Results = make([]CodeQLResult, 0)
	//扫描结果,长度就是漏洞数量
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

		// locations是个数组，但貌似长度都是1
		resultCodeQLLocation := CodeQLLocation{}
		resultCodeQLLocation.Uri = result.Get("locations.0.physicalLocation.artifactLocation.uri").String()
		resultCodeQLLocation.StartLine = result.Get("locations.0.physicalLocation.region.startLine").Int()
		resultCodeQLLocation.StartColumn = result.Get("locations.0.physicalLocation.region.startColumn").Int()
		resultCodeQLLocation.EndLine = result.Get("locations.0.physicalLocation.region.endLine").Int()
		resultCodeQLLocation.EndColumn = result.Get("locations.0.physicalLocation.region.endColumn").Int()
		//上下文代码及位置
		resultCodeQLLocation.ContextStartLine = result.Get("locations.0.physicalLocation.contextRegion.startLine").Int()
		resultCodeQLLocation.ContextEndLine = result.Get("locations.0.physicalLocation.contextRegion.endLine").Int()
		resultCodeQLLocation.ContextSnippet = result.Get("locations.0.physicalLocation.contextRegion.snippet.text").String()
		codeQLResult.Location = resultCodeQLLocation

		//一个漏洞所有的路径（终点都一样），没有路径的话就是problem
		codeFlows := result.Get("codeFlows").Array()
		for _, codeFlow := range codeFlows {
			codeQLCodeFlow := CodeQLCodeFlow{}
			//单条路径的位置序列
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
