package main

import (
	// "bufio"
  "fmt"
  "math/rand"
  "os"
  "os/exec"
  "strconv"
  "time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type Pos struct {
  X, Y int
}

type FloatPos struct {
  X, Y float64
}

func (p FloatPos) pos() Pos {
  return Pos{int(p.X), int(p.Y)}
}

type Block int

const (
  Air Block = iota
  Ground
  Spikes
)

type State struct {
  pos, vel FloatPos
  cam FloatPos
  world map[Pos]Block
  border int
  score int
}

func main() {
  keys, _ := keyboard.GetKeys(10)
  defer keyboard.Close()

  s := &State{
    world: make(map[Pos]Block),
  }

  seed := rand.Int63()

  if len(os.Args) > 1 {
    seedraw, _ := strconv.Atoi(os.Args[1])
    seed = int64(seedraw)
  }

  rng := rand.New(rand.NewSource(seed))

  main: for {
    s.cam.X += (s.pos.X - s.cam.X)*.2
    // s.cam.Y += (s.pos.Y - s.cam.Y)*.2

    curblock := s.world[s.pos.pos()]

    if curblock == Spikes {
      s.pos = FloatPos{X: 0}
      continue main
    }

    if curblock == Ground {
      s.pos.Y -= 1
    }

    if s.pos.X > float64(s.border) {
      for range 20 {
        X := rng.Intn(90)+10+s.border
        Y := 1-rng.Intn(3)
        for x := -1; x < 2; x++ {
          for y := -1; y < 2; y++ {
            if Y+y > 0 { continue }
            s.world[Pos{
             X+x, Y+y,
            }] = Ground
          }
        }
      }
      for range 30 {
        s.world[Pos{
          rng.Intn(90)+10+s.border,
          0-rng.Intn(3),
        }] = Spikes
      }
      s.border += 100
    }

    r := s.world[Pos{int(s.pos.X)+1, int(s.pos.Y)}] == Ground
    l := s.world[Pos{int(s.pos.X)-1, int(s.pos.Y)}] == Ground
    d := s.world[Pos{int(s.pos.X), int(s.pos.Y)+1}] == Ground
    u := s.world[Pos{int(s.pos.X), int(s.pos.Y)-1}] == Ground

    select {
    case key := <-keys:
      if key.Key == keyboard.KeyEsc { break main }
      if key.Rune == 't' { s.pos.X += 100 }
      if key.Rune == 'd' && !r { s.pos.X += 1 }
      if key.Rune == 'a' && !l { s.pos.X -= 1 }
      if key.Rune == 'w' && (s.pos.Y == 0 || d) && !u {
        s.pos.Y -= 1
        s.vel.Y = -.8
        d = false
      }
    default:
    }

    s.vel.Y += .1
    if s.pos.Y >= 0 || d {
      if s.pos.Y >= 0 {
        s.pos.Y = 0
      }
      s.vel.Y = 0
    }
    s.pos.Y += s.vel.Y

    if int(s.pos.X) > s.score {
      s.score = int(s.pos.X)
    }

    screen := ""
    for y := int(s.cam.Y)-10; y < int(s.cam.Y)+10; y++ {
      row := ""
      for x := int(s.cam.X)-40; x < int(s.cam.X)+40; x++ {
        cur := Pos{x, y}
        if s.pos.pos() == cur {
          row += color.MagentaString("O")
        } else if s.world[cur] == Spikes {
          row += color.RedString("A")
        } else if s.world[cur] == Ground || y == 1 {
          if x == s.border-100 {
            row += color.YellowString("M")
          } else {
            row += "X"
          }
        } else { row += " " }
      }
      screen += row + "\n"
    }
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
    fmt.Printf("%v  score: %v  max: %v\n", seed, s.score, int(s.pos.X))
    fmt.Print(screen)

    time.Sleep(time.Millisecond * 30)
  }
}
