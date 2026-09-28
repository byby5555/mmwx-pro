package license

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type PlanInfo struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description,omitempty"`
	MaxServers  int      `json:"max_servers"`
	MaxNodes    int      `json:"max_nodes"`
	MaxUsers    int      `json:"max_users"`
	Features    []string `json:"features"`
	FeatureTokens map[string]string `json:"feature_tokens,omitempty"`
}

type Status struct {
	Valid      bool      `json:"valid"`
	Error      string    `json:"error,omitempty"`
	MaxServers int       `json:"max_servers"`
	ExpiresAt  string    `json:"expires_at,omitempty"`
	Plan       *PlanInfo `json:"plan,omitempty"`
	LastCheck  time.Time `json:"last_check"`
	HardRevoked bool    `json:"hard_revoked,omitempty"`
}

func (s *Status) HasFeature(_ string) bool {
	return true
}

var defaultStatus = Status{
	Valid:      true,
	MaxServers: 9999,
	Plan: &PlanInfo{
		Name:        "PRO",
		DisplayName: "专业版(已解锁)",
		MaxServers:  9999,
		MaxNodes:    99999,
		MaxUsers:    99999,
		Features: []string{
			"node_speed_test",
			"node_rate_limit",
			"limiter",
			"server_share",
			"embed_xray",
			"reality_domain_pool",
			"reality_pool",
			"speed_test",
			"rate_limit",
			"embedded",
			"custom_branding",
			"premium_theme",
		},
		FeatureTokens: map[string]string{
			"node_speed_test":     "bypassed",
			"node_rate_limit":     "bypassed",
			"limiter":             "bypassed",
			"server_share":        "bypassed",
			"embed_xray":          "bypassed",
			"reality_domain_pool": "bypassed",
			"reality_pool":        "bypassed",
			"speed_test":          "bypassed",
			"rate_limit":          "bypassed",
			"embedded":            "bypassed",
			"custom_branding":     "bypassed",
			"premium_theme":       "bypassed",
		},
	},
}

type SettingsGetter interface {
	GetSystemSetting(ctx context.Context, key string) (string, error)
}

type SettingsStore interface {
	GetSystemSetting(ctx context.Context, key string) (string, error)
	SetSystemSetting(ctx context.Context, key, value string) error
}

type UsageReporter interface {
	LicenseUsage(ctx context.Context) (servers, nodes, users int, err error)
}

type Manager struct {
	mu        sync.RWMutex
	status    Status
	serverURL string
	key       string
	machineID string
	settings  SettingsStore
	usage     UsageReporter
	client    *http.Client
	cancel    context.CancelFunc
	onRecover func()
	onQuotaChange func()
}

func (m *Manager) SetOnRecover(cb func()) {
	m.mu.Lock()
	m.onRecover = cb
	m.mu.Unlock()
}

func (m *Manager) SetOnQuotaChange(cb func()) {
	m.mu.Lock()
	m.onQuotaChange = cb
	m.mu.Unlock()
}

func (m *Manager) SetUsageReporter(r UsageReporter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usage = r
}

const DefaultServerURL = "https://license.miaomiaowux.com"

func NewManager(settings SettingsStore, machineID string) *Manager {
	return &Manager{
		status:    defaultStatus,
		serverURL: DefaultServerURL,
		machineID: machineID,
		settings:  settings,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *Manager) Start(ctx context.Context) {
	ctx, m.cancel = context.WithCancel(ctx)
	log.Printf("[license] 许可证验证已绕过，所有 PRO 功能已解锁")
	m.mu.Lock()
	m.status = defaultStatus
	m.mu.Unlock()
}

func (m *Manager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *Manager) GetStatus() Status {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

func (m *Manager) IsValid() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return true
}

func (m *Manager) isValidLocked() bool {
	return true
}

func (m *Manager) QuotaEnforced() bool {
	return false
}

func (m *Manager) EffectiveServerQuota() int {
	return 99999
}

func (m *Manager) HasFeature(name string) bool {
	return true
}

func (m *Manager) Refresh(ctx context.Context) {}

func (m *Manager) withinGracePeriod() bool {
	return true
}

func (m *Manager) loadSettings(ctx context.Context) {
	if m.settings == nil {
		return
	}
	if key, err := m.settings.GetSystemSetting(ctx, "license_key"); err == nil && key != "" {
		m.key = key
	}
	if url, err := m.settings.GetSystemSetting(ctx, "license_server_url"); err == nil && url != "" {
		m.serverURL = url
	}
}

func (m *Manager) loadPersistedStatus(ctx context.Context) {
	if m.settings == nil {
		return
	}
	raw, err := m.settings.GetSystemSetting(ctx, "license_status")
	if err != nil || raw == "" {
		return
	}
	var status Status
	if err := json.Unmarshal([]byte(raw), &status); err != nil {
		log.Printf("[license] failed to load persisted status: %v", err)
		return
	}
	m.mu.Lock()
	m.status = status
	m.mu.Unlock()
	log.Printf("[license] restored status from database: valid=%v plan=%s", status.Valid, status.Plan.Name)
}

func (m *Manager) persistStatus(ctx context.Context) {
	if m.settings == nil {
		return
	}
	m.mu.RLock()
	data, err := json.Marshal(m.status)
	m.mu.RUnlock()
	if err != nil {
		return
	}
	if err := m.settings.SetSystemSetting(ctx, "license_status", string(data)); err != nil {
		log.Printf("[license] failed to persist status: %v", err)
	}
}

func (m *Manager) activate(ctx context.Context) {}

func (m *Manager) heartbeat(ctx context.Context) {}

func (m *Manager) parseResponse(ctx context.Context, resp *http.Response, nonce string) {
	if resp.StatusCode != http.StatusOK {
		log.Printf("[license] non-200 from server: %d (treated as transient, grace in effect)", resp.StatusCode)
		return
	}
	var result struct {
		Valid      bool            `json:"valid"`
		Error      string          `json:"error,omitempty"`
		MaxServers int             `json:"max_servers"`
		ExpiresAt  string          `json:"expires_at,omitempty"`
		Plan       json.RawMessage `json:"plan,omitempty"`
		Sig        string          `json:"sig,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[license] parse response error: %v", err)
		return
	}
	var plan PlanInfo
	var hasPlan bool
	if result.Valid && result.Plan != nil {
		if err := json.Unmarshal(result.Plan, &plan); err == nil {
			if plan.Features == nil {
				var raw struct {
					Features json.RawMessage `json:"features"`
				}
				_ = json.Unmarshal(result.Plan, &raw)
				if raw.Features != nil {
					_ = json.Unmarshal(raw.Features, &plan.Features)
				}
			}
			hasPlan = true
		}
	}
	if result.Valid {
		if !verifyLicenseSig(nonce, m.machineID, true, result.MaxServers, result.ExpiresAt, plan.Features, result.Sig) {
			log.Printf("[license] response signature verification FAILED — ignoring response (possible forged license server / MITM)")
			return
		}
	}
	wasValid := m.IsValid()
	oldQuota := m.EffectiveServerQuota()
	m.mu.Lock()
	m.status.Valid = result.Valid
	m.status.Error = result.Error
	m.status.LastCheck = time.Now()
	if result.Valid {
		m.status.HardRevoked = false
		m.status.MaxServers = result.MaxServers
		m.status.ExpiresAt = result.ExpiresAt
		if hasPlan {
			m.status.Plan = &plan
		}
	} else {
		m.status.HardRevoked = true
		log.Printf("[license] HARD REVOKED by server: %s", result.Error)
	}
	cb := m.onRecover
	quotaCb := m.onQuotaChange
	m.mu.Unlock()
	m.persistStatus(ctx)
	if result.Valid && !wasValid && cb != nil {
		go cb()
	}
	if quotaCb != nil && m.EffectiveServerQuota() != oldQuota {
		go quotaCb()
	}
}

func (m *Manager) heartbeatLoop(ctx context.Context) {
	<-ctx.Done()
}

func (m *Manager) StatusForAgent() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := map[string]any{
		"valid":       true,
		"max_servers": m.status.MaxServers,
	}
	if m.status.ExpiresAt != "" {
		result["expires_at"] = m.status.ExpiresAt
	}
	if m.status.Plan != nil {
		result["plan"] = map[string]any{
			"name":         m.status.Plan.Name,
			"display_name": m.status.Plan.DisplayName,
			"description":  m.status.Plan.Description,
			"max_servers":  m.status.Plan.MaxServers,
			"max_nodes":    m.status.Plan.MaxNodes,
			"max_users":    m.status.Plan.MaxUsers,
			"features":     m.status.Plan.Features,
		}
	}
	return result
}

func GetMachineID() string {
	return persistentMachineID()
}