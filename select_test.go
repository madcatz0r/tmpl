package templates

import (
	"github.com/madcatz0r/tmpl/conditions"
	"github.com/madcatz0r/tmpl/order"
	"testing"
)

type Accounts struct{}
type Users struct{}
type Tariffs struct{}
type TariffVersions struct{}

func init() {
	_ = ParseTags(Accounts{})
	_ = ParseTags(Users{})
	_ = ParseTags(Tariffs{})
	_ = ParseTags(TariffVersions{})
}

func TestSyn(t *testing.T) {
	type args struct {
		syn    string
		fields []string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "simple",
			args: args{
				syn:    "t1",
				fields: []string{"accounts.type", "accounts.balance", "accounts.currency_id"},
			},
			expected: "t1.type,t1.balance,t1.currency_id",
		},
		{
			name: "empty",
			args: args{
				syn:    "t1",
				fields: []string{},
			},
			expected: "",
		},
		{
			name: "single",
			args: args{
				syn:    "t1",
				fields: []string{"accounts.balance"},
			},
			expected: "t1.balance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Syn(tt.args.syn, tt.args.fields...); got != tt.expected {
				t.Errorf("Syn() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_sel_String(t *testing.T) {
	tests := []struct {
		name    string
		sel     *sel
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			sel:     Select("1"),
			want:    "SELECT 1",
			wantErr: false,
		},
		{
			name:    "+ from",
			sel:     Select("accounts.balance").From(Accounts{}),
			want:    "SELECT accounts.balance FROM accounts",
			wantErr: false,
		},
		{
			name:    "+ where",
			sel:     Select("accounts.balance").From(Accounts{}).Where(conditions.Eq("accounts.id", conditions.String("1"))),
			want:    "SELECT accounts.balance FROM accounts\n WHERE accounts.id = '1'",
			wantErr: false,
		},
		{
			name:    "+ group by",
			sel:     Select("accounts.balance", "accounts.type").From(Accounts{}).Where(conditions.Eq("accounts.id", conditions.String("1"))).GroupBy("accounts.type"),
			want:    "SELECT accounts.balance,accounts.type FROM accounts\n WHERE accounts.id = '1'\n GROUP BY accounts.type",
			wantErr: false,
		},
		{
			name:    "order by",
			sel:     Select("accounts.balance", "accounts.type").From(Accounts{}).Where(conditions.Eq("accounts.id", conditions.String("1"))).OrderBy(order.Asc("accounts.type"), order.Desc("accounts.balance")),
			want:    "SELECT accounts.balance,accounts.type FROM accounts\n WHERE accounts.id = '1'\n ORDER BY accounts.type ASC,accounts.balance DESC",
			wantErr: false,
		},
		{
			name:    "join",
			sel:     Select("accounts.balance", "accounts.type").From(Accounts{}).InnerJoin(Users{}).On(conditions.Eq("accounts.owner_id", "users.id")).Where(conditions.Eq("accounts.type", conditions.String("1"))),
			want:    "SELECT accounts.balance,accounts.type FROM accounts\n INNER JOIN users ON accounts.owner_id = users.id\n WHERE accounts.type = '1'",
			wantErr: false,
		},
		{
			name:    "double join",
			sel:     Select("accounts.balance", "accounts.type", "tariffs.user_id").From(Accounts{}).InnerJoin(Users{}).On(conditions.Eq("accounts.owner_id", "users.id")).LeftJoin(Tariffs{}).On(conditions.Eq("users.id", "tariffs.user_id")).Where(conditions.Eq("accounts.type", conditions.String("1"))),
			want:    "SELECT accounts.balance,accounts.type,tariffs.user_id FROM accounts\n INNER JOIN users ON accounts.owner_id = users.id\n LEFT OUTER JOIN tariffs ON users.id = tariffs.user_id\n WHERE accounts.type = '1'",
			wantErr: false,
		},
		{
			name: "complex",
			sel: Select("user_tariff_version.*").From(TariffVersions{}).InnerJoin(Tariffs{}).On(
				conditions.Eq("tariff_versions.tariff_id", "tariffs.outer_id"),
			).Where(conditions.Eq("tariffs.user_id", "$1").And().IsNull("tariffs.deleted_at").And().Le("tariffs.valid_from", "now() AT TIME ZONE 'UTC'").And().Par(
				conditions.Ge("tariffs.valid_to", "now() AT TIME ZONE 'UTC'").Or().IsNull("tariffs.valid_to"),
			).And().Le("tariff_versions.valid_from_date", "(now() AT TIME ZONE 'UTC')::date"),
			).OrderBy("tariff_versions.valid_from_date").Limit(1),
			want: `SELECT user_tariff_version.* FROM tariff_versions
 INNER JOIN tariffs ON tariff_versions.tariff_id = tariffs.outer_id
 WHERE tariffs.user_id = $1 AND tariffs.deleted_at IS NULL AND tariffs.valid_from <= now() AT TIME ZONE 'UTC' AND (tariffs.valid_to >= now() AT TIME ZONE 'UTC' OR tariffs.valid_to IS NULL) AND tariff_versions.valid_from_date <= (now() AT TIME ZONE 'UTC')::date
 ORDER BY tariff_versions.valid_from_date
 LIMIT 1`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sel.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String() got = %v, want %v", got, tt.want)
			}
		})
	}
}
