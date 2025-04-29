package helper

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
	"golang.org/x/sync/errgroup"
)

// BatchRequest 并发发送命令到 lgproxy 节点，并返回按输入顺序的响应和错误信息
func BatchRequest(
	parentCtx context.Context,
	serverIDs []string,
	endpoint, command string,
) ([]string, []error) {
	g, ctx := errgroup.WithContext(parentCtx)
	responses := make([]string, len(serverIDs))
	errors := make([]error, len(serverIDs))
	cfg := GetConfig()
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.singleTimeout) * time.Second,
	}

	for idx, serverID := range serverIDs {
		idx, serverID := idx, serverID

		g.Go(func() error {
			// 校验节点合法性
			if !slices.Contains(cfg.servers, serverID) {
				responses[idx] = "request failed: invalid server"
				return nil
			}

			// 构造请求URL
			scheme := "http"
			if source.AppConfig.LG.SSL {
				scheme = "https"
			}

			serverAddr := serverID + "." + source.AppConfig.LG.Domain
			u := url.URL{
				Scheme:   scheme,
				Host:     serverAddr + ":" + strconv.Itoa(source.AppConfig.LG.ProxyPort),
				Path:     endpoint,
				RawQuery: "q=" + url.QueryEscape(command),
			}

			// 发起请求
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				errors[idx] = err
				return nil
			}

			resp, err := client.Do(req)
			if err != nil {
				errors[idx] = err
				return nil
			}
			defer resp.Body.Close()

			// 读取响应
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errors[idx] = err
				return nil
			}

			if len(body) == 0 {
				responses[idx] = "node returned empty response, please refresh to try again."
			} else {
				responses[idx] = string(body)
			}

			return nil
		})
	}

	g.Wait() // 等待所有请求完成
	return responses, errors
}
