package raid3_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/rclone/rclone/backend/raid3"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/object"
	"github.com/rclone/rclone/lib/readers"
	"github.com/stretchr/testify/require"
)

func newRaid3ForContextTests(t *testing.T) *raid3.Fs {
	t.Helper()

	ctx := context.Background()
	evenDir := t.TempDir()
	oddDir := t.TempDir()
	parityDir := t.TempDir()

	fInt, err := raid3.NewFs(ctx, "raid3", "", configmap.Simple{
		"even":   evenDir,
		"odd":    oddDir,
		"parity": parityDir,
	})
	require.NoError(t, err)

	f, ok := fInt.(*raid3.Fs)
	require.True(t, ok, "expected *raid3.Fs")

	t.Cleanup(func() {
		_ = f.Features().Shutdown(ctx)
	})

	return f
}

func TestPutCanceledContext(t *testing.T) {
	f := newRaid3ForContextTests(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	payload := []byte("raid3 context cancel test")
	src := object.NewStaticObjectInfo("cancel-put.txt", time.Now(), int64(len(payload)), true, nil, f)
	in := readers.NewContextReader(ctx, bytes.NewReader(payload))

	_, err := f.Put(ctx, in, src)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestCommandStatusCanceledContext(t *testing.T) {
	f := newRaid3ForContextTests(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := f.Command(ctx, "status", nil, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestCommandHealCanceledContext(t *testing.T) {
	f := newRaid3ForContextTests(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := f.Command(ctx, "heal", nil, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestListCanceledContext(t *testing.T) {
	f := newRaid3ForContextTests(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := f.List(ctx, "")
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

