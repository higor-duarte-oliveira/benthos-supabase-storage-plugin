package output

import (
	"context"
	"fmt"
	"os"

	"github.com/redpanda-data/benthos/v4/public/service"
	"github.com/supabase-community/supabase-go"
)

func supabaseOutputConfigSpec() *service.ConfigSpec {
	x := service.NewConfigSpec().
		Summary("Executes a request in the supabase to write a query").
		Categories("Services").
		Field(service.NewStringField("host").Default("localhost:8000")).
		Field(service.NewStringField("apiKey").Default("")).
		Field(service.NewStringField("file").Default(""))
	return x
}

func init() {
	err := service.RegisterOutput(
		"supabase_stdout", supabaseOutputConfigSpec(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (out service.Output, mif int, err error) {
			return newSupabaseOutput(conf)
		})
	if err != nil {
		panic(err)
	}
}

func newSupabaseOutput(conf *service.ParsedConfig) (
	out service.Output,
	mif int,
	err error) {
	
	mif = 1
	
	host, err := conf.FieldString("host")
	if err != nil {
		return
	}

	apiKey, err := conf.FieldString("apiKey")
	if err != nil {
		return
	}

	file, err := conf.FieldString("file")
	if err != nil {
		return
	}
	
	out = &supabaseOutput{host: host, apiKey: apiKey, file: file}
	return
}

//------------------------------------------------------------------------------

type supabaseOutput struct{
	host string
	apiKey string
	file string
	client *supabase.Client
}

func (s *supabaseOutput) Connect(ctx context.Context) error {
	var err error

	s.client, err = supabase.NewClient(s.host, s.apiKey, &supabase.ClientOptions{})

	if err != nil {
	  fmt.Println("cannot initalize client", err)
	}

	return nil
}

func (s *supabaseOutput) Write(ctx context.Context, msg *service.Message) error {
	if s.client == nil {
		return service.ErrNotConnected
	}

	file, err := os.Open(s.file)

    if err != nil {
        panic(err)
    }

	result, err := s.client.Storage.UploadFile("test", file.Name(), file)

	if err != nil {
		return err
	}

	fmt.Println(result)
	file.Close()

	return nil
}

func (s *supabaseOutput) Close(ctx context.Context) error {
	return nil
}
