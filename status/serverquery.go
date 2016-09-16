package status

import (
	"time"
)

type ServerQuery struct {
	ip   string
	port string
}

type CachedServerQuery struct {
	ServerQuery
	expireTime     time.Duration
	CacheTimestamp time.Time
	cachedStatus   ServerStatus
	cachedError    error
}

func NewServerQuery(ip string, port string) ServerQuery {
	server := ServerQuery{}
	server.ip = ip
	server.port = port
	return server
}

func NewCachedServerQuery(ip string, port string) CachedServerQuery {
	cachedServer := CachedServerQuery{}
	cachedServer.ip = ip
	cachedServer.port = port
	cachedServer.expireTime = time.Duration(10 * time.Second)
	return cachedServer
}

func (query ServerQuery) IsPresent() bool {
	return query.ip != "" && query.port != ""
}

func (query ServerQuery) Matches(ip string, port string) bool {
	return query.ip == ip && query.port == port
}

func (serverQuery ServerQuery) Poll() (ServerStatus, error) {
	return query(serverQuery.ip, serverQuery.port)
}

func (serverQuery CachedServerQuery) ExpireAfter(expireTime time.Duration) CachedServerQuery {
	serverQuery.expireTime = expireTime
	return serverQuery
}

func (serverQuery *CachedServerQuery) Expired() bool {
	return time.Now().Sub(serverQuery.CacheTimestamp) > serverQuery.expireTime
}

func (serverQuery *CachedServerQuery) Poll() (ServerStatus, error) {
	if !serverQuery.Expired() {
		// return the cached status
		return serverQuery.cachedStatus, serverQuery.cachedError
	}

	// query server for new status
	serverQuery.cachedStatus, serverQuery.cachedError = query(serverQuery.ip, serverQuery.port)
	serverQuery.CacheTimestamp = time.Now()
	return serverQuery.cachedStatus, serverQuery.cachedError
}