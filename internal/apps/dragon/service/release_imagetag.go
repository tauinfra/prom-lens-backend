package service

import (
	"fmt"
	"time"
	"valyria-backend/internal/apps/dragon/model"

	"github.com/google/uuid"
)

type Image struct {
	Registry string
	Name     string
	Tag      string
}

func GenerateImage(p *model.Pipeline, ref, commit string) Image {
	ts := time.Now().Format("20060102150405")
	return Image{
		Registry: fmt.Sprintf(
			"%s/%s",
			p.Environment.Harbor.BaseURL,
			p.Environment.Project.Name,
		),
		Name: p.Name,
		Tag:  fmt.Sprintf("%s-%s-%s", ts, ref, commit),
	}
}

func GenerateTaskID(p *model.Pipeline) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		p.Environment.Project.Name,
		p.Environment.Name,
		p.Name,
		uuid.NewString()[:8],
	)
}
