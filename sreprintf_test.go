package sreprintf

import (
	"testing"
)

func TestSreprintf(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		message     string
		translation string
		expected    string
		wantErr     bool
	}{
		// Basic format tests
		{
			name:        "basic string",
			template:    "Hello %s",
			message:     "Hello World",
			translation: "こんにちは %s",
			expected:    "こんにちは World",
		},
		{
			name:        "basic integer",
			template:    "You have %d messages",
			message:     "You have 5 messages",
			translation: "%d件のメッセージがあります",
			expected:    "5件のメッセージがあります",
		},
		{
			name:        "basic float",
			template:    "Price: $%f",
			message:     "Price: $19.99",
			translation: "価格: ¥%.2f",
			expected:    "価格: ¥19.99",
		},
		{
			name:        "basic bool",
			template:    "Success: %t",
			message:     "Success: true",
			translation: "成功: %t",
			expected:    "成功: true",
		},
		{
			name:        "generic value",
			template:    "Value: %v",
			message:     "Value: test123",
			translation: "値: %v",
			expected:    "値: test123",
		},
		{
			name:        "literal percent",
			template:    "100%% complete",
			message:     "100% complete",
			translation: "100%% 完了",
			expected:    "100% 完了",
		},

		// Multiple placeholders
		{
			name:        "multiple placeholders",
			template:    "Hello %s, you have %d new messages",
			message:     "Hello Alice, you have 5 new messages",
			translation: "こんにちは %sさん、あなたには%d件の新しいメッセージがあります",
			expected:    "こんにちは Aliceさん、あなたには5件の新しいメッセージがあります",
		},
		{
			name:        "mixed types",
			template:    "%s: %d items at $%f each",
			message:     "Cart: 3 items at $9.99 each",
			translation: "%s: %d個 各¥%.2f",
			expected:    "Cart: 3個 各¥9.99",
		},

		// Format options tests
		{
			name:        "width specification",
			template:    "Name: %10s",
			message:     "Name:       John",
			translation: "名前: %10s",
			expected:    "名前:       John",
		},
		{
			name:        "left align",
			template:    "Name: %-10s!",
			message:     "Name: John      !",
			translation: "名前: %-10s！",
			expected:    "名前: John      ！",
		},
		{
			name:        "zero padding",
			template:    "ID: %05d",
			message:     "ID: 00042",
			translation: "ID番号: %05d",
			expected:    "ID番号: 00042",
		},
		{
			name:        "precision float",
			template:    "Pi: %.2f",
			message:     "Pi: 3.14",
			translation: "円周率: %.2f",
			expected:    "円周率: 3.14",
		},
		{
			name:        "precision string",
			template:    "Short: %.5s",
			message:     "Short: Hello",
			translation: "短縮: %.5s",
			expected:    "短縮: Hello",
		},
		{
			name:        "signed number",
			template:    "Change: %+d",
			message:     "Change: +10",
			translation: "変化: %+d",
			expected:    "変化: +10",
		},
		{
			name:        "complex format",
			template:    "Value: %+10.2f",
			message:     "Value:     +12.35",
			translation: "値: %+10.2f",
			expected:    "値:     +12.35",
		},

		// Edge cases
		{
			name:        "empty strings",
			template:    "Name: %s",
			message:     "Name: ",
			translation: "名前: %s",
			expected:    "名前: ",
		},
		{
			name:        "consecutive placeholders",
			template:    "%s%s",
			message:     "HelloWorld",
			translation: "%s%s",
			expected:    "HelloWorld",
		},
		{
			name:        "more placeholders in translation",
			template:    "Hello %s",
			message:     "Hello World",
			translation: "こんにちは %sさん、今日は%sです",
			expected:    "こんにちは Worldさん、今日はです",
		},
		{
			name:        "fewer placeholders in translation",
			template:    "Hello %s, age %d",
			message:     "Hello John, age 30",
			translation: "こんにちは %s",
			expected:    "こんにちは John",
		},
		{
			name:        "no placeholders",
			template:    "Hello World",
			message:     "Hello World",
			translation: "こんにちは世界",
			expected:    "こんにちは世界",
		},
		{
			name:        "percent in text",
			template:    "Loading 50%% done",
			message:     "Loading 50% done",
			translation: "読み込み中 50%% 完了",
			expected:    "読み込み中 50% 完了",
		},

		// Error cases
		{
			name:        "message doesn't match template",
			template:    "Hello %s",
			message:     "Goodbye World",
			translation: "こんにちは %s",
			wantErr:     true,
		},
		{
			name:        "type mismatch ignored",
			template:    "Count: %d",
			message:     "Count: abc",
			translation: "カウント: %s",
			expected:    "カウント: abc",
		},
		{
			name:        "negative numbers",
			template:    "Temperature: %d°C",
			message:     "Temperature: -5°C",
			translation: "温度: %d°C",
			expected:    "温度: -5°C",
		},
		{
			name:        "float with no decimal",
			template:    "Price: %f",
			message:     "Price: 100",
			translation: "価格: %.0f",
			expected:    "価格: 100",
		},
		{
			name:        "special characters",
			template:    "Hello %s!",
			message:     "Hello 世界!",
			translation: "こんにちは %s！",
			expected:    "こんにちは 世界！",
		},
		{
			name:        "hex format as string",
			template:    "Color: %#x",
			message:     "Color: 0xff00ff",
			translation: "色: %#x",
			expected:    "色: 0xff00ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sreprintf(tt.template, tt.message, tt.translation)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sreprintf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Sreprintf() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExtractValues(t *testing.T) {
	tests := []struct {
		name     string
		template string
		message  string
		expected []interface{}
		wantErr  bool
	}{
		{
			name:     "extract string",
			template: "Hello %s",
			message:  "Hello World",
			expected: []interface{}{"World"},
		},
		{
			name:     "extract integer",
			template: "Count: %d",
			message:  "Count: 42",
			expected: []interface{}{42},
		},
		{
			name:     "extract float",
			template: "Price: %f",
			message:  "Price: 19.99",
			expected: []interface{}{19.99},
		},
		{
			name:     "extract bool",
			template: "Success: %t",
			message:  "Success: false",
			expected: []interface{}{false},
		},
		{
			name:     "multiple values",
			template: "%s has %d items",
			message:  "Cart has 3 items",
			expected: []interface{}{"Cart", 3},
		},
		{
			name:     "with format options",
			template: "ID: %05d",
			message:  "ID: 00042",
			expected: []interface{}{42},
		},
		{
			name:     "literal percent",
			template: "100%% done",
			message:  "100% done",
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractValues(tt.template, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("extractValues() returned %d values, want %d", len(got), len(tt.expected))
				return
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("extractValues()[%d] = %v (%T), want %v (%T)",
						i, v, v, tt.expected[i], tt.expected[i])
				}
			}
		})
	}
}

func TestApplyValues(t *testing.T) {
	tests := []struct {
		name        string
		translation string
		values      []interface{}
		expected    string
		wantErr     bool
	}{
		{
			name:        "apply single value",
			translation: "Hello %s",
			values:      []interface{}{"World"},
			expected:    "Hello World",
		},
		{
			name:        "apply multiple values",
			translation: "%s: %d items",
			values:      []interface{}{"Cart", 3},
			expected:    "Cart: 3 items",
		},
		{
			name:        "more values than placeholders",
			translation: "Hello %s",
			values:      []interface{}{"World", "Extra"},
			expected:    "Hello World",
		},
		{
			name:        "fewer values than placeholders",
			translation: "Hello %s %s",
			values:      []interface{}{"World"},
			expected:    "Hello World ",
		},
		{
			name:        "with format options",
			translation: "ID: %05d",
			values:      []interface{}{42},
			expected:    "ID: 00042",
		},
		{
			name:        "literal percent",
			translation: "100%% complete",
			values:      []interface{}{},
			expected:    "100% complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyValues(tt.translation, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("applyValues() = %q, want %q", got, tt.expected)
			}
		})
	}
}
