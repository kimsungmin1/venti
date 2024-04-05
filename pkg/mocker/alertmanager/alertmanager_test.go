package alertmanager_test

import (
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/stretchr/testify/assert"
)

var (
	server *mocker.Server
	client *mockerClient.Client
)

func init() {
	var err error
	server, err = alertmanager.New()
	if err != nil {
		panic(err)
	}
	client = mockerClient.New(server.URL())
}

func Test_api_v2_status(t *testing.T) {
	code, body, err := client.GET("/api/v2/status", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"cluster":{"name":"01HT1VJ2WY533RPFJ32JY1WDX7","peers":[{"address":"172.17.0.4:9094","name":"01HT1VJ2WY533RPFJ32JY1WDX7"}],"status":"ready"},"config":{"original":"global:\n  resolve_timeout: 5m\n  http_config:\n    follow_redirects: true\n    enable_http2: true\n  smtp_hello: localhost\n  smtp_require_tls: true\n  pagerduty_url: https://events.pagerduty.com/v2/enqueue\n  opsgenie_api_url: https://api.opsgenie.com/\n  wechat_api_url: https://qyapi.weixin.qq.com/cgi-bin/\n  victorops_api_url: https://alert.victorops.com/integrations/generic/20131114/alert/\n  telegram_api_url: https://api.telegram.org\n  webex_api_url: https://webexapis.com/v1/messages\nroute:\n  receiver: web.hook\n  group_by:\n  - alertname\n  continue: false\n  group_wait: 30s\n  group_interval: 5m\n  repeat_interval: 1h\ninhibit_rules:\n- source_match:\n    severity: critical\n  target_match:\n    severity: warning\n  equal:\n  - alertname\n  - dev\n  - instance\nreceivers:\n- name: web.hook\n  webhook_configs:\n  - send_resolved: true\n    http_config:\n      follow_redirects: true\n      enable_http2: true\n    url: <secret>\n    url_file: \"\"\n    max_alerts: 0\ntemplates: []\n"},"uptime":"2024-03-28T06:22:06.240Z","versionInfo":{"branch":"HEAD","buildDate":"20240228-11:51:20","buildUser":"root@22cd11f671e9","goVersion":"go1.21.7","revision":"0aa3c2aad14cff039931923ab16b26b7481783b5","version":"0.27.0"}}`, body)
}

func Test_api_v2_alerts(t *testing.T) {
	code, body, err := client.GET("/api/v2/alerts", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"status":"success","data":{}}`, body)
}
