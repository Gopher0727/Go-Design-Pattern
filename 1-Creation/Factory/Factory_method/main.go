package main

// 具体工厂方法
type RuleConfigParser interface {
	Parse(data []byte) (string, error)
}

type jsonParser struct{}

func (j *jsonParser) Parse(data []byte) (string, error) {
	return "json 规则配置解析完成", nil
}

type xmlParser struct{}

func (x *xmlParser) Parse(data []byte) (string, error) {
	return "xml 规则配置解析完成", nil
}

// type yamlParser struct{}

// func (y *yamlParser) Parse(data []byte) (string, error) {
// 	return "yaml 规则配置解析完成", nil
// }

// 工厂方法接口
type RuleConfigParserFactory interface {
	CreateParser() RuleConfigParser
}

type jsonParserFactory struct{}

func (j *jsonParserFactory) CreateParser() RuleConfigParser {
	return &jsonParser{}
}

type xmlParserFactory struct{}

func (x *xmlParserFactory) CreateParser() RuleConfigParser {
	return &xmlParser{}
}

// type yamlParserFactory struct{}

// func (y *yamlParserFactory) CreateParser() RuleConfigParser {
// 	return &yamlParser{}
// }

func NewRuleConfigParserFactory(parserType string) RuleConfigParserFactory {
	switch parserType {
	case "json":
		return &jsonParserFactory{}
	case "xml":
		return &xmlParserFactory{}
	default:
		return nil
	}
}
