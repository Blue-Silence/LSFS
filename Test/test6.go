package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
	"LSF/UserInterface/FileGeneral"
	"fmt"
	//"UserInterface/FileGeneral"
)

func Test6() {
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(DiskLayer.VirtualDisk{})

	_, hRoot := FileGeneral.GetHandler(&testFS, "/")
	_, hF1 := FileGeneral.Create(&testFS, "/Folder1", BlockLayer.Folder)
	_, _ = FileGeneral.Create(&testFS, "/Folder1/F11", BlockLayer.Folder)
	_, _ = FileGeneral.Create(&testFS, "/Folder1/N12", BlockLayer.Folder)
	_, hF2 := FileGeneral.Create(&testFS, "/Folder2", BlockLayer.Folder)
	FileGeneral.Create(&testFS, "/Folder2/F21", BlockLayer.Folder)
	_, _ = FileGeneral.Create(&testFS, "/Folder1/F11/F111", BlockLayer.Folder)
	_, _ = FileGeneral.Create(&testFS, "/Folder1/F11/F111/N1", BlockLayer.NormalFile)
	_, hF211 := FileGeneral.Create(&testFS, "/Folder2/F21/F211", BlockLayer.Folder)
	_, _ = FileGeneral.Create(&testFS, "/Folder2/F21/F211/N2", BlockLayer.NormalFile)
	errN3, hN3 := FileGeneral.Create(&testFS, "/Folder2/F21/F211/N3", BlockLayer.NormalFile)
	FileGeneral.Create(&testFS, "/Folder2/F21/F211/N4", BlockLayer.NormalFile)

	//dataB := createDataBlock().Data
	indexL := []int{1, 2, 3, 6, 9}
	ds := []BlockLayer.DataBlock{createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock()}
	fmt.Println("ErrN3:", errN3)
	fmt.Println(FileGeneral.GetInfo(&testFS, hN3))
	for in, i := range indexL {
		fmt.Println("Err:", FileGeneral.Write(&testFS, hN3, i, ds[in].Data), "   $")
	}

	fmt.Println(FileGeneral.Delete(&testFS, "/Folder1/F11"))
	testFS.GC(-1)
	fmt.Println(testFS.ReadInodeUnsafe(FileGeneral.GetIUnsafe(hF1)))
	printSingleBlock(testFS, 254)
	printSingleBlock(testFS, 255)
	printSingleBlock(testFS, 256)
	printSingleBlock(testFS, 257)
	printSingleBlock(testFS, 258)
	printSingleBlock(testFS, 259)
	printSingleBlock(testFS, 260)

	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF1))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF2))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hRoot))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF211))
	//fmt.Println(FileGeneral.GetFolderContent(&testFS, hF2))

	for iD, v := range indexL {
		testF := true
		err, arr := FileGeneral.Read(&testFS, hN3, v)
		fmt.Println("Err:", err)
		for i := 0; i < Setting.BlockSize; i++ {
			if ds[iD].Data[i] != arr[i] {
				testF = false
			}
		}
		fmt.Println("Index:", v, "  test passed?: ", testF)
	}
}
