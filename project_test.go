package redmine

import "testing"

func Test_validateAdditionalFields(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should pass for single valid field", args{fields: []string{"trackers"}}, false},
		{"should pass for multiple valid fields", args{fields: []string{"enabled_modules", "trackers"}}, false},
		{"should pass for all fields", args{fields: []string{
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities"}}, false},
		{"should pass for repeated fields", args{fields: []string{
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities",
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities"}}, false},
		{"should fail for single unsupported field", args{fields: []string{"invalid value"}}, true},
		{"should fail for single unsupported field among valid fields", args{fields: []string{"enabled_modules", "invalidvalue"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAdditionalFields(tt.args.fields...); (err != nil) != tt.wantErr {
				t.Errorf("validateAdditionalFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("should pass for non-existing parameter", func(t *testing.T) {
		if err := validateAdditionalFields(); (err != nil) != false {
			t.Errorf("expected error for empty string")
		}
	})
	t.Run("should fail for empty string", func(t *testing.T) {
		if err := validateAdditionalFields(""); (err != nil) != true {
			t.Errorf("expected error for empty string")
		}
	})
}

func Test_additionalFieldsParam(t *testing.T) {
	type args struct {
		additionalFields []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should create empty string for nil", args{nil}, ""},
		{"should create empty string for empty slice", args{[]string{}}, ""},
		{"should create parameter without comma for single field", args{[]string{"enabled_modules"}}, "enabled_modules"},
		{"should create comma delimited list without trailing comma for multiple fields", args{[]string{
			"enabled_modules", "enabled_modules", "trackers"}}, "enabled_modules,enabled_modules,trackers"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := additionalFieldsParam(tt.args.additionalFields...); got != tt.want {
				t.Errorf("additionalFieldsParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
