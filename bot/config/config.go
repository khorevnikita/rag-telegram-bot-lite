package config

import (
	"gopkg.in/yaml.v3"
	"gorag-telegram-bot/models"
	"log"
	"os"
)

// Config структура для хранения конфигураций приложения
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type LikesModule struct {
	Enabled         bool   `yaml:"enabled"`
	LikeLabel       string `yaml:"like_label"`
	LikeResponse    string `yaml:"like_response"`
	DislikeLabel    string `yaml:"dislike_label"`
	DislikeResponse string `yaml:"dislike_response"`
}

type FilesModule struct {
	Enabled bool `yaml:"enabled"`
}

type AudioModule struct {
	Enabled bool `yaml:"enabled"`
}

type QuestionOptionConfig struct {
	Text                  string `yaml:"text"`
	RequireAdditionalText bool   `yaml:"require_additional_text,omitempty"`
}

type QuestionConfig struct {
	Text                   string                 `yaml:"text"`
	Order                  int                    `yaml:"order"`
	IsRequired             bool                   `yaml:"is_required"`
	Type                   models.QuestionType    `yaml:"type"`
	SelectableOptionsCount int                    `yaml:"selectable_options_count,omitempty"`
	Hint                   *string                `yaml:"hint,omitempty"`
	Options                []QuestionOptionConfig `yaml:"options,omitempty"`
}

type FormModule struct {
	Enabled               bool             `yaml:"enabled"`
	ShowOnStart           bool             `yaml:"show_on_start"`
	CanSkip               bool             `yaml:"can_skip"`
	DisclaimerText        string           `yaml:"disclaimer_label"`
	StartLabel            string           `yaml:"start_label"`
	LaterLabel            string           `yaml:"later_label"`
	LaterMessage          string           `yaml:"later_message"`
	AllowEdit             bool             `yaml:"allow_edit"`
	EditLabel             string           `yaml:"edit_label"`
	SelectQuestionMessage string           `yaml:"select_question_message"`
	ViewLabel             string           `yaml:"view_label"`
	CompletedMessage      string           `yaml:"completed_message"`
	WrongOptionMessage    string           `yaml:"wrong_option_message"`
	CustomOptionMessage   string           `yaml:"custom_option_message"`
	NextQuestionLabel     string           `yaml:"next_question_label"`
	MoreOptionLabel       string           `yaml:"more_option_label"`
	OptionSavedMessage    string           `yaml:"option_saved_message"`
	AddOptionMessage      string           `yaml:"add_option_message"`
	Questions             []QuestionConfig `yaml:"questions"`
	ContextPrefix         string           `yaml:"context_prefix"`
}

type ModuleConfig struct {
	Likes   LikesModule   `yaml:"likes"`
	Files   FilesModule   `yaml:"files"`
	Audio   AudioModule   `yaml:"audio"`
	Form    FormModule    `yaml:"form"`
	Billing BillingConfig `yaml:"billing"`
}

type CommandAction struct {
	ActUnique string `yaml:"act_unique"`
	ActData   string `yaml:"act_data"`
	Label     string `yaml:"label"`
	Prompt    string `yaml:"prompt"`
}

type CommandConfig struct {
	Name       string          `yaml:"name"`
	Label      string          `yaml:"label"`
	Enabled    bool            `yaml:"enabled"`
	ShowInMenu bool            `yaml:"show_in_menu"`
	Message    string          `yaml:"message"`
	Replies    []string        `yaml:"replies"`
	Actions    []CommandAction `yaml:"actions"`
}

type CloudPaymentsProvider struct {
	Enabled    bool    `yaml:"enabled"`
	PublicKey  string  `yaml:"public_key"`
	SecretKey  string  `yaml:"secret_key"`
	Price      float64 `yaml:"price"`
	Period     int     `yaml:"period"`
	PeriodUnit string  `yaml:"period_unit"`
}
type YouKassaProvider struct {
	Enabled    bool    `yaml:"enabled"`
	Token      string  `yaml:"token"`
	Price      int     `yaml:"price"`
	Period     float64 `yaml:"period"`
	PeriodUnit string  `yaml:"period_unit"`
}

type BillingProvider struct {
	CloudPayments CloudPaymentsProvider `yaml:"cloud_payments"`
	YooKassa      YouKassaProvider      `yaml:"yoo_kassa"`
}

type BillingConfig struct {
	Enabled                         bool            `yaml:"enabled"`
	SubscriptionAlert               string          `yaml:"subscription_alert"`
	SubscriptionMessage             string          `yaml:"subscription_message"`
	SubscribeBtn                    string          `yaml:"subscribe_btn"`
	UnsubscribeBtn                  string          `yaml:"unsubscribe_btn"`
	UnsubscribeConfirmationMessage  string          `yaml:"unsubscribe_confirmation_message"`
	UnsubscribeCancel               string          `yaml:"unsubscribe_cancel"`
	UnsubscribeConfirm              string          `yaml:"unsubscribe_confirm"`
	ByeMessage                      string          `yaml:"bye_message"`
	NotEnoughMoneyNotification      string          `yaml:"not_enough_amount_notification"`
	SubscriptionGrantedNotification string          `yaml:"subscription_granted_notification"`
	Providers                       BillingProvider `yaml:"providers"`
}

type MDSchool struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

type IntegrationsConfig struct {
	MD MDSchool `yaml:"md"`
}

type MenuItem struct {
	Key         string          `yaml:"key"`
	Enabled     bool            `yaml:"enabled"`
	ButtonLabel string          `yaml:"button_label"`
	Context     string          `yaml:"context"`
	Message     string          `yaml:"message"`
	Actions     []CommandAction `yaml:"actions"`
}

type MenuConfig struct {
	Enabled bool       `yaml:"enabled"`
	Label   string     `yaml:"label"`
	Items   []MenuItem `yaml:"items"`
}

type Config struct {
	AppName          string             `yaml:"app_name"`
	AppURL           string             `yaml:"app_url"`
	BotToken         string             `yaml:"bot_token"`
	BotUsername      string             `yaml:"bot_username"`
	RagApiEndpoint   string             `yaml:"rag_api_endpoint"`
	RagApiToken      string             `yaml:"rag_api_token"`
	AIModeMessage    string             `yaml:"ai_mode_message"`
	TemporaryMessage string             `yaml:"temporary_message"`
	Database         DatabaseConfig     `yaml:"database"`
	Menu             MenuConfig         `yaml:"menu"`
	Modules          ModuleConfig       `yaml:"modules"`
	Commands         []CommandConfig    `yaml:"commands"`
	Integrations     IntegrationsConfig `yaml:"integrations"`
}

var AppConfig Config

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig(path string) {
	// Загружаем переменные из .env файла (если он есть)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Error decoding YAML file: %v", err)
	}

	log.Println("Configuration loaded successfully")
}
