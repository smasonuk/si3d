package gosie3d

import (
	"image/color"
	"testing"
)

func TestNewPolygonBatcher(t *testing.T) {
	pb := NewPolygonBatcher(100)
	if pb == nil {
		t.Fatal("NewPolygonBatcher returned nil")
	}
	if cap(pb.commands) != 100 {
		t.Errorf("NewPolygonBatcher initial capacity expected 100, got %d", cap(pb.commands))
	}
}

func TestPolygonBatcher_AddPolygon(t *testing.T) {
	pb := NewPolygonBatcher(10)
	xp := []float32{1, 2, 3}
	yp := []float32{4, 5, 6}
	col := color.RGBA{255, 0, 0, 255}

	pb.AddPolygon(xp, yp, col)

	if len(pb.commands) != 1 {
		t.Errorf("AddPolygon command count expected 1, got %d", len(pb.commands))
	}

	cmd := pb.commands[0]
	if !cmd.hasFill {
		t.Error("AddPolygon should have fill")
	}
	if cmd.hasStroke {
		t.Error("AddPolygon should not have stroke")
	}
	if cmd.fillClr != col {
		t.Error("AddPolygon color mismatch")
	}
}

func TestPolygonBatcher_AddPolygonOutline(t *testing.T) {
	pb := NewPolygonBatcher(10)
	xp := []float32{1, 2, 3}
	yp := []float32{4, 5, 6}
	col := color.RGBA{0, 255, 0, 255}
	width := float32(2.0)

	pb.AddPolygonOutline(xp, yp, width, col)

	if len(pb.commands) != 1 {
		t.Errorf("AddPolygonOutline command count expected 1, got %d", len(pb.commands))
	}

	cmd := pb.commands[0]
	if cmd.hasFill {
		t.Error("AddPolygonOutline should not have fill")
	}
	if !cmd.hasStroke {
		t.Error("AddPolygonOutline should have stroke")
	}
	if cmd.strokeClr != col {
		t.Error("AddPolygonOutline color mismatch")
	}
	if cmd.strokeW != width {
		t.Errorf("AddPolygonOutline width mismatch, expected %f, got %f", width, cmd.strokeW)
	}
}

func TestPolygonBatcher_AddPolygonAndOutline(t *testing.T) {
	pb := NewPolygonBatcher(10)
	xp := []float32{1, 2, 3}
	yp := []float32{4, 5, 6}
	fillCol := color.RGBA{0, 0, 255, 255}
	strokeCol := color.RGBA{255, 255, 0, 255}
	width := float32(1.5)

	pb.AddPolygonAndOutline(xp, yp, fillCol, strokeCol, width)

	if len(pb.commands) != 1 {
		t.Errorf("AddPolygonAndOutline command count expected 1, got %d", len(pb.commands))
	}

	cmd := pb.commands[0]
	if !cmd.hasFill {
		t.Error("AddPolygonAndOutline should have fill")
	}
	if !cmd.hasStroke {
		t.Error("AddPolygonAndOutline should have stroke")
	}
	if cmd.fillClr != fillCol {
		t.Error("AddPolygonAndOutline fill color mismatch")
	}
	if cmd.strokeClr != strokeCol {
		t.Error("AddPolygonAndOutline stroke color mismatch")
	}
}
