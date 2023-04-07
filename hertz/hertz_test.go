package hertz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/transport/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
)

func EncodeResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	data := map[string]interface{}{
		"code": 200,
		"msg":  "succeed",
		"data": v,
	}
	return http.DefaultResponseEncoder(w, r, data)
}

func ErrEncoder(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{
		"code": -1,
		"msg":  err.Error(),
	}
	if ierr, ok := err.(interface{ Code() int }); ok {
		data["code"] = ierr.Code()
	}
	if ierr, ok := err.(interface{ Msg() string }); ok {
		data["msg"] = ierr.Msg()
	}
	body, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "json")
	_, _ = w.Write(body)
}

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(func(s *Server) {
		s.addr = "127.0.0.1:8001"
		s.ene = ErrEncoder
		s.enc = EncodeResponse
	})

	r := srv.Group("/api")
	r.GET("/hello", func(ctx context.Context, c *app.RequestContext) {
		rsp := struct {
			Uid      uint64 `json:"uid"'`
			Nickname string `json:"nickname"`
			Gender   int8   `json:"gender"`
		}{
			Uid:      100002323,
			Nickname: "jager",
			Gender:   2,
		}
		c.JSON(200, &rsp)
	})

	http.NewServer(func(s *http.Server) {

	})

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
}
