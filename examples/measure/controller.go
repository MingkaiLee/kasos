package main

import (
	"crypto/md5"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomHash(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		s := randomString(1000)
		md5Hasher := md5.New()
		_, err := io.WriteString(md5Hasher, s)
		if err != nil {
			LogErrorf("error: %v", err)
		}
	}
	w.WriteHeader(http.StatusOK)
	diff := float64(time.Since(startTime).Milliseconds()) / float64(1000)
	LatencyMetric.WithLabelValues("measure").Observe(diff)
	ServiceQPSMetric.WithLabelValues("on", "measure").Inc()
}

func randomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return strings.TrimRight(string(b), "\x00") // 去除末尾可能出现的空字符
}
