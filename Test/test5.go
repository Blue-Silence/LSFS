package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/MemoryDisk"
	"LSF/Setting"
	"LSF/UserInterface/File"
	"LSF/UserInterface/Helper"
	"fmt"
	//"UserInterface/FileGeneral"
)

func Test5() {
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(&MemoryDisk.RamDisk{})

	_, hRoot := File.GetHandler(&testFS, "/")
	_, hF1 := Helper.CreateByPath(&testFS, "/Folder1", BlockLayer.Folder)
	_, _ = Helper.CreateByPath(&testFS, "/Folder1/F11", BlockLayer.Folder)
	_, _ = Helper.CreateByPath(&testFS, "/Folder1/N12", BlockLayer.Folder)
	_, hF2 := Helper.CreateByPath(&testFS, "/Folder2", BlockLayer.Folder)
	Helper.CreateByPath(&testFS, "/Folder2/F21", BlockLayer.Folder)
	_, _ = Helper.CreateByPath(&testFS, "/Folder1/F11/F111", BlockLayer.Folder)
	_, _ = Helper.CreateByPath(&testFS, "/Folder1/F11/F111/N1", BlockLayer.NormalFile)
	_, hF211 := Helper.CreateByPath(&testFS, "/Folder2/F21/F211", BlockLayer.Folder)
	_, _ = Helper.CreateByPath(&testFS, "/Folder2/F21/F211/N2", BlockLayer.NormalFile)
	errN3, hN3 := Helper.CreateByPath(&testFS, "/Folder2/F21/F211/N3", BlockLayer.NormalFile)
	Helper.CreateByPath(&testFS, "/Folder2/F21/F211/N4", BlockLayer.NormalFile)

	//dataB := createDataBlock().Data
	indexL := []int{1, 2, 3, 6, 9}
	ds := []BlockLayer.DataBlock{createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock()}
	fmt.Println("ErrN3:", errN3)
	fmt.Println(File.GetInfo(&testFS, hN3))
	for in, i := range indexL {
		fmt.Println("Err:", File.Write(&testFS, hN3, i, ds[in].Data), "   $")
	}

	fmt.Println(Helper.DeleteByPath(&testFS, "/Folder1/F11"))
	fmt.Println(File.GetFolderContent(&testFS, hF1))
	fmt.Println(File.GetFolderContent(&testFS, hF2))
	fmt.Println(File.GetFolderContent(&testFS, hRoot))
	fmt.Println(File.GetFolderContent(&testFS, hF211))
	//fmt.Println(File.GetFolderContent(&testFS, hF2))

	for iD, v := range indexL {
		testF := true
		err, arr := File.Read(&testFS, hN3, v)
		fmt.Println("Err:", err)
		for i := 0; i < Setting.BlockSize; i++ {
			if ds[iD].Data[i] != arr[i] {
				testF = false
			}
		}
		fmt.Println("Index:", v, "  test passed?: ", testF)
	}
}
