# zabbix-exporter-3000
 just another zabbix exporter for prometheus

 like the other exporters it use Zabbix API and represent response as prometheus metrics.

### Limitations:

Main limitation - one instance = one query.

### Configuration
| Variable |	Description |	Default |
| --- | ----- | ---- |
| ZABBIX_API_ENDPOINT | full url to Zabbix API | http://zabbix/api_jsonrpc.php |
| ZABBIX_USER | Zabbix user | admin |
| ZABBIX_PASSWORD | Zabbix password | admin |
| ZABBIX_SKIP_SSL | Skip Zabbix endpoint SSL check | true |
| ZE3000_STRICT_METRIC_REG | May be useful when you have an error about metric duplicate on registration - set this to false. On this case, you highly likely have a problem with query, but this may help you investigate. Don't set this to 'false' on real environment | true |
| ZE3000_SINGLE_METRIC | If you, for some reason, won't use Default mechanics with mapping metric name and field from Zabbix response.  | true |
| ZE3000_SINGLE_METRIC_HELP | Hardcoded HELP field for Single metric mechanics | single description |
| ZE3000_HOST_PORT | which host and port exporter should listening. Supported notations - 0.0.0.0:9080 or :9080 | localhost:8080 |
| ZE3000_METRIC_NAMESPACE | Metric namespace (part of metric name in Prometheus) | zbx |
| ZE3000_METRIC_SUBSYSTEM | Metric subsystem (part of metric name in Prometheus) | subsystem |
| ZE3000_METRIC_NAME_PREFIX | Metric name prefix | prefix |
| ZE3000_METRIC_NAME_FIELD | `Mapping field.` Which field form Zabbix response use as part of a name. Please note - this field will be trimmed, set to lower case and rid off of all symbols except A-z and 0-9. `Only top level Zabbix response fields supported`  | key_ |
| ZE3000_METRIC_VALUE | `Mapping field.` Which field form Zabbix response use as value of metric. `Only top level Zabbix response fields supported`| lastvalue |
| ZE3000_METRIC_HELP | `Mapping field.` Which field form Zabbix response use as help field of metric. `Only top level Zabbix response fields supported` | description |
| ZE3000_ZABBIX_METRIC_LABELS | `Mapping field.` Which field form Zabbix response use as labels. `This field supported first level and second level fields ` | name,itemid,key_,hosts>host,hosts>name,interfaces>ip,interface>dns |
| ZE3000_ZABBIX_REFRESH_DELAY_SEC | How frequent Zabbix exporter will be query Zabbix. In seconds | 10 |
| ZE3000_ZABBIX_QUERY  | any Zabbix query, with field "auth" with value "%auth-token%" - yes, literally "%auth-token%" | ```{     "jsonrpc": "2.0",     "method": "item.get",     "params": {     	"application":"My Valuable Application",         "output": ["itemid","key_","description","lastvalue"],         "selectDependencies": "extend",         "selectHosts": ["name","status","host"],         "selectInterfaces": ["ip","dns"],         "sortfield":"key_" },     "auth": "%auth-token%",     "id": 1 }``` |


### How-to use
#### requirements
 - zabbix
 - prometheus
 - docker or k8s

Make some query to zabbix server over [Insomnia](https://insomnia.rest/download/), [Postman](https://www.postman.com/), [curl](https://curl.haxx.se/), you name it. Let's say this query is:
``` json
{
    "jsonrpc": "2.0",
    "method": "item.get",
    "params": {
    	"application":"My Super Application",
        "output": ["itemid","key_","description","lastvalue"],
        "selectDependencies": "extend",
        "selectHosts": ["name","status","host"],
        "selectInterfaces": ["ip","dns"],
        "sortfield":"key_"
    },
    "auth": "1234ml34kl3f4mk4gkl680klfmkl3fml",
    "id": 1
}
```

and response of this query is:
``` json

"jsonrpc": "2.0",
    "result": [
        {
            "itemid": "452345",
            "key_": "concurrencyConnections",
            "description": "The number of current concurrency connections.",
            "hosts": [
                {
                    "hostid": "54637",
                    "name": "Mighty Frontend",
                    "status": "2",
                    "host": "mighty.fronend"
                }
            ],
            "interfaces": [],
            "lastvalue": "9"
        },
        {
            "itemid": "902934",
            "key_": "numbeOfConnections",
            "description": "The number of currently active connections.",
            "hosts": [
                {
                    "hostid": "42092",
                    "name": "Mega Application",
                    "status": "0",
                    "host": "mega.application"
                }
            ],
            "interfaces": [
                {
                    "interfaceid": "1900",
                    "ip": "10.4.4.3",
                    "dns": ""
                }
            ],
            "lastvalue": "10987"
        },
      ],
  "id": 1
}

```
Since we know the query and know what is return - let's configure and start Zabbix Exporter 3000:

``` bash
docker run -d \
      -p 8080:8080 \
      -e ZABBIX_API_ENDPOINT=https://zabbix.example.com/zabbix/api_jsonrpc.php \
      -e ZABBIX_USER=someuser \
      -e ZABBIX_PASSWORD=str0nGpA5sw0rd \
      -e ZABBIX_SKIP_SSL=true \
      -e ZE3000_STRICT_METRIC_REG=true \
      -e ZE3000_METRIC_NAME_FIELD="key_" \
      -e ZABBIX_SKIP_SSL=true \
      -e ZE3000_SINGLE_METRIC=false \
      -e ZE3000_METRIC_NAMESPACE="megacompany" \
      -e ZE3000_METRIC_SUBSYSTEM="frontend" \
      -e ZE3000_METRIC_NAME_PREFIX="nginx" \
      -e ZE3000_METRIC_NAME_FIELD="key_" \
      -e ZE3000_METRIC_VALUE="lastvalue" \
      -e ZE3000_METRIC_HELP="description" \
      -e ZE3000_ZABBIX_REFRESH_DELAY_SEC=20 \
      -e ZE3000_ZABBIX_METRIC_LABELS="itemid,key_,hosts>host,hosts>name,interfaces>ip,interface>dns" \
      -e ZE3000_HOST_PORT=localhost:8080 \
      -e ZE3000_ZABBIX_QUERY="{     "jsonrpc": "2.0",     "method": "item.get",     "params": {     	"application":"My Super Application",         "output": ["itemid","key_","description","lastvalue"],         "selectDependencies": "extend",         "selectHosts": ["name","status","host"],         "selectInterfaces": ["ip","dns"],         "sortfield":"key_"     },     "auth": "%auth-token%",     "id": 1 }"
      rzrbld/adminio-api:latest

```
:boom: let's suppose everything running ok, and you don't have any error messages from ze3000 <br/><br/>
ze3000 brings up next endpoints:
- `/metrics` - main and exported metrics
- `/ready` - readiness probe for k8s monitoring
- `/live` - liveness probe for k8s monitoring

Let's se at `/metrics`
``` bash
$ curl http://localhost:8080/metrics
...
megacompany_frontend_nginx_concurrencyconnections{hosts_host="mighty.fronend",hosts_name="Mighty Frontend",interface_dns="NA",interfaces_ip="10.4.4.3",itemid="452345",key_="concurrencyConnections"} 9
...
megacompany_frontend_nginx_numbeofconnections{hosts_host="mega.application",hosts_name="Mega Application",interface_dns="NA",interfaces_ip="NA",itemid="902934",key_="numbeOfConnections"} 10987


```
