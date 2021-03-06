package main

import (
  "cache"
  "client"
  "fmt"
  "log"
  "net/http"
  "time"
)

var Cache = cache.New()

func handler(w http.ResponseWriter, r *http.Request) {

  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in handler:", re, " at ", r.URL.String())
      for hk, _ := range w.Header() {
        if hk == "Server" {
          continue
        }
        w.Header().Del(hk)
      }
      http.Error(w, fmt.Sprintf("BackenServer Error,%s", r.URL.String()), http.StatusInternalServerError)
    }
  }()

  var (
    cache_key        = cache.GenKey(r)
    cache_storage    *cache.Storage
    cache_exists     bool
    resp_body        []byte
    err              error
    resp_status_code int
  )

  if r.Header.Get("ACS_RELOAD") == "true" {
    Cache.Remove(cache_key)
  }

  cache_storage, cache_exists = Cache.Get(cache_key)
  if cache_exists {
    for hk, _ := range cache_storage.Response.Header {
      w.Header().Set(hk, cache_storage.Response.Header.Get(hk))
    }
    w.Header().Set("aup", fmt.Sprintf("%s,%d,%d", cache_key, cache_storage.ClientAccessCount, Cache.Size()))
    w.Write(cache_storage.Response.Body)
    return
  }

  resp_body, resp_status_code, _, err = client.HttpRequest(w, r)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  if resp_status_code != 200 {
    return
  }

  cache_storage = &cache.Storage{
    InitAt:             time.Now(),
    UpdatedAt:          time.Now(),
    UpdateDuration:     cache.ENTITY_UPDATE_DURATION,
    Duration:           cache.ENTITY_DURATION,
    ClientLastAccessAt: time.Now(),
    ClientAccessCount:  1,
    CurrentStatus:      cache.STATUS_NORMAL,
    Request:            r,
    Response: &cache.ResponseStorage{
      Header:     w.Header(),
      Body:       resp_body,
      StatusCode: resp_status_code,
    },
  }
  Cache.Set(cache_key, cache_storage)
}
