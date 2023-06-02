package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/LogLayer"
	"log"
)

func (afs *AppFS) createFileWithSpecINodeType(fType int, name string, level int, isRoot bool) int {
	newInodeN := afs.findFreeINode()
	if afs.isINodeInLog(newInodeN) {
		afs.LogCommit()
		//fmt.Println("Ha?")
		newInodeN = afs.findFreeINode()
	} //Avoid reallocating a inode.

	if newInodeN == -1 {
		log.Fatal("No inode number available.")
	}
	afs.tryLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN, level, isRoot)}, []LogLayer.DataBlockMem{})
	return newInodeN
}

type InodeTrace struct {
	inode  BlockLayer.INode
	offset int
}

func (afs *AppFS) findBlockFromStart(allocateWhenNeed bool, inodeN int, index int) (bool, []InodeTrace, []InodeTrace) {
	return afs.findBlockInTree(allocateWhenNeed, inodeN, index, 0)
}

/*
func linkINodes(rootTrs []InodeTrace, treeTrs []InodeTrace) []BlockLayer.INode {
	trees := reverseTrace(linkTreeINodes(reverseTrace(treeTrs)))
	roots := reverseTrace(linkTreeINodes(reverseTrace(append(rootTrs, trees[0]))))
	traces :=  append(roots, trees[1:]...)
	re := []BlockLayer.INode{}


}

func linkTreeINodes(treeTrsO []InodeTrace) []InodeTrace {
	treeTrs := make([]InodeTrace, len(treeTrsO))
	if len(treeTrs) < 2 {
		return treeTrs
	}
	if !testNonStartValid(treeTrs[0].inode) {
		treeTrs[1].inode.Pointers[treeTrs[1].offset] = -1
		treeTrs[0].inode.Valid = false
	} else {
		treeTrs[1].inode.Pointers[treeTrs[1].offset] = treeTrs[0].inode.InodeN
	}
	return append([]InodeTrace{treeTrs[0]}, linkTreeINodes(treeTrs[1:])...)
}

func linkRootINodes(rootTrsO []InodeTrace) []InodeTrace {
	rootTrs := make([]InodeTrace, len(rootTrsO))
	if len(rootTrs) < 2 {
		return rootTrs
	}
	if !testNonStartValid(rootTrs[0].inode) {
		rootTrs[1].inode.PointerToNextINode = -1
		rootTrs[1].inode.Valid = false
	} else {
		rootTrs[1].inode.PointerToNextINode = rootTrs[0].inode.InodeN
	}
	return append([]InodeTrace{rootTrs[0]}, linkTreeINodes(rootTrs[1:])...)
}

func testNonStartValid(n BlockLayer.INode) bool {
	if !n.Valid {
		return false
	}
	c := 0
	for _, v := range n.Pointers {
		if v >= 0 {
			c++
		}
	}
	if n.PointerToNextINode >= 0 {
		c++
	}
	if c <= 0 {
		return false
	}
	return true
}

func reverseTrace(tr []InodeTrace) []InodeTrace {
	r := []InodeTrace{}
	for _, v := range tr {
		r = append([]InodeTrace{v}, tr...)
	}
	return r
}*/

func (afs *AppFS) findBlockInTree(allocateWhenNeed bool, inodeN int, index int, level int) (bool, []InodeTrace, []InodeTrace) {
	inode := afs.GetFileINfo(inodeN)
	if !inode.Valid {
		if allocateWhenNeed {
			inodeN := afs.createFileWithSpecINodeType(BlockLayer.NormalFile, "//", level, true)
			inode = afs.GetFileINfo(inodeN)
		} else {
			return false, []InodeTrace{}, []InodeTrace{}
		}
	}

	if index < blockInInodeLevel(level) {
		b, trs := afs.findBlockInTreeLeaf(allocateWhenNeed, inodeN, index, level)
		if b {
			return true, []InodeTrace{}, trs
		} else {
			return false, []InodeTrace{}, []InodeTrace{}
		}
	} else {
		b, rootTrs, treeTrs := afs.findBlockInTree(allocateWhenNeed, inode.PointerToNextINode, index-blockInInodeLevel(level), level+1)
		return b, append([]InodeTrace{InodeTrace{inode, -1}}, rootTrs...), treeTrs
	}

}

func (afs *AppFS) findBlockInTreeLeaf(allocateWhenNeed bool, inodeN int, index int, level int) (bool, []InodeTrace) {
	inode := afs.GetFileINfo(inodeN)
	if !inode.Valid {
		if allocateWhenNeed {
			inodeN := afs.createFileWithSpecINodeType(BlockLayer.NormalFile, "///////", level, false)
			inode = afs.GetFileINfo(inodeN)
		} else {
			return false, []InodeTrace{}
		}
	}
	if inode.CurrentLevel == 0 {
		//if allocateWhenNeed && inode.Pointers[index] <0 {}
		return true, []InodeTrace{InodeTrace{inode, index}}
	} else {
		offset := index / blockInInodeLevel(level-1)
		trace := InodeTrace{inode: inode, offset: offset}

		b, trs := afs.findBlockInTreeLeaf(allocateWhenNeed, inode.Pointers[offset], index%blockInInodeLevel(level-1), level-1)
		return b, append([]InodeTrace{trace}, trs...)
	}
}

func blockInInodeLevel(level int) int {
	r := 1
	for i := 0; i <= level; i++ {
		r = r * BlockLayer.DirectPointerPerINode
	}
	return r
}
