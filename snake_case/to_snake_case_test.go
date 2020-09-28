package snake_case

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		have     string
		expected string
	}{
		// TODO: Add test cases.
		{
			name:     "ID",
			have:     "ID",
			expected: "id",
		},
		{
			name:     "OuterID",
			have:     "OuterID",
			expected: "outer_id",
		},
		{
			name:     "FreeCashoutAmount",
			have:     "FreeCashoutAmount",
			expected: "free_cashout_amount",
		},
		{
			name:     "FreeRemittancesAmount",
			have:     "FreeRemittancesAmount",
			expected: "free_remittances_amount",
		},
		{
			name:     "FreeC2COutAmount",
			have:     "FreeC2COutAmount",
			expected: "free_c2c_out_amount",
		},
		{
			name:     "FreeC2CInAmount",
			have:     "FreeC2CInAmount",
			expected: "free_c2c_in_amount",
		},
		{
			name:     "SalaryFreeCashoutAmount",
			have:     "SalaryFreeCashoutAmount",
			expected: "salary_free_cashout_amount",
		},
		{
			name:     "FriendTransferOutFeeMinAmount",
			have:     "FriendTransferOutFeeMinAmount",
			expected: "friend_transfer_out_fee_min_amount",
		},
		{
			name:     "FreeFriendTransferOutAmount",
			have:     "FreeFriendTransferOutAmount",
			expected: "free_friend_transfer_out_amount",
		},
		{
			name:     "C2COutFeePercentage",
			have:     "C2COutFeePercentage",
			expected: "c2c_out_fee_percentage",
		},
		{
			name:     "FeePercentage",
			have:     "FeePercentage",
			expected: "fee_percentage",
		},
		{
			name:     "CreatedAt",
			have:     "CreatedAt",
			expected: "created_at",
		},
		{
			name:     "UpdatedAt",
			have:     "UpdatedAt",
			expected: "updated_at",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.have); got != tt.expected {
				t.Errorf("ToSnakeCase() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
