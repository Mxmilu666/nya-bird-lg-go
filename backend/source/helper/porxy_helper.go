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

// BatchRequest 并发发送命令到 lgproxy 节点，并返回按输入顺序的响应
func BatchRequest(
	parentCtx context.Context,
	serverIDs []string,
	endpoint, command string,
) ([]string, error) {
	g, ctx := errgroup.WithContext(parentCtx)
	responses := make([]string, len(serverIDs))
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

			// 构造完整的服务器地址
			serverAddr := serverID + "." + source.AppConfig.LG.Domain
			// 构造 URL
			u := url.URL{
				Scheme: "http",
				Host:   serverAddr + ":" + strconv.Itoa(source.AppConfig.LG.ProxyPort),
				Path:   endpoint,
			}
			q := u.Query()
			q.Set("q", command)
			u.RawQuery = q.Encode()

			// 发起请求
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// 读取全部
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if len(body) == 0 {
				responses[idx] = "node returned empty response, please refresh to try again."
			} else {
				responses[idx] = string(body)
			}
			return nil
		})
	}

	// 等待所有 goroutine 完成或出错
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return responses, nil
}
