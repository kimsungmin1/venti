package prometheus

import (
	"fmt"

	"github.com/kuoss/venti/pkg/mocker"
)

func New() (*mocker.Server, error) {
	s := NewWithPort(0)
	if err := s.Start(); err != nil {
		return nil, fmt.Errorf("start err: %w", err)
	}
	return s, nil
}

func NewWithPort(port int) *mocker.Server {
	s := mocker.NewWithPort(port)
	s.GET("/api/v1/status/buildinfo", handleBuildInfo)
	s.GET("/api/v1/metadata", handleMetadata)
	s.GET("/api/v1/query", handleQuery)
	s.GET("/api/v1/query_range", handleQueryRange)
	return s
}

func handleBuildInfo(c *mocker.Context) {
	c.JSONString(200, `{"status":"success","data":{"version":"2.41.0-prometheus","revision":"c0d8a56c69014279464c0e15d8bfb0e153af0dab","branch":"HEAD","buildUser":"root@d20a03e77067","buildDate":"20221220-10:40:45","goVersion":"go1.19.4"}}`)
}

func handleMetadata(c *mocker.Context) {
	c.JSONString(200, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`)
}

func handleQuery(c *mocker.Context) {
	query := c.Query("query")

	// 405
	if query == "" {
		c.JSONString(405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`)
		return
	}

	// 200
	if query == "up" {
		c.JSONString(200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]},{"metric":{"__name__":"up","job":"prometheus2","instance2":"localhost:9092"},"value":[1435781451.781,"1"]}]}}`)
		return
	}
	// 200
	if query == "unmarshalable" {
		c.JSONString(200, `{"status":"success","data":{"a":"b":"c"}}`)
		return
	}
	// 200 metric_not_exists
	c.JSONString(200, `{"status":"success","data":{"resultType":"vector","result":[]}}`)
}

func handleQueryRange(c *mocker.Context) {
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	step := c.Query("step")

	// 405
	if start == "" {
		c.JSONString(405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`)
		return
	}
	if end == "" {
		c.JSONString(405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`)
		return
	}
	if step == "" {
		c.JSONString(405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`)
		return
	}
	if query == "" {
		c.JSONString(405, `{"status":"error","errorType":"bad_data","error":"1:1: parse error: no expression found in input"}`)
		return
	}

	// 200
	if query == "up" {
		c.JSONString(200, `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`)
		return
	}
	// 200 metric_not_exists
	c.JSONString(200, `{"status":"success","data":{"resultType":"matrix","result":[]}}`)
}
