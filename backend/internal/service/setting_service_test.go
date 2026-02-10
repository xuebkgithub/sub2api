//go:build unit

package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// settingRepoTestStub 用于测试的 SettingRepository stub
type settingRepoTestStub struct {
	values map[string]string
}

func (s *settingRepoTestStub) Get(ctx context.Context, key string) (*Setting, error) {
	if v, ok := s.values[key]; ok {
		return &Setting{Key: key, Value: v}, nil
	}
	return nil, ErrSettingNotFound
}

func (s *settingRepoTestStub) GetValue(ctx context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
	}
	return "", ErrSettingNotFound
}

func (s *settingRepoTestStub) Set(ctx context.Context, key, value string) error {
	s.values[key] = value
	return nil
}

func (s *settingRepoTestStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, key := range keys {
		if v, ok := s.values[key]; ok {
			result[key] = v
		}
	}
	return result, nil
}

func (s *settingRepoTestStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	for k, v := range settings {
		s.values[k] = v
	}
	return nil
}

func (s *settingRepoTestStub) GetAll(ctx context.Context) (map[string]string, error) {
	return s.values, nil
}

func (s *settingRepoTestStub) Delete(ctx context.Context, key string) error {
	delete(s.values, key)
	return nil
}

// ldapConfigRepoTestStub 用于测试的 LdapConfigRepository stub
type ldapConfigRepoTestStub struct {
	enabled bool
}

func (s *ldapConfigRepoTestStub) Get(ctx context.Context) (*LdapConfig, error) {
	if s.enabled {
		return &LdapConfig{Enabled: true}, nil
	}
	return nil, ErrLdapConfigNotFound
}

func (s *ldapConfigRepoTestStub) GetEnabled(ctx context.Context) (*LdapConfig, error) {
	if s.enabled {
		return &LdapConfig{Enabled: true}, nil
	}
	return nil, ErrLdapConfigNotFound
}

func (s *ldapConfigRepoTestStub) Create(ctx context.Context, config *LdapConfig) error {
	return nil
}

func (s *ldapConfigRepoTestStub) Update(ctx context.Context, config *LdapConfig) error {
	return nil
}

func (s *ldapConfigRepoTestStub) Exists(ctx context.Context) (bool, error) {
	return s.enabled, nil
}

// setupSettingService 创建测试用的 SettingService
func setupSettingService(t *testing.T) *SettingService {
	t.Helper()

	// 创建测试配置
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireHour: 1,
		},
	}

	// 创建测试数据
	settingRepo := &settingRepoTestStub{
		values: map[string]string{
			SettingKeyRegistrationEnabled:         "true",
			SettingKeyEmailVerifyEnabled:          "false",
			SettingKeyTotpEnabled:                 "true",
			SettingKeyTurnstileEnabled:            "true",
			SettingKeyTurnstileSiteKey:            "test-site-key",
			SettingKeySiteName:                    "Test Site",
			SettingKeyLinuxDoConnectEnabled:       "true",
			SettingKeyPromoCodeEnabled:            "false",
			SettingKeyPasswordResetEnabled:        "true",
			SettingKeyInvitationCodeEnabled:       "false",
			SettingKeyHideCcsImportButton:         "false",
			SettingKeyPurchaseSubscriptionEnabled: "false",
		},
	}

	ldapConfigRepo := &ldapConfigRepoTestStub{
		enabled: true, // LDAP 启用
	}

	return NewSettingService(settingRepo, ldapConfigRepo, cfg)
}

// TestPublicSettingsConsistency 测试配置注入与 API 配置的一致性
func TestPublicSettingsConsistency(t *testing.T) {
	ctx := context.Background()
	service := setupSettingService(t)

	// 获取 API 配置
	apiSettings, err := service.GetPublicSettings(ctx)
	require.NoError(t, err)

	// 获取注入配置
	injectedSettings, err := service.GetPublicSettingsForInjection(ctx)
	require.NoError(t, err)

	// 将注入配置转换为 map
	injectedMap := make(map[string]interface{})
	jsonBytes, err := json.Marshal(injectedSettings)
	require.NoError(t, err)
	err = json.Unmarshal(jsonBytes, &injectedMap)
	require.NoError(t, err)

	// 比较关键字段
	assert.Equal(t, apiSettings.LdapEnabled, injectedMap["ldap_enabled"], "ldap_enabled should match")
	assert.Equal(t, apiSettings.TotpEnabled, injectedMap["totp_enabled"], "totp_enabled should match")
	assert.Equal(t, apiSettings.TurnstileEnabled, injectedMap["turnstile_enabled"], "turnstile_enabled should match")
	assert.Equal(t, apiSettings.RegistrationEnabled, injectedMap["registration_enabled"], "registration_enabled should match")
	assert.Equal(t, apiSettings.EmailVerifyEnabled, injectedMap["email_verify_enabled"], "email_verify_enabled should match")
	assert.Equal(t, apiSettings.LinuxDoOAuthEnabled, injectedMap["linuxdo_oauth_enabled"], "linuxdo_oauth_enabled should match")
}
