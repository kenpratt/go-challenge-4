package main

import "fmt"

// A repacker repacks trucks.
type repacker struct {
	pallets [][]pallet
}

// There are 10 possible sizes of boxes
const (
	Dim1x1 = iota
	Dim2x1
	Dim2x2
	Dim3x1
	Dim3x2
	Dim3x3
	Dim4x1
	Dim4x2
	Dim4x3
	Dim4x4
)

// Translate a width and length onto a box size index.
func dimIndex(w, l uint8) uint8 {
	switch {
	case w == 1 && l == 1:
		return Dim1x1
	case w == 2 && l == 1:
		return Dim2x1
	case w == 2 && l == 2:
		return Dim2x2
	case w == 3 && l == 1:
		return Dim3x1
	case w == 3 && l == 2:
		return Dim3x2
	case w == 3 && l == 3:
		return Dim3x3
	case w == 4 && l == 1:
		return Dim4x1
	case w == 4 && l == 2:
		return Dim4x2
	case w == 4 && l == 3:
		return Dim4x3
	case w == 4 && l == 4:
		return Dim4x4
	default:
		panic(fmt.Sprintf("Invalid box size: w=%v l=%v", w, l))
	}
}

func (r *repacker) unload(t *truck) {
	r.pallets = append(r.pallets, t.pallets)
}

func (r *repacker) sortPallets() (boxes [10][]box) {
	for _, pallets := range r.pallets {
		for _, p := range pallets {
			for _, b := range p.boxes {
				if b.w < b.l {
					b.l, b.w = b.w, b.l
				}
				i := dimIndex(b.w, b.l)
				boxes[i] = append(boxes[i], b)
			}
		}
	}
	return
}

func (r *repacker) packEverything(id int) *truck {
	out := &truck{id: id}
	boxesByDim := r.sortPallets()

	// Pack 4x4s, one per pallet
	for _, b := range boxesByDim[Dim4x4] {
		b.x, b.y = 0, 0
		out.pallets = append(out.pallets, pallet{boxes: []box{b}})
	}
	boxesByDim[Dim4x4] = []box{}

	// Pack 4x2s, two at a time
	boxes := []box{}
	for i, b := range boxesByDim[Dim4x2] {
		switch i % 2 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)
			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim4x2] = []box{}

	// Pack 2x2s, four at a time
	boxes = []box{}
	for i, b := range boxesByDim[Dim2x2] {
		switch i % 4 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 0, 2
			boxes = append(boxes, b)
		case 2:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)
		case 3:
			b.x, b.y = 2, 2
			boxes = append(boxes, b)
			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim2x2] = []box{}

	// Pack 4x3s, one per pallet, filling the left-over row
	for _, b := range boxesByDim[Dim4x3] {
		b.x, b.y = 0, 0
		boxes = []box{b}

		// fill the last row
		if len(boxesByDim[Dim4x1]) > 0 {
			b = boxesByDim[Dim4x1][0]
			b.x, b.y = 3, 0
			boxes = append(boxes, b)
			boxesByDim[Dim4x1] = boxesByDim[Dim4x1][1:]
		}

		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim4x3] = []box{}

	// Pack 3x3s, one per pallet, filling the left-over row and column
	for _, b := range boxesByDim[Dim3x3] {
		b.x, b.y = 0, 0
		boxes = []box{b}

		// fill the last row
		if len(boxesByDim[Dim4x1]) > 0 {
			b = boxesByDim[Dim4x1][0]
			b.x, b.y = 3, 0
			boxes = append(boxes, b)
			boxesByDim[Dim4x1] = boxesByDim[Dim4x1][1:]
		}

		// fill the last column
		if len(boxesByDim[Dim3x1]) > 0 {
			b = boxesByDim[Dim3x1][0]
			b.l, b.w = b.w, b.l // flip vertical
			b.x, b.y = 0, 3
			boxes = append(boxes, b)
			boxesByDim[Dim3x1] = boxesByDim[Dim3x1][1:]
		}

		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim3x3] = []box{}

	// Pack 3x2s, two per pallet, filling the left-over column
	boxes = []box{}
	for i, b := range boxesByDim[Dim3x2] {
		switch i % 2 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)

			// fill the last column
			if len(boxesByDim[Dim4x1]) > 0 {
				b = boxesByDim[Dim4x1][0]
				b.l, b.w = b.w, b.l // flip vertical
				b.x, b.y = 0, 3
				boxes = append(boxes, b)
				boxesByDim[Dim4x1] = boxesByDim[Dim4x1][1:]
			}

			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim3x2] = []box{}

	// Pack 3x1s, four per pallet, filling the left-over column
	boxes = []box{}
	for i, b := range boxesByDim[Dim3x1] {
		switch i % 4 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 1, 0
			boxes = append(boxes, b)
		case 2:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)
		case 3:
			b.x, b.y = 3, 0
			boxes = append(boxes, b)

			// fill the last column
			if len(boxesByDim[Dim4x1]) > 0 {
				b = boxesByDim[Dim4x1][0]
				b.l, b.w = b.w, b.l // flip vertical
				b.x, b.y = 0, 3
				boxes = append(boxes, b)
				boxesByDim[Dim4x1] = boxesByDim[Dim4x1][1:]
			}

			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim3x1] = []box{}

	// Pack 2x1s, eight per pallet
	boxes = []box{}
	for i, b := range boxesByDim[Dim2x1] {
		switch i % 8 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 0, 2
			boxes = append(boxes, b)
		case 2:
			b.x, b.y = 1, 0
			boxes = append(boxes, b)
		case 3:
			b.x, b.y = 1, 2
			boxes = append(boxes, b)
		case 4:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)
		case 5:
			b.x, b.y = 2, 2
			boxes = append(boxes, b)
		case 6:
			b.x, b.y = 3, 0
			boxes = append(boxes, b)
		case 7:
			b.x, b.y = 3, 2
			boxes = append(boxes, b)
			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim2x1] = []box{}

	// Pack 1x1s, sixteen per pallet
	boxes = []box{}
	for i, b := range boxesByDim[Dim1x1] {
		switch i % 16 {
		case 0:
			b.x, b.y = 0, 0
			boxes = append(boxes, b)
		case 1:
			b.x, b.y = 0, 1
			boxes = append(boxes, b)
		case 2:
			b.x, b.y = 0, 2
			boxes = append(boxes, b)
		case 3:
			b.x, b.y = 0, 3
			boxes = append(boxes, b)
		case 4:
			b.x, b.y = 1, 0
			boxes = append(boxes, b)
		case 5:
			b.x, b.y = 1, 1
			boxes = append(boxes, b)
		case 6:
			b.x, b.y = 1, 2
			boxes = append(boxes, b)
		case 7:
			b.x, b.y = 1, 3
			boxes = append(boxes, b)
		case 8:
			b.x, b.y = 2, 0
			boxes = append(boxes, b)
		case 9:
			b.x, b.y = 2, 1
			boxes = append(boxes, b)
		case 10:
			b.x, b.y = 2, 2
			boxes = append(boxes, b)
		case 11:
			b.x, b.y = 2, 3
			boxes = append(boxes, b)
		case 12:
			b.x, b.y = 3, 0
			boxes = append(boxes, b)
		case 13:
			b.x, b.y = 3, 1
			boxes = append(boxes, b)
		case 14:
			b.x, b.y = 3, 2
			boxes = append(boxes, b)
		case 15:
			b.x, b.y = 3, 3
			boxes = append(boxes, b)
			out.pallets = append(out.pallets, pallet{boxes: boxes})
			boxes = []box{}
		}
	}
	if len(boxes) > 0 {
		out.pallets = append(out.pallets, pallet{boxes: boxes})
	}
	boxesByDim[Dim1x1] = []box{}

	return out
}

func newRepacker(in <-chan *truck, out chan<- *truck) *repacker {
	r := &repacker{}
	go func() {
		for t := range in {
			// The last truck is indicated by its id. You might
			// need to do something special here to make sure you
			// send all the boxes.
			r.unload(t)
			if t.id == idLastTruck {
				out <- r.packEverything(t.id)
			} else {
				out <- &truck{id: t.id}
			}

		}
		// The repacker must close channel out after it detects that
		// channel in is closed so that the driver program will finish
		// and print the stats.
		close(out)
	}()
	return r
}
