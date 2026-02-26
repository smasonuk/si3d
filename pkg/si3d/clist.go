package si3d

type Clist struct {
	points []Vector3
	count  int
	max    int
	end    int
}

func NewClist(size int) *Clist {
	return &Clist{
		points: make([]Vector3, size),
		max:    size,
	}
}
func (c *Clist) AddPoint(p Vector3) {
	if c.end < c.max {
		c.points[c.end] = p
		c.end++
	}
}
func (c *Clist) Back() {
	c.count--
	if c.count < 0 {
		c.count = c.end - 1
	}
}
func (c *Clist) NextPoint() Vector3 {
	old := c.count
	c.count++
	if c.count >= c.end {
		c.count = 0
	}
	return c.points[old]
}
