package ipcheck

import (
	"net"
	"net/http"
)

func WithIPCheck(next http.Handler, subnet string) http.Handler {
	checkFn := func(w http.ResponseWriter, r *http.Request) {
		if subnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			http.Error(w, "need to specify X-Real-IP", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(realIP)
		_, subn, err := net.ParseCIDR(subnet)
		if err != nil {
			http.Error(w, "wrong trusted subnet format", http.StatusBadRequest)
			return
		}

		if !subn.Contains(ip) {
			http.Error(w, "ip not in trusted subnet", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(checkFn)
}
