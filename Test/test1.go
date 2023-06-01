package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
	"fmt"
)

func Test1() {
	fmt.Println("Test FS format:")
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(&DiskLayer.VirtualDisk{})
	printSuperBlock(testFS)
	printSingleBlock(testFS, 0)
	printSingleBlock(testFS, 1)
	//printBlock(testFS, 100)
	node1N := testFS.CreateFile(BlockLayer.NormalFile, "test1")

	testFS.LogCommit()
	//printBlock(testFS, 100)
	fmt.Println("\n\n\n\n\nAnd the new superBlock\n")
	printSuperBlock(testFS)
	printSingleBlock(testFS, 130)
	printSingleBlock(testFS, 131)
	printSingleBlock(testFS, 132)
	printSingleBlock(testFS, 133)
	fmt.Println("My dear world!")

	testFS.DeleteFile(node1N)
	testFS.LogCommit()
	printSuperBlock(testFS)
	printSingleBlock(testFS, 133)
	printSingleBlock(testFS, 134)
	printSingleBlock(testFS, 135)
	printSingleBlock(testFS, 136)
	printSingleBlock(testFS, 2)

	node2N := testFS.CreateFile(BlockLayer.NormalFile, "test2")
	//testFS.LogCommit()
	node3N := testFS.CreateFile(BlockLayer.NormalFile, "test3")
	testFS.LogCommit()
	testFileWrite(testFS)
	fmt.Println(node2N)
	fmt.Println(node3N)
}

func testFileWrite(afs AppFSLayer.AppFS) {
	indexL := []int{1, 2, 3, 6, 9}
	fmt.Println("-------------------Test write start----------------------------")
	inodeN := afs.CreateFile(BlockLayer.NormalFile, "TestWrite")
	afs.LogCommit()
	ds := []DiskLayer.Block{createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock()}
	afs.WriteFile(inodeN, indexL, ds)
	fmt.Println("-------------------Log after write-----------------------")
	//fmt.Println(afs.fLog.inodeByImap)
	afs.LogCommit()
	inode := afs.GetFileINfo(inodeN)
	fmt.Println("-------------------Inode after write-----------------------")
	fmt.Println(inode)
	fmt.Println("-------------------Data check-----------------------")
	printSingleBlock(afs, 145)
	printSingleBlock(afs, 146)
	printSingleBlock(afs, 147)
	printSingleBlock(afs, 148)
	for iD, v := range indexL {
		testF := true
		arr := BlockLayer.DataBlock{}.FromBlock(afs.ReadFile(inodeN, v)).(BlockLayer.DataBlock).Data
		for i := 0; i < Setting.BlockSize; i++ {
			if ds[iD].(BlockLayer.DataBlock).Data[i] != arr[i] {
				testF = false
			}
		}
		fmt.Println("Index:", v, "  test passed?: ", testF)
	}
	fmt.Println("-------------------Test read after other write----------------------------")
	ds2 := []DiskLayer.Block{createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock()}
	inodeN2 := afs.CreateFile(BlockLayer.NormalFile, "TestWrite")
	afs.WriteFile(inodeN2, indexL, ds2)
	afs.LogCommit()
	inode2 := afs.GetFileINfo(inodeN2)
	fmt.Println(inode2)
	afs.DeleteFile(inodeN2)
	afs.LogCommit()

	for iD, v := range indexL {
		testF := true
		arr := BlockLayer.DataBlock{}.FromBlock(afs.ReadFile(inodeN, v)).(BlockLayer.DataBlock).Data
		for i := 0; i < Setting.BlockSize; i++ {
			if ds[iD].(BlockLayer.DataBlock).Data[i] != arr[i] {
				testF = false
			}
		}
		fmt.Println("Index:", v, "  test passed?: ", testF)
	}

	fmt.Println("-------------------Test write done----------------------------")
}
