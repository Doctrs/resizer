package resize


func (r *Resizer) Cover(width uint, height uint) (newWidth uint, newHeight uint){
	koefOld := float32(r.img.Bounds().Max.X) / float32(r.img.Bounds().Max.Y)
	koefNew := float32(width) / float32(height)
	switch true {
	case width == 0:
		return uint(float32(height) * koefOld), height
	case height == 0:
		return width, uint(float32(width) / koefOld)
	}

	switch true {
	case koefNew > koefOld:
		return uint(float32(height) * koefOld), height
	case koefNew < koefOld:
		return width, uint(float32(width) / koefOld)
	default:
		return width, height
	}
}

func (r *Resizer) Contain(width uint, height uint) (newWidth uint, newHeight uint){
	koefOld := float32(r.img.Bounds().Max.X) / float32(r.img.Bounds().Max.Y)
	koefNew := float32(width) / float32(height)
	switch true {
	case width == 0:
		return uint(float32(height) * koefOld), height
	case height == 0:
		return width, uint(float32(width) / koefOld)
	}

	switch true {
	case koefNew > koefOld:
		return width, uint(float32(width) * koefOld)
	case koefNew < koefOld:
		return uint(float32(height) / koefOld), height
	default:
		return width, height
	}
}
