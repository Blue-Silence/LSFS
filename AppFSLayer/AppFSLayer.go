package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/LogLayer"
	"LSF/Setting"

	//"fmt"
	"log"
)

type AppFS struct {
	blockFs BlockLayer.BlockFS
	fLog    LogLayer.FSLog
}

func (afs *AppFS) FormatFS(VD DiskLayer.VirtualDisk) {
	afs.blockFs.VD = VD
	afs.fLog.InitLog()
	initINodes := []BlockLayer.INode{createInode(BlockLayer.Folder, "", true, 0, 0, true)} //Adding root
	for i := 1; i < Setting.MaxInodeN; i++ {
		initINodes = append(initINodes, createInode(BlockLayer.NormalFile, "", false, i, 0, true)) //Adding invalid inodes to init imap
	}
	afs.fLog.ConstructLog(initINodes, []LogLayer.DataBlockMem{})
	_, _, _, initSegLen := afs.fLog.LenInBlock()
	initStart := afs.blockFs.FindSpaceForSeg(initSegLen)
	blocks, imapLs := afs.fLog.Log2DiskBlock(initStart, make(map[int]BlockLayer.INodeMap))
	afs.blockFs.ApplyUpdate(initStart, blocks, imapLs)
	afs.fLog.InitLog()
}

func createInode(fType int, name string, valid bool, inodeN int, level int, isRoot bool) BlockLayer.INode {
	// isRoot and level is used  to ad support for inode tree.
	in := BlockLayer.INode{Valid: valid, FileType: fType, Name: name, InodeN: inodeN, PointerToNextINode: -1, CurrentLevel: level, IsRoot: true}
	for i, _ := range in.Pointers {
		in.Pointers[i] = -1 //Init to invalid pointers
	}
	return in
}

func (afs *AppFS) findFreeINode() int {
	for i := 0; i < Setting.MaxInodeN; i++ {
		if afs.blockFs.INodeN2iNode(i).Valid == false {
			return i
		}
	}
	return -1
}

func (afs *AppFS) LogCommitWithINMap(imapNeeded map[int]BlockLayer.INodeMap) {
	for _, v := range afs.fLog.ImapNeeded() {
		imapNeeded[v] = BlockLayer.INodeMap{}.FromBlock(afs.blockFs.VD.ReadBlock(BlockLayer.SuperBlock{}.FromBlocks(afs.blockFs.VD.ReadSuperBlock()).INodeMaps[v])).(BlockLayer.INodeMap)
	} //Get inaodmap needed
	_, _, _, logSegLen := afs.fLog.LenInBlock()
	start := afs.blockFs.FindSpaceForSeg(logSegLen)
	if start < 0 {
		//WE will add GC later. TO BE DONE
		afs.GC(-1)
		log.Fatal("No space!")
	}
	bs, newIMap := afs.fLog.Log2DiskBlock(start, imapNeeded)
	afs.blockFs.ApplyUpdate(start, bs, newIMap)
	afs.fLog.InitLog()
}

func (afs *AppFS) LogCommit() {
	afs.LogCommitWithINMap(make(map[int]BlockLayer.INodeMap))
}

func (afs *AppFS) isINodeInLog(n int) bool {
	return afs.fLog.IsINodeInLog(n)
}

func (afs *AppFS) GetFileINfo(inodeN int) BlockLayer.INode {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	return afs.blockFs.INodeN2iNode(inodeN)
}

func (afs *AppFS) CreateFile(fType int, name string) int {
	return afs.createFileWithSpecINodeType(fType, name, 0, true)
}

func (afs *AppFS) WriteFile(inodeN int, index []int, data []DiskLayer.Block) {
	inode := afs.GetFileINfo(inodeN)
	if inode.Valid == false {
		log.Fatal("Invalid write to non-exsistent inode:", inodeN, "  get inode:", inode)
	}
	ds := []LogLayer.DataBlockMem{}
	ins := []BlockLayer.INode{}
	for i, v := range index {
		_, _, traces := afs.findBlockFromStart(true, inodeN, v)
		traceTail := traces[len(traces)-1]
		//ds = append(ds, LogLayer.DataBlockMem{Inode: inodeN, Index: v, Data: data[i].ToBlock()})
		ins = append(ins, afs.GetFileINfo(traceTail.inode.InodeN))
		ds = append(ds, LogLayer.DataBlockMem{Inode: traceTail.inode.InodeN, Index: traceTail.offset, Data: data[i].ToBlock()})
	}
	//afs.fLog.ConstructLog([]BlockLayer.INode{inode}, ds)
	afs.tryLog(ins, ds)
}

func (afs *AppFS) ReadFile(inodeN int, index int) DiskLayer.RealBlock {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	var emptyB DiskLayer.RealBlock
	b, rs, trsTree := afs.findBlockFromStart(false, inodeN, index)
	if b {
		leaf := trsTree[len(trsTree)-1]
		return afs.blockFs.ReadFile(leaf.inode.InodeN, leaf.offset)
	}
	log.Println(rs)
	log.Println(trsTree)
	log.Fatal("GOT EMPTY BLOCK!")
	return emptyB
}

func (afs *AppFS) DeleteFile(inodeN int) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	inode := BlockLayer.INode{InodeN: inodeN, Valid: false}
	//afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
	afs.tryLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
}

func (afs *AppFS) DeleteBlockInFile(inodeN int, index []int) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	for _, v := range index {
		b, _, treeTrs := afs.findBlockFromStart(false, inodeN, v)
		if b {
			inode := afs.GetFileINfo(treeTrs[len(treeTrs)-1].inode.InodeN)
			//treeTrs[len(treeTrs)-1].inode.Pointers[treeTrs[len(treeTrs)-1].offset] = -1
			inode.Pointers[treeTrs[len(treeTrs)-1].offset] = -1
			afs.tryLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
			//afs.LogCommit() //We may do GC later.
		}
		//inode.Pointers[v] = -1
	}
}

func (afs *AppFS) tryLog(inodes []BlockLayer.INode, ds []LogLayer.DataBlockMem) {
	if afs.fLog.NeedCommit() {
		afs.LogCommit()
	}
	if !afs.fLog.ConstructLog(inodes, ds) {
		afs.LogCommit()
		if !afs.fLog.ConstructLog(inodes, ds) {
			log.Fatal("No space!")
		}
	}
}

//////////////////////
/////////////    the function bellow is to get debug info. Don't use these!

func (afs *AppFS) ReadBlockUnsafe(a int) DiskLayer.Block {
	//return afs.blockFs.VD.ReadBlock(a)
	return nil
}
func (afs *AppFS) ReadSuperUnsafe() BlockLayer.SuperBlock {
	return BlockLayer.SuperBlock{}.FromBlocks(afs.blockFs.VD.ReadSuperBlock())
}

func (afs *AppFS) ReadInodeUnsafe(n int) BlockLayer.INode {
	return afs.blockFs.INodeN2iNode(n)
}

/*func (afs *AppFS) PrintLogUnsafe() {
	afs.fLog.PrintLog()
}*/
