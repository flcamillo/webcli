package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// arquivos html do servidor web
//
//go:embed root/templates/*
var content embed.FS

// Define a estrutura de um alerta para o usuário.
type Alert struct {
	Title   string
	Message string
}

// Define a estrutura para o servidor web.
type Server struct {
	config     *Config
	templates  *template.Template
	httpServer *http.Server
}

// Cria uma instancia do servidor web.
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

// Incia o servidor.
func (p *Server) Run(ctx context.Context) error {
	templates, err := template.ParseFS(content, "root/templates/*.html")
	if err != nil {
		return err
	}
	p.templates = templates
	router := http.NewServeMux()
	router.HandleFunc("/", p.handleHome)
	p.httpServer = &http.Server{
		Addr:         p.config.ServerAddress,
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      router,
	}
	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Go(func() {
		errChan <- p.httpServer.ListenAndServe()
	})
	defer wg.Wait()
	select {
	case <-ctx.Done():
		p.httpServer.Shutdown(context.Background())
	case err := <-errChan:
		if err != nil {
			return err
		}
	}
	return nil
}

// Função para auxiliar a adicionar alertas no objeto de dados que será retornado ao usuário.
func (p *Server) addAlert(data map[string]interface{}, group string, alert *Alert) {
	alerts, ok := data[group]
	if !ok {
		data[group] = []*Alert{alert}
	} else {
		alerts = append(alerts.([]*Alert), alert)
	}
}

// Retorna os dados básicos.
func (p *Server) loadDefaultData(r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	// adiciona o arquivo de configuração
	data["Config"] = p.config
	// adiciona a variavel de ambiente em uso
	environment := r.FormValue("environment")
	if environment == "" {
		p.addAlert(data, "Alerts", &Alert{Title: "Variável ausente", Message: "Ambiente não selecionado."})
		return data
	}
	found := false
	for k := range p.config.Environments {
		if k == environment {
			found = true
			break
		}
	}
	if !found {
		p.addAlert(data, "Errors", &Alert{Title: "Variável inválida", Message: "Ambiente inválido informado."})
		return data
	}
	data["Environment"] = environment
	return data
}

// Processa requisições para a pagina principal.
func (p *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	data := p.loadDefaultData(r)
	defer func() {
		err := p.templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Printf("failed to parse template, %s", err)
		}
	}()
}
