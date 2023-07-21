package app

import (
	"app/downloader"
	"app/extractors"
	_ "app/extractors/bilibili"
	"app/request"
	"encoding/json"
	"os"
	"reflect"
)

type (
	Context struct {
		Playlist         bool
		Items            string
		Start            uint
		End              uint
		EpisodeTitleOnly bool
		Json             bool
		Cookie           string
		Silent           bool
		Info             bool
		StreamFormat     string
		Refer            string
		Debug            bool
		OutputPath       string
		OutputName       string
		FileNameLength   uint
		Caption          bool
		MultiThread      bool
		Thread           uint
		Retry            uint
		UserAgent        string
		ChunkSize        uint
	}
	Engine struct {
		ctx       *Context
		requester *request.Requester
	}
)

var DefaultContext = &Context{
	FileNameLength: 255,
	Start:          1,
	End:            0,
	Retry:          10,
	ChunkSize:      1,
	Thread:         10,
}

func New(ctx *Context) *Engine {
	engine := &Engine{ctx: ctx}
	engine.requester = request.New(request.Options{
		RetryTimes: int(engine.ctx.Retry),
		Cookie:     engine.ctx.Cookie,
		UserAgent:  engine.ctx.UserAgent,
		Refer:      engine.ctx.Refer,
		Debug:      engine.ctx.Debug,
		Silent:     engine.ctx.Silent,
	})
	return engine
}

func Default(userCtx ...*Context) *Engine {
	ctx := &Context{}
	*ctx = *DefaultContext

	if len(userCtx) > 0 {
		userCtxVal := reflect.ValueOf(userCtx[0]).Elem()
		defaultCtxVal := reflect.ValueOf(ctx).Elem()

		for i := 0; i < userCtxVal.NumField(); i++ {
			userField := userCtxVal.Field(i)
			defaultField := defaultCtxVal.Field(i)

			if !userField.IsZero() {
				defaultField.Set(userField)
			}
		}
	}

	return New(ctx)
}

func (engine *Engine) Download(videoURL string) error {
	data, err := extractors.Extract(videoURL, extractors.Options{
		Playlist:         engine.ctx.Playlist,
		Items:            engine.ctx.Items,
		ItemStart:        int(engine.ctx.Start),
		ItemEnd:          int(engine.ctx.End),
		ThreadNumber:     int(engine.ctx.Thread),
		EpisodeTitleOnly: engine.ctx.EpisodeTitleOnly,
		Cookie:           engine.ctx.Cookie,
	}, engine.requester)
	if err != nil {
		// if this error occurs, it means that an error occurred before actually starting to extract data
		// (there is an error in the preparation step), and the data list is empty.
		return err
	}

	if engine.ctx.Json {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "\t")
		e.SetEscapeHTML(false)
		if err := e.Encode(data); err != nil {
			return err
		}

		return nil
	}

	defaultDownloader := downloader.New(downloader.Options{
		Silent:         engine.ctx.Silent,
		InfoOnly:       engine.ctx.Info,
		Stream:         engine.ctx.StreamFormat,
		Refer:          engine.ctx.Refer,
		OutputPath:     engine.ctx.OutputPath,
		OutputName:     engine.ctx.OutputName,
		FileNameLength: int(engine.ctx.FileNameLength),
		Caption:        engine.ctx.Caption,
		MultiThread:    engine.ctx.MultiThread,
		ThreadNumber:   int(engine.ctx.Thread),
		RetryTimes:     int(engine.ctx.Retry),
		ChunkSizeMB:    int(engine.ctx.ChunkSize),
	}, engine.requester)

	errors := make([]error, 0)
	for _, item := range data {
		if item.Err != nil {
			// if this error occurs, the preparation step is normal, but the data extraction is wrong.
			// the data is an empty struct.
			errors = append(errors, item.Err)
			continue
		}
		if err = defaultDownloader.Download(item); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		return errors[0]
	}
	return nil
}
