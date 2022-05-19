/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    traefik
 * @Date:    2022/3/24 6:34 PM
 * @package: traefik
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package traefik

const (
	namespace = "traefik" // root key

	// Routers
	kRule            = "/http/routers/%s/rule"                   // value: PathPrefix(`/ts`)
	kEntryPoints     = "/http/routers/%s/entrypoints/%d"         // value: web
	kMiddlewares     = "/http/routers/%s/middlewares/%d"         // value: auth
	kService         = "/http/routers/%s/service"                // value: myservice
	kTls             = "/http/routers/%s/tls"                    // value: true
	kTlsCertResolver = "/http/routers/%s/tls/certresolver"       // value: myresolver
	kTlsDomain       = "/http/routers/%s/tls/domains/%d/main"    // value: example.org
	kTlsSubDomain    = "/http/routers/%s/tls/domains/%d/sans/%d" // value: test.example.org
	kTlsOption       = "/http/routers/%s/tls/options"            // value: foobar
	kPriority        = "/http/routers/%s/priority"               // value: 1/2/3
	// === service ===

	// Services
	kServerUrl            = "/http/services/%s/loadbalancer/servers/%d/url"                   // value: http://<ip-server-1>:<port-server-1>/
	kServersTransport     = "/http/services/%s/loadbalancer/serverstransport"                 // value: foobar@file
	kPassHostHeader       = "/http/services/%s/loadbalancer/passhostheader"                   // value: true
	kHealthCheckHeader    = "/http/services/%s/loadbalancer/healthcheck/headers/%s"           // value: foobar
	kHealthCheckHostname  = "/http/services/%s/loadbalancer/healthcheck/hostname"             // value: example.org
	kHealthCheckInterval  = "/http/services/%s/loadbalancer/healthcheck/interval"             // value: 10
	kHealthCheckPath      = "/http/services/%s/loadbalancer/healthcheck/path"                 // value: /foo
	kHealthCheckPort      = "/http/services/%s/loadbalancer/healthcheck/port"                 // value: 8000
	kHealthCheckScheme    = "/http/services/%s/loadbalancer/healthcheck/scheme"               // value: http
	kHealthCheckTimeout   = "/http/services/%s/loadbalancer/healthcheck/timeout"              // value: 10
	kSticky               = "/http/services/%s/loadbalancer/sticky"                           // value: true
	kHttpOnly             = "/http/services/%s/loadbalancer/sticky/cookie/httponly"           // value: true
	kCookieName           = "/http/services/%s/loadbalancer/sticky/cookie/name"               // value: foobar
	kCookieSecure         = "/http/services/%s/loadbalancer/sticky/cookie/secure"             // value: true
	kCookieSameSite       = "/http/services/%s/loadbalancer/sticky/cookie/samesite"           // value: none
	kFlushInterval        = "/http/services/%s/loadbalancer/responseforwarding/flushinterval" // value: 10
	kMirroringService     = "/http/services/%s/mirroring/service"                             // value: foobar
	kMirrorsName          = "/http/services/%s/mirroring/mirrors/%d/name"                     // value: foobar
	kMirrorsPercent       = "/http/services/%s/mirroring/mirrors/%d/percent"                  // value: 42
	kServicesName         = "/http/services/%s/weighted/services/%d/name"                     // value: foobar
	kServicesWeight       = "/http/services/%s/weighted/services/%d/weight"                   // value: 42
	kWeightedCookieName   = "/http/services/%s/weighted/sticky/cookie/name"                   // value: foobar
	kWeightedCookieSecure = "/http/services/%s/weighted/sticky/cookie/secure"                 // value: true
	kWeightSameSite       = "/http/services/%s/weighted/sticky/cookie/samesite"               // value: none
	kWeightedHttpOnly     = "/http/services/%s/weighted/sticky/cookie/httpOnly"               // value: true
)
