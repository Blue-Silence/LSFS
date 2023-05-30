package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/UserInterface/Folder"
	"fmt"
)

func Test2() {
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(DiskLayer.VirtualDisk{})

	fmt.Println("-----------------------Before adding file.----------------------------------")
	fmt.Println(Folder.GetFolderContent(&testFS, 0))
	fmt.Println("-----------------------After adding 3 file.----------------------------------")
	iNF1 := testFS.CreateFile(BlockLayer.NormalFile, "Normal file 1")
	iF2 := testFS.CreateFile(BlockLayer.NormalFile, "Folder 2")
	iNF3 := testFS.CreateFile(BlockLayer.NormalFile, "Normal file 3")
	testFS.LogCommit()

	Folder.AddFileToFolder(&testFS, 0, iNF1)
	Folder.AddFileToFolder(&testFS, 0, iF2)
	Folder.AddFileToFolder(&testFS, 0, iNF3)
	testFS.LogCommit()
	fmt.Println(Folder.GetFolderContent(&testFS, 0))

	iF4 := testFS.CreateFile(BlockLayer.NormalFile, "Folder 4")
	testFS.LogCommit()
	Folder.AddFileToFolder(&testFS, iF2, iF4)
	testFS.LogCommit()

	fmt.Println("Root:", Folder.GetFolderContent(&testFS, 0))
	fmt.Println("Secondary:", Folder.GetFolderContent(&testFS, iF2))

	inFs := []int{}
	for i := 4; i < 73; i++ {
		inFs = append(inFs, testFS.CreateFile(BlockLayer.NormalFile, (fmt.Sprint("-Normalfile-", i))))
	}
	testFS.LogCommit()

	fmt.Println("-----------------------Massive test.----------------------------------")
	for _, v := range inFs {
		Folder.AddFileToFolder(&testFS, iF2, v)
		//testFS.LogCommit()
		//fmt.Println("Secondary:", Folder.ConcatFolderUnsafe(&testFS, iF2)) //Folder.GetFolderContent(&testFS, iF2))
		//fmt.Println(i, ": Pass")
	}
	testFS.LogCommit()
	fmt.Println("Secondary:", Folder.GetFolderContent(&testFS, iF2))
	fmt.Println("-----------------------Delete test.----------------------------------")
	//fmt.Println("Long Before apply:", testFS.GetFileINfo(iF2).Pointers)

	Folder.DeleteFileToFolder(&testFS, iF2, inFs[1])
	//fmt.Println("Secondary:", Folder.GetFolderContent(&testFS, iF2))

	//Folder.DeleteFileToFolder(&testFS, iF2, inFs[10])
	//Folder.DeleteFileToFolder(&testFS, iF2, inFs[39])

	//fmt.Println("Before apply:", testFS.GetFileINfo(iF2).Pointers)
	testFS.LogCommit()
	//fmt.Println("After apply:", testFS.GetFileINfo(iF2).Pointers)
	/*printSingleBlock(testFS, 1135)
	printSingleBlock(testFS, 1136)
	printSingleBlock(testFS, 1137)
	printSingleBlock(testFS, 1138)*/
	fmt.Println("Tag1:", inFs[1])

	fmt.Println("Secondary:", Folder.GetFolderContent(&testFS, iF2))
}
