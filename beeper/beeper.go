package beeper

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
	frequency  = 440
)

type sineWaveStream struct {
	freq       float64
	sampleRate int
	position   int64
}

type Beeper struct {
	audioContext *audio.Context
	player       *audio.Player
}

func Init() (*Beeper, error) {
	audioContext := audio.NewContext(sampleRate)

	beeper := &Beeper{
		audioContext: audioContext,
	}

	return beeper, nil
}

func (b *Beeper) Play() {
	sineWave := &sineWaveStream{
		freq:       frequency,
		sampleRate: sampleRate,
	}

	player, err := b.audioContext.NewPlayer(sineWave)
	if err != nil {
		return
	}
	b.player = player
	b.player.Play()

	go func() {
		time.Sleep(time.Second / 10)
		if b.player != nil {
			b.player.Pause()
			b.player.Close()
			b.player = nil
		}
	}()
}

func (s *sineWaveStream) Read(p []byte) (int, error) {
	// remember that slices are passed as references, not as values

	if len(p) == 0 {
		return 0, nil
	}

	sampleCount := len(p) / 2 // 16-bit mono audio

	for i := 0; i < sampleCount; i++ {
		pos := float64(s.position) / float64(s.sampleRate)
		phase := 2 * math.Pi * s.freq * pos

		amplitude := int16(math.Sin(phase) * 10000)

		offset := i * 2
		p[offset] = byte(amplitude & 0xFF)
		p[offset+1] = byte((amplitude >> 8) & 0xFF)

		s.position++
	}

	return len(p), nil
}

func (s *sineWaveStream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0: // SEEK_SET
		s.position = offset
	case 1: // SEEK_CUR
		s.position += offset
	case 2: // SEEK_END
		s.position = 0
	}
	return s.position, nil
}
