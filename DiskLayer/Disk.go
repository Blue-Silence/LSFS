package DiskLayer

import (
	"LSF/Setting"
	"log"
)

type VirtualDisk struct {
	//blocks     [Setting.BlockN]Block
	blocks [Setting.BlockN]RealBlock
	//superBlock Block
	//superBlock RealBlock
	superBlock []RealBlock
}

type Block interface {
	CanBeBlock()
	ToBlock() RealBlock
	FromBlock(RealBlock) Block
}

type RealBlock = [Setting.BlockSize]byte

func BytesToBlock(d []byte) RealBlock {
	var b RealBlock
	if len(d) > Setting.BlockSize {
		log.Fatal("Too big to be a block.")
	}
	copy(b[:], d)
	return b
}

func BytesToBlocks(d []byte) []RealBlock {
	dN := d[:]
	bs := []RealBlock{}
	for len(dN) > Setting.BlockSize {
		bs = append(bs, BytesToBlock(dN[:Setting.BlockSize]))
		dN = dN[Setting.BlockSize:]
	}
	if len(dN) > 0 {
		bs = append(bs, BytesToBlock(dN[:Setting.BlockSize]))
	}
	return bs
}

func BlockToBytes(b RealBlock) []byte {
	return b[:]
}

func (d *VirtualDisk) ReadBlock(index int) RealBlock {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk read access at ", index)
	}
	return d.blocks[index]
}

func (d *VirtualDisk) WriteBlock(index int, b Block) {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk write access at ", index)
	}
	d.blocks[index] = b.ToBlock()
}

func (d *VirtualDisk) ReadSuperBlock() []RealBlock {
	return d.superBlock
}

func (d *VirtualDisk) WriteSuperBlock(b []RealBlock) {
	d.superBlock = b
}
