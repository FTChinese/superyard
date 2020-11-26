package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/guregu/null"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_SaveMemberSnapshot(t *testing.T) {
	p := test.NewPersona()

	env := NewEnv(test.DBX, zaptest.NewLogger(t))

	type args struct {
		s reader.MemberSnapshot
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save member snapshot",
			args: args{
				s: reader.NewSnapshot(p.Membership()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SaveMemberSnapshot(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SaveMemberSnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListMemberSnapshots(t *testing.T) {
	env := NewEnv(test.DBX, zaptest.NewLogger(t))

	type args struct {
		ids reader.IDs
		p   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List membership snapshot",
			args: args{
				ids: reader.IDs{
					FtcID:   null.StringFrom("216b2a94-e140-40d1-94e1-0f1de0fb3320"),
					UnionID: null.StringFrom("ogfvwjn5kmva3hRz4_SvRujh4mJM"),
				},
				p: gorest.NewPagination(1, 20),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListMemberSnapshots(tt.args.ids, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMemberSnapshots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
