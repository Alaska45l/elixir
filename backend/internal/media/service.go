package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

const maxImagePixels = 24_000_000

var ErrInvalidImage = errors.New("media: imagen inválida")

type StorageService struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewStorageService(accountID, accessKey, secretKey, bucketName, publicURL string) (*StorageService, error) {
	accountID = strings.TrimSpace(accountID)
	accessKey = strings.TrimSpace(accessKey)
	secretKey = strings.TrimSpace(secretKey)
	bucketName = strings.TrimSpace(bucketName)
	publicURL = strings.TrimRight(strings.TrimSpace(publicURL), "/")

	values := []string{accountID, accessKey, secretKey, bucketName, publicURL}
	configured := 0
	for _, value := range values {
		if value != "" {
			configured++
		}
	}
	if configured == 0 {
		return nil, nil
	}
	if configured != len(values) {
		return nil, errors.New("media: configuración R2 incompleta")
	}

	cfg, err := awscfg.LoadDefaultConfig(context.Background(),
		awscfg.WithRegion("auto"),
		awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("media: aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
		o.UsePathStyle = true
	})

	return &StorageService{client: client, bucketName: bucketName, publicURL: publicURL}, nil
}

func (s *StorageService) IsConfigured() bool {
	return s != nil && s.client != nil && s.bucketName != "" && s.publicURL != ""
}

func (s *StorageService) ProcessAndUpload(ctx context.Context, src io.Reader, folder string) (string, error) {
	if !s.IsConfigured() {
		return "", errors.New("media: storage no configurado")
	}

	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("media: read: %w", err)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("%w: decode config: %v", ErrInvalidImage, err)
	}
	if format != "jpeg" && format != "png" && format != "webp" {
		return "", fmt.Errorf("%w: formato no soportado: %s", ErrInvalidImage, format)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxImagePixels/cfg.Height {
		return "", fmt.Errorf("%w: dimensiones inválidas o excesivas", ErrInvalidImage)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("%w: decode: %v", ErrInvalidImage, err)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 80}); err != nil {
		return "", fmt.Errorf("media: webp encode: %w", err)
	}

	folder = sanitizeFolder(folder)
	key := fmt.Sprintf("%s/%s.webp", folder, uuid.NewString())

	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.bucketName),
		Key:          aws.String(key),
		Body:         bytes.NewReader(buf.Bytes()),
		ContentType:  aws.String("image/webp"),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	}); err != nil {
		return "", fmt.Errorf("media: r2 put: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func sanitizeFolder(folder string) string {
	folder = strings.Trim(strings.ReplaceAll(folder, "\\", "/"), "/")
	if folder == "" {
		return "products"
	}

	parts := strings.Split(folder, "/")
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		part = sanitizeSegment(part)
		if part != "" && part != "." && part != ".." {
			clean = append(clean, part)
		}
	}
	if len(clean) == 0 {
		return "products"
	}
	return strings.Join(clean, "/")
}

func sanitizeSegment(segment string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(segment)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
