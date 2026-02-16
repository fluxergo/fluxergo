package voice

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"time"
)

const (
	// OpusFrameSizeMs is the size of an opus frame in milliseconds.
	OpusFrameSizeMs = 20
)

type (
	// AudioSenderCreateFunc is used to create a new AudioSender sending audio to the given Conn.
	AudioSenderCreateFunc func(logger *slog.Logger, provider OpusFrameProvider, conn Conn) AudioSender

	// AudioSender is used to send audio to a Conn.
	AudioSender interface {
		Open()
		Close()
	}

	// OpusFrameProvider is used to provide opus frames to an AudioSender.
	OpusFrameProvider interface {
		// ProvideOpusFrame provides an opus frame to the AudioSender.
		ProvideOpusFrame() ([]byte, error)

		// Close closes the OpusFrameProvider.
		Close()
	}
)

// NewAudioSender creates a new AudioSender sending audio from the given OpusFrameProvider to the given Conn.
func NewAudioSender(logger *slog.Logger, opusProvider OpusFrameProvider, conn Conn) AudioSender {
	return &defaultAudioSender{
		logger:       logger,
		opusProvider: opusProvider,
		conn:         conn,
	}
}

type defaultAudioSender struct {
	logger       *slog.Logger
	cancelFunc   context.CancelFunc
	opusProvider OpusFrameProvider
	conn         Conn
	trackWriter  io.Writer
}

func (s *defaultAudioSender) Open() {
	go s.open()
}

func (s *defaultAudioSender) open() {
	w, err := s.conn.LiveKit().VideoWriter()
	if err != nil {
		s.logger.Error("error creating audio writer", slog.Any("err", err))
		return
	}
	defer w.Close()
	s.trackWriter = w

	defer s.logger.Debug("closing audio sender")
	lastFrameSent := time.Now().UnixMilli()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		default:
			s.send()
			sleepTime := time.Duration(int64(OpusFrameSizeMs) - (time.Now().UnixMilli() - lastFrameSent))
			if sleepTime > 0 {
				time.Sleep(sleepTime * time.Millisecond)
			}
			if time.Now().UnixMilli() < lastFrameSent+int64(OpusFrameSizeMs)*3 {
				lastFrameSent += int64(OpusFrameSizeMs)
			} else {
				lastFrameSent = time.Now().UnixMilli()
			}
		}
	}
}

func (s *defaultAudioSender) send() {
	if s.opusProvider == nil {
		return
	}
	opus, err := s.opusProvider.ProvideOpusFrame()
	if err != nil && !errors.Is(err, io.EOF) {
		s.logger.Error("error while reading opus frame", slog.Any("err", err))
		return
	}
	if len(opus) == 0 {
		return
	}

	if _, err = s.trackWriter.Write(opus); err != nil {
		s.handleErr(err)
	}
}

func (s *defaultAudioSender) handleErr(err error) {
	if errors.Is(err, net.ErrClosed) {
		s.Close()
		return
	}
	s.logger.Error("failed to send audio", slog.Any("err", err))
}

func (s *defaultAudioSender) Close() {
	s.cancelFunc()
}
