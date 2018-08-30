// DO NOT EDIT; automatically generated by brio-gen

package hbook

import (
	"encoding/binary"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (o *H1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.Binning.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.Ann.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *H1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Binning.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Ann.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *H2D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.Binning.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.Ann.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *H2D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Binning.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Ann.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *P1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.bng.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.ann.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *P1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.bng.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.ann.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *S2D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.pts)))
	data = append(data, buf[:8]...)
	for i := range o.pts {
		o := &o.pts[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.ann.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *S2D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.pts = make([]Point2D, n)
		data = data[8:]
		for i := range o.pts {
			oi := &o.pts[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.ann.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return err
}
