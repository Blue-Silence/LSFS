package Test

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/UserInterface/File"
	"LSF/UserInterface/Helper"
	"fmt"
)

func Test4() {
	testFS := AppFSLayer.AppFS{}
	testFS.FormatFS(&DiskLayer.VirtualDisk{})

	_, hRoot := File.GetHandler(&testFS, "/")
	_, hF1 := Helper.CreateByPath(&testFS, "/Folder1", BlockLayer.Folder)
	_, hF0 := Helper.CreateByPath(&testFS, "/Folder1/F11", BlockLayer.Folder)
	_, hF4 := Helper.CreateByPath(&testFS, "/Folder1/N12", BlockLayer.Folder)
	_, hF2 := Helper.CreateByPath(&testFS, "/Folder2", BlockLayer.Folder)
	_, hF3 := Helper.CreateByPath(&testFS, "/Folder1/F11/F111", BlockLayer.Folder)
	_, hN1 := Helper.CreateByPath(&testFS, "/Folder1/F11/F111/N1", BlockLayer.NormalFile)
	//File.Write()

	fmt.Println(File.GetInfo(&testFS, hRoot))
	fmt.Println(File.GetInfo(&testFS, hF1))
	fmt.Println(File.GetInfo(&testFS, hF0))
	fmt.Println(File.GetInfo(&testFS, hF4))
	fmt.Println(File.GetInfo(&testFS, hF2))
	fmt.Println(File.GetInfo(&testFS, hF3))
	fmt.Println(File.GetInfo(&testFS, hN1))

	fmt.Println("--------------------------------Let's see whats inside-------------------")

	fmt.Println(File.GetFolderContent(&testFS, hRoot))
	fmt.Println(File.GetFolderContent(&testFS, hF0))
	fmt.Println(File.GetFolderContent(&testFS, hF1))
	fmt.Println(File.GetFolderContent(&testFS, hF2))
	fmt.Println(File.GetFolderContent(&testFS, hF3))
	fmt.Println(File.GetFolderContent(&testFS, hF4))

	fmt.Println("--------------------------------Start dleleting files!-------------------")
	fmt.Println(Helper.DeleteByPath(&testFS, "/Folder1/F11"))
	fmt.Println(Helper.DeleteByPath(&testFS, "/Folder1/F11/F111/N1"))
	fmt.Println("Ok?")
	fmt.Println(File.GetInfo(&testFS, hF0))
	fmt.Println(File.GetInfo(&testFS, hN1))
	fmt.Println(File.GetInfo(&testFS, hF3))
	fmt.Println(File.GetInfo(&testFS, hRoot))
	fmt.Println(File.GetInfo(&testFS, hF1))
	fmt.Println("And the info:")
	fmt.Println(File.GetFolderContent(&testFS, hF1))
	fmt.Println(File.GetFolderContent(&testFS, hF3))
	fmt.Println(File.GetFolderContent(&testFS, hRoot))
	fmt.Println(File.GetFolderContent(&testFS, hF0))
	fmt.Println(File.GetFolderContent(&testFS, hF2))
	fmt.Println("Is every one that deserve the daeth truly dead?:")
	fmt.Println(File.GetInfo(&testFS, hF0))
	fmt.Println(File.GetInfo(&testFS, hF3))
	fmt.Println(File.GetInfo(&testFS, hN1))
}
