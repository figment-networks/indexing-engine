package postgreshealth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/figment-networks/indexing-engine/health"
	"go.uber.org/zap"
)

func Example() {
	ctx := context.Background()
	logger := zap.NewNop()

	// Your database
	db := &sql.DB{}

	dbMonitor := NewPostgresMonitorWithMetrics(db, logger)
	monitor := &health.Monitor{}
	monitor.AddProber(ctx, dbMonitor)
	go monitor.RunChecks(ctx, 1*time.Second)

	// Attach endpoints
	mux := http.NewServeMux()
	monitor.AttachHttp(mux)
}
