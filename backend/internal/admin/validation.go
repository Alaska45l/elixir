package admin

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	slugRE     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	skuRE      = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	usernameRE = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
)

const (
	defaultHeroImageMode          = "product_covers"
	defaultHeroRotationIntervalMS = 8000
	minHeroRotationIntervalMS     = 1000
	maxHeroRotationIntervalMS     = 60000
)

func normalizeProductPayload(p *productPayload) error {
	p.Name = strings.TrimSpace(p.Name)
	p.Slug = strings.ToLower(strings.TrimSpace(p.Slug))
	p.Tagline = strings.TrimSpace(p.Tagline)
	p.Description = strings.TrimSpace(p.Description)
	p.ScentFamily = strings.TrimSpace(p.ScentFamily)
	p.GenderTag = strings.TrimSpace(p.GenderTag)
	p.Concentration = strings.TrimSpace(p.Concentration)

	if len(p.Name) < 2 || len(p.Name) > 120 {
		return errors.New("el nombre del producto debe tener entre 2 y 120 caracteres")
	}
	if !slugRE.MatchString(p.Slug) {
		return errors.New("el slug solo puede usar minúsculas, números y guiones")
	}
	if len(p.Tagline) > 180 {
		return errors.New("la frase corta no puede superar 180 caracteres")
	}
	if len(p.Description) > 4000 {
		return errors.New("la descripción no puede superar 4000 caracteres")
	}
	if p.ScentFamily == "" {
		return errors.New("indicá una familia olfativa")
	}
	if p.GenderTag == "" {
		p.GenderTag = "Unisex"
	}
	if p.Concentration == "" {
		p.Concentration = "EDP"
	}
	p.TopNotes = cleanNoteList(p.TopNotes)
	p.HeartNotes = cleanNoteList(p.HeartNotes)
	p.BaseNotes = cleanNoteList(p.BaseNotes)

	if len(p.Variants) == 0 {
		return errors.New("agregá al menos una variante con precio y stock")
	}
	if len(p.Variants) > 20 {
		return errors.New("un producto no puede tener más de 20 variantes")
	}
	seenSKU := map[string]bool{}
	seenSize := map[int]bool{}
	for i := range p.Variants {
		v := &p.Variants[i]
		v.SKU = strings.TrimSpace(v.SKU)
		if v.SKU == "" {
			v.SKU = fmt.Sprintf("%s-%d", p.Slug, v.SizeML)
		}
		if !skuRE.MatchString(v.SKU) || len(v.SKU) > 80 {
			return fmt.Errorf("el SKU %q solo puede usar letras, números, punto, guion o guion bajo", v.SKU)
		}
		if seenSKU[strings.ToUpper(v.SKU)] {
			return fmt.Errorf("el SKU %q está repetido", v.SKU)
		}
		seenSKU[strings.ToUpper(v.SKU)] = true
		if v.SizeML <= 0 || v.SizeML > 10000 {
			return errors.New("cada variante necesita un tamaño válido en ml")
		}
		if seenSize[v.SizeML] {
			return fmt.Errorf("ya existe una variante de %d ml", v.SizeML)
		}
		seenSize[v.SizeML] = true
		if v.PriceARSCents <= 0 || v.PriceARSCents > 100_000_000_000 {
			return errors.New("cada variante necesita un precio válido")
		}
		if v.Stock < 0 || v.Stock > 100000 {
			return errors.New("el stock no puede ser negativo ni excesivo")
		}
		if v.WeightGrams <= 0 {
			v.WeightGrams = 200
		}
		if v.WeightGrams > 50000 {
			return errors.New("el peso debe estar expresado en gramos y ser razonable")
		}
	}

	cleanImages := make([]imageForm, 0, len(p.Images))
	primaryIndex := -1
	for _, img := range p.Images {
		img.URL = strings.TrimSpace(img.URL)
		img.AltText = strings.TrimSpace(img.AltText)
		if img.URL == "" {
			continue
		}
		if !validImageURL(img.URL) {
			return errors.New("cada imagen debe usar una URL pública http o https")
		}
		if len(img.AltText) > 160 {
			return errors.New("el texto alternativo de imagen no puede superar 160 caracteres")
		}
		if img.AltText == "" {
			img.AltText = p.Name
		}
		if img.IsPrimary && primaryIndex == -1 {
			primaryIndex = len(cleanImages)
		}
		img.SortOrder = len(cleanImages)
		cleanImages = append(cleanImages, img)
	}
	if len(cleanImages) == 0 {
		return errors.New("agregá al menos una imagen pública del producto")
	}
	for i := range cleanImages {
		cleanImages[i].IsPrimary = i == primaryIndex
	}
	if primaryIndex == -1 {
		cleanImages[0].IsPrimary = true
	}
	p.Images = cleanImages
	return nil
}

func cleanNoteList(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.Trim(strings.TrimSpace(value), ",")
		if value == "" || len(value) > 40 {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
		if len(out) == 12 {
			break
		}
	}
	return out
}

func validImageURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "https" || parsed.Scheme == "http"
}

func normalizeHomepagePayload(p *homepageRequest) error {
	p.HeroHeading = strings.TrimSpace(p.HeroHeading)
	p.HeroSubheading = strings.TrimSpace(p.HeroSubheading)
	p.HeroImageURL = strings.TrimSpace(p.HeroImageURL)
	p.HeroImageMode = strings.TrimSpace(p.HeroImageMode)
	p.HeroCTALabel = strings.TrimSpace(p.HeroCTALabel)
	p.HeroCTAURL = strings.TrimSpace(p.HeroCTAURL)
	p.EditorialHeading = strings.TrimSpace(p.EditorialHeading)
	p.EditorialBody = strings.TrimSpace(p.EditorialBody)
	p.EditorialImageURL = strings.TrimSpace(p.EditorialImageURL)
	if p.HeroImageMode == "" {
		p.HeroImageMode = defaultHeroImageMode
	}
	if p.HeroImageMode != "static" && p.HeroImageMode != "product_covers" {
		return errors.New("el modo de imagen del hero no es válido")
	}
	if p.HeroRotationIntervalMS <= 0 {
		p.HeroRotationIntervalMS = defaultHeroRotationIntervalMS
	}
	if p.HeroRotationIntervalMS < minHeroRotationIntervalMS || p.HeroRotationIntervalMS > maxHeroRotationIntervalMS {
		return errors.New("el intervalo del hero debe estar entre 1000 y 60000 milisegundos")
	}
	if p.HeroHeading == "" {
		return errors.New("el título principal es obligatorio")
	}
	if err := validateURLField(p.HeroImageURL, "la imagen principal", false); err != nil {
		return err
	}
	if err := validateURLField(p.EditorialImageURL, "la imagen editorial", false); err != nil {
		return err
	}
	if p.HeroCTAURL != "" && !strings.HasPrefix(p.HeroCTAURL, "/") {
		if err := validateURLField(p.HeroCTAURL, "el enlace del botón", true); err != nil {
			return err
		}
	}
	return nil
}

func normalizeSiteSettingsPayload(s *siteSettings) error {
	s.FooterInstagramURL = strings.TrimSpace(s.FooterInstagramURL)
	s.FooterTikTokURL = strings.TrimSpace(s.FooterTikTokURL)
	s.FooterWhatsAppURL = strings.TrimSpace(s.FooterWhatsAppURL)
	s.AnnouncementBarText = strings.TrimSpace(s.AnnouncementBarText)
	s.AboutTitle = strings.TrimSpace(s.AboutTitle)
	s.AboutDescription = strings.TrimSpace(s.AboutDescription)
	s.AboutLocation = strings.TrimSpace(s.AboutLocation)
	s.AboutPhone = strings.TrimSpace(s.AboutPhone)
	s.ReturnPolicyHTML = strings.TrimSpace(s.ReturnPolicyHTML)
	if err := validateURLField(s.FooterInstagramURL, "Instagram", false); err != nil {
		return err
	}
	if err := validateURLField(s.FooterTikTokURL, "TikTok", false); err != nil {
		return err
	}
	if err := validateURLField(s.FooterWhatsAppURL, "WhatsApp", false); err != nil {
		return err
	}
	if s.LowStockThreshold <= 0 {
		s.LowStockThreshold = 5
	}
	if s.LowStockThreshold > 1000 {
		return errors.New("el aviso de stock bajo debe ser menor a 1000")
	}
	s.FAQItems = cleanFAQItems(s.FAQItems)
	s.NavbarProductCategories = cleanNavItems(s.NavbarProductCategories)
	return nil
}

func validateURLField(raw, label string, required bool) error {
	if raw == "" {
		if required {
			return fmt.Errorf("%s es obligatorio", label)
		}
		return nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return fmt.Errorf("%s debe ser un enlace http o https válido", label)
	}
	return nil
}

func cleanFAQItems(values []faqItem) []faqItem {
	out := make([]faqItem, 0, len(values))
	for _, item := range values {
		item.Question = strings.TrimSpace(item.Question)
		item.Answer = strings.TrimSpace(item.Answer)
		if item.Question == "" && item.Answer == "" {
			continue
		}
		out = append(out, item)
		if len(out) == 20 {
			break
		}
	}
	return out
}

func cleanNavItems(values []navItem) []navItem {
	out := make([]navItem, 0, len(values))
	for _, item := range values {
		item.Label = strings.TrimSpace(item.Label)
		item.Href = strings.TrimSpace(item.Href)
		if item.Label == "" && item.Href == "" {
			continue
		}
		if item.Href == "" {
			item.Href = "/fragrances"
		}
		out = append(out, item)
		if len(out) == 12 {
			break
		}
	}
	return out
}

func normalizeDiscountPayload(req *discountWriteRequest, creating bool) error {
	req.Code = CleanCode(req.Code)
	req.DiscountType = strings.TrimSpace(req.DiscountType)
	if creating && req.Code == "" {
		return errors.New("el código es obligatorio")
	}
	if req.Code != "" && (!skuRE.MatchString(req.Code) || len(req.Code) > 40) {
		return errors.New("el código solo puede usar letras, números, punto, guion o guion bajo")
	}
	switch req.DiscountType {
	case "percent":
		if req.DiscountValue <= 0 || req.DiscountValue > 100 {
			return errors.New("el descuento porcentual debe estar entre 1 y 100")
		}
	case "fixed":
		if req.DiscountValue <= 0 {
			return errors.New("el descuento fijo debe ser mayor a cero")
		}
	default:
		return errors.New("elegí descuento por porcentaje o monto fijo")
	}
	if req.MinOrderCents < 0 {
		return errors.New("el mínimo de compra no puede ser negativo")
	}
	if req.MaxUses != nil && *req.MaxUses <= 0 {
		return errors.New("el máximo de usos debe ser mayor a cero")
	}
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now().Add(-time.Minute)) {
		return errors.New("la fecha de vencimiento no puede estar en el pasado")
	}
	return nil
}

func validateOrderStatus(status string) bool {
	switch status {
	case "pending", "paid", "failed", "cancelled", "shipped", "delivered":
		return true
	default:
		return false
	}
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 40 {
		return errors.New("el usuario debe tener entre 3 y 40 caracteres")
	}
	if !usernameRE.MatchString(username) {
		return errors.New("el usuario solo puede usar letras, números, punto, guion o guion bajo")
	}
	return nil
}

func validatePassword(password, username string) error {
	if len(password) < 10 {
		return errors.New("la contraseña debe tener al menos 10 caracteres")
	}
	if len(password) > 200 {
		return errors.New("la contraseña es demasiado larga")
	}
	lower := strings.ToLower(password)
	if username != "" && strings.Contains(lower, strings.ToLower(username)) {
		return errors.New("la contraseña no debe incluir el nombre de usuario")
	}
	hasLetter := false
	hasDigit := false
	for _, ch := range password {
		if ch >= '0' && ch <= '9' {
			hasDigit = true
		}
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			hasLetter = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("la contraseña debe combinar letras y números")
	}
	return nil
}

func validateLoginPayload(username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if password == "" {
		return errors.New("la contraseña es obligatoria")
	}
	if len(password) > 200 {
		return errors.New("la contraseña es demasiado larga")
	}
	return nil
}

func auditActor(username string) string {
	username = strings.TrimSpace(username)
	if username == "" {
		return "unknown"
	}
	return username
}
