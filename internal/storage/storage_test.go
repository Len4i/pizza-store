package storage

import (
	"testing"

	"github.com/Len4i/pizza-store/internal/storage/sqlite"
	_ "modernc.org/sqlite"
)

func TestStorage_SaveOrder(t *testing.T) {
	type args struct {
		order Order
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "working",
			args: args{
				order: Order{
					Size:      "family",
					Amount:    1,
					PizzaType: "margherita",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db, err := sqlite.New(":memory:")
		if err != nil {
			t.Fatalf("Failed to ping db: %v", err)
		}
		storage := New(db)
		if err != nil {
			t.Fatalf("Failed to create storage: %v", err)
		}
		defer storage.Close()
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.SaveOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.SaveOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.SaveOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetOrder(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    Order
		wantErr bool
	}{
		{
			name: "working",
			args: args{
				id: 1,
			},
			want: Order{
				Size:      "family",
				Amount:    1,
				PizzaType: "margherita",
			},
			wantErr: false,
		},
		{
			name: "not working",
			args: args{
				id: 2,
			},
			want:    Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		db, err := sqlite.New(":memory:")
		if err != nil {
			t.Fatalf("Failed to ping db: %v", err)
		}
		storage := New(db)
		if err != nil {
			t.Fatalf("Failed to create storage: %v", err)
		}
		defer storage.Close()
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.SaveOrder(Order{
				Size:      "family",
				Amount:    1,
				PizzaType: "margherita",
			})
			if err != nil {
				t.Fatalf("Failed to save order: %v", err)
			}

			got, err := storage.GetOrder(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.GetOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
