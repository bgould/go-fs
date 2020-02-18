package fat

import "io"

type File struct {
	chain *ClusterChain
	dir   *Directory
	entry *DirectoryClusterEntry

	readLen int
	readErr bool
}

func (f *File) Read(p []byte) (n int, err error) {

	fileSize := int(f.entry.fileSize)
	if f.readErr || f.readLen >= fileSize {
		return 0, io.EOF
	}
	n, err = f.chain.read(p)
	if err != nil {
		f.readErr = true
		return n, err
	}
	bytesLeft := fileSize - f.readLen
	f.readLen += n
	/*
		return n, err
	*/
	if bytesLeft > n {
		return n, nil
	} else {
		return bytesLeft, nil
	}

}

func (f *File) Write(p []byte) (n int, err error) {
	lastByte := f.chain.writeOffset + uint32(len(p))
	if lastByte > f.entry.fileSize {
		// Increase the file size since we're writing past the end of the file
		f.entry.fileSize = lastByte

		// Write the entry out
		if err := f.dir.dirCluster.WriteToDevice(f.dir.device, f.dir.fat); err != nil {
			return 0, err
		}
	}

	return f.chain.Write(p)
}

func (f *File) FileSize() uint32 {
	return f.entry.fileSize
}
