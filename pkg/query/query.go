package query

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	"fmt"
	"regexp"
	"time"

	"go.uber.org/zap"
)

type QueryType int
type MathOp int
type FileType string

const (
	SELECT QueryType = 0
	DELETE QueryType = 1
)

const (
	GT   MathOp = 0
	LT   MathOp = 1
	EQ   MathOp = 2
	NE   MathOp = 3
	LIKE MathOp = 4
)

type Query struct {
	Node *node.Node
}

func NewQuery(node *node.Node) *Query {
	return &Query{
		Node: node,
	}
}

// tokenise uses regex to split the query string
func (query *Query) tokenise(queryStr string) []string {
	regex := "(?i)([a-zA-Z-_.]*[><=/ ]*[0-9]+)|[a-z-_.,]*[a-z-_.,]*([ = /LIKE]*\"[a-z0-9-_.\\[\\]\\* ]+\")*"
	return regexp.MustCompile(regex).FindAllString(queryStr, 100)
}

// Execute executes the query
func (query *Query) Execute(querystr string) ([]map[string]interface{}, int64, int, []error) {
	start := time.Now()
	decodedData, errs := query.decodeQuery(querystr)
	if len(errs) > 0 {
		logger.Logger.Warn("Found errors in query", zap.Errors("errs", errs))
	}
	elapsed := time.Since(start).Milliseconds()
	logger.Logger.Info(fmt.Sprintf("%s executed in %vms", querystr, elapsed))

	return decodedData, elapsed, len(decodedData), errs
}
