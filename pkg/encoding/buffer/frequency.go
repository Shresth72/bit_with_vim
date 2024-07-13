package buffer

import "github.com/Shresth72/ascii/pkg/encoding/utils"

type Frequency struct {
  PointMap map[int]*FreqPoint
  count int
  Points []*FreqPoint
}

type FreqPoint struct {
  Val int
  Count int
}

func NewFrequency() Frequency {
  return Frequency{
    PointMap: map[int]*FreqPoint{},
    Points: make([]*FreqPoint, 0),
    count: 0,
  }
}

func (f *Frequency) Freq(iterator utils.ByteIterator) {
  for {
    val := iterator.Next()
    point, ok := f.PointMap[val.Value]

    if !ok {
      f.count++

      point = &FreqPoint{
        Count: 0,
        Val:  val.Value,
      }
      f.Points = append(f.Points, point)
      f.PointMap[val.Value] = point
    }
    point.Count++

    if val.Done {
      break
    }
  }
}

func (f *Frequency) Length() int {
  return f.count
}

func (f *Frequency) Reset() {
  f.count = 0
  f.PointMap = map[int]*FreqPoint{}
  f.Points = make([]*FreqPoint, 0)
}
