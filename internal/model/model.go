package model

type ResizedImage struct {
	Data   []byte
	Format string
}

type ImageData struct {
	Data        []byte
	ContentType string
	URL         string
	Width       int
	Height      int
}
