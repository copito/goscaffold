package entity

type Prompt struct {
	Items map[string]PromptItem `mapstructure:"prompt"`
}

type PromptItem struct {
	Key string

	OrderID      int           `mapstructure:"order"`
	DefaultValue any           `mapstructure:"default"`
	Options      []interface{} `mapstructure:"options"`
	AllowEdit    bool          `mapstructure:"allow_edit"`
	HideEntered  bool          `mapstructure:"hide_entered"`
	// Validation   []string `mapstructure:"validation"`
}
