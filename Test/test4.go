package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/UserInterface/FileGeneral"
	"fmt"
)

func Test4() {
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(DiskLayer.VirtualDisk{})

	_, hRoot := FileGeneral.GetHandler(&testFS, "/")
	_, hF1 := FileGeneral.Create(&testFS, "/Folder1", BlockLayer.Folder)
	_, hF0 := FileGeneral.Create(&testFS, "/Folder1/F11", BlockLayer.Folder)
	_, hF4 := FileGeneral.Create(&testFS, "/Folder1/N12", BlockLayer.Folder)
	_, hF2 := FileGeneral.Create(&testFS, "/Folder2", BlockLayer.Folder)
	_, hF3 := FileGeneral.Create(&testFS, "/Folder1/F11/F111", BlockLayer.Folder)
	_, hN1 := FileGeneral.Create(&testFS, "/Folder1/F11/F111/N1", BlockLayer.NormalFile)
	//FileGeneral.Write()

	fmt.Println(FileGeneral.GetInfo(&testFS, hRoot))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF1))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF0))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF4))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF2))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF3))
	fmt.Println(FileGeneral.GetInfo(&testFS, hN1))

	fmt.Println("--------------------------------Let's see whats inside-------------------")

	fmt.Println(FileGeneral.GetFolderContent(&testFS, hRoot))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF0))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF1))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF2))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF3))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF4))

	fmt.Println("--------------------------------Start dleleting files!-------------------")
	fmt.Println(FileGeneral.Delete(&testFS, "/Folder1/F11"))
	fmt.Println(FileGeneral.Delete(&testFS, "/Folder1/F11/F111/N1"))
	fmt.Println("Ok?")
	fmt.Println(FileGeneral.GetInfo(&testFS, hF0))
	fmt.Println(FileGeneral.GetInfo(&testFS, hN1))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF3))
	fmt.Println(FileGeneral.GetInfo(&testFS, hRoot))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF1))
	fmt.Println("And the info:")
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF1))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF3))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hRoot))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF0))
	fmt.Println(FileGeneral.GetFolderContent(&testFS, hF2))
	fmt.Println("Is every one that deserve the daeth truly dead?:")
	fmt.Println(FileGeneral.GetInfo(&testFS, hF0))
	fmt.Println(FileGeneral.GetInfo(&testFS, hF3))
	fmt.Println(FileGeneral.GetInfo(&testFS, hN1))
}
